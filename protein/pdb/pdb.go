package pdb

import (
	"fmt"
	"time"
	"varq/http"
	"varq/utils"
)

type PDB struct {
	ID         string
	URL        string
	PDBURL     string
	CIFURL     string
	RawPDB     []byte `json:"-"`
	RawCIF     []byte `json:"-"`
	Title      string
	Date       *time.Time
	Method     string
	Resolution float64
	FromPos    int64
	ToPos      int64
	Length     int64
	Atoms      []*Atom
	Chains     map[string][]*utils.Aminoacid
}

func (pdb *PDB) Fetch() error {
	url := "https://www.rcsb.org/structure/" + pdb.ID
	urlCIF := "https://files.rcsb.org/download/" + pdb.ID + ".cif"
	rawCIF, err := http.Get(urlCIF)
	if err != nil {
		return fmt.Errorf("download CIF file: %v", err)
	}

	urlPDB := "https://files.rcsb.org/download/" + pdb.ID + ".pdb"
	rawPDB, err := http.Get(urlPDB)
	if err != nil {
		return fmt.Errorf("download PDB file: %v", err)
	}

	// Mandatory data
	pdb.URL = url
	pdb.PDBURL = urlPDB
	pdb.CIFURL = urlCIF
	pdb.RawPDB = rawPDB
	pdb.RawCIF = rawCIF
	pdb.Atoms, err = parsePDBAtoms(pdb.RawPDB)
	if err != nil {
		return fmt.Errorf("parsing PDB atoms: %v", err)
	}
	pdb.Chains, err = parsePDBChains(pdb.Atoms)
	if err != nil {
		return fmt.Errorf("parsing chains: %v", err)
	}

	// Optional data, but can be nice to have
	if t, err := parseCIFTitle(pdb.RawCIF); err == nil {
		pdb.Title = t
	}
	if d, err := parseCIFDate(pdb.RawCIF); err == nil {
		pdb.Date = d
	}

	return nil
}
