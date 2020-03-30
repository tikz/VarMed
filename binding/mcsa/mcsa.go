package mcsa

import (
	"encoding/json"
	"net/url"
	"varq/http"
	"varq/pdb"
)

// CatalyticResidues holds the protein's residues that have catalytic activity according to M-CSA.
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

// GetCSA queries M-CSA for the catalytic positions of the given UniProt ID
func GetCSA(pdb *pdb.PDB) (*Catalytic, error) {
	url := "https://www.ebi.ac.uk/thornton-srv/m-csa/api/entries/?" + url.Values{
		"format":                                 {"json"},
		"entries.proteins.sequences.uniprot_ids": {pdb.UniProtID},
	}.Encode()

	raw, err := http.Get(url)
	if err != nil {
		return nil, nil // TODO: handle 404
	}

	response := searchAPIResponse{}
	err = json.Unmarshal(raw, &response)
	if err != nil {
		return nil, err
	}

	if response.Count == 0 {
		return nil, nil
	}

	cs := Catalytic{}
	for _, res := range response.Results[0].Residues {
		for _, seq := range res.ResidueSequences {
			cs.UniProtPositions = append(cs.UniProtPositions, seq.ResID)
			if resInPos, ok := pdb.UniProtPositions[seq.ResID]; ok {
				for _, res := range resInPos {
					cs.Residues = append(cs.Residues, res)
				}
			}
		}
	}

	return &cs, nil
}
