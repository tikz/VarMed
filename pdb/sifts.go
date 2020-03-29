package pdb

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"varq/http"
)

// https://www.ebi.ac.uk/pdbe/api/doc/sifts.html

type SIFTS struct {
	UniProtIDs map[string]*SIFTSUniProt
}

type SIFTSUniProt struct {
	Chains map[string]*SIFTSMapping
}

type SIFTSMapping struct {
	UniProtStart int64
	UniProtEnd   int64
	PDBStart     int64
	PDBEnd       int64
}

// JSON response unmarshaling structs

type ResponseUniProtAccession struct {
	Identifier string            `json:"identifier"`
	Mappings   []ResponseMapping `json:"mappings"`
	Name       string            `json:"name"`
}

type ResponseMapping struct {
	Start        ResponsePos `json:"start"`
	EntityID     int64       `json:"entity_id"`
	End          ResponsePos `json:"end"`
	UnpStart     int64       `json:"unp_start"`
	UnpEnd       int64       `json:"unp_end"`
	ChainID      string      `json:"chain_id"`
	StructAsymID string      `json:"struct_asym_id"`
}

type ResponsePos struct {
	ResidueNumber int64 `json:"residue_number"`
}

func (pdb *PDB) GetSIFTSMappings() error {
	pdbID := strings.ToLower(pdb.ID)
	raw, _ := http.Get("https://www.ebi.ac.uk/pdbe/api/mappings/uniprot_segments/" + pdbID)
	pdbs := make(map[string]json.RawMessage)
	err := json.Unmarshal(raw, &pdbs)
	if err != nil {
		return fmt.Errorf("unmarshal: %v", err)
	}

	databases := make(map[string]json.RawMessage)
	err = json.Unmarshal(pdbs[pdbID], &databases)
	if err != nil {
		return fmt.Errorf("unmarshal databases keys: %v", err)
	}

	accessions := make(map[string]*ResponseUniProtAccession)
	err = json.Unmarshal(databases["UniProt"], &accessions)
	if err != nil {
		return fmt.Errorf("unmarshal UniProt key: %v", err)
	}

	if len(accessions) == 0 {
		return errors.New("UniProt accession not found")
	}

	s := make(map[string]*SIFTSUniProt)
	for accessionID, accession := range accessions {
		s[accessionID] = &SIFTSUniProt{Chains: make(map[string]*SIFTSMapping)}
		for _, mapping := range accession.Mappings {
			s[accessionID].Chains[mapping.StructAsymID] = &SIFTSMapping{
				UniProtStart: mapping.UnpStart,
				UniProtEnd:   mapping.UnpEnd,
				PDBStart:     mapping.Start.ResidueNumber,
				PDBEnd:       mapping.End.ResidueNumber,
			}
		}
	}

	pdb.SIFTS = &SIFTS{UniProtIDs: s}
	return nil
}
