package mcsa

import (
	"encoding/json"
	"fmt"
	"net/url"
	"varq/http"
	"varq/pdb"
	"varq/uniprot"
)

// Catalytic holds the protein's residues that have catalytic activity according to M-CSA.
type Catalytic struct {
	UniProtPositions []int64
	Residues         []*pdb.Residue // Pointers to original residues in the requested structure.
}

type searchAPIResponse struct {
	Count    int64       `json:"count"`
	Next     interface{} `json:"next"`
	Previous interface{} `json:"previous"`
	Results  []result    `json:"results"`
}

type result struct {
	Residues []residue `json:"residues"`
}

type residue struct {
	MCSAID              int64             `json:"mcsa_id"`
	RolesSummary        string            `json:"roles_summary"`
	FunctionLocationAbv string            `json:"function_location_abv"`
	ResidueChains       []residueChain    `json:"residue_chains"`
	ResidueSequences    []residueSequence `json:"residue_sequences"`
}

type residueChain struct {
	ChainName         string `json:"chain_name"`
	PdbID             string `json:"pdb_id"`
	AssemblyChainName string `json:"assembly_chain_name"`
	Assembly          int64  `json:"assembly"`
	Code              string `json:"code"`
	ResID             int64  `json:"resid"`
	AuthResid         int64  `json:"auth_resid"`
	IsReference       bool   `json:"is_reference"`
	DomainName        string `json:"domain_name"`
	DomainCathID      string `json:"domain_cath_id"`
}
type residueSequence struct {
	UniProtID   string `json:"uniprot_id"`
	Code        string `json:"code"`
	IsReference bool   `json:"is_reference"`
	ResID       int64  `json:"resid"`
}

// GetPositions queries M-CSA and fetches catalytic residue UniProt positions for the given PDB.
func GetPositions(unp *uniprot.UniProt, pdb *pdb.PDB, msg func(string)) (*Catalytic, error) {
	url := "https://www.ebi.ac.uk/thornton-srv/m-csa/api/entries/?" + url.Values{
		"format":                                 {"json"},
		"entries.proteins.sequences.uniprot_ids": {unp.ID},
	}.Encode()

	raw, err := http.Get(url)
	if err != nil {
		return nil, nil // TODO: handle 404
	}
	msg("downloaded M-CSA data")

	response := searchAPIResponse{}
	err = json.Unmarshal(raw, &response)
	if err != nil {
		return nil, err
	}

	cs := Catalytic{}
	if response.Count == 0 {
		msg("no M-CSA residues found")
		return &cs, nil
	}

	for _, res := range response.Results[0].Residues {
		for _, seq := range res.ResidueSequences {
			cs.UniProtPositions = append(cs.UniProtPositions, seq.ResID)
			if resInPos, ok := pdb.UniProtPositions[unp.ID][seq.ResID]; ok {
				for _, res := range resInPos {
					cs.Residues = append(cs.Residues, res)
				}
			}
		}
	}

	msg(fmt.Sprintf("%d M-CSA residues found", len(cs.Residues)))

	return &cs, nil
}
