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

type Pocket struct {
	Name      string
	DrugScore float64
	Chains    map[string]map[int64]*pdb.Aminoacid
}

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

	// Run FPocket
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
		// For each pocket PDB file
		if strings.Contains(file, "_atm.pdb") {
			data, err := ioutil.ReadFile(file)
			if err != nil {
				return err
			}
			// Extract drug score from FPocket result PDB
			regexScore, _ := regexp.Compile("Drug Score.*: ([0-9.]*)")
			drugScore, err := strconv.ParseFloat(regexScore.FindAllStringSubmatch(string(data), -1)[0][1], 64)
			if err != nil {
				return err
			}

			// Druggability score threshold as VarQ spec
			if drugScore > 0.00001 {
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
