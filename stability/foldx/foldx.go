package foldx

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"varq/pdb"
	"varq/uniprot"
)

// SASEnergyDiff represents the ddG energy difference between an original
// and a mutated structure with a single aminoacid substitution.
type SASEnergyDiff struct {
	SAS *uniprot.SAS
	ddG float64 // kcal/mol
}

func fileNotExist(path string) bool {
	_, err := os.Stat(path)
	return os.IsNotExist(err)
}

func writeFile(path string, contents string) {
	err := ioutil.WriteFile(path, []byte(contents), 0644)
	if err != nil {
		panic(err)
	}
}

// Repair runs the RepairPDB FoldX command on a PDB file specified by pdb.LocalPath
// and stores the resulting file in foldxDir, only if the file doesn't exist already.
func Repair(p *pdb.PDB, foldxDir string, msg func(string)) error {
	path := foldxDir + p.ID + "_Repair.pdb"
	if fileNotExist(path) {
		msg("running FoldX RepairPDB")

		cmd := exec.Command("./foldx",
			"--command=RepairPDB",
			"--pdb="+p.ID+".pdb",
			"--output-dir=../"+foldxDir)
		cmd.Dir = "bin/"

		out, err := cmd.CombinedOutput()
		if err != nil {
			return err
		}

		if !strings.Contains(string(out), "run OK") || fileNotExist(path) {
			fmt.Println(string(out))
			return errors.New("RepairPDB failed")
		}
	} else {
		msg("found existing FoldX PDB")
	}

	return nil
}

// formatMut receives a given mutation in UniProt position and returns
// the corresponding PDB positions in FoldX format, i.e.: KA42I,KB42I;
func formatMut(unpID string, p *pdb.PDB, pos int64, aa string) string {
	var muts []string
	for _, res := range p.UniProtPositions[unpID][int64(pos)] {
		muts = append(muts, res.Name1+res.Chain+strconv.FormatInt(pos, 64)+aa)
	}
	return strings.Join(muts, ",") + ";"
}

func Run(sasList []*uniprot.SAS, unpID string,
	p *pdb.PDB, foldxDir string, msg func(string)) ([]*SASEnergyDiff, error) {
	err := Repair(p, foldxDir, msg)
	if err != nil {
		return nil, fmt.Errorf("repair: %v", err)
	}

	var results []*SASEnergyDiff
	for _, sas := range sasList {
		pos := sas.Position
		aa := sas.ToAa
		name := p.ID + strconv.FormatInt(pos, 64) + aa

		// Create FoldX job output dir
		os.MkdirAll("bin/"+name, os.ModePerm)

		// Create file containing individual list of mutations
		mutantFile := "individual_list_" + name
		writeFile("bin/"+mutantFile, formatMut(unpID, p, pos, aa))

		// Remove job files on scope exit
		defer func() {
			os.RemoveAll("bin/" + name)
			os.RemoveAll("bin/" + mutantFile)
		}()

		cmd := exec.Command("./foldx",
			"--command=BuildModel",
			"--pdb="+p.ID+".pdb",
			"--mutant-file="+mutantFile,
			"--output-dir="+name)
		cmd.Dir = "bin/"

		out, err := cmd.CombinedOutput()
		if err != nil {
			return nil, err
		}

		if !strings.Contains(string(out), "run OK") {
			fmt.Println(string(out))
			return nil, errors.New("BuildModel failed")
		}

		ddG, err := extractddG("bin/" + name + "/Dif_" + p.ID + ".fxout")
		if err != nil {
			return nil, fmt.Errorf("extract results: %v", err)
		}

		results = append(results, &SASEnergyDiff{SAS: sas, ddG: ddG})
		fmt.Println(p.ID, ddG)
	}

	return results, nil
}

func extractddG(path string) (ddG float64, err error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return ddG, err
	}

	// First column contains the ddG
	// https://regex101.com/r/BGcps6/1
	r, _ := regexp.Compile("pdb\t(.*?)\t")
	ddG, err = strconv.ParseFloat(r.FindAllStringSubmatch(string(data), -1)[0][1], 64)
	return ddG, err
}
