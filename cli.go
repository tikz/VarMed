package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"strings"
)

type arrayFlags []string

func (i *arrayFlags) String() string {
	return ""
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, strings.ToUpper(value))
	return nil
}

func cliRun(uniprotID string, pdbFlags arrayFlags) {
	p, err := NewPipeline(uniprotID, pdbFlags, nil)
	if err != nil {
		log.Fatal(err)
	}
	err = p.RunPipeline()
	if err != nil {
		log.Fatal(err)
	}
	// var analyses []*Analysis

	// a, err := RunPipeline(uniprotID, pdbFlags)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// analyses = a

	// dumpJSON(analyses)
}

func dumpJSON(analyses []*Results) {
	for _, analysis := range analyses {
		out, _ := json.MarshalIndent(analysis, "", "\t")
		err := ioutil.WriteFile(analysis.PDB.ID+".json", out, 0644)
		if err != nil {
			log.Fatal(err)
		}
	}
}
