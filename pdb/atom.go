package pdb

type Atom struct {
	// PDB fields
	Number        int64
	Residue       string
	Chain         string
	ResidueNumber int64
	X             float64
	Y             float64
	Z             float64

	// Parent ref
	Aminoacid *Aminoacid `json:"-"`
}
