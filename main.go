package main

import (
	"encoding/gob"
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
	// Load config.yaml
	c, err := config.LoadFile("config.yaml")
	if err != nil {
		log.Fatalf("Cannot open and parse config.yaml: %v", err)
	}
	cfg = c
	http.Cfg = c
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

	// fmt.Println("Gob Example")
	// p, _ := pdb.NewPDBFromID("3CON", "P01111")
	// // fmt.Println(p)
	// // b := bytes.Buffer{}
	// // e := gob.NewEncoder(&b)
	// // e.Encode(p)
	// // fmt.Println(string(b.Bytes()))
	// b := bytes.Buffer{}
	// enc := gob.NewEncoder(&b)
	// err := enc.Encode(p)
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println(p)
	// err := writeGob("./pdb.gob", p)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// var pdbRead = new(pdb.PDB)
	// err = readGob("./student.gob", pdbRead)
	// if err != nil {
	// 	fmt.Println(err)
	// } else {
	// 	fmt.Println(pdbRead.SIFTS, "\t", pdbRead.SIFTS)
	// }

}

func writeGob(filePath string, object interface{}) error {
	file, err := os.Create(filePath)
	if err == nil {
		encoder := gob.NewEncoder(file)
		encoder.Encode(object)
	}
	file.Close()
	return err
}

func readGob(filePath string, object interface{}) error {
	file, err := os.Open(filePath)
	if err == nil {
		decoder := gob.NewDecoder(file)
		err = decoder.Decode(object)
	}
	file.Close()
	return err
}
