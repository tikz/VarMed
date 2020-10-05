package main

import (
	"encoding/gob"
	"os"
)

func makeDirs() {
	os.MkdirAll(cfg.Paths.UniProt, os.ModePerm)
	os.MkdirAll(cfg.Paths.PDB, os.ModePerm)
	os.MkdirAll(cfg.Paths.Jobs, os.ModePerm)
	os.MkdirAll(cfg.Paths.Fpocket, os.ModePerm)
	os.MkdirAll(cfg.Paths.ClinVar, os.ModePerm)
	os.MkdirAll(cfg.Paths.Pfam, os.ModePerm)
	os.MkdirAll(cfg.Paths.Abswitch, os.ModePerm)

	os.MkdirAll(cfg.Paths.FoldXRepair, os.ModePerm)
	os.MkdirAll(cfg.Paths.FoldXMutations, os.ModePerm)
}

func writeJob(j *Job) error {
	return write(cfg.Paths.Jobs+j.ID+cfg.Paths.FileExt, j)
}

func loadJob(id string) (*Job, error) {
	j := Job{}
	err := read(cfg.Paths.Jobs+id+cfg.Paths.FileExt, &j)
	if err != nil {
		return nil, err
	}

	return &j, nil
}

// func loadPDB(pdbID string) (*pdb.PDB, error) {
// 	path := pdbDir + pdbID + fileExt
// 	_, err := os.Stat(path)
// 	if os.IsNotExist(err) {
// 		p, err := pdb.NewPDBFromID(pdbID)
// 		if err != nil {
// 			return nil, err
// 		}

// 		err = write(path, &p)
// 		if err != nil {
// 			return nil, fmt.Errorf("write PDB: %v", err)
// 		}
// 		return &p, nil
// 	}

// 	return readPDB(pdbID)
// }

// func readPDB(pdbID string) (*pdb.PDB, error) {
// 	path := pdbDir + pdbID + fileExt
// 	p := new(pdb.PDB)
// 	err := read(path, &p)
// 	if err != nil {
// 		return nil, fmt.Errorf("load file: %v", err)
// 	}

// 	err = p.Parse()
// 	return p, err
// }

// func loadUniProt(unpID string) (*uniprot.UniProt, error) {
// 	path := unpDir + unpID + fileExt
// 	_, err := os.Stat(path)
// 	if os.IsNotExist(err) {
// 		u, err := uniprot.NewUniProt(unpID)
// 		if err != nil {
// 			return nil, err
// 		}

// 		err = write(path, &u)
// 		if err != nil {
// 			return nil, fmt.Errorf("write UniProt: %v", err)
// 		}
// 		return u, nil
// 	}

// 	u := new(uniprot.UniProt)
// 	err = read(path, &u)
// 	if err != nil {
// 		return nil, fmt.Errorf("load file: %v", err)
// 	}

// 	return u, nil
// }

func write(filePath string, object interface{}) error {
	file, err := os.Create(filePath)
	if err == nil {
		encoder := gob.NewEncoder(file)
		err = encoder.Encode(object)
		if err != nil {
			return err
		}
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
