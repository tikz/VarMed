package interaction

import (
	"math"
	"time"
	"varq/pdb"
)

// InteractionAnalysis holds the collected data in the interaction analysis step
type InteractionAnalysis struct {
	ChainsInteractions []*ChainsPair
	Residues           []*pdb.Residue
	Duration           time.Duration
	Error              error
}

// ChainsPair holds all residue interactions between two chains.
type ChainsPair struct {
	Chain1               string
	Chain2               string
	ResiduesInteractions []*ResiduesPair
}

// ResiduesPair holds all interaction parameters between two residues.
type ResiduesPair struct {
	Distance float64
	Residue1 *pdb.Residue
	Residue2 *pdb.Residue
}

// RunInteractionAnalysis starts the interaction analysis step
func RunInteractionAnalysis(p *pdb.PDB, results chan<- *InteractionAnalysis) {
	start := time.Now()
	interactions := calculateChainsInteraction(p.Chains)
	results <- &InteractionAnalysis{ChainsInteractions: interactions,
		Residues: getInteractionResidues(interactions),
		Duration: time.Since(start)}
}

// calculateDistance returns the distance between two atoms
func calculateDistance(atom1 *pdb.Atom, atom2 *pdb.Atom) float64 {
	return math.Sqrt(math.Pow(atom1.X-atom2.X, 2) + math.Pow(atom1.Y-atom2.Y, 2) + math.Pow(atom1.Z-atom2.Z, 2))
}

func calculateChainsInteraction(chains map[string]map[int64]*pdb.Residue) (chainInters []*ChainsPair) {
	var i1, i2 int

	for chainName1, chain1 := range chains {
		for chainName2, chain2 := range chains {
			if i2 > i1 && chainName1 != chainName2 {
				chainsInteraction := &ChainsPair{
					Chain1: chainName1,
					Chain2: chainName2,
				}
				chainAtoms1 := flatten(chain1)
				chainAtoms2 := flatten(chain2)

				chainsInteraction.ResiduesInteractions = calculateResiduesInteraction(chainAtoms1, chainAtoms2)

				chainInters = append(chainInters, chainsInteraction)
			}
			i2++
		}
		i1++
	}

	return chainInters
}

func calculateResiduesInteraction(chain1 []*pdb.Atom, chain2 []*pdb.Atom) (resInteracts []*ResiduesPair) {
	for _, atom1 := range chain1 {
		for _, atom2 := range chain2 {
			if dist := calculateDistance(atom1, atom2); dist < 5 {
				resInteracts = append(resInteracts, &ResiduesPair{
					Distance: dist,
					Residue1: atom1.Aminoacid,
					Residue2: atom2.Aminoacid,
				})
			}
		}
	}
	return resInteracts
}

// getInteractionResidues converts the slice of chain interaction pairs to a flat slice of residues
func getInteractionResidues(interactions []*ChainsPair) (residues []*pdb.Residue) {
	for _, chainPair := range interactions {
		for _, interaction := range chainPair.ResiduesInteractions {
			residues = append(residues, interaction.Residue1)
			residues = append(residues, interaction.Residue2)
		}
	}

	return residues
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
