package pdb

import (
	"fmt"
	"log"
	"time"
	"varq/http"
)

type PDB struct {
	UniProtID        string
	ID               string
	URL              string
	PDBURL           string
	CIFURL           string
	LocalPath        string
	LocalFilename    string
	RawPDB           []byte `json:"-"`
	RawCIF           []byte `json:"-"`
	Title            string
	Date             *time.Time
	Method           string
	Resolution       float64
	TotalLength      int64
	Unpseq           string                          // TODO: delete
	Chains           map[string]map[int64]*Aminoacid `json:"-"` // PDB ATOM chain name and position to Aminoacid pointer.
	SeqRes           map[string][]*Aminoacid         `json:"-"`
	SeqResChains     map[string]map[int64]*Aminoacid `json:"-"` // PDB SEQRES chain name and position to Aminoacid pointer (same instance as PDB).
	UniProtPositions map[int64][]*Aminoacid          `json:"-"` // UniProt primary sequence position to Aminoacid pointers (same instances as PDB).
	ChainsOffsets    map[string]int64                `json:"-"` // ATOM residue number to SEQRES position offsets.
	SIFTS            *SIFTS
	Error            error
}

// Fetch populates the instance with parsed data retrieved from RCSB
func (pdb *PDB) Fetch() {
	start := time.Now()
	log.Printf("Downloading PDB and CIF files for %s", pdb.ID)
	url := "https://www.rcsb.org/structure/" + pdb.ID
	urlCIF := "https://files.rcsb.org/download/" + pdb.ID + ".cif"
	rawCIF, err := http.Get(urlCIF)
	if err != nil {
		pdb.Error = fmt.Errorf("download CIF file: %v", err)
	}

	urlPDB := "https://files.rcsb.org/download/" + pdb.ID + ".pdb"
	rawPDB, err := http.Get(urlPDB)
	if err != nil {
		pdb.Error = fmt.Errorf("download PDB file: %v", err)
	}

	pdb.URL = url
	pdb.PDBURL = urlPDB
	pdb.CIFURL = urlCIF
	pdb.RawPDB = rawPDB
	pdb.RawCIF = rawCIF

	err = pdb.GetSIFTSMappings()
	if err != nil {
		pdb.Error = fmt.Errorf("SIFTS: %v", err)
	}

	err = pdb.ExtractSeqRes()
	if err != nil {
		pdb.Error = fmt.Errorf("extracting SEQRES: %v", err)
	}

	err = pdb.ExtractChains()
	if err != nil {
		pdb.Error = fmt.Errorf("extracting PDB atoms to chains: %v", err)
	}

	err = pdb.ExtractCIFData()
	if err != nil {
		pdb.Error = fmt.Errorf("extracting CIF data: %v", err)
	}

	pdb.calculateChainsOffset()
	pdb.makeMappings()

	end := time.Since(start)
	log.Printf("PDB %s obtained in %.3f secs", pdb.ID, end.Seconds())
}
