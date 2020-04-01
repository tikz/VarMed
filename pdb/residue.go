package pdb

import (
	"errors"
	"fmt"
	"regexp"
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

// ExtractSeqRes parses the raw PDB for SEQRES records containing the primary sequence.
func (pdb *PDB) ExtractSeqRes() error {
	regex, _ := regexp.Compile("SEQRES[ ]*.*?[ ]+(.*?)[ ]+([0-9]*)[ ]*([A-Z ]*)")
	matches := regex.FindAllStringSubmatch(string(pdb.RawPDB), -1)
	if len(matches) == 0 {
		return errors.New("SEQRES not found")
	}

	pdb.SeqRes = make(map[string][]*Residue)
	for _, match := range matches {
		chain := match[1]
		resSplit := strings.Split(match[3], " ")
		for i, resStr := range resSplit {
			if resStr != "" {
				res := NewResidue(chain, int64(i), resStr)
				pdb.SeqRes[chain] = append(pdb.SeqRes[chain], res)
			}
		}
	}

	return nil
}

// ExtractResidues extracts data from the ATOM and HETATM records and parses them accordingly.
func (pdb *PDB) ExtractResidues() error {
	atoms, err := pdb.extractPDBATMRecords("ATOM")
	if err != nil {
		return fmt.Errorf("extract ATOM records: %v", err)
	}

	hetatms, _ := pdb.extractPDBATMRecords("HETATM")

	pdb.Atoms = atoms
	pdb.HetAtoms = hetatms

	err = pdb.ExtractPDBChains()
	if err != nil {
		return fmt.Errorf("extract PDB chains: %v", err)
	}

	return nil
}

// ExtractPDBChains parses the residue chains from raw PDB contents.
func (pdb *PDB) ExtractPDBChains() error {
	atoms := pdb.Atoms
	if len(atoms) == 0 {
		return errors.New("empty atoms list")
	}

	chains := make(map[string]map[int64]*Residue)

	var res *Residue
	for _, atom := range atoms {
		chain, chainOk := chains[atom.Chain]
		pos, posOk := chain[atom.ResidueNumber]

		if !chainOk {
			chains[atom.Chain] = make(map[int64]*Residue)
		}
		if !posOk {
			res = NewResidue(atom.Chain, atom.ResidueNumber, atom.Residue)
			res.Atoms = []*Atom{atom}
			chains[atom.Chain][atom.ResidueNumber] = res
		} else {
			pos.Atoms = append(pos.Atoms, atom)
		}

		atom.Aminoacid = res // parent ref
	}

	pdb.Chains = chains

	for _, chain := range pdb.Chains {
		pdb.TotalLength += int64(len(chain))
	}

	return nil
}
