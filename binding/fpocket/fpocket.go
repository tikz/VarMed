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
func Run(p *pdb.PDB, msg func(string)) (pockets []*Pocket, err error) {
	outPath := strings.Split(p.LocalPath, ".")[0] + "_out"

	// Delete Fpocket result files on function exit
	defer func() {
		os.RemoveAll(outPath)
	}()

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

	msg("retrieving Fpocket results")
	// Walk created folder containing pocket analysis files
	dir := outPath + "/pockets"
	fmt.Println(dir)
	pockets, err = walkPocketDir(p, dir, msg)
	if err != nil {
		return nil, err
	}

	return pockets, nil
}

func walkPocketDir(crystal *pdb.PDB, dir string, msg func(string)) (pockets []*Pocket, err error) {
	err = filepath.Walk(dir, func(file string, info os.FileInfo, err error) error {
		n := 0
		// For each Fpocket result PDB file
		if strings.Contains(file, "_atm.pdb") {
			msg(fmt.Sprintf("parsing Fpocket %s", file))
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
			n++
		}

		msg(fmt.Sprintf("found %d pockets, %d suitable", n, len(pockets)))
		return nil
	})

	return pockets, err
}
