package interaction

import (
	"fmt"
	"math"
	"time"
	"varq/pdb"
)

type InteractionAnalysis struct {
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

func calculateInteractions(chains map[string]map[int64]*pdb.Aminoacid) (interactions []*AminoacidsInteraction) {
	start := time.Now()

	var i1, i2 int
	for _, chain1 := range chains {
		for _, chain2 := range chains {
			if i2 > i1 {
				chainAtoms1 := flatten(chain1)
				chainAtoms2 := flatten(chain2)
				_ = calculateChainAAsDistance(chainAtoms1, chainAtoms2)

			}
			i2++
		}
		i1++
	}

	end := time.Since(start)
	fmt.Println("finished in", end)
	return interactions
}

func calculateChainAAsDistance(c1 []*pdb.Atom, c2 []*pdb.Atom) (ai []*AminoacidsInteraction) {
	for i1, atom1 := range c1 {
		for i2, atom2 := range c2 {
			if i2 > i1 {
				if dist := calculateDistance(atom1, atom2); dist < 5 {
					fmt.Println(dist, atom1, atom2)
					ai = append(ai, &AminoacidsInteraction{
						Distance:   dist,
						Aminoacid1: atom1.Aminoacid,
						Aminoacid2: atom2.Aminoacid,
					})
				}
			}
		}
	}
	return ai
}

// flatten converts the aminoacid slice datatype to a flat slice of atom pointers
func flatten(chain map[int64]*pdb.Aminoacid) (atoms []*pdb.Atom) {
	for _, aminoacid := range chain {
		for _, atom := range aminoacid.Atoms {
			atoms = append(atoms, atom)
		}
	}
	return atoms
}

func NewInteraction(p *pdb.PDB) {
	calculateInteractions(p.Chains)
}
