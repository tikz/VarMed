package uniprot

import (
	"errors"
	"fmt"
	"regexp"
	"varq/http"
	"varq/pdb"
)

// UniProt contains general protein data retrieved from UniProt
type UniProt struct {
	ID       string
	URL      string
	TXTURL   string
	Raw      []byte     `json:"-"`
	Crystals []*pdb.PDB `json:"-"`
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

func (u *UniProt) extract() error {
	crystals, err := u.extractCrystals()
	if err != nil {
		return fmt.Errorf("extracting crystals from UniProt TXT: %v", err)
	}
	u.Crystals = crystals

	return nil
}

func (u *UniProt) extractCrystals() (crystals []*pdb.PDB, err error) {
	// Regex match all PDB IDs in the UniProt TXT entry. X-ray only, ignore others (NMR, etc).
	// https://regex101.com/r/BpJ3QB/1
	regexPDB, _ := regexp.Compile("PDB;[ ]*(.*?);[ ]*(X.*?ray);[ ]*([0-9\\.]*).*?;.*?\n")
	matches := regexPDB.FindAllStringSubmatch(string(u.Raw), -1)
	if len(matches) == 0 {
		return nil, errors.New("UniProt entry has no associated crystal PDB entries")
	}

	// Parse each PDB match in TXT
	for _, match := range matches {
		crystal := pdb.PDB{
			ID: match[1],
		}
		crystals = append(crystals, &crystal)
	}

	return crystals, nil
}
