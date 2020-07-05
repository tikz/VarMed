package secondary

import (
	"fmt"
	"time"
	"varq/pdb"
	"varq/secondary/tango"
	"varq/uniprot"
)

// Results holds the collected data in the secondary structure analysis step
type Results struct {
	Tango    []*MappedTangoResidue `json:"tango"`
	Duration time.Duration         `json:"duration"`
	Error    error                 `json:"error"`
}

// Run starts the secondary structure analysis step
func Run(unp *uniprot.UniProt, pdb *pdb.PDB, results chan<- *Results, msg func(string)) {
	start := time.Now()

	tangoResidues, err := RunTango(unp, pdb)
	if err != nil {
		results <- &Results{Error: fmt.Errorf("Tango: %v", err)}
	}

	results <- &Results{
		Tango:    tangoResidues,
		Duration: time.Since(start),
	}
}

type MappedTangoResidue struct {
	Position int64               `json:"position"`
	Results  *tango.TangoResidue `json:"results"`
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
			tangoResidues, err := tango.Run(unp.ID+chain.ChainID, seq)
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
