package uniprot

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"varq/http"
)

// UniProt contains relevant protein data for a single accession.
type UniProt struct {
	ID       string          `json:"id"`       // accession ID
	URL      string          `json:"url"`      // page URL for the entry
	TXTURL   string          `json:"txtUrl"`   // TXT API URL for the entry.
	Name     string          `json:"name"`     // protein name
	Gene     string          `json:"gene"`     // gene code
	Organism string          `json:"organism"` // organism
	Sequence string          `json:"sequence"` // canonical sequence
	PDBIDs   []string        `json:"pdbIds"`   // PDB IDs
	Variants []*VariantEntry `json:"variants"` // dbSNP variants
	Raw      []byte          `json:"-"`        // TXT API raw bytes.
}

// VariantEntry represents a single variant entry extracted from the TXT.
type VariantEntry struct {
	Position int64  `json:"position"`
	Note     string `json:"note"`
	Evidence string `json:"evidence"`
	ID       string `json:"id"`
	SAS      *SAS   `json:"sas"`
}

// SAS represents a single aminoacid substitution.
type SAS struct {
	Position int64  `json:"position"`
	FromAa   string `json:"fromAa"`
	ToAa     string `json:"toAa"`
}

// NewUniProt constructs an instance from an UniProt accession ID and a list of target PDB IDs
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

// extract parses the TXT response.
func (u *UniProt) extract() error {
	err := u.extractSequence()
	if err != nil {
		return fmt.Errorf("get seq: %v", err)
	}

	err = u.extractPDBs()
	if err != nil {
		return fmt.Errorf("extracting crystals from UniProt TXT: %v", err)
	}

	err = u.extractNames()
	if err != nil {
		return fmt.Errorf("extracting names from UniProt TXT: %v", err)
	}

	err = u.extractVariants()
	if err != nil {
		return fmt.Errorf("extracting variants from UniProt TXT: %v", err)
	}

	return nil
}

// extractPDBs parses the TXT for PDB IDs and populates UniProt.PDBs
func (u *UniProt) extractPDBs() error {
	// Regex match all PDB IDs in the UniProt TXT entry. X-ray only, ignore others (NMR, etc).
	// https://regex101.com/r/BpJ3QB/1
	r, _ := regexp.Compile("PDB;[ ]*(.*?);[ ]*(X.*?ray);[ ]*([0-9\\.]*).*?;.*?\n")
	matches := r.FindAllStringSubmatch(string(u.Raw), -1)

	// Parse each PDB match in TXT
	for _, m := range matches {
		u.PDBIDs = append(u.PDBIDs, m[1])
	}

	return nil
}

// extractSequence parses the canonical sequence.
func (u *UniProt) extractSequence() error {
	r, _ := regexp.Compile("(?ms)^SQ.*?$(.*?)//") // https://regex101.com/r/ZTOYaJ/1
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

// extractNames parses protein, gene and organism names
func (u *UniProt) extractNames() error {
	r, _ := regexp.Compile("(?m)^DE.*?Name.*?Full=(.*?)(;| {)")
	matches := r.FindAllStringSubmatch(string(u.Raw), -1)

	if len(matches) == 0 {
		return errors.New("protein name not found")
	}
	u.Name = matches[0][1]

	r, _ = regexp.Compile("(?m)^GN.*?=(.*?);")
	matches = r.FindAllStringSubmatch(string(u.Raw), -1)

	if len(matches) == 0 {
		return errors.New("gene name not found")
	}
	u.Gene = matches[0][1]

	r, _ = regexp.Compile("(?m)^OS[ ]+(.*?)\\.")
	matches = r.FindAllStringSubmatch(string(u.Raw), -1)

	if len(matches) == 0 {
		return errors.New("organism name not found")
	}
	u.Organism = matches[0][1]

	return nil
}

// extractVariants parses for variant references
func (u *UniProt) extractVariants() error {
	var variants []*VariantEntry

	// https://regex101.com/r/BpJ3QB/1
	r, _ := regexp.Compile("(?ms)^FT[ ]*VARIANT[ ]*([0-9]*)$(.*?)id=\"(.*?)\"")
	matches := r.FindAllStringSubmatch(string(u.Raw), -1)

	for _, variant := range matches {
		pos, err := strconv.ParseInt(variant[1], 10, 64)
		if err != nil {
			return fmt.Errorf("cannot parse variant position int: %s", variant[1])
		}

		data := variant[2]
		s := regexp.MustCompile("\nFT \\s+")
		d := s.ReplaceAllString(data, " ")

		r, _ := regexp.Compile("(?s)/note=\"(.*?)\"")
		n := r.FindAllStringSubmatch(d, -1)
		note := n[0][1]

		r, _ = regexp.Compile("(?s)/evidence=\"(.*?)\"")
		e := r.FindAllStringSubmatch(d, -1)
		var evidence string
		if len(e) > 0 {
			evidence = e[0][1]
		}

		id := variant[3]

		variants = append(variants, &VariantEntry{
			Position: pos,
			ID:       id,
			Note:     note,
			Evidence: evidence,
		})
	}

	u.Variants = variants

	return nil
}

// PDBIDExists returns true if the given PDB ID is included in this
// UniProt entry, false otherwise.
func (u *UniProt) PDBIDExists(pdbID string) bool {
	for _, id := range u.PDBIDs {
		if id == pdbID {
			return true
		}
	}
	return false
}
