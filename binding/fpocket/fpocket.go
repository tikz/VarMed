package fpocket

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
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
	Name      string         `json:"name"`
	DrugScore float64        `json:"drugScore"`
	Residues  []*pdb.Residue `json:"residues"` // pointers to original residues in the requested structure
}

// Run runs Fpocket on a PDB and parses the results
func Run(p *pdb.PDB, msg func(string)) (pockets []*Pocket, err error) {
	dirName := p.ID + "_out"
	dirPath := "data/fpocket/" + dirName

	_, err = os.Stat(dirPath)
	if os.IsNotExist(err) {
		// Run Fpocket
		msg("running Fpocket")
		out, err := exec.Command("fpocket", "-f", p.LocalPath).CombinedOutput()
		if err != nil {
			return nil, err
		}
		if strings.Contains(string(out), "failed") {
			fmt.Println(string(out))
			return nil, errors.New("FPocket failed")
		}
	}

	err = os.Rename("bin/"+dirName, "data/fpocket/"+dirName)
	if err != nil {
		log.Fatal(err)
	}

	msg("retrieving Fpocket results")
	// Walk created folder containing pocket analysis files
	dir := dirPath + "/pockets"
	pockets, err = walkPocketDir(p, dir, msg)
	if err != nil {
		return nil, err
	}

	return pockets, nil
}

func walkPocketDir(crystal *pdb.PDB, dir string, msg func(string)) (pockets []*Pocket, err error) {
	n := 0
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
				pocketPDB, err := pdb.NewPDBFromRaw(data)
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
			n++
		}

		return nil
	})
	msg(fmt.Sprintf("found %d pockets, %d suitable", n, len(pockets)))

	return pockets, err
}
