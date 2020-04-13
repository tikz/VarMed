package interaction

import (
	"math"
	"time"
	"varq/pdb"
)

// InteractionAnalysis holds the collected data in the interaction analysis step
type Results struct {
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

// Run starts the interaction analysis step
func Run(p *pdb.PDB, results chan<- *Results) {
	start := time.Now()
	interactions := calculateChainsInteraction(p.Chains)
	results <- &Results{ChainsInteractions: interactions,
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
				res1, a1 := flatten(chain1)
				res2, a2 := flatten(chain2)

				chainsInteraction.ResiduesInteractions = calculateResiduesInteraction(res1, a1, res2, a2)

				chainInters = append(chainInters, chainsInteraction)
			}
			i2++
		}
		i1++
	}

	return chainInters
}

func calculateResiduesInteraction(res1 []*pdb.Residue, a1 []*pdb.Atom,
	res2 []*pdb.Residue, a2 []*pdb.Atom) (resInteracts []*ResiduesPair) {
	for i1, atom1 := range a1 {
		for i2, atom2 := range a2 {
			if dist := calculateDistance(atom1, atom2); dist < 5 {
				resInteracts = append(resInteracts, &ResiduesPair{
					Distance: dist,
					Residue1: res1[i1],
					Residue2: res2[i2],
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

// flatten converts the residue map datatype to flat slices of residues and atom pointers
func flatten(chain map[int64]*pdb.Residue) (residues []*pdb.Residue, atoms []*pdb.Atom) {
	for _, residue := range chain {
		for _, atom := range residue.Atoms {
			atoms = append(atoms, atom)
			residues = append(residues, residue)
		}
	}
	return residues, atoms
}
