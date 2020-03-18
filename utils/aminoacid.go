package utils

import (
	"errors"
	"strings"
)

var aas = [...][3]string{
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

// Aminoacid represents a single aminoacid. Holds all the ways that can be represented as a string.
type Aminoacid struct {
	Position int64
	Name     string
	Abbrv1   string
	Abbrv3   string
}

// NewAminoacid returns a new *Aminoacid from a string that can be either full name, one or three letter abbreviation.
func NewAminoacid(pos int64, input string) (*Aminoacid, error) {
	r, err := matchName(input)
	if err != nil {
		return nil, err
	}

	aminoacid := Aminoacid{
		Position: pos,
		Name:     r[0],
		Abbrv1:   r[1],
		Abbrv3:   r[2],
	}

	return &aminoacid, err
}

// AminoacidExists checks if the given aminoacid position already exists in a given slice (list) of aminoacids
func AminoacidExists(aminoacids []*Aminoacid, pos int64) bool {
	for _, aa := range aminoacids {
		if aa.Position == pos {
			return true
		}
	}
	return false
}

func matchName(input string) (*[3]string, error) {
	s := strings.Title(strings.ToLower(input))
	for _, aa := range aas {
		if aa[0] == s || aa[1] == s || aa[2] == s {
			return &aa, nil
		}
	}

	return nil, errors.New("unknown aminoacid name or abbreviation: " + s)
}
