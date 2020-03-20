package fpocket

import (
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"varq/pdb"
)

//Pocket represents a single pocket found with Fpocket (pocketN_atm.pdb)
type Pocket struct {
	Name      string
	DrugScore float64
	Chains    map[string]map[int64]*pdb.Aminoacid
}

// Run creates a temp file of the specified PDB structure, runs Fpocket on it and parses the results
func Run(crystal *pdb.PDB) (pockets []*Pocket, err error) {
	// Create tmp file with the PDB data
	fileName := "varq_" + crystal.ID
	path := "/tmp/" + fileName + ".pdb"
	outPath := "/tmp/" + fileName + "_out"

	// Delete tmp files on function exit
	defer func() {
		os.Remove(path)
		os.RemoveAll(outPath)
	}()

	err = ioutil.WriteFile(path, crystal.RawPDB, 0644)
	if err != nil {
		return nil, err
	}

	// Run Fpocket
	out, err := exec.Command("fpocket", "-f", path).CombinedOutput()
	if err != nil {
		return nil, err
	}
	if strings.Contains(string(out), "failed") {
		return nil, errors.New("FPocket failed")
	}

	// Walk created folder containing pocket analysis files
	dir := outPath + "/pockets"
	pockets, err = walkPocketDir(dir)
	if err != nil {
		return nil, err
	}

	return pockets, nil
}

func walkPocketDir(dir string) (pockets []*Pocket, err error) {
	err = filepath.Walk(dir, func(file string, info os.FileInfo, err error) error {
		// For each Fpocket result PDB file
		if strings.Contains(file, "_atm.pdb") {
			data, err := ioutil.ReadFile(file)
			if err != nil {
				return err
			}
			// Extract drug score
			regexScore, _ := regexp.Compile("Drug Score.*: ([0-9.]*)")
			drugScore, err := strconv.ParseFloat(regexScore.FindAllStringSubmatch(string(data), -1)[0][1], 64)
			if err != nil {
				return err
			}

			// Druggability score threshold as VarQ spec
			if drugScore > 0.5 {
				// then load PDB
				pocketPDB := &pdb.PDB{RawPDB: data}
				err := pocketPDB.ExtractChains()
				if err != nil {
					return err
				}
				pocket := &Pocket{
					Name:      file,
					DrugScore: drugScore,
					Chains:    pocketPDB.Chains,
				}
				pockets = append(pockets, pocket)
			}

		}
		return nil
	})

	return pockets, err
}
