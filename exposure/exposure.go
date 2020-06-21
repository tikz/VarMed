package exposure

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"
	"varq/pdb"
)

type Results struct {
	Residues []*pdb.Residue `json:"residues"`
	Duration time.Duration  `json:"duration"`
	Error    error          `json:"error"`
}
type ResidueExposure struct {
	BFactor   float64      `json:"bFactor"`
	ExposureP float64      `json:"exposureP"`
	Residue   *pdb.Residue `json:"-"`
}

type PyMOLResults struct {
	Error error
	Lines []string
}

func Run(pdb *pdb.PDB, results chan<- *Results, msg func(string)) {
	start := time.Now()
	RunPyMOL(pdb, msg)
	chains, err := RunPyMOL(pdb, msg)
	if err != nil {
		results <- &Results{Error: fmt.Errorf("running PyMOL: %v", err)}
	}

	buried := buriedResidues(pdb, chains)
	results <- &Results{Residues: buried, Duration: time.Since(start)}
}

// RunPyMOL creates a temp file of the specified PDB structure, runs the PyMOL script on it and parses the results
func RunPyMOL(crystal *pdb.PDB, msg func(string)) (map[string]map[int64]*ResidueExposure, error) {
	pymolWorkers := 16

	length := totalLength(crystal.Chains)
	chunkSize := length / pymolWorkers

	var totalJobs int
	if chunkSize == 0 { // more workers than length
		chunkSize = 1
		totalJobs = length
	} else {
		if length%chunkSize != 0 {
			totalJobs++
		}
		totalJobs += length / chunkSize
	}

	jobs := make(chan [2]int, totalJobs)
	results := make(chan *PyMOLResults, totalJobs)

	for w := 0; w < pymolWorkers; w++ {
		go pymolWorker(crystal.LocalPath, jobs, results)
	}

	for j := 0; j < totalJobs; j++ {
		start := j * chunkSize
		end := (j + 1) * chunkSize
		msg(fmt.Sprintf("do SASA for residues %d-%d", start, end))
		pos := [2]int{start, end}
		jobs <- pos
	}
	close(jobs)

	exposureChains := make(map[string]map[int64]*ResidueExposure)
	for a := 0; a < totalJobs; a++ {
		res := <-results
		msg("SASA worker done")

		for _, line := range res.Lines {
			cols := strings.Split(line, " ")
			if len(cols) > 1 {
				chain, pos, bFactor, exposureP, err := parseLine(cols)
				if err != nil {
					return nil, fmt.Errorf("parsing line: %v", err)
				}

				if _, ok := exposureChains[chain]; !ok {
					exposureChains[chain] = make(map[int64]*ResidueExposure)
				}
				if crystal.Chains[chain][pos] != nil {
					exposureChains[chain][pos] = &ResidueExposure{
						BFactor:   bFactor,
						ExposureP: exposureP,
						Residue:   crystal.Chains[chain][pos],
					}
				}
			}

		}
	}

	return exposureChains, nil
}

func pymolWorker(path string, jobs <-chan [2]int, results chan<- *PyMOLResults) {
	for j := range jobs {
		start := strconv.Itoa(j[0])
		end := strconv.Itoa(j[1])
		out, _ := exec.Command("../pymol/bin/python3", "exposure/run_pymol.py", path, start, end).CombinedOutput()

		res := PyMOLResults{Error: nil, Lines: strings.Split(string(out), "\n")}
		results <- &res
	}
}

func totalLength(chains map[string]map[int64]*pdb.Residue) (totalLength int) {
	for _, chain := range chains {
		totalLength += len(chain)
	}
	return totalLength
}

func parseLine(line []string) (chain string, pos int64, bFactor float64, exposureP float64, err error) {
	chain = line[0]
	pos, err = strconv.ParseInt(line[1], 10, 64)
	bFactor, err = strconv.ParseFloat(line[2], 64)
	exposureP, err = strconv.ParseFloat(line[3], 64)

	return chain, pos, bFactor, exposureP, err
}

func buriedResidues(pdb *pdb.PDB, chains map[string]map[int64]*ResidueExposure) (buriedResidues []*pdb.Residue) {
	for chainName, chain := range chains {
		for pos, resExp := range chain {
			if resExp.ExposureP < 0.5 && resExp.Residue.Name1 != "G" {
				buriedResidues = append(buriedResidues, pdb.Chains[chainName][pos])
			}
		}
	}
	return buriedResidues
}
