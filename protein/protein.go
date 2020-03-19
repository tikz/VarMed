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
	UniProt  *UniProt
	Crystals []*pdb.PDB
	// BestCrystal *pdb.PDB
}

// UniProt contains general protein data retrieved from UniProt
type UniProt struct {
	ID     string
	URL    string
	TXTURL string
	Raw    []byte `json:"-"`
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
			ID:     uniprotID,
			URL:    url,
			TXTURL: txtURL,
			Raw:    raw,
		},
	}

	// Parse UniProt TXT
	err = p.extract()
	if err != nil {
		return nil, fmt.Errorf("extract PDB crystals %v: %v", uniprotID, err)
	}

	return &p, nil
}

func (p *Protein) extract() error {
	crystals, err := p.extractCrystals()
	if err != nil {
		return fmt.Errorf("extracting crystals from UniProt TXT: %v", err)
	}
	p.Crystals = crystals

	// bestCrystal, err := pickBestCrystal(p.Crystals)
	// if err != nil {
	// 	return fmt.Errorf("choosing best crystal: %v", err)
	// }
	// p.BestCrystal = bestCrystal

	return nil
}

func (p *Protein) extractCrystals() (crystals []*pdb.PDB, err error) {
	// Regex match all PDB IDs in the UniProt TXT entry. X-ray only, ignore others (NMR, etc).
	// https://regex101.com/r/BpJ3QB/1
	regexPDB, _ := regexp.Compile("PDB;[ ]*(.*?);[ ]*(X.*?ray);[ ]*([0-9\\.]*).*?;.*?\n")
	matches := regexPDB.FindAllStringSubmatch(string(p.UniProt.Raw), -1)
	if len(matches) == 0 {
		return nil, errors.New("UniProt entry has no associated crystal PDB entries")
	}

	// Parse each PDB match in TXT
	for _, match := range matches {
		resolution, err := strconv.ParseFloat(match[3], 64)
		if err != nil {
			return nil, fmt.Errorf("parsing resolution %v as float: %v", resolution, err)
		}

		// Extract start and end positions for each chain, and calculate length sum
		regexChains, _ := regexp.Compile("=([0-9]*)-([0-9]*)")
		chains := regexChains.FindAllStringSubmatch(match[0], -1)
		var totalLength int64
		for _, chain := range chains {
			fromPos, err := strconv.ParseInt(chain[1], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("parsing from position %v as int: %v", fromPos, err)
			}

			toPos, err := strconv.ParseInt(chain[2], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("parsing to position %v as int: %v", toPos, err)
			}
			totalLength += toPos - fromPos + 1
		}

		crystal := pdb.PDB{
			UniProtID:  p.UniProt.ID,
			ID:         match[1],
			Method:     match[2],
			Resolution: resolution,
			Length:     totalLength,
		}
		crystals = append(crystals, &crystal)
	}

	return crystals, nil
}

// pickBestCrystal picks the best crystal to our criteria from the available ones
func pickBestCrystal(crystals []*pdb.PDB) (*pdb.PDB, error) {
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
