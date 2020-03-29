package mcsa

import (
	"encoding/json"
	"net/url"
	"varq/http"
	"varq/pdb"
)

// CatalyticResidues holds the protein's residues that have catalytic activity according to M-CSA.
type Catalytic struct {
	UniProtChains map[string]map[int64][]*pdb.Aminoacid
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
	MCSAID              int64          `json:"mcsa_id"`
	RolesSummary        string         `json:"roles_summary"`
	FunctionLocationAbv string         `json:"function_location_abv"`
	ResidueChains       []residueChain `json:"residue_chains"`
}

type residueChain struct {
	ChainName         string `json:"chain_name"`
	PdbID             string `json:"pdb_id"`
	AssemblyChainName string `json:"assembly_chain_name"`
	Assembly          int64  `json:"assembly"`
	Code              string `json:"code"`
	Resid             int64  `json:"resid"`
	AuthResid         int64  `json:"auth_resid"`
	IsReference       bool   `json:"is_reference"`
	DomainName        string `json:"domain_name"`
	DomainCathID      string `json:"domain_cath_id"`
}

// GetCSA queries M-CSA for the catalytic positions of the given UniProt ID
func GetCSA(uniprotID string) (*Catalytic, error) {
	url := "https://www.ebi.ac.uk/thornton-srv/m-csa/api/entries/?" + url.Values{
		"format":                                 {"json"},
		"entries.proteins.sequences.uniprot_ids": {uniprotID},
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

	chains := make(map[string]map[int64][]*pdb.Aminoacid)
	cs := Catalytic{UniProtChains: chains}

	for _, res := range response.Results[0].Residues {
		for _, resC := range res.ResidueChains {
			if _, ok := chains[resC.ChainName]; !ok {
				chains[resC.ChainName] = make(map[int64][]*pdb.Aminoacid)
			}
			aa := pdb.NewAminoacid(resC.Resid, resC.Code)
			chains[resC.ChainName][resC.Resid] = append(chains[resC.ChainName][resC.Resid], aa)
		}
	}

	return &cs, nil
}
