package exposure

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"
	"varq/pdb"
)

type Results struct {
	Residues []*pdb.Residue `json:"residues"`
	Duration time.Duration  `json:"duration"`
	Error    error          `json:"error"`
}

func Run(pdb *pdb.PDB, results chan<- *Results, msg func(string)) {
	start := time.Now()
	buried, err := BuriedResidues(pdb)
	if err != nil {
		results <- &Results{Error: fmt.Errorf("calculate rSASA: %v", err)}
	}

	results <- &Results{Residues: buried, Duration: time.Since(start)}
}

func BuriedResidues(p *pdb.PDB) ([]*pdb.Residue, error) {
	var buried []*pdb.Residue

	cmd := exec.Command("freesasa",
		p.ID+".pdb",
		"--format=rsa")
	cmd.Dir = "bin/"

	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	for _, line := range strings.Split(string(out), "\n") {
		f := strings.Fields(line)
		if len(f) > 0 && f[0] == "RES" {
			aa := f[1]
			chain := f[2]
			pos, _ := strconv.ParseInt(f[3], 10, 64)
			rsasa, _ := strconv.ParseFloat(f[7], 64)
			if aa != "GLY" && rsasa < 50 {
				buried = append(buried, p.Chains[chain][pos])
			}
		}
	}

	return buried, nil
}
