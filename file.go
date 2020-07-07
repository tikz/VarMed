package main

import (
	"encoding/gob"
	"fmt"
	"os"
	"varq/pdb"
	"varq/uniprot"
)

const (
	dataDir     = "data/"
	unpDir      = dataDir + "uniprot/"
	pdbDir      = dataDir + "pdb/"
	jobDir      = dataDir + "jobs/"
	fpocketDir  = dataDir + "fpocket/"
	clinvarDir  = dataDir + "clinvar/"
	pfamDir     = dataDir + "pfam/"
	foldxDir    = dataDir + "foldx/"
	abswitchDir = dataDir + "abswitch/"
	fileExt     = ".varq"
)

func makeDirs() {
	os.MkdirAll(unpDir, os.ModePerm)
	os.MkdirAll(pdbDir, os.ModePerm)
	os.MkdirAll(jobDir, os.ModePerm)
	os.MkdirAll(fpocketDir, os.ModePerm)
	os.MkdirAll(clinvarDir, os.ModePerm)
	os.MkdirAll(pfamDir, os.ModePerm)
	os.MkdirAll(abswitchDir, os.ModePerm)

	os.MkdirAll(foldxDir, os.ModePerm)
	os.MkdirAll(foldxDir+"repair", os.ModePerm)
	os.MkdirAll(foldxDir+"mutations", os.ModePerm)
}

func writeJob(j *Job) error {
	return write(jobDir+j.ID+fileExt, j)
}

func loadJob(id string) (*Job, error) {
	j := Job{}
	err := read(jobDir+id+fileExt, &j)
	if err != nil {
		return nil, err
	}

	return &j, nil
}

func loadPDB(pdbID string) (*pdb.PDB, error) {
	path := pdbDir + pdbID + fileExt
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		p, err := pdb.NewPDBFromID(pdbID)
		if err != nil {
			return nil, err
		}

		err = write(path, &p)
		if err != nil {
			return nil, fmt.Errorf("write PDB: %v", err)
		}
		return &p, nil
	}

	return readPDB(pdbID)
}

func readPDB(pdbID string) (*pdb.PDB, error) {
	path := pdbDir + pdbID + fileExt
	p := new(pdb.PDB)
	err := read(path, &p)
	if err != nil {
		return nil, fmt.Errorf("load file: %v", err)
	}

	err = p.Parse()
	return p, err
}

func loadUniProt(unpID string) (*uniprot.UniProt, error) {
	path := unpDir + unpID + fileExt
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		u, err := uniprot.NewUniProt(unpID)
		if err != nil {
			return nil, err
		}

		err = write(path, &u)
		if err != nil {
			return nil, fmt.Errorf("write UniProt: %v", err)
		}
		return u, nil
	}

	u := new(uniprot.UniProt)
	err = read(path, &u)
	if err != nil {
		return nil, fmt.Errorf("load file: %v", err)
	}

	return u, nil
}

func write(filePath string, object interface{}) error {
	file, err := os.Create(filePath)
	if err == nil {
		encoder := gob.NewEncoder(file)
		encoder.Encode(object)
	}
	file.Close()

	return err
}

func read(filePath string, object interface{}) error {
	file, err := os.Open(filePath)
	if err == nil {
		decoder := gob.NewDecoder(file)
		err = decoder.Decode(object)
	}
	file.Close()

	return err
}
