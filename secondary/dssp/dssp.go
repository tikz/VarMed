package dssp

import (
	"os/exec"
	"strconv"
	"strings"
	"varq/pdb"
)

// RunDSSP calculates secondary structure for a given PDB.
func RunDSSP(p *pdb.PDB) error {
	cmd := exec.Command("mkdssp", "-i", p.ID+".pdb")
	cmd.Dir = "bin/"

	out, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	// https://swift.cmbi.umcn.nl/gv/dssp/
	start := false
	for _, l := range strings.Split(string(out), "\n") {
		if len(l) > 17 {
			if start {
				posStr := strings.TrimSpace(l[5:10])
				if len(posStr) > 0 {
					chain := string(l[11])
					pos, _ := strconv.ParseInt(posStr, 10, 64)
					if chain, ok := p.Chains[chain]; ok {
						if res, ok := chain[pos]; ok {
							res.DSSP = string(l[16])
						}
					}
				}
			}
			if string(l[2]) == "#" {
				start = true
			}
		}
	}

	return nil
}
