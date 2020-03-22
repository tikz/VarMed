package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type arrayFlags []string

func (i *arrayFlags) String() string {
	return ""
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func cliRun(uniprotFlag arrayFlags, pdbFlag arrayFlags) {
	var analyses []*Analysis

	analyses, err := RunPipelineForPDBs(pdbFlag)
	if err != nil {
		log.Fatal(err)
	}

	for _, uniprotID := range uniprotFlag {
		a, err := RunPipelineForUniProt(uniprotID)
		if err != nil {
			log.Fatal(err)
		}
		analyses = append(analyses, a...)
	}

	dumpJSON(analyses)
}

func dumpJSON(analyses []*Analysis) {
	for _, analysis := range analyses {
		out, _ := json.MarshalIndent(analysis, "", "\t")
		err := ioutil.WriteFile(analysis.PDB.ID+".json", out, 0644)
		if err != nil {
			log.Fatal(err)
		}
	}
}