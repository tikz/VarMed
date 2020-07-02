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
	SAS *uniprot.SAS `json:"sas"`
	DdG float64      `json:"ddG"` // kcal/mol
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

// repair runs the RepairPDB FoldX command on a PDB file specified by pdb.LocalPath
// and stores the resulting file in foldxDir, only if the file doesn't exist already.
// Returns the path where the repaired PDB is located.
func repair(p *pdb.PDB, msg func(string)) (string, error) {
	repairDir := "data/foldx/repair/"
	path := repairDir + p.ID + "_Repair.pdb"
	if fileNotExist(path) {
		msg("running FoldX RepairPDB")

		cmd := exec.Command("./foldx",
			"--command=RepairPDB",
			"--pdb="+p.ID+".pdb",
			"--output-dir=../"+repairDir)
		cmd.Dir = "bin/"

		out, err := cmd.CombinedOutput()
		if err != nil {
			return path, err
		}

		if !strings.Contains(string(out), "run OK") || fileNotExist(path) {
			fmt.Println(string(out))
			return path, errors.New("RepairPDB failed")
		}
	} else {
		msg("found existing FoldX PDB")
	}

	return path, nil
}

// formatMutant receives a given mutation in UniProt position and returns
// the corresponding PDB positions in FoldX format, i.e.: KA42I,KB42I;
func formatMutant(unpID string, p *pdb.PDB, pos int64, aa string) (string, error) {
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
	p *pdb.PDB, msg func(string)) ([]*SASDiff, error) {
	pdbPath, err := repair(p, msg)
	if err != nil {
		return nil, fmt.Errorf("repair: %v", err)
	}

	var results []*SASDiff
	for _, sas := range sasList {
		mut, err := formatMutant(unpID, p, sas.Position, sas.ToAa)
		if err == nil {
			diff, err := buildModel(p.ID, pdbPath, sas, mut)
			if err != nil {
				return nil, err
			}
			results = append(results, diff)
		}
	}

	return results, nil
}

func buildModel(pdbID string, pdbPath string, sas *uniprot.SAS, mut string) (*SASDiff, error) {
	diff := &SASDiff{SAS: sas}

	change := sas.FromAa + strconv.FormatInt(sas.Position, 10) + sas.ToAa
	destDirPath := "data/foldx/mutations/" + pdbID + "/" + change
	diffPath := destDirPath + "/Dif_" + pdbID + "_Repair.fxout"

	if fileNotExist(diffPath) {
		// Create FoldX job output dir
		os.MkdirAll(destDirPath, os.ModePerm)

		// Create file containing individual list of mutations
		mutantFile := "individual_list_" + pdbID + change
		writeFile("bin/"+mutantFile, mut)

		// Create hardlink
		// (FoldX seems to only look for PDBs in the same dir)
		os.Link(pdbPath, "bin/"+pdbID+"_Repair.pdb")

		// Remove files on scope exit
		defer func() {
			os.RemoveAll("bin/" + pdbPath + "_Repair.pdb") // hardlink
			os.RemoveAll("bin/" + pdbID)
			os.RemoveAll("bin/" + mutantFile)
		}()

		cmd := exec.Command("./foldx",
			"--command=BuildModel",
			"--pdb="+pdbID+"_Repair.pdb",
			"--mutant-file="+mutantFile,
			"--output-dir=../"+destDirPath)
		cmd.Dir = "bin/"

		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println(string(out))
			return nil, err
		}

		if !strings.Contains(string(out), "run OK") {
			fmt.Println(string(out))
			return nil, errors.New("BuildModel failed")
		}
	}

	ddG, err := extractddG(diffPath)
	if err != nil {
		return nil, fmt.Errorf("extract results: %v", err)
	}
	diff.DdG = ddG

	return diff, nil
}

func extractddG(path string) (ddG float64, err error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return ddG, err
	}

	r, _ := regexp.Compile("pdb\t(.*?)\t")
	ddG, err = strconv.ParseFloat(r.FindAllStringSubmatch(string(data), -1)[0][1], 64)
	return ddG, err
}
