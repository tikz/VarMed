package interaction

import (
	"fmt"
	"math"
	"time"
	"varq/protein/pdb"
)

// TODO: better names
type InteractionData struct {
	Interactions []*Interaction
}

type Interaction struct {
	Distance   float64
	Atom1      *pdb.Atom
	Atom2      *pdb.Atom
	Aminoacid1 *pdb.Aminoacid
	Aminoacid2 *pdb.Aminoacid
}

func calculateDistance(atom1 *pdb.Atom, atom2 *pdb.Atom) float64 {
	return math.Sqrt(math.Pow(atom1.X-atom2.X, 2) + math.Pow(atom1.Y-atom2.Y, 2) + math.Pow(atom1.Z-atom2.Z, 2))
}

func calculateInteractions(atoms []*pdb.Atom) (interactions []*Interaction) {
	start := time.Now()
	for i1, atom1 := range atoms {
		for i2, atom2 := range atoms {
			if i2 > i1 && atom1.Chain != atom2.Chain {
				if dist := calculateDistance(atom1, atom2); dist < 5 {
					// aa1, _ := pdb.NewAminoacid(atom1.Chain, atom1.ResidueNumber, atom1.Residue)
					// aa2, _ := pdb.NewAminoacid(atom2.Chain, atom2.ResidueNumber, atom2.Residue)
					// interactions = append(interactions, &Interaction{
					// 	Distance:   dist,
					// 	Atom1:      atom1,
					// 	Atom2:      atom2,
					// 	Aminoacid1: aa1,
					// 	Aminoacid2: aa2,
					// })
				}
			}
		}
	}
	end := time.Since(start)
	fmt.Println("finished in", end)
	fmt.Printf("%v+", interactions)
	return interactions
}

func NewInteraction(atoms []*pdb.Atom) {
	calculateInteractions(atoms)
}
