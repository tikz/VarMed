package ligand

import (
	"fmt"
	"math"
	"respdb/pdb"
)

// calculateDistance returns the distance between two atoms
func calculateDistance(atom1 *pdb.Atom, atom2 *pdb.Atom) float64 {
	return math.Sqrt(math.Pow(atom1.X-atom2.X, 2) + math.Pow(atom1.Y-atom2.Y, 2) + math.Pow(atom1.Z-atom2.Z, 2))
}

// ResiduesNearLigands returns a map of ligand IDs to near residues.
func ResiduesNearLigands(p *pdb.PDB, msg func(string)) (map[string][]*pdb.Residue, error) {
	// TODO: just load one time, somewhere
	// pdbBind, err := LoadPDBBind()
	// if err != nil {
	// 	return nil, fmt.Errorf("cannot load PDBBind data: %v", err)
	// }

	ligands := make(map[string][]*pdb.Residue)

	for _, hetatm := range p.HetAtoms {
		ln := hetatm.Residue

		// ignore ligand IDs not present in PDBBind
		// if _, exists := pdbBind[ln]; !exists {
		// 	continue
		// }

		for _, chain := range p.Chains {
			for _, res := range chain {
				var hasCloseAtom bool
				var i int

				for !hasCloseAtom && i < len(res.Atoms) {
					hasCloseAtom = calculateDistance(hetatm, res.Atoms[i]) < 5
					i++
				}

				if hasCloseAtom && (len(ligands[ln]) == 0 || ligands[ln][len(ligands[ln])-1] != res) {
					var exists bool
					for _, r := range ligands[ln] {
						if res == r {
							exists = true
							break
						}
					}
					if !exists {
						ligands[ln] = append(ligands[ln], res)
					}
				}

			}
		}
	}

	for name, res := range ligands {
		msg(fmt.Sprintf("%d residues near ligand %s", len(res), name))
	}

	return ligands, nil
}
