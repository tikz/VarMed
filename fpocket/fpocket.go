package fpocket

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"respdb/file"
	"respdb/pdb"
	"strconv"
	"strings"
)

//Pocket represents a single pocket found with Fpocket (pocketN_atm.pdb)
type Pocket struct {
	Name      string  `json:"name"`
	DrugScore float64 `json:"drugScore"`

	// residues in pocket, chain name to original PDB positions
	ChainPos map[string][]int64 `json:"chain_pos"`
}

// Run runs Fpocket on a PDB and parses the results
func Run(path string) (pockets []*Pocket, err error) {
	dir, fileName := filepath.Split(path)
	outDirName := strings.Split(fileName, ".")[0] + "_out"
	dirPath := "data/fpocket/" + outDirName

	_, err = os.Stat(dirPath)
	if os.IsNotExist(err) {
		// Run Fpocket
		out, err := exec.Command("fpocket", "-f", path).CombinedOutput()
		if err != nil || strings.Contains(string(out), "failed") {
			return nil, fmt.Errorf("fpocket: %v %s", err, string(out))
		}

		err = file.Copy(dir+outDirName, dirPath)
		if err != nil {
			return nil, err
		}
	}

	// Walk created folder containing pocket analysis files
	pocketsDir := dirPath + "/pockets"
	pockets, err = walkPocketDir(pocketsDir)
	if err != nil {
		return nil, err
	}

	return pockets, nil
}

func walkPocketDir(dir string) (pockets []*Pocket, err error) {
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

			pocketPDB, err := pdb.NewPDBFromRaw(data)
			if err != nil {
				return err
			}

			residues := make(map[string][]int64)
			for chain, chainPos := range pocketPDB.Chains {
				for pos := range chainPos {
					residues[chain] = append(residues[chain], pos)
				}
			}

			pocket := &Pocket{
				Name:      file,
				DrugScore: drugScore,
				ChainPos:  residues,
			}
			pockets = append(pockets, pocket)
			n++
		}

		return nil
	})

	return pockets, err
}
