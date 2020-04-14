package main

import (
	"flag"
	"log"
	"varq/config"
	"varq/http"
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
	http.Cfg = c

	MakeDirs()
}

func main() {
	pdbsFlag := arrayFlags{}
	uniprotID := flag.String("u", "", "UniProt accession.")
	flag.Var(&pdbsFlag, "p", "PDB ID(s) to analyse, can repeat this flag.")
	flag.Parse()

	if len(*uniprotID) > 0 {
		cliRun(*uniprotID, pdbsFlag)
	} else {
		httpServe()
	}
}
