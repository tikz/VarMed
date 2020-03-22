package main

import (
	"flag"
	"log"
	"varq/config"
)

var (
	cfg *config.Config
)

func init() {
	// Load config.yaml
	c, err := config.LoadFile("config.yaml")
	if err != nil {
		log.Fatalf("Cannot open and parse config.yaml: %v", err)
	}
	cfg = c
}

func main() {
	uniprotsFlag := arrayFlags{}
	pdbsFlag := arrayFlags{}
	flag.Var(&uniprotsFlag, "u", "UniProt ID to analyse. Can pass multiple flags.")
	flag.Var(&pdbsFlag, "p", "PDB ID to analyse. Can pass multiple flags.")
	flag.Parse()

	if len(uniprotsFlag) > 0 || len(pdbsFlag) > 0 {
		cliRun(uniprotsFlag, pdbsFlag)
	} else {
		httpServe()
	}
}
