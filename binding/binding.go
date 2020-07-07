package binding

import (
	"fmt"
	"math"
	"time"
	"varq/binding/fpocket"
	"varq/binding/ligand"
	"varq/pdb"
	"varq/uniprot"
)

// Results holds the collected data in the binding analysis step
type Results struct {
	Pockets  []*fpocket.Pocket         `json:"pockets"`
	Ligands  map[string][]*pdb.Residue `json:"ligands"` // ligand ID to near residues
	Residues []*pdb.Residue            `json:"residues"`
	Duration time.Duration             `json:"duration"`
	Error    error                     `json:"error"`
}

// Run starts the binding analysis step
func Run(unp *uniprot.UniProt, pdb *pdb.PDB, results chan<- *Results, msg func(string)) {
	start := time.Now()
	ligand, err := ligand.ResiduesNearLigands(pdb, msg)
	if err != nil {
		results <- &Results{Error: fmt.Errorf("Ligands: %v", err)}
	}

	pockets, err := fpocket.Run(pdb, msg)
	if err != nil {
		results <- &Results{Error: fmt.Errorf("running Fpocket: %v", err)}
	}

	results <- &Results{
		Pockets:  pockets,
		Ligands:  ligand,
		Residues: BindingResidues(unp, pdb, pockets),
		Duration: time.Since(start),
	}
}

func BindingResidues(unp *uniprot.UniProt, p *pdb.PDB, pockets []*fpocket.Pocket) (bindingResidues []*pdb.Residue) {
	// UniProt site positions to residues in structure
	var siteResidues []*pdb.Residue
	for _, site := range unp.Sites {
		siteResidues = append(siteResidues, p.UniProtPositions[unp.ID][site.Position]...)
	}

	for _, siteRes := range siteResidues {
		for _, res := range residuesFromNearPockets(siteRes, pockets) {
			if !residueExists(bindingResidues, res) {
				bindingResidues = append(bindingResidues, res)
			}
		}
	}

	if len(bindingResidues) == 0 {
		for _, siteRes := range siteResidues {
			for _, chain := range p.Chains {
				for _, res := range chain {
					if isNear(siteRes, res) && !residueExists(bindingResidues, res) {
						bindingResidues = append(bindingResidues, res)
					}
				}
			}
		}
	}

	return bindingResidues
}

// residuesFromNearPockets receives a residue and a slice of pockets, and returns all residues from
// pockets that are < 5 A from the passed residue.
func residuesFromNearPockets(residue *pdb.Residue, pockets []*fpocket.Pocket) (nearRes []*pdb.Residue) {
	for _, pocket := range pockets {
		if residueExists(pocket.Residues, residue) {
			nearRes = append(nearRes, pocket.Residues...)
		}
	}

	return nearRes
}

// isNear returns true if a pair of atoms between residues is closer than 5 angstroms, false otherwise.
func isNear(r1 *pdb.Residue, r2 *pdb.Residue) bool {
	for _, a1 := range r1.Atoms {
		for _, a2 := range r2.Atoms {
			distance := math.Sqrt(math.Pow(a1.X-a2.X, 2) + math.Pow(a1.Y-a2.Y, 2) + math.Pow(a1.Z-a2.Z, 2))
			if distance < 5 {
				return true
			}
		}
	}
	return false
}

func residueExists(rs []*pdb.Residue, r *pdb.Residue) bool {
	for _, res := range rs {
		if res == r {
			return true
		}
	}
	return false
}
