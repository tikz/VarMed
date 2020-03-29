package pdb

import (
	"fmt"
	"log"
	"time"
	"varq/http"
)

type PDB struct {
	ID     string
	URL    string
	PDBURL string
	CIFURL string

	Title       string
	Date        *time.Time
	Method      string
	Resolution  float64
	TotalLength int64

	UniProtID       string // UniProt accession. If the PDB is a complex of multiple proteins, this defines the chains of interest.
	UniProtSequence string // UniProt choosen canonical sequence.

	SIFTS            *SIFTS                        // EBI SIFTS data for residue position mapping between UniProt and PDB.
	Chains           map[string]map[int64]*Residue `json:"-"` // PDB ATOM chain name and position to Residue pointer.
	SeqRes           map[string][]*Residue         `json:"-"`
	SeqResChains     map[string]map[int64]*Residue `json:"-"` // PDB SEQRES chain name and position to Residue pointer in structure.
	SeqResOffsets    map[string]int64              `json:"-"` // PDB ATOM residue number to SEQRES position offsets.
	UniProtPositions map[int64][]*Residue          `json:"-"` // UniProt sequence position to Residue pointer(s) in structure. Multiple chains can come from same positions in the sequence.

	RawPDB        []byte `json:"-"`
	RawCIF        []byte `json:"-"`
	LocalPath     string
	LocalFilename string

	Error error
}

type Chain struct {
	UniProtName string
	PDBName     string

	Residues map[int64]*Residue
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
		return
	}

	err = pdb.ExtractSeqRes()
	if err != nil {
		pdb.Error = fmt.Errorf("extracting SEQRES: %v", err)
		return
	}

	err = pdb.ExtractChains()
	if err != nil {
		pdb.Error = fmt.Errorf("extracting PDB atoms to chains: %v", err)
		return
	}

	err = pdb.ExtractCIFData()
	if err != nil {
		pdb.Error = fmt.Errorf("extracting CIF data: %v", err)
		return
	}

	pdb.calculateChainsOffset()
	pdb.makeMappings()

	end := time.Since(start)
	log.Printf("PDB %s obtained in %.3f secs", pdb.ID, end.Seconds())
}
