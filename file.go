package main

import (
	"encoding/gob"
	"fmt"
	"os"
	"varq/pdb"
	"varq/uniprot"
)

func makeDirs() {
	os.MkdirAll("data/uniprot", os.ModePerm)
	os.MkdirAll("data/pdb", os.ModePerm)
	os.MkdirAll("data/jobs", os.ModePerm)
}

func WriteJob(j *Job) error {
	return write("data/jobs/"+j.ID+".varq", j)
}

func LoadJob(id string) (*Job, error) {
	j := Job{}
	err := read("data/jobs/"+id+".varq", &j)
	if err != nil {
		return nil, err
	}

	return &j, nil
}

func LoadPDB(pdbID string) (*pdb.PDB, error) {
	path := "data/pdb/" + pdbID + ".varq"
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

	return ReadPDB(pdbID)
}

func ReadPDB(pdbID string) (*pdb.PDB, error) {
	path := "data/pdb/" + pdbID + ".varq"
	p := new(pdb.PDB)
	err := read(path, &p)
	if err != nil {
		return nil, fmt.Errorf("load file: %v", err)
	}

	err = p.Parse()
	return p, err
}

func LoadUniProt(unpID string) (*uniprot.UniProt, error) {
	path := "data/uniprot/" + unpID + ".varq"
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
