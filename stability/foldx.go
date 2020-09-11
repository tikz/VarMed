package stability

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"respdb/glyco"
	"respdb/pdb"
	"respdb/sasa"
	"respdb/uniprot"
	"strconv"
	"strings"
)

// Mutation represents the parameters between an original
// and a mutated structure with a single aminoacid substitution.
type Mutation struct {
	SAS         *uniprot.SAS     `json:"sas"`
	DdG         float64          `json:"ddG"` // kcal/mol
	SASA        float64          `json:"sasa"`
	SASAApolar  float64          `json:"sasa_apolar"`
	SASAPolar   float64          `json:"sasa_polar"`
	SASAUnknown float64          `json:"sasa_unknown"`
	GlycoDist   *glyco.GlycoDist `json:"glyco_dist"`
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
	residues := p.UniProtPositions[unpID][int64(pos)]
	if len(residues) == 0 {
		return "", errors.New("no coverage")
	}
	res := residues[0]
	return res.Name1 + res.Chain + strconv.FormatInt(res.StructPosition, 10) + aa + ";", nil
}

func FoldXRun(sasList []*uniprot.SAS, unpID string,
	p *pdb.PDB, msg func(string)) ([]*Mutation, error) {
	pdbPath, err := repair(p, msg)
	if err != nil {
		return nil, fmt.Errorf("repair: %v", err)
	}

	// // Wild type SASAs
	// wtS, wtSA, _, err := sasa.SASA("bin/" + p.ID + ".pdb")

	var results []*Mutation
	for _, sas := range sasList {
		mut, err := formatMutant(unpID, p, sas.Position, sas.ToAa)
		if err == nil {
			mutation, err := buildModel(p.ID, pdbPath, sas, mut)
			if err != nil {
				return nil, err
			}
			results = append(results, mutation)
		}
	}

	return results, nil
}

func buildModel(pdbID string, pdbPath string, sas *uniprot.SAS, mut string) (*Mutation, error) {
	mutation := &Mutation{SAS: sas}

	change := sas.FromAa + strconv.FormatInt(sas.Position, 10) + sas.ToAa
	destDirPath := "data/foldx/mutations/" + pdbID + "/" + change
	diffPath := destDirPath + "/Dif_" + pdbID + "_Repair.fxout"
	mutatedPDBPath := destDirPath + "/" + pdbID + "_Repair_1.pdb"

	ddG, err := extractddG(diffPath)
	if fileNotExist(diffPath) || fileNotExist(mutatedPDBPath) || err != nil {
		// Create FoldX job output dir
		os.MkdirAll(destDirPath, os.ModePerm)

		// Create file containing individual list of mutations
		mutantFile := "individual_list_" + pdbID + change
		mutantPath := "bin/" + mutantFile
		writeFile(mutantPath, mut)

		// Create hardlink
		// (FoldX seems to only look for PDBs in the same dir)
		linkPath := "bin/" + pdbID + "_Repair.pdb"
		os.Link(pdbPath, linkPath)

		// Remove files on scope exit
		defer func() {
			// hardlink
			os.RemoveAll(linkPath)
			os.RemoveAll(mutantPath)

			// Duplicate of the original repaired PDB copied to mutation folder by FoldX
			os.RemoveAll(destDirPath + "/WT_" + pdbID + "_Repair_1.pdb")

			// Mutated PDB
			// os.RemoveAll(destDirPath + "/" + pdbID + "_Repair_1.pdb")
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

	ddG, err = extractddG(diffPath)
	if err != nil {
		return nil, fmt.Errorf("extract results: %v", err)
	}

	mutation.DdG = ddG

	// SASA of mutated structure
	total, apolar, polar, unk, err := sasa.SASA(destDirPath + "/" + pdbID + "_Repair_1.pdb")
	if err != nil {
		return nil, fmt.Errorf("mutated SASA: %v", err)
	}
	mutation.SASA = total
	mutation.SASAApolar = apolar
	mutation.SASAPolar = polar
	mutation.SASAUnknown = unk

	return mutation, nil
}

func extractddG(path string) (ddG float64, err error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return ddG, err
	}

	r, _ := regexp.Compile("pdb\t(.*?)\t")
	m := r.FindAllStringSubmatch(string(data), -1)
	if len(m) == 0 {
		return ddG, errors.New("ddG not found")
	}
	ddG, err = strconv.ParseFloat(m[0][1], 64)
	return ddG, err
}
