package stability

import (
	"fmt"
	"time"
	"varq/pdb"
	"varq/sasa"
	"varq/stability/foldx"
	"varq/uniprot"
)

// Results holds the collected data in the stability analysis step
type Results struct {
	RepairedStructure *RepairedStructure `json:"repaired_structure"`
	FoldX             []*foldx.Mutation  `json:"foldx"`
	Duration          time.Duration      `json:"duration"`
	Error             error              `json:"error"`
}

type RepairedStructure struct {
	SASA       float64 `json:"sasa"`
	SASAApolar float64 `json:"sasa_apolar"`
	SASAPolar  float64 `json:"sasa_polar"`
}

// Run starts the stability analysis step
func Run(sasList []*uniprot.SAS, unp *uniprot.UniProt, pdb *pdb.PDB,
	results chan<- *Results, msg func(string)) {
	start := time.Now()

	// Run FoldX Repair + BuildModel
	foldxResults, err := foldx.Run(sasList, unp.ID, pdb, msg)
	if err != nil {
		results <- &Results{Error: fmt.Errorf("FoldX: %v", err)}
	}

	// SASA of repaired structure
	total, apolar, polar, err := sasa.SASA("data/foldx/repair/" + pdb.ID + "_Repair.pdb")
	if err != nil {
		results <- &Results{Error: fmt.Errorf("repaired SASA: %v", err)}
	}

	// repairDir := "data/foldx/repair/"
	// path := repairDir + pdb.ID + "_Repair.pdb"
	// fpocket.Run(path)

	results <- &Results{
		FoldX: foldxResults,
		RepairedStructure: &RepairedStructure{
			SASA:       total,
			SASAApolar: apolar,
			SASAPolar:  polar,
		},
		Duration: time.Since(start),
	}
}
