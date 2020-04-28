package main

import (
	"flag"
	"log"
	"os"
	"strings"
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
		cliRun(strings.ToUpper(*uniprotID), pdbsFlag)
	} else {
		httpServe()
	}
}

func makeSampleResults() {
	_, err := os.Stat(jobDir + "aa2725a483568c283274c6e551b83ac1c34548736c0dbb2581ba770bb0de21eb" + fileExt)
	if os.IsNotExist(err) {
		log.Println("Running pipeline to populate sample results...")
		j := NewJob(&JobRequest{UniProtID: "P01112", PDBIDs: []string{"1LFD"}})
		j.Process()
	}
}
