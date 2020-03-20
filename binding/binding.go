package binding

import (
	"fmt"
	"varq/binding/fpocket"
	"varq/pdb"
)

type Binding struct {
	Pockets []*fpocket.Pocket
}

func NewBinding(pdb *pdb.PDB) (*Binding, error) {
	pockets, err := fpocket.Run(pdb)
	if err != nil {
		return nil, fmt.Errorf("running Fpocket: %v", err)
	}

	return &Binding{Pockets: pockets}, nil
}
