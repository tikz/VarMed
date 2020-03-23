package exposure

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"
	"varq/pdb"
)

type ExposureAnalysis struct {
	Chains   map[string]map[int64]*AminoacidExposure
	Duration time.Duration
	Error    error
}
type AminoacidExposure struct {
	BFactor   float64
	ExposureP float64
	Aminoacid *pdb.Aminoacid `json:"-"`
}

type PyMOLResults struct {
	Error error
	Lines []string
}

func RunExposureAnalysis(pdb *pdb.PDB, results chan<- *ExposureAnalysis) {
	start := time.Now()
	Run(pdb)
	chains, err := Run(pdb)
	if err != nil {
		results <- &ExposureAnalysis{Error: fmt.Errorf("running PyMOL: %v", err)}
	}

	results <- &ExposureAnalysis{Chains: chains, Duration: time.Since(start)}
}

// Run creates a temp file of the specified PDB structure, runs the PyMOL script on it and parses the results
func Run(crystal *pdb.PDB) (map[string]map[int64]*AminoacidExposure, error) {
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
		pos := [2]int{j * chunkSize, (j + 1) * chunkSize}
		jobs <- pos
	}
	close(jobs)

	exposureChains := make(map[string]map[int64]*AminoacidExposure)
	for a := 0; a < totalJobs; a++ {
		res := <-results

		for _, line := range res.Lines {
			cols := strings.Split(line, " ")
			if len(cols) > 1 {
				chain, pos, bFactor, exposureP, err := parseLine(cols)
				if err != nil {
					return nil, fmt.Errorf("parsing line: %v", err)
				}

				if _, ok := exposureChains[chain]; !ok {
					exposureChains[chain] = make(map[int64]*AminoacidExposure)
				}
				exposureChains[chain][pos] = &AminoacidExposure{
					BFactor:   bFactor,
					ExposureP: exposureP,
					Aminoacid: crystal.Chains[chain][pos],
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
		out, _ := exec.Command("python3", "exposure/run_pymol.py", path, start, end).CombinedOutput()

		res := PyMOLResults{Error: nil, Lines: strings.Split(string(out), "\n")}
		results <- &res
	}
}

func totalLength(chains map[string]map[int64]*pdb.Aminoacid) (totalLength int) {
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
