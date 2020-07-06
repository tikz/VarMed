package exposure

import (
	"fmt"
	"time"
	"varq/pdb"
	"varq/sasa"
)

type Results struct {
	Residues []*pdb.Residue `json:"residues"`
	Duration time.Duration  `json:"duration"`
	Error    error          `json:"error"`
}

func Run(pdb *pdb.PDB, results chan<- *Results, msg func(string)) {
	start := time.Now()
	buried, err := sasa.BuriedResidues(pdb)
	if err != nil {
		results <- &Results{Error: fmt.Errorf("calculate rSASA: %v", err)}
	}

	results <- &Results{Residues: buried, Duration: time.Since(start)}
}
