package uniprot

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"varq/http"
	"varq/pdb"
)

// UniProt contains relevant protein data for a single accession.
type UniProt struct {
	ID       string     // accession ID
	URL      string     // page URL for the entry
	TXTURL   string     // TXT API URL for the entry.
	Sequence string     // canonical sequence
	Raw      []byte     `json:"-"` // TXT API raw bytes.
	PDBIDs   []string   // associated PDB IDs
	PDBs     []*pdb.PDB `json:"-"` // associated PDB entries
}

// NewUniProt constructs a Protein instance from an UniProt accession ID
func NewUniProt(uniprotID string) (*UniProt, error) {
	url := "https://www.uniprot.org/uniprot/" + uniprotID
	txtURL := url + ".txt"
	raw, err := http.Get(txtURL)
	if err != nil {
		return nil, fmt.Errorf("get UniProt accession %v: %v", uniprotID, err)
	}

	u := &UniProt{
		ID:     uniprotID,
		URL:    url,
		TXTURL: txtURL,
		Raw:    raw,
	}

	// Parse UniProt TXT
	err = u.extract()
	if err != nil {
		return nil, fmt.Errorf("extract PDB crystals %v: %v", uniprotID, err)
	}

	return u, nil
}

// FilterPDBs removes PDBIDs that aren't in the given slice.
func (u *UniProt) FilterPDBs(PDBIDs []string) {
	var newCrystals []*pdb.PDB
	for _, filterID := range PDBIDs {
		for _, crystal := range u.PDBs {
			if strings.ToLower(filterID) == strings.ToLower(crystal.ID) {
				newCrystals = append(newCrystals, crystal)
			}
		}
	}

	u.PDBs = newCrystals
}

// extract launches all the parsing to be done in the TXT response.
func (u *UniProt) extract() error {
	err := u.extractPDBs()
	if err != nil {
		return fmt.Errorf("extracting crystals from UniProt TXT: %v", err)
	}

	err = u.extractSequence()
	if err != nil {
		return fmt.Errorf("get seq: %v", err)
	}

	return nil
}

// extractPDBs parses the TXT for PDB IDs, and populates UniProt.PDBs
func (u *UniProt) extractPDBs() error {
	// Regex match all PDB IDs in the UniProt TXT entry. X-ray only, ignore others (NMR, etc).
	// https://regex101.com/r/BpJ3QB/1
	r, _ := regexp.Compile("PDB;[ ]*(.*?);[ ]*(X.*?ray);[ ]*([0-9\\.]*).*?;.*?\n")
	matches := r.FindAllStringSubmatch(string(u.Raw), -1)
	if len(matches) == 0 {
		return errors.New("UniProt entry has no associated crystal PDB entries")
	}

	// Parse each PDB match in TXT
	var PDBs []*pdb.PDB
	for _, m := range matches {
		pdb := pdb.PDB{
			ID:              m[1],
			UniProtID:       u.ID,
			UniProtSequence: u.Sequence,
		}
		PDBs = append(PDBs, &pdb)
	}
	u.PDBs = PDBs

	return nil
}

// extractSequence parses the TXT for the canonical sequence.
func (u *UniProt) extractSequence() error {
	r, _ := regexp.Compile("(?ms)SQ.*?$(.*?)//") // https://regex101.com/r/ZTOYaJ/1
	matches := r.FindAllStringSubmatch(string(u.Raw), -1)

	if len(matches) == 0 {
		return errors.New("canonical sequence not found")
	}

	seqGroup := matches[0][1]
	sequence := strings.ReplaceAll(seqGroup, " ", "")
	sequence = strings.ReplaceAll(sequence, "\n", "")

	u.Sequence = sequence

	return nil
}
