package main

import (
	"flag"
	"log"
	"os"
	"varq/config"
	"varq/http"
)

var (
	cfg *config.Config
)

func init() {
	c, err := config.LoadFile("config.yaml")
	if err != nil {
		log.Fatalf("Cannot open and parse config.yaml: %v", err)
	}
	cfg = c
	http.Cfg = c

	makeDirs()
	makeSampleResults()
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

func makeSampleResults() {
	_, err := os.Stat("data/jobs/fe2423053f1a75a300e4074b1609ed972e3e2eaeae149f21d9c5fd79b4ef3d5c.varq")
	if os.IsNotExist(err) {
		j := NewJob(&JobRequest{UniProtID: "P00390", PDBIDs: []string{"2GH5"}})
		j.Process()
	}
}
