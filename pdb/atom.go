package pdb

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

// Atom represents a single atom in the structure.
// It contains all the columns from an ATOM or HETATM record in a PDB file.
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

// extractPDBATMRecords extracts either ATOM or HETATM records.
func (pdb *PDB) extractPDBATMRecords(recordName string) ([]*Atom, error) {
	var atoms []*Atom

	r, _ := regexp.Compile("(?m)^" + recordName + ".*$")
	matches := r.FindAllString(string(pdb.RawPDB), -1)
	if len(matches) == 0 {
		return atoms, errors.New("atoms not found")
	}

	for _, match := range matches {
		var atom Atom

		// https://www.wwpdb.org/documentation/file-format-content/format23/sect9.html#ATOM
		atom.Number, _ = strconv.ParseInt(strings.TrimSpace(match[6:11]), 10, 64)
		atom.Residue = strings.TrimSpace(match[17:20])
		atom.Chain = match[21:22]
		atom.ResidueNumber, _ = strconv.ParseInt(strings.TrimSpace(match[22:26]), 10, 64)
		atom.X, _ = strconv.ParseFloat(strings.TrimSpace(match[30:38]), 64)
		atom.Y, _ = strconv.ParseFloat(strings.TrimSpace(match[38:46]), 64)
		atom.Z, _ = strconv.ParseFloat(strings.TrimSpace(match[46:54]), 64)

		atoms = append(atoms, &atom)
	}

	return atoms, nil
}
