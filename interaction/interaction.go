package interaction

import (
	"math"
	"time"
	"varq/pdb"
)

// InteractionAnalysis holds the collected data in the interaction analysis step
type InteractionAnalysis struct {
	Interactions []*AminoacidsInteraction
	Duration     time.Duration
	Error        error
}

// ChainsInteraction holds all aminoacid interactions between two chains
type ChainsInteraction struct {
	Chain1                 string
	Chain2                 string
	AminoacidsInteractions []*AminoacidsInteraction
}

// AminoacidsInteraction holds all interaction parameters between two aminoacids
type AminoacidsInteraction struct {
	Distance   float64
	Aminoacid1 *pdb.Aminoacid
	Aminoacid2 *pdb.Aminoacid
}

// RunInteractionAnalysis starts the interaction analysis step
func RunInteractionAnalysis(p *pdb.PDB, results chan<- *InteractionAnalysis) {
	start := time.Now()
	interactions := calculateChainsInteraction(p.Chains)
	results <- &InteractionAnalysis{Interactions: interactions, Duration: time.Since(start)}
}

// calculateDistance returns the distance between two atoms
func calculateDistance(atom1 *pdb.Atom, atom2 *pdb.Atom) float64 {
	return math.Sqrt(math.Pow(atom1.X-atom2.X, 2) + math.Pow(atom1.Y-atom2.Y, 2) + math.Pow(atom1.Z-atom2.Z, 2))
}

func calculateChainsInteraction(chains map[string]map[int64]*pdb.Aminoacid) (aaInteracts []*AminoacidsInteraction) {
	var i1, i2 int

	for chainName1, chain1 := range chains {
		for chainName2, chain2 := range chains {
			if i2 > i1 && chainName1 != chainName2 {
				chainAtoms1 := flatten(chain1)
				chainAtoms2 := flatten(chain2)

				aaInteracts = calculateAminoacidsInteraction(chainAtoms1, chainAtoms2)
			}
			i2++
		}
		i1++
	}

	return aaInteracts
}

func calculateAminoacidsInteraction(chain1 []*pdb.Atom, chain2 []*pdb.Atom) (aaInteracts []*AminoacidsInteraction) {
	for _, atom1 := range chain1 {
		for _, atom2 := range chain2 {
			if dist := calculateDistance(atom1, atom2); dist < 5 {
				aaInteracts = append(aaInteracts, &AminoacidsInteraction{
					Distance:   dist,
					Aminoacid1: atom1.Aminoacid,
					Aminoacid2: atom2.Aminoacid,
				})
			}
		}
	}
	return aaInteracts
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
