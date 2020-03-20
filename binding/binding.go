package binding

import (
	"fmt"
	"varq/binding/fpocket"
	"varq/pdb"
)

type BindingAnalysis struct {
	Pockets []*fpocket.Pocket
}

func NewBindingAnalysis(pdb *pdb.PDB) (*BindingAnalysis, error) {
	pockets, err := fpocket.Run(pdb)
	if err != nil {
		return nil, fmt.Errorf("running Fpocket: %v", err)
	}

	return &BindingAnalysis{Pockets: pockets}, nil
}
