package main

import (
	"flag"
	"log"
	"os"
	"strings"
	"varmed/config"

	"github.com/tikz/bio/abswitch"
	"github.com/tikz/bio/clinvar"
	"github.com/tikz/bio/conservation"
	"github.com/tikz/bio/foldx"
	"github.com/tikz/bio/tango"
)

var (
	cfg       *config.Config
	instances *Instances
)

type Instances struct {
	FoldX    *foldx.FoldX
	Pfam     *conservation.Pfam
	ClinVar  *clinvar.ClinVar
	AbSwitch *abswitch.AbSwitch
	Tango    *tango.Tango
}

func init() {
	c, err := config.LoadFile("config.yaml")
	if err != nil {
		log.Fatalf("Cannot open and parse config.yaml: %v", err)
	}

	cfg = c

	makeDirs()

	instances = &Instances{}
	if instances.ClinVar, err = clinvar.NewClinVar(cfg.Paths.ClinVar); err != nil {
		log.Fatalf("Cannot instance ClinVar dir: %v", err)
	}

	if instances.Pfam, err = conservation.NewPfam(cfg.Paths.Pfam); err != nil {
		log.Fatalf("Cannot instance Pfam dir: %v", err)
	}

	if instances.FoldX, err = foldx.NewFoldX(cfg.Paths.FoldXBin,
		cfg.Paths.FoldXRepair,
		cfg.Paths.FoldXMutations); err != nil {
		log.Fatalf("Cannot instance FoldX: %v", err)
	}

	if instances.AbSwitch, err = abswitch.NewAbSwitch(cfg.Paths.AbSwitchBin, cfg.Paths.AbSwitch); err != nil {
		log.Fatalf("Cannot instance abSwitch: %v", err)
	}

	if instances.Tango, err = tango.NewTango(cfg.Paths.TangoBin, cfg.Paths.Tango); err != nil {
		log.Fatalf("Cannot instance Tango: %v", err)
	}
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
	_, err := os.Stat(cfg.Paths.Jobs + "15e20e5f18326d264b60eeaa07c9af8d04b0a6c70f037b7f69b6d40d22fb590b" + cfg.Paths.FileExt)
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
