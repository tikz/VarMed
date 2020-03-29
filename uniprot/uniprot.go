package uniprot

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"varq/http"
	"varq/pdb"
)

// UniProt contains general protein data retrieved from UniProt
type UniProt struct {
	ID       string
	URL      string
	TXTURL   string
	Sequence string
	Raw      []byte     `json:"-"`
	Crystals []*pdb.PDB `json:"-"`
}

type SeqPosMap struct {
	Chains []string // Chain name in PDB
	Start  int      // Start pos compared to UniProt canonical seq
	End    int      // End pos compared to UniProt canonical seq
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

	// Get canonical sequence
	err = u.getSequence()
	if err != nil {
		return nil, fmt.Errorf("get seq %v: %v", uniprotID, err)
	} // TODO: see if its not in the TXT

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
			ID:              match[1],
			UniProtID:       u.ID,
			UniProtSequence: u.Sequence,
		}
		crystals = append(crystals, &crystal)
	}

	return crystals, nil
}

func (u *UniProt) getSequence() error {
	url := "https://www.uniprot.org/uniprot/" + u.ID + ".fasta"
	raw, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("get UniProt FASTA seq: %v", err)
	}

	regexSeq, _ := regexp.Compile(">(?ms).*?$(.*)")
	matches := regexSeq.FindAllStringSubmatch(string(raw), -1)
	if len(matches) == 0 {
		return errors.New("cannot parse FASTA")
	}
	seq := strings.ReplaceAll(matches[0][1], " ", "")
	seq = strings.ReplaceAll(seq, "\n", "")

	u.Sequence = seq
	return nil
}

func (u *UniProt) CleanCrystals() {
	var newCrystals []*pdb.PDB
	for _, crystal := range u.Crystals {
		_, uniprotExistsInSIFTS := crystal.SIFTS.UniProtIDs[u.ID]
		if crystal.SIFTS != nil && // Has SIFTS data
			uniprotExistsInSIFTS { // Requested UniProt ID in SIFTS
			newCrystals = append(newCrystals, crystal)
		}
	}
	u.Crystals = newCrystals
}
