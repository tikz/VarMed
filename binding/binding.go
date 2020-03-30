package binding

import (
	"fmt"
	"time"
	"varq/binding/fpocket"
	"varq/binding/mcsa"
	"varq/pdb"
)

// BindingAnalysis holds the collected data in the binding analysis step
type Step struct {
	Pockets   []*fpocket.Pocket // Only pockets with drug score >0.5
	Catalytic *mcsa.Catalytic
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

	csa, err := mcsa.GetCSA(pdb)
	if err != nil {
		results <- &Step{Error: fmt.Errorf("M-CSA: %v", err)}
	}

	results <- &Step{Pockets: pockets, Catalytic: csa, Duration: time.Since(start)}
}
