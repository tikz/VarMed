package interaction

import (
	"fmt"
	"math"
	"respdb/pdb"
	"time"
)

// Results holds the collected data in the interaction analysis step
type Results struct {
	ChainsInteractions []*ChainsPair  `json:"chainsInteractions"`
	Residues           []*pdb.Residue `json:"residues"`
	Duration           time.Duration  `json:"duration"`
	Error              error          `json:"error"`
}

// ChainsPair holds all residue interactions between two chains.
type ChainsPair struct {
	Chain1               string          `json:"chain1"`
	Chain2               string          `json:"chain2"`
	ResiduesInteractions []*ResiduesPair `json:"residuesInteractions"`
}

// ResiduesPair holds all interaction parameters between two residues.
type ResiduesPair struct {
	Distance float64      `json:"distance"`
	Residue1 *pdb.Residue `json:"residue1"`
	Residue2 *pdb.Residue `json:"residue2"`
}

// Run starts the interaction analysis step
func Run(p *pdb.PDB, results chan<- *Results, msg func(string)) {
	start := time.Now()
	interactions := calculateChainsInteraction(p.Chains)
	residues := getInteractionResidues(interactions)
	msg(fmt.Sprintf("%d interface residues by distance", len(residues)))
	results <- &Results{ChainsInteractions: interactions,
		Residues: residues,
		Duration: time.Since(start)}
}

// Distance returns the distance between two atoms
func Distance(atom1 *pdb.Atom, atom2 *pdb.Atom) float64 {
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
			if dist := Distance(atom1, atom2); dist < 5 {
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
			var r1Exists, r2Exists bool
			for _, r := range residues {
				if r == interaction.Residue1 {
					r1Exists = true
					break
				}
			}
			for _, r := range residues {
				if r == interaction.Residue2 {
					r2Exists = true
					break
				}
			}
			if !r1Exists {
				residues = append(residues, interaction.Residue1)
			}
			if !r2Exists {
				residues = append(residues, interaction.Residue2)
			}
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
