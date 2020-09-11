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

func cliRun(uniprotID string, pdbFlags arrayFlags, sas []string) {
	j := NewJob(&JobRequest{
		UniProtID: uniprotID,
		PDBIDs:    pdbFlags,
		SAS:       sas,
	})

	fmt.Println("RespDB CLI")
	fmt.Println()
	fmt.Printf("UniProt ID: \t %s\n", uniprotID)
	fmt.Printf("PDB IDs: \t %s\n", pdbFlags)
	fmt.Printf("SAS: \t\t %s\n", sas)
	fmt.Printf("Job hash: \t %s...\n", j.ID[:10])
	fmt.Println()

	j.Process(true)
	if j.Error != nil {
		log.Fatal(j.Error)
	}
	out, _ := json.MarshalIndent(j.Pipeline.Results, "", "\t")
	ioutil.WriteFile("output.json", out, 0644)
}
