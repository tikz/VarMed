package pdb

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"varq/http"
)

// Reference: https://www.ebi.ac.uk/pdbe/api/doc/sifts.html

// SIFTS represents a valid response from the SIFTS mapping project.
type SIFTS struct {
	UniProtIDs map[string]*SIFTSUniProt
}

// SIFTSUniProt represents the available UniProt sequence mappings for protein chains.
type SIFTSUniProt struct {
	Chains map[string]*SIFTSMapping
}

// SIFTSMapping represents the position offsets between UniProt and PDB for a specific chain.
type SIFTSMapping struct {
	UniProtStart int64
	UniProtEnd   int64
	PDBStart     int64
	PDBEnd       int64
}

// Private structs for JSON unmarshaling

type responseUniProtAccession struct {
	Identifier string            `json:"identifier"`
	Mappings   []responseMapping `json:"mappings"`
	Name       string            `json:"name"`
}

type responseMapping struct {
	Start        responsePos `json:"start"`
	EntityID     int64       `json:"entity_id"`
	End          responsePos `json:"end"`
	UnpStart     int64       `json:"unp_start"`
	UnpEnd       int64       `json:"unp_end"`
	ChainID      string      `json:"chain_id"`
	StructAsymID string      `json:"struct_asym_id"`
}

type responsePos struct {
	ResidueNumber int64 `json:"residue_number"`
}

// getSIFTSMappings retrieves the UniProt<->PDB position mappings from the SIFTS project.
func (pdb *PDB) getSIFTSMappings() error {
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

	accessions := make(map[string]*responseUniProtAccession)
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
			s[accessionID].Chains[mapping.ChainID] = &SIFTSMapping{
				UniProtStart: mapping.UnpStart,
				UniProtEnd:   mapping.UnpEnd,
				PDBStart:     mapping.Start.ResidueNumber,
				PDBEnd:       mapping.End.ResidueNumber,
			}
		}
	}

	pdb.SIFTS = &SIFTS{UniProtIDs: s}
	if _, ok := pdb.SIFTS.UniProtIDs[pdb.UniProtID]; !ok {
		return fmt.Errorf("no mappings available between %s and %s", pdb.ID, pdb.UniProtID)
	}

	return nil
}
