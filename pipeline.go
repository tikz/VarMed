package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"
	"varq/binding"
	"varq/exposure"
	"varq/interaction"
	"varq/pdb"
	"varq/uniprot"
)

// Analysis contains all pipeline steps results for a single PDB entry.
type Analysis struct {
	PDB         *pdb.PDB
	Binding     *binding.Step
	Interaction *interaction.Step
	Exposure    *exposure.Step
	Error       error `json:"-"`
}

// pipelinePDBWorker fetches and loads a single PDB file.
func pipelinePDBWorker(pdbChan <-chan *pdb.PDB, aChan chan<- *Analysis) {
	for pdb := range pdbChan {
		analysis := Analysis{PDB: pdb}

		start := time.Now()
		log.Printf("Loading PDB %s...", pdb.ID)
		err := pdb.Load()
		if err != nil {
			analysis.Error = fmt.Errorf("load PDB %s: %v", pdb.ID, err)
			aChan <- &analysis
			continue
		}
		end := time.Since(start)
		log.Printf("PDB %s loaded in %.3f secs", pdb.ID, end.Seconds())

		aChan <- analysePDB(&analysis)
	}
}

// analysePDB runs each available analysis in parallel for a single structure.
func analysePDB(a *Analysis) *Analysis {
	// Create temp PDB on filesystem for analysis with external tools
	a.PDB.LocalFilename = "varq_" + a.PDB.ID
	a.PDB.LocalPath = "/tmp/" + a.PDB.LocalFilename + ".pdb"

	err := ioutil.WriteFile(a.PDB.LocalPath, a.PDB.RawPDB, 0644)
	if err != nil {
		a.Error = fmt.Errorf("create tmp PDB: %v", err)
		return a
	}

	defer func() {
		os.Remove(a.PDB.LocalPath)
	}()

	bindingChan := make(chan *binding.Step)
	interactionChan := make(chan *interaction.Step)
	exposureChan := make(chan *exposure.Step)

	go binding.RunBindingStep(a.PDB, bindingChan)
	go interaction.RunInteractionStep(a.PDB, interactionChan)
	go exposure.RunExposureStep(a.PDB, exposureChan)

	// TODO: refactor these repeated patterns
	bindingRes := <-bindingChan
	if bindingRes.Error != nil {
		a.Error = fmt.Errorf("binding analysis: %v", bindingRes.Error)
		return a
	}
	a.Binding = bindingRes
	log.Printf("PDB %s binding analysis done in %.3f secs", a.PDB.ID, bindingRes.Duration.Seconds())

	interactionRes := <-interactionChan
	if interactionRes.Error != nil {
		a.Error = fmt.Errorf("interaction analysis: %v", interactionRes.Error)
		return a
	}
	a.Interaction = interactionRes
	log.Printf("PDB %s interaction analysis done in %.3f secs", a.PDB.ID, interactionRes.Duration.Seconds())

	exposureRes := <-exposureChan
	if exposureRes.Error != nil {
		a.Error = fmt.Errorf("exposure analysis: %v", exposureRes.Error)
		return a
	}
	a.Exposure = exposureRes
	log.Printf("PDB %s exposure analysis done in %.3f secs", a.PDB.ID, exposureRes.Duration.Seconds())

	return a
}

func runPipelinePDBs(pdbs []*pdb.PDB) (analyses []*Analysis, err error) {
	length := len(pdbs)
	pdbChan := make(chan *pdb.PDB, length)
	analysisChan := make(chan *Analysis, length)

	for w := 1; w <= 20; w++ {
		go pipelinePDBWorker(pdbChan, analysisChan)
	}

	for _, pdb := range pdbs {
		pdbChan <- pdb
	}
	close(pdbChan)

	for a := 1; a <= length; a++ {
		analysis := <-analysisChan
		if analysis.Error != nil {
			log.Printf("ignoring crystal: %v", analysis.Error)
		} else {
			analyses = append(analyses, analysis)
			debugPrintChains(analysis)
		}
	}

	return analyses, nil
}

// RunPipeline grabs and analyses all structures from a given UniProt ID.
func RunPipeline(uniprotID string, filterPDBIDs []string) ([]*Analysis, error) {
	start := time.Now()

	u, err := uniprot.NewUniProt(uniprotID)
	if err != nil {
		return nil, fmt.Errorf("run pipeline: %v", err)
	}

	if len(filterPDBIDs) > 0 {
		u.FilterPDBs(filterPDBIDs)
	}

	analyses, err := runPipelinePDBs(u.PDBs)
	if err != nil {
		return nil, fmt.Errorf("analyzing crystals: %v", err)
	}

	end := time.Since(start)
	log.Printf("Finished UniProt %s in %.3f secs", u.ID, end.Seconds())
	return analyses, nil
}
