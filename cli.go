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

func cliRun(uniprotID string, pdbFlags arrayFlags) {
	msgs := make(chan string)
	go func() {
		for m := range msgs {
			fmt.Println(m)
		}
	}()

	variants := make(map[int]string)
	p, err := NewPipeline(uniprotID, pdbFlags, variants, msgs)
	if err != nil {
		log.Fatal(err)
	}

	if len(pdbFlags) == 0 {
		p.pdbIDs = p.UniProt.PDBIDs
	}

	err = p.Run()
	if err != nil {
		log.Fatal(err)
	}
	close(msgs)
}
