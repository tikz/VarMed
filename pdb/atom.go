package pdb

type Atom struct {
	// PDB columns for the ATOM tag
	Number        int64
	Residue       string
	Chain         string
	ResidueNumber int64
	X             float64
	Y             float64
	Z             float64
	// TODO: add remaining

	// Parent ref
	Aminoacid *Residue `json:"-"`
}
