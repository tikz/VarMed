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
	msgs := make(chan string)
	go func() {
		for m := range msgs {
			fmt.Println(m)
		}
	}()

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

	j.Process()
	if j.Error != nil {
		log.Fatal(j.Error)
	}

	// if len(pdbFlags) == 0 {
	// 	p.pdbIDs = p.UniProt.PDBIDs
	// }

	close(msgs)
}
