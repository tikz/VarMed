package interaction

import (
	"fmt"
	"math"
	"time"
	"varq/protein/pdb"
)

type InteractionData struct {
	Interactions []*AminoacidsInteraction
}

type AminoacidsInteraction struct {
	Distance   float64
	Aminoacid1 *pdb.Aminoacid
	Aminoacid2 *pdb.Aminoacid
}

func calculateDistance(atom1 *pdb.Atom, atom2 *pdb.Atom) float64 {
	return math.Sqrt(math.Pow(atom1.X-atom2.X, 2) + math.Pow(atom1.Y-atom2.Y, 2) + math.Pow(atom1.Z-atom2.Z, 2))
}

func calculateInteractions(chains map[string][]*pdb.Aminoacid) (interactions []*AminoacidsInteraction) {
	atoms := flatten(chains)
	start := time.Now()

	for i1, atom1 := range atoms {
		for i2, atom2 := range atoms {
			if i2 > i1 && atom1.Chain != atom2.Chain && atom1.Aminoacid != atom2.Aminoacid {
				if dist := calculateDistance(atom1, atom2); dist < 5 {
					interactions = append(interactions, &AminoacidsInteraction{
						Distance:   dist,
						Aminoacid1: atom1.Aminoacid,
						Aminoacid2: atom2.Aminoacid,
					})
					fmt.Println(dist, "aa1 pos", atom1.Aminoacid.Position, "aa2 pos", atom2.Aminoacid.Position, "aa1 chain", atom1.Chain, "aa2 chain", atom2.Chain)
				}
			}
		}
	}
	end := time.Since(start)
	fmt.Println("finished in", end)
	return interactions
}

func flatten(chains map[string][]*pdb.Aminoacid) (atoms []*pdb.Atom) {
	for _, aminoacids := range chains {
		for _, aminoacid := range aminoacids {
			for _, atom := range aminoacid.Atoms {
				atoms = append(atoms, atom)
			}
		}
	}
	return atoms
}

func NewInteraction(chains map[string][]*pdb.Aminoacid) {
	calculateInteractions(chains)
}
