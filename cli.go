package main

import (
	"encoding/json"
	"fmt"
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
	msgs := make(chan string, 100) // TODO: check
	go func() {
		for {
			fmt.Println(<-msgs)
		}
	}()
	p, err := NewPipeline(uniprotID, pdbFlags, msgs)
	if err != nil {
		log.Fatal(err)
	}

	err = p.Run()
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
