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
	FoldX    []*foldx.Variant `json:"foldx"`
	Duration time.Duration    `json:"duration"`
	Error    error            `json:"error"`
}

// Run starts the energy analysis step
func Run(variants map[int]string, unp *uniprot.UniProt, pdb *pdb.PDB, foldxDir string,
	results chan<- *Results, msg func(string)) {
	start := time.Now()

	err := foldx.Run(variants, unp.ID, pdb, foldxDir, msg)
	if err != nil {
		results <- &Results{Error: fmt.Errorf("FoldX: %v", err)}
	}

	// pockets, err := fpocket.Run(pdb, msg)
	// if err != nil {
	// 	results <- &Results{Error: fmt.Errorf("running FoldX: %v", err)}
	// }

	// ligand, err := ligand.ResiduesNearLigands(pdb, msg)
	// if err != nil {
	// 	results <- &Results{Error: fmt.Errorf("Ligands: %v", err)}
	// }

	results <- &Results{
		Duration: time.Since(start),
	}
}
