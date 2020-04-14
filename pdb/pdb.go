package pdb

import (
	"fmt"
	"time"
	"varq/http" // TOOD: decouple this
)

// PDB represents a single PDB entry.
type PDB struct {
	ID     string // PDB ID
	URL    string // RCSB web page URL
	PDBURL string // RCSB download URL for the PDB file
	CIFURL string // RCSB download URL for the CIF file

	Title       string     // publication title
	Date        *time.Time // publication date
	Method      string     // experimental method used
	Resolution  float64    // method resolution
	TotalLength int64      // total length as sum of residues of all chains in the structure

	Atoms    []*Atom // ATOM records in the structure
	HetAtoms []*Atom // HETATM records in the structure

	// Position mapping
	SIFTS            *SIFTS                          // EBI SIFTS data for residue position mapping
	Chains           map[string]map[int64]*Residue   // PDB ATOM chain name and position to pointer in structure
	SeqRes           map[string][]*Residue           // PDB SEQRES chain name to residue pointers
	SeqResChains     map[string]map[int64]*Residue   // PDB SEQRES chain name and PDB ATOM position to residue in structure
	SeqResOffsets    map[string]int64                // PDB ATOM residue number to SEQRES position offsets
	UniProtPositions map[string]map[int64][]*Residue // UniProt ID to sequence position to residue(s) (multiple chains) in structure

	// Extra data
	// SITE records
	BindingSite map[string][]*Residue // binding site identifier to residues compromising it

	// REMARK 800 site descriptions
	BindingSiteDesc map[string]string // binding site identifier to description

	RawPDB        []byte // PDB file raw data
	RawCIF        []byte // CIF file raw data
	LocalPath     string // local path for the PDB file
	LocalFilename string // local filename for the PDB file
}

// NewPDBFromID constructs a new instance from a UniProt accession ID and PDB ID, fetching and parsing the data.
func NewPDBFromID(pdbID string) (PDB, error) {
	pdb := PDB{ID: pdbID}

	err := pdb.Load()
	return pdb, err
}

// NewPDBNoMetadata constructs a new instance from raw bytes, and only extracts ATOM records.
// This is useful for parsing PDB output files generated from external tools.
func NewPDBNoMetadata(raw []byte) (*PDB, error) {
	pdb := PDB{RawPDB: raw}

	err := pdb.ExtractResidues()
	if err != nil {
		return nil, fmt.Errorf("parse: %v", err)
	}

	return &pdb, nil
}

// Load fetches and parses the necessary data.
func (pdb *PDB) Load() error {
	err := pdb.Fetch()
	if err != nil {
		return fmt.Errorf("fetch data: %v", err)
	}

	err = pdb.Parse()
	if err != nil {
		return fmt.Errorf("parse: %v", err)
	}

	return nil
}

func (pdb *PDB) Parse() error {
	err := pdb.Extract()
	if err != nil {
		return fmt.Errorf("extract: %v", err)
	}

	pdb.makeMappings()

	pdb.extractSites()
	return nil
}

// Fetch downloads all external data for the entry.
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

	pdb.URL = url
	pdb.PDBURL = urlPDB
	pdb.CIFURL = urlCIF
	pdb.RawPDB = rawPDB
	pdb.RawCIF = rawCIF

	err = pdb.getSIFTSMappings()
	if err != nil {
		return fmt.Errorf("SIFTS: %v", err)
	}

	return nil
}

// Extract parses data from the raw PDB, raw CIF, SIFTS, and populates the entry.
func (pdb *PDB) Extract() error {
	err := pdb.ExtractSeqRes()
	if err != nil {
		return fmt.Errorf("extract SEQRES: %v", err)
	}

	err = pdb.ExtractResidues()
	if err != nil {
		return fmt.Errorf("extract PDB residues: %v", err)
	}

	err = pdb.ExtractCIFData()
	if err != nil {
		return fmt.Errorf("extract CIF data: %v", err)
	}

	return nil
}

// SeqExactMatchInUniProt returns true if the crystal primary sequence is contained
// and exactly matched per each residue in the canonical UniProt sequence range, false otherwise.
// func (pdb *PDB) SeqExactMatchInUniProt() bool {
// 	for _, m := range pdb.SIFTS.UniProt[pdb.UniProtID].Mappings {
// 		var i int64
// 		for i = m.PDBStart.ResidueNumber; i < m.PDBEnd.ResidueNumber; i++ {
// 			if pdb.Chains[m.ChainID][i].Name != string(pdb.UniProtSequence[i]) {
// 				return false
// 			}
// 		}
// 	}

// 	return true
// }
