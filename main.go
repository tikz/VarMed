package main

import (
	"flag"
	"log"
	"os"
	"strings"
	"varq/clinvar"
	"varq/config"
	"varq/http"
	"varq/uniprot"
)

var (
	cfg *config.Config
)

func init() {
	c, err := config.LoadFile("config.yaml")
	if err != nil {
		log.Fatalf("Cannot open and parse config.yaml: %v", err)
	}

	makeDirs()

	cfg = c
	http.Cfg = c
	uniprot.DbSNP = clinvar.NewDbSNP()
}

func main() {
	pdbsFlag := arrayFlags{}
	uniprotID := flag.String("u", "", "UniProt accession.")
	flag.Var(&pdbsFlag, "p", "PDB ID(s) to analyse, can repeat this flag.")
	flag.Parse()

	if len(*uniprotID) > 0 {
		cliRun(strings.ToUpper(*uniprotID), pdbsFlag, flag.Args())
	} else {
		makeSampleResults()
		httpServe()
	}
}

func makeSampleResults() {
	_, err := os.Stat(jobDir + "15e20e5f18326d264b60eeaa07c9af8d04b0a6c70f037b7f69b6d40d22fb590b" + fileExt)
	if os.IsNotExist(err) {
		log.Println("Running pipeline to populate sample results...")
		j := NewJob(&JobRequest{
			Name:      "Sample Job - AGAL",
			UniProtID: "P06280",
			PDBIDs:    []string{"1R47"},
		})
		j.Process(false)
	}
}
