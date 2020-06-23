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
)

type Variant struct {
	Position int64   `json:"position"`
	FromAa   string  `json:"fromAa"`
	ToAa     string  `json:"toAa"`
	DeltadG  float64 `json:"ddG"`
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
			return errors.New("FoldX RepairPDB failed")
		}
	} else {
		msg("found existing FoldX PDB")
	}

	return nil
}

// formatMut receives a given mutation in UniProt position and returns
// the corresponding PDB positions in FoldX format, i.e.: KA42I,KB42I;
func formatMut(unpID string, p *pdb.PDB, pos int, aa string) string {
	var muts []string
	for _, res := range p.UniProtPositions[unpID][int64(pos)] {
		muts = append(muts, res.Name1+res.Chain+strconv.Itoa(pos)+aa)
	}
	return strings.Join(muts, ",") + ";"
}

func Run(variants map[int]string, unpID string,
	p *pdb.PDB, foldxDir string, msg func(string)) error {
	err := Repair(p, foldxDir, msg)
	if err != nil {
		return fmt.Errorf("FoldX repair: %v", err)
	}

	variants[108] = "G"

	for pos, aa := range variants {
		name := p.ID + strconv.Itoa(pos) + aa

		// Create FoldX job output dir
		os.MkdirAll("bin/"+name, os.ModePerm)

		// Create file containing individual list of mutations
		mutantFile := "individual_list_" + name
		writeFile("bin/"+mutantFile, formatMut(unpID, p, pos, aa))

		// Remove job files on scope exit
		defer func() {
			// os.RemoveAll("bin/" + name)
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
			return err
		}

		if !strings.Contains(string(out), "run OK") {
			fmt.Println(string(out))
			return errors.New("FoldX BuildModel failed")
		}

		ddG, err := extractddG("bin/" + name + "/Dif_" + p.ID + ".fxout")
		fmt.Println(p.ID, ddG)
	}

	return nil
}

func extractddG(path string) (ddG float64, err error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return ddG, err
	}

	// Parse first column
	// https://regex101.com/r/BGcps6/1
	r, _ := regexp.Compile("pdb\t(.*?)\t")
	ddG, err = strconv.ParseFloat(r.FindAllStringSubmatch(string(data), -1)[0][1], 64)
	return ddG, err
}
