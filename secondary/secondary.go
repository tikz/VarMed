package secondary

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sync"
	"time"
	"varq/pdb"
	"varq/secondary/abswitch"
	"varq/secondary/dssp"
	"varq/secondary/tango"
	"varq/uniprot"
)

var fileMux sync.Mutex

// Results holds the collected data in the secondary structure analysis step
type Results struct {
	Tango    []*MappedTangoResidue    `json:"tango"`
	AbSwitch []*MappedAbSwitchResidue `json:"abswitch"`
	Duration time.Duration            `json:"duration"`
	Error    error                    `json:"error"`
}

// Run starts the secondary structure analysis step
func Run(unp *uniprot.UniProt, pdb *pdb.PDB, results chan<- *Results, msg func(string)) {
	start := time.Now()

	err := dssp.RunDSSP(pdb)
	if err != nil {
		results <- &Results{Error: fmt.Errorf("DSSP: %v", err)}
	}

	tangoResidues, err := RunTango(unp, pdb)
	if err != nil {
		results <- &Results{Error: fmt.Errorf("tango: %v", err)}
	}

	abswitchResidues, err := RunAbSwitch(unp, pdb)
	if err != nil {
		results <- &Results{Error: fmt.Errorf("abswitch: %v", err)}
	}

	results <- &Results{
		Tango:    tangoResidues,
		AbSwitch: abswitchResidues,
		Duration: time.Since(start),
	}
}

type MappedTangoResidue struct {
	Position int64          `json:"position"`
	Results  *tango.Residue `json:"results"`
}

func RunTango(unp *uniprot.UniProt, pdb *pdb.PDB) ([]*MappedTangoResidue, error) {
	var seqs []string
	var results []*MappedTangoResidue

	isUnique := func(seqs []string, seq string) bool {
		for _, s := range seqs {
			if s == seq {
				return false
			}
		}
		return true
	}

	// Run tango for each chain in structure
	// with unique sequence (don't rerun it)
	for _, chain := range pdb.SIFTS.UniProt[unp.ID].Mappings {
		seq := unp.Sequence[chain.UnpStart:chain.UnpEnd]
		if isUnique(seqs, seq) {
			name := fmt.Sprintf("%s-%s-%s", unp.ID, pdb.ID, chain.ChainID)
			tangoResidues, err := tango.Run(name, seq)
			if err != nil {
				return nil, err
			}

			for i, residue := range tangoResidues {
				results = append(results, &MappedTangoResidue{
					Position: chain.UnpStart + int64(i),
					Results:  residue,
				})
			}
		}
		seqs = append(seqs, seq)
	}
	return results, nil
}

type MappedAbSwitchResidue struct {
	Position int64             `json:"position"`
	Results  *abswitch.Residue `json:"results"`
}

func generateID(seq string) string {
	hash := sha256.Sum256([]byte(seq))
	return hex.EncodeToString(hash[:])
}

func RunAbSwitch(unp *uniprot.UniProt, pdb *pdb.PDB) ([]*MappedAbSwitchResidue, error) {
	var seqs []string
	var results []*MappedAbSwitchResidue

	isUnique := func(seqs []string, seq string) bool {
		for _, s := range seqs {
			if s == seq {
				return false
			}
		}
		return true
	}

	// Run abswitch for each chain in structure
	// with unique sequence (don't rerun it)
	for _, chain := range pdb.SIFTS.UniProt[unp.ID].Mappings {
		seq := unp.Sequence[chain.UnpStart:chain.UnpEnd]
		if isUnique(seqs, seq) {
			name := generateID(seq)
			fileMux.Lock()
			abswitchResidues, err := abswitch.Run(name, seq)
			fileMux.Unlock()
			if err != nil {
				return nil, err
			}

			for i, residue := range abswitchResidues {
				results = append(results, &MappedAbSwitchResidue{
					Position: chain.UnpStart + int64(i),
					Results:  residue,
				})
			}
		}
		seqs = append(seqs, seq)
	}
	return results, nil
}
