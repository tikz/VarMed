package binding

import (
	"fmt"
	"time"
	"varq/binding/fpocket"
	"varq/binding/ligand"
	"varq/binding/mcsa"
	"varq/pdb"
	"varq/uniprot"
)

// Results holds the collected data in the binding analysis step
type Results struct {
	Pockets   []*fpocket.Pocket         `json:"pockets"`   // pockets with Fpocket drug score of >0.5
	Catalytic *mcsa.Catalytic           `json:"catalytic"` // catalytic residues in M-CSA
	Ligands   map[string][]*pdb.Residue `json:"ligands"`   // ligand ID to near residues
	Duration  time.Duration             `json:"duration"`
	Error     error                     `json:"error"`
}

// Run starts the binding analysis step
func Run(unp *uniprot.UniProt, pdb *pdb.PDB, results chan<- *Results, msg func(string)) {
	start := time.Now()
	pockets, err := fpocket.Run(pdb, msg)
	if err != nil {
		results <- &Results{Error: fmt.Errorf("running Fpocket: %v", err)}
	}

	csa, err := mcsa.GetPositions(unp, pdb, msg)
	if err != nil {
		results <- &Results{Error: fmt.Errorf("M-CSA: %v", err)}
	}

	ligand, err := ligand.ResiduesNearLigands(pdb, msg)
	if err != nil {
		results <- &Results{Error: fmt.Errorf("Ligands: %v", err)}
	}

	results <- &Results{
		Pockets:   pockets,
		Catalytic: csa,
		Ligands:   ligand,
		Duration:  time.Since(start),
	}
}
