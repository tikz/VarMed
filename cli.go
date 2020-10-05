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

func cliRun(uniprotID string, pdbFlags arrayFlags, variants []string) {
	j := NewJob(&JobRequest{
		UniProtID: uniprotID,
		PDBIDs:    pdbFlags,
		Variants:  variants,
	})

	fmt.Println("RespDB CLI")
	fmt.Println()
	fmt.Printf("UniProt ID: \t %s\n", uniprotID)
	fmt.Printf("PDB IDs: \t %s\n", pdbFlags)
	fmt.Printf("Variants: \t %s\n", variants)
	fmt.Printf("Job hash: \t %s...\n", j.ID[:10])
	fmt.Println()

	j.Process(true)
	if j.Error != nil {
		log.Fatal(j.Error)
	}

	out, _ := json.MarshalIndent(j.Pipeline.Results["1R47"], "", "\t")
	ioutil.WriteFile("output.json", out, 0644)
}
