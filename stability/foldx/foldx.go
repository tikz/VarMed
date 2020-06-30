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

// SASDiff represents the ddG energy difference between an original
// and a mutated structure with a single aminoacid substitution.
type SASDiff struct {
	SAS         *uniprot.SAS `json:"sas"`
	InStructure bool         `json:"inStructure"`
	DdG         float64      `json:"ddG"` // kcal/mol
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

// formatMutants receives a given mutation in UniProt position and returns
// the corresponding PDB positions in FoldX format, i.e.: KA42I,KB42I;
func formatMutants(unpID string, p *pdb.PDB, pos int64, aa string) (string, error) {
	var muts []string
	residues := p.UniProtPositions[unpID][int64(pos)]
	if len(residues) == 0 {
		return "", errors.New("no coverage")
	}

	for _, res := range residues {
		muts = append(muts, res.Name1+res.Chain+strconv.FormatInt(pos, 10)+aa)
	}
	return strings.Join(muts, ",") + ";", nil
}

func Run(sasList []*uniprot.SAS, unpID string,
	p *pdb.PDB, foldxDir string, msg func(string)) ([]*SASDiff, error) {
	err := Repair(p, foldxDir, msg)
	if err != nil {
		return nil, fmt.Errorf("repair: %v", err)
	}

	var results []*SASDiff
	for _, sas := range sasList {
		diff := &SASDiff{SAS: sas}
		results = append(results, diff)

		pos := sas.Position
		aa := sas.ToAa
		name := p.ID + strconv.FormatInt(pos, 10) + aa

		muts, err := formatMutants(unpID, p, pos, aa)
		if err != nil {
			continue
		}
		diff.InStructure = true

		// Create FoldX job output dir
		os.MkdirAll("bin/"+name, os.ModePerm)

		// Create file containing individual list of mutations
		mutantFile := "individual_list_" + name
		writeFile("bin/"+mutantFile, muts)

		// Create hardlink
		// (FoldX seems to only look for PDBs in the same dir)
		os.Link(foldxDir+p.ID+"_Repair.pdb", "bin/"+p.ID+"_Repair.pdb")

		// Remove job files on scope exit
		defer func() {
			os.RemoveAll("bin/" + name)
			os.RemoveAll("bin/" + mutantFile)
			os.RemoveAll("bin/" + p.ID + "_Repair.pdb")
		}()

		cmd := exec.Command("./foldx",
			"--command=BuildModel",
			"--pdb="+p.ID+"_Repair.pdb",
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

		ddG, err := extractddG("bin/" + name + "/Dif_" + p.ID + "_Repair.fxout")
		if err != nil {
			return nil, fmt.Errorf("extract results: %v", err)
		}

		diff.DdG = ddG
	}

	return results, nil
}

func extractddG(path string) (ddG float64, err error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return ddG, err
	}

	// First column contains ddG
	// https://regex101.com/r/BGcps6/1
	r, _ := regexp.Compile("pdb\t(.*?)\t")
	ddG, err = strconv.ParseFloat(r.FindAllStringSubmatch(string(data), -1)[0][1], 64)
	return ddG, err
}
