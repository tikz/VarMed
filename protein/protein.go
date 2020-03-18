package protein

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"varq/http"
	"varq/protein/pdb"
)

// Protein contains all the raw and parsed data for a protein
type Protein struct {
	UniProt     *UniProt
	Crystals    []*pdb.PDB `json:"-"`
	BestCrystal *pdb.PDB
}

// UniProt contains general protein data retrieved from UniProt
type UniProt struct {
	AccessionID string
	URL         string
	TXTURL      string
	Raw         []byte `json:"-"`
}

// NewProtein constructs a Protein instance from an UniProt accession ID
func NewProtein(uniprotID string) (*Protein, error) {
	url := "https://www.uniprot.org/uniprot/" + uniprotID
	txtURL := url + ".txt"
	raw, err := http.Get(txtURL)
	if err != nil {
		return nil, fmt.Errorf("get UniProt accession %v: %v", uniprotID, err)
	}

	p := Protein{
		UniProt: &UniProt{
			AccessionID: uniprotID,
			URL:         url,
			TXTURL:      txtURL,
			Raw:         raw,
		},
	}

	// Parse UniProt TXT
	err = p.extractPDBIDs()
	if err != nil {
		return nil, fmt.Errorf("extract PDB crystals %v: %v", uniprotID, err)
	}

	return &p, nil
}

func (p *Protein) extractPDBIDs() error {
	// Regex match all PDB IDs in the UniProt TXT entry. X-ray only, ignore others (NMR, etc).
	// https://regex101.com/r/QCI3cu/3
	regexpPDB, _ := regexp.Compile("PDB;[ ]*(.*?);[ ]*(X.*?ray);[ ]*([0-9\\.]*).*?;[ ]*(.*?)=([0-9]*)-([0-9]*)")
	matches := regexpPDB.FindAllStringSubmatch(string(p.UniProt.Raw), -1)
	if len(matches) == 0 {
		return errors.New("UniProt entry has no associated crystal PDB entries")
	}

	// Parse each match
	for _, match := range matches {
		resolution, err := strconv.ParseFloat(match[3], 64)
		if err != nil {
			return fmt.Errorf("parsing resolution %v as float: %v", resolution, err)
		}

		from, err := strconv.ParseInt(match[5], 10, 64)
		if err != nil {
			return fmt.Errorf("parsing from position %v as int: %v", from, err)
		}

		to, err := strconv.ParseInt(match[6], 10, 64)
		if err != nil {
			return fmt.Errorf("parsing to position %v as int: %v", to, err)
		}

		crystal := pdb.PDB{
			ID:         match[1],
			Method:     match[2],
			Resolution: resolution,
			FromPos:    from,
			ToPos:      to,
			Length:     to - from + 1,
		}
		p.Crystals = append(p.Crystals, &crystal)
	}

	bestCrystal, err := decideBestCrystal(p.Crystals)
	if err != nil {
		return fmt.Errorf("choosing best crystal: %v", err)
	}
	p.BestCrystal = bestCrystal

	// TODO: for now just fetch the PDB data for the best crystal.
	// decide if it's worth it to grab all crystals for use in other projects.
	err = p.BestCrystal.Fetch()
	if err != nil {
		return fmt.Errorf("fetching crystal: %v", err)
	}

	return nil
}

// decideBestCrystal picks the best crystal to our criteria from the available ones
func decideBestCrystal(crystals []*pdb.PDB) (*pdb.PDB, error) {
	bestCovLength := crystals[0].Length
	var bestCovCrystals []*pdb.PDB

	// One or more crystals with the same best coverage
	for _, crystal := range crystals {
		if crystal.Length >= 20 {
			if crystal.Length > bestCovLength {
				bestCovLength = crystal.Length
				bestCovCrystals = []*pdb.PDB{crystal}
			} else {
				if crystal.Length == bestCovLength {
					bestCovCrystals = append(bestCovCrystals, crystal)
				}
			}
		}
	}

	if len(bestCovCrystals) == 0 {
		return nil, errors.New("no suitable crystal with >= 20 aminoacids found")
	}

	// From the crystals with the best coverage, pick the one with the best resolution
	var bestResCrystal *pdb.PDB = bestCovCrystals[0]
	for _, crystal := range bestCovCrystals {
		if crystal.Resolution < bestResCrystal.Resolution {
			bestResCrystal = crystal
		}
	}

	// TODO: !! a reasonable structure for having multiple PDB files with different ligands,
	// also complexes (maybe Protein{} pointers)

	return bestResCrystal, nil
}
