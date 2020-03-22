package pdb

import (
	"fmt"
	"log"
	"time"
	"varq/http"
)

type PDB struct {
	UniProtID  string
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
	Length     int64 // Sum of chains length
	Chains     map[string]map[int64]*Aminoacid
}

// Fetch populates the instance with parsed data retrieved from RCSB
func (pdb *PDB) Fetch() error {
	start := time.Now()
	log.Printf("[UniProt %s] Downloading PDB file for %s", pdb.UniProtID, pdb.ID)
	url := "https://www.rcsb.org/structure/" + pdb.ID
	urlCIF := "https://files.rcsb.org/download/" + pdb.ID + ".cif"
	rawCIF, err := http.Get(urlCIF)
	if err != nil {
		return fmt.Errorf("download CIF file: %v", err)
	}

	log.Printf("[UniProt %s] Downloading CIF file for %s", pdb.UniProtID, pdb.ID)
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

	err = pdb.ExtractChains()
	if err != nil {
		return err
	}

	// Optional data, but can be nice to have
	if t, err := extractCIFTitle(pdb.RawCIF); err == nil {
		pdb.Title = t
	}
	if d, err := extractCIFDate(pdb.RawCIF); err == nil {
		pdb.Date = d
	}
	end := time.Since(start)
	log.Printf("[UniProt %s] PDB %s loaded in %d msecs", pdb.UniProtID, pdb.ID, end.Milliseconds())

	return nil
}

func (pdb *PDB) ExtractChains() error {
	chains, err := extractPDBChains(pdb.RawPDB)
	if err != nil {
		return fmt.Errorf("parsing chains: %v", err)
	}
	pdb.Chains = chains
	return nil
}