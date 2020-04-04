package binding

import (
	"fmt"
	"time"
	"varq/binding/fpocket"
	"varq/binding/ligand"
	"varq/binding/mcsa"
	"varq/pdb"
)

// BindingAnalysis holds the collected data in the binding analysis step
type Step struct {
	Pockets   []*fpocket.Pocket         // pockets with Fpocket drug score of >0.5
	Catalytic *mcsa.Catalytic           // catalytic residues in M-CSA
	Ligands   map[string][]*pdb.Residue // ligand ID to near residues
	Duration  time.Duration
	Error     error
}

// RunBindingAnalysis starts the binding analysis step
func RunBindingStep(pdb *pdb.PDB, results chan<- *Step) {
	start := time.Now()
	pockets, err := fpocket.Run(pdb)
	if err != nil {
		results <- &Step{Error: fmt.Errorf("running Fpocket: %v", err)}
	}

	csa, err := mcsa.GetPositions(pdb)
	if err != nil {
		results <- &Step{Error: fmt.Errorf("M-CSA: %v", err)}
	}

	ligand, err := ligand.ResiduesNearLigands(pdb)
	if err != nil {
		results <- &Step{Error: fmt.Errorf("Ligands: %v", err)}
	}

	results <- &Step{
		Pockets:   pockets,
		Catalytic: csa,
		Ligands:   ligand,
		Duration:  time.Since(start),
	}
}
