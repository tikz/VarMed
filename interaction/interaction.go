package interaction

import (
	"math"
	"time"
	"varq/pdb"
)

// InteractionAnalysis holds the collected data in the interaction analysis step
type InteractionAnalysis struct {
	Interactions []*ResiduesInteraction
	Duration     time.Duration
	Error        error
}

// ChainsInteraction holds all residue interactions between two chains.
type ChainsInteraction struct {
	Chain1               string
	Chain2               string
	ResiduesInteractions []*ResiduesInteraction
}

// ResiduesInteraction holds all interaction parameters between two residues.
type ResiduesInteraction struct {
	Distance float64
	Residue1 *pdb.Residue
	Residue2 *pdb.Residue
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

func calculateChainsInteraction(chains map[string]map[int64]*pdb.Residue) (resInteracts []*ResiduesInteraction) {
	var i1, i2 int

	for chainName1, chain1 := range chains {
		for chainName2, chain2 := range chains {
			if i2 > i1 && chainName1 != chainName2 {
				chainAtoms1 := flatten(chain1)
				chainAtoms2 := flatten(chain2)

				resInteracts = calculateResiduesInteraction(chainAtoms1, chainAtoms2)
			}
			i2++
		}
		i1++
	}

	return resInteracts
}

func calculateResiduesInteraction(chain1 []*pdb.Atom, chain2 []*pdb.Atom) (resInteracts []*ResiduesInteraction) {
	for _, atom1 := range chain1 {
		for _, atom2 := range chain2 {
			if dist := calculateDistance(atom1, atom2); dist < 5 {
				resInteracts = append(resInteracts, &ResiduesInteraction{
					Distance: dist,
					Residue1: atom1.Aminoacid,
					Residue2: atom2.Aminoacid,
				})
			}
		}
	}
	return resInteracts
}

// flatten converts the residue map datatype to a flat slice of atom pointers
func flatten(chain map[int64]*pdb.Residue) (atoms []*pdb.Atom) {
	for _, residue := range chain {
		for _, atom := range residue.Atoms {
			atoms = append(atoms, atom)
		}
	}
	return atoms
}
