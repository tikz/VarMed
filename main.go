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
		// makeSampleResults()
		httpServe()
	}
}

func makeSampleResults() {
	_, err := os.Stat(jobDir + "236ff01847ea475576e3c7d972c489b673c30d8990ff52248d970fbcc467b605" + fileExt)
	if os.IsNotExist(err) {
		log.Println("Running pipeline to populate sample results...")
		j := NewJob(&JobRequest{
			Name:      "Sample Job - GTPase HRas",
			UniProtID: "P01112",
			PDBIDs:    []string{"6D5H"},
		})
		j.Process(false)
	}
}
