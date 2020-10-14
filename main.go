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
	_, err := os.Stat(cfg.Paths.Jobs + "f44f3c3880d25872ce61cfa3612e2d683fec7c8296924337a24a28da26ee4d8d" + cfg.Paths.FileExt)
	if os.IsNotExist(err) {
		log.Println("Running pipeline to populate sample results...")
		j := NewJob(&JobRequest{
			Name:      "Sample Job - AGAL",
			UniProtID: "P06280",
			PDBIDs:    []string{"1R47"},
			Variants:  []string{"A121T", "A135V", "A143P", "A143T", "A156T", "A156V", "A20D", "A20P", "A230T", "A285P", "A288D", "A309V", "A31V", "A352G", "A377D", "A97V", "C142R", "C142Y", "C172R", "C172Y", "C202W", "C202Y", "C223G", "C378Y", "C52R", "C52S", "C56F", "C56G", "C56Y", "C94S", "C94Y", "D165V", "D170V", "D231N", "D234E", "D244H", "D244N", "D264V", "D264Y", "D266H", "D266N", "D266V", "D313N", "D313Y", "D315N", "D33G", "D92H", "D92Y", "D93G", "D93N", "E338K", "E341K", "E358A", "E358K", "E48D", "E59K", "E66Q", "E71G", "F113I", "F113L", "F113S", "F396Y", "G128E", "G138R", "G144V", "G163V", "G171D", "G183D", "G258R", "G260A", "G261D", "G328A", "G328R", "G328V", "G35E", "G35R", "G360C", "G360S", "G361R", "G373D", "G373S", "G375A", "G43R", "G80D", "G85D", "H46P", "H46R", "H46Y", "I154T", "I198T", "I219M", "I219N", "I219T", "I242N", "I242V", "I253T", "I289F", "I289V", "I317S", "I64F", "I91N", "I91T", "K213R", "L120V", "L131P", "L166V", "L167Q", "L180F", "L21P", "L243F", "L300F", "L32P", "L36W", "L3P", "L3V", "L414S", "L45P", "L89P", "L89R", "M187I", "M187V", "M267I", "M284T", "M296I", "M296V", "M42L", "M42T", "M42V", "M72V", "N215S", "N224D", "N224S", "N228S", "N249K", "N263S", "N272K", "N272S", "N298H", "N298K", "N298S", "N320K", "N320Y", "N34S", "P146S", "P205T", "P214L", "P259L", "P259R", "P265R", "P323R", "P409A", "P409T", "P40L", "P40S", "P60L", "Q279E", "Q279H", "Q280H", "Q321E", "Q327K", "Q327L", "Q327R", "Q330R", "R100K", "R100T", "R112C", "R112H", "R112S", "R196S", "R227P", "R227Q", "R301Q", "R342P", "R342Q", "R356P", "R356Q", "R356W", "R363H", "R392S", "R49L", "R49P", "R49S", "S148N", "S148R", "S201F", "S235C", "S247P", "S276G", "S297F", "S65T", "T410A", "V164G", "V164L", "V254A", "V269A", "V269G", "V316A", "V316E", "W162C", "W162R", "W204R", "W226R", "W236C", "W236L", "W262R", "W287C", "W287G", "W340R", "W399S", "W47G", "W47R", "W95S", "Y134S", "Y216D", "Y86C", "Y86H"},
		})
		j.Process(false)
	}
}
