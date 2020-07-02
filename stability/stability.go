package stability

import (
	"fmt"
	"time"
	"varq/pdb"
	"varq/stability/foldx"
	"varq/uniprot"
)

// Results holds the collected data in the stability analysis step
type Results struct {
	FoldX    []*foldx.SASDiff `json:"foldx"`
	Duration time.Duration    `json:"duration"`
	Error    error            `json:"error"`
}

// Run starts the stability analysis step
func Run(sasList []*uniprot.SAS, unp *uniprot.UniProt, pdb *pdb.PDB,
	results chan<- *Results, msg func(string)) {
	start := time.Now()

	foldxResults, err := foldx.Run(sasList, unp.ID, pdb, msg)
	if err != nil {
		results <- &Results{Error: fmt.Errorf("FoldX: %v", err)}
	}

	results <- &Results{
		FoldX:    foldxResults,
		Duration: time.Since(start),
	}
}
