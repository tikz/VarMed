package pdb

import (
	"strings"
)

var residueNames = [...][3]string{
	[3]string{"Alanine", "Ala", "A"},
	[3]string{"Arginine", "Arg", "R"},
	[3]string{"Asparagine", "Asn", "N"},
	[3]string{"Aspartic acid", "Asp", "D"},
	[3]string{"Cysteine", "Cys", "C"},
	[3]string{"Glutamic acid", "Glu", "E"},
	[3]string{"Glutamine", "Gln", "Q"},
	[3]string{"Glycine", "Gly", "G"},
	[3]string{"Histidine", "His", "H"},
	[3]string{"Isoleucine", "Ile", "I"},
	[3]string{"Leucine", "Leu", "L"},
	[3]string{"Lysine", "Lys", "K"},
	[3]string{"Methionine", "Met", "M"},
	[3]string{"Phenylalanine", "Phe", "F"},
	[3]string{"Proline", "Pro", "P"},
	[3]string{"Serine", "Ser", "S"},
	[3]string{"Threonine", "Thr", "T"},
	[3]string{"Tryptophan", "Trp", "W"},
	[3]string{"Tyrosine", "Tyr", "Y"},
	[3]string{"Valine", "Val", "V"},
}

// Residue represents a single residue from the PDB structure.
type Residue struct {
	Chain    string
	Position int64
	Name     string
	Abbrv1   string
	Abbrv3   string
	Atoms    []*Atom
}

// NewResidue constructs a new residue given a chain, position and aminoacid name.
// The name is case-insensitive and can be either a full aminoacid name, one or three letter abbreviation.
func NewResidue(chain string, pos int64, input string) *Residue {
	name, abbrv3, abbrv1 := matchName(input)

	res := &Residue{
		Chain:    chain,
		Position: pos,
		Name:     name,
		Abbrv3:   abbrv3,
		Abbrv1:   abbrv1,
	}

	return res
}

// matchName receives a residue name and returns a 3-sized array of all the possible representations as a string.
func matchName(input string) (string, string, string) {
	s := strings.Title(strings.ToLower(input))
	for _, res := range residueNames {
		for _, n := range res {
			if n == s {
				return res[0], res[1], res[2]
			}
		}
	}

	return input, "Unk", "X"
}
