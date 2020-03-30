package pdb

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func (pdb *PDB) ExtractSeqRes() error {
	regex, _ := regexp.Compile("SEQRES[ ]*.*?[ ]+(.*?)[ ]+([0-9]*)[ ]*([A-Z ]*)") // https://regex101.com/r/9vwbyc/1
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

func (pdb *PDB) ExtractChains() error {
	chains, err := extractPDBChains(pdb.RawPDB)
	if err != nil {
		return fmt.Errorf("parsing chains: %v", err)
	}
	pdb.Chains = chains

	for _, chain := range pdb.Chains {
		pdb.TotalLength += int64(len(chain))
	}

	return nil
}

func (pdb *PDB) ExtractCIFData() error {
	title, err := extractCIFLine("title", "_struct.title", pdb.RawCIF)
	if err != nil {
		return err
	}

	method, err := extractCIFLine("method", "_refine.pdbx_refine_id", pdb.RawCIF)
	if err != nil {
		return err
	}

	resolutionStr, err := extractCIFLine("resolution", "_refine.ls_d_res_high", pdb.RawCIF)
	if err != nil {
		return err
	}
	resolution, err := strconv.ParseFloat(resolutionStr, 64)
	if err != nil {
		return err
	}

	date, err := extractCIFDate(pdb.RawCIF)
	if err != nil {
		return err
	}

	pdb.Title = title
	pdb.Method = method
	pdb.Resolution = resolution
	pdb.Date = date

	return nil
}

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

		// Positions reference: https://www.cgl.ucsf.edu/chimera/docs/UsersGuide/tutorials/pdbintro.html
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

// extractPDBChains extracts the residue chains from a slice of atoms
func extractPDBChains(raw []byte) (map[string]map[int64]*Residue, error) {
	atoms, err := extractPDBAtoms(raw)
	if err != nil {
		return nil, fmt.Errorf("parsing PDB atoms: %v", err)
	}
	if len(atoms) == 0 {
		return nil, errors.New("empty atoms slice")
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
		// Parent ref
		atom.Aminoacid = res
	}

	return chains, nil
}

// CIF contains additional data that in PDB files is included under the REMARK tag, which is not standarized and hard to parse.

func extractCIFLine(name string, pattern string, raw []byte) (string, error) {
	regex, _ := regexp.Compile("(?s)" + pattern + "[ ]*(.*?)_")
	matches := regex.FindAllStringSubmatch(string(raw), -1)
	if len(matches) == 0 {
		return "", errors.New("CIF " + name + " not found")
	}
	match := matches[0][1]
	match = strings.TrimSpace(match)
	match = strings.Replace(match, "'", "", -1)

	return match, nil
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
