package energy

import (
	"fmt"
	"time"
	"varq/energy/foldx"
	"varq/pdb"
	"varq/uniprot"
)

// Results holds the collected data in the energy analysis step
type Results struct {
	FoldX    []*foldx.SASEnergyDiff `json:"foldx"`
	Duration time.Duration          `json:"duration"`
	Error    error                  `json:"error"`
}

// Run starts the energy analysis step
func Run(sasList []*uniprot.SAS, unp *uniprot.UniProt, pdb *pdb.PDB, foldxDir string,
	results chan<- *Results, msg func(string)) {
	start := time.Now()

	foldxResults, err := foldx.Run(sasList, unp.ID, pdb, foldxDir, msg)
	if err != nil {
		results <- &Results{Error: fmt.Errorf("FoldX: %v", err)}
	}

	results <- &Results{
		FoldX:    foldxResults,
		Duration: time.Since(start),
	}
}
