package main

import (
	"fmt"
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

	fmt.Println("VarQ CLI")
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
}
