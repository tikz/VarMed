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
	os.MkdirAll(cfg.Paths.AbSwitch, os.ModePerm)
	os.MkdirAll(cfg.Paths.Tango, os.ModePerm)

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
	if err != nil {
		return err
	}

	decoder := gob.NewDecoder(file)
	err = decoder.Decode(object)
	file.Close()

	return err
}
