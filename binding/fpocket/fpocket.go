package fpocket

import (
	"errors"
	"fmt"
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
	Residues  []*pdb.Residue // pointers to original residues in the requested structure
}

// Run creates a temp file of the specified PDB structure, runs Fpocket on it and parses the results
func Run(crystal *pdb.PDB) (pockets []*Pocket, err error) {
	outPath := "/tmp/" + crystal.LocalFilename + "_out"

	// Delete Fpocket result files on function exit
	defer func() {
		os.RemoveAll(outPath)
	}()

	// Run Fpocket
	out, err := exec.Command("fpocket", "-f", crystal.LocalPath).CombinedOutput()
	if err != nil {
		return nil, err
	}
	if strings.Contains(string(out), "failed") {
		fmt.Println(string(out))
		return nil, errors.New("FPocket failed")
	}

	// Walk created folder containing pocket analysis files
	dir := outPath + "/pockets"
	pockets, err = walkPocketDir(crystal, dir)
	if err != nil {
		return nil, err
	}

	return pockets, nil
}

func walkPocketDir(crystal *pdb.PDB, dir string) (pockets []*Pocket, err error) {
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
				pocketPDB, err := pdb.NewPDBNoMetadata(data)
				if err != nil {
					return err
				}

				var pocketResidues []*pdb.Residue
				for chain, chainPos := range pocketPDB.Chains {
					for pos := range chainPos {
						pocketResidues = append(pocketResidues, crystal.Chains[chain][pos])
					}
				}

				pocket := &Pocket{
					Name:      file,
					DrugScore: drugScore,
					Residues:  pocketResidues,
				}
				pockets = append(pockets, pocket)
			}

		}
		return nil
	})

	return pockets, err
}
