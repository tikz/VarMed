package pdb

import (
	"fmt"
	"log"
	"strconv"
	"time"
	"varq/http"
)

type PDB struct {
	ID          string
	URL         string
	PDBURL      string
	CIFURL      string
	RawPDB      []byte `json:"-"`
	RawCIF      []byte `json:"-"`
	Title       string
	Date        *time.Time
	Method      string
	Resolution  float64
	TotalLength int64
	Chains      map[string]map[int64]*Aminoacid `json:"-"`
}

// Fetch populates the instance with parsed data retrieved from RCSB
func (pdb *PDB) Fetch() error {
	start := time.Now()
	log.Printf("Downloading PDB and CIF files for %s", pdb.ID)
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

	pdb.URL = url
	pdb.PDBURL = urlPDB
	pdb.CIFURL = urlCIF
	pdb.RawPDB = rawPDB
	pdb.RawCIF = rawCIF

	err = pdb.ExtractCIFData()
	if err != nil {
		return fmt.Errorf("extracting CIF data: %v", err)
	}

	err = pdb.ExtractChains()
	if err != nil {
		return fmt.Errorf("extracting PDB atoms to chains: %v", err)
	}

	end := time.Since(start)
	log.Printf("PDB %s loaded in %d msecs", pdb.ID, end.Milliseconds())

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
