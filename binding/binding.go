package binding

import (
	"fmt"
	"time"
	"varq/binding/fpocket"
	"varq/pdb"
)

// BindingAnalysis holds the collected data in the binding analysis step
type BindingAnalysis struct {
	Pockets  []*fpocket.Pocket
	Duration time.Duration
	Error    error
}

// RunBindingAnalysis starts the binding analysis step
func RunBindingAnalysis(pdb *pdb.PDB, results chan<- *BindingAnalysis) {
	start := time.Now()
	pockets, err := fpocket.Run(pdb)
	if err != nil {
		results <- &BindingAnalysis{Error: fmt.Errorf("running Fpocket: %v", err)}
	}

	results <- &BindingAnalysis{Pockets: pockets, Duration: time.Since(start)}
}
