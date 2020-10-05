package config

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// Config holds constant parameters for the app
type Config struct {
	HTTPClient struct {
		UserAgent string `yaml:"user-agent"`
		Timeout   int    `yaml:"timeout"`
	} `yaml:"http-client"`

	HTTPServer struct {
		Port string `yaml:"port"`
	} `yaml:"http-server"`

	RespDB struct {
		JobWorkers int `yaml:"job-workers"`
		Pipeline   struct {
			StructureWorkers int `yaml:"structure-workers"`
			EnableSteps      struct {
				Binding      bool `yaml:"binding"`
				Exposure     bool `yaml:"exposure"`
				Interaction  bool `yaml:"interaction"`
				Secondary    bool `yaml:"secondary"`
				Conservation bool `yaml:"conservation"`
				Stability    bool `yaml:"stability"`
			} `yaml:"enable-steps"`
		} `yaml:"pipeline"`
	} `yaml:"respdb"`

	DebugPrint struct {
		Enabled bool `yaml:"enabled"`
		Rulers  struct {
			UniProt bool `yaml:"uniprot"`
			PDB     bool `yaml:"pdb"`
		} `yaml:"rulers"`
	} `yaml:"debug-print"`

	Paths struct {
		Data           string `yaml:"data"`
		UniProt        string `yaml:"uniprot"`
		PDB            string `yaml:"pdb"`
		Jobs           string `yaml:"jobs"`
		Fpocket        string `yaml:"fpocket"`
		ClinVar        string `yaml:"clinvar"`
		Pfam           string `yaml:"pfam"`
		FoldXBin       string `yaml:"foldx-bin"`
		FoldXRepair    string `yaml:"foldx-repair"`
		FoldXMutations string `yaml:"foldx-mutations"`
		Abswitch       string `yaml:"abswitch"`
		FileExt        string `yaml:"file-ext"`
	} `yaml:"paths"`
}

// LoadFile opens and parses the YAML config file
func LoadFile(path string) (*Config, error) {
	f, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		return nil, fmt.Errorf("opening config.yaml: %v", err)
	}

	cfg := Config{}
	err = yaml.Unmarshal([]byte(f), &cfg)
	if err != nil {
		return nil, fmt.Errorf("parsing config.yaml: %v", err)
	}
	return &cfg, nil
}
