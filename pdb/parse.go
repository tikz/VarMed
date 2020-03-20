package pdb

// Just a bunch of auxiliary regexes and parsing functions.
import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// extractPDBAtoms extracts the atoms from raw PDB contents
func extractPDBAtoms(raw []byte) ([]*Atom, error) {
	var atoms []*Atom

	regex, _ := regexp.Compile("(?m)^ATOM.*$")
	matches := regex.FindAllString(string(raw), -1)
	if len(matches) == 0 {
		return atoms, errors.New("atoms not found")
	}

	for _, match := range matches {
		var atom Atom

		// Columns position spec: https://www.cgl.ucsf.edu/chimera/docs/UsersGuide/tutorials/pdbintro.html
		atom.Number, _ = strconv.ParseInt(strings.TrimSpace(match[6:11]), 10, 64)
		atom.Residue = match[17:20]
		atom.Chain = match[21:22]
		atom.ResidueNumber, _ = strconv.ParseInt(strings.TrimSpace(match[22:26]), 10, 64)
		atom.X, _ = strconv.ParseFloat(strings.TrimSpace(match[30:38]), 64)
		atom.Y, _ = strconv.ParseFloat(strings.TrimSpace(match[38:46]), 64)
		atom.Z, _ = strconv.ParseFloat(strings.TrimSpace(match[46:54]), 64)

		atoms = append(atoms, &atom)
	}

	return atoms, nil
}

// extractPDBChains extracts the aminoacid chains from a slice of atoms
func extractPDBChains(raw []byte) (map[string]map[int64]*Aminoacid, error) {
	atoms, err := extractPDBAtoms(raw)
	if err != nil {
		return nil, fmt.Errorf("parsing PDB atoms: %v", err)
	}
	if len(atoms) == 0 {
		return nil, errors.New("empty atoms slice")
	}

	chains := make(map[string]map[int64]*Aminoacid)

	for _, atom := range atoms {
		chain, chainOk := chains[atom.Chain]
		pos, posOk := chain[atom.ResidueNumber]

		if !chainOk {
			chains[atom.Chain] = make(map[int64]*Aminoacid)
		}
		if !posOk {
			aa, err := NewAminoacid(atom.ResidueNumber, atom.Residue, []*Atom{atom})
			if err != nil {
				return nil, fmt.Errorf("cannot parse aminoacid: %v", atom.Residue)
			}
			chains[atom.Chain][atom.ResidueNumber] = aa
		} else {
			pos.Atoms = append(pos.Atoms, atom)
		}
	}

	return chains, nil
}

// CIF contains additional data that in PDB files is included under the REMARK tag, which is not standarized and hard to parse.

// extractCIFTitle extracts the main publication title from the CIF file
func extractCIFTitle(raw []byte) (string, error) {
	regex, _ := regexp.Compile("(?s)_struct.title.*?'(.*?)'")
	matches := regex.FindAllStringSubmatch(string(raw), -1)
	if len(matches) == 0 {
		return "", errors.New("CIF title not found")
	}
	return matches[0][1], nil
}

// extractCIFDate extracts the main publication date from the CIF file
func extractCIFDate(raw []byte) (*time.Time, error) {
	regex, _ := regexp.Compile("_pdbx_database_status.recvd_initial_deposition_date[ ]*([0-9]*-[0-9]*-[0-9]*)")
	matches := regex.FindAllStringSubmatch(string(raw), -1)
	if len(matches) == 0 {
		return nil, errors.New("CIF date not found")
	}

	t, err := time.Parse("2006-01-02", string(matches[0][1]))
	if err != nil {
		return nil, fmt.Errorf("parse CIF date: %v", err)
	}
	return &t, nil
}
