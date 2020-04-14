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

type Results struct {
	UniProt     *uniprot.UniProt
	PDB         *pdb.PDB
	Binding     *binding.Results
	Interaction *interaction.Results
	Exposure    *exposure.Results
	Error       error `json:"-"`
}

type Pipeline struct {
	UniProt  *uniprot.UniProt
	Results  map[string]*Results // PDB ID to results
	Duration time.Duration

	pdbIDs []string
	msg    func(string) // private callback for passing messages
}

func NewPipeline(unpID string, pdbIDs []string, msgChan chan<- string) (*Pipeline, error) {
	uniprot, err := LoadUniProt(unpID)
	if err != nil {
		return nil, err
	}

	msgHook := func(m string) {
		log.Println(m)
		// msgChan <- m
	} // TODO: remove this crap, let the caller manage the channel

	p := Pipeline{
		UniProt: uniprot,
		pdbIDs:  pdbIDs,
		msg:     msgHook,
		Results: make(map[string]*Results),
	}

	return &p, nil
}

// RunPipeline grabs and analyses all structures from a given UniProt ID.
func (p *Pipeline) RunPipeline() error {
	start := time.Now()

	pdbIDChan := make(chan string)
	resultsChan := make(chan *Results)

	for w := 1; w <= cfg.VarQ.Pipeline.StructureWorkers; w++ {
		go p.pdbWorker(pdbIDChan, resultsChan)
	}

	for _, id := range p.pdbIDs {
		if !p.UniProt.PDBIDExists(id) {
			return fmt.Errorf("PDB ID %s not found", id)
		}
		pdbIDChan <- id
	}

	close(pdbIDChan)

	for a := 1; a <= len(p.pdbIDs); a++ {
		result := <-resultsChan
		if result.Error != nil {
			return fmt.Errorf("step error: %v", result.Error)
		}

		p.Results[result.PDB.ID] = result
	}

	p.Duration = time.Since(start)
	p.msg(fmt.Sprintf("Finished UniProt %s in %.3f secs", p.UniProt.ID, p.Duration.Seconds()))
	return nil
}

// pdbWorker fetches and loads a single PDB file.
func (p *Pipeline) pdbWorker(pdbIDChan <-chan string, resChan chan<- *Results) {
	for pdbID := range pdbIDChan {
		results := Results{}

		start := time.Now()
		p.msg(fmt.Sprintf("Loading PDB %s...", pdbID))
		pdb, err := LoadPDB(pdbID)
		if err != nil {
			results.Error = fmt.Errorf("load PDB %s: %v", pdbID, err)
			resChan <- &results
			continue
		}
		results.PDB = pdb
		results.UniProt = p.UniProt
		end := time.Since(start)
		p.msg(fmt.Sprintf("PDB %s loaded in %.3f secs", pdbID, end.Seconds()))

		resChan <- p.analysePDB(&results)
	}
}

// analysePDB runs each available analysis in parallel for a single structure.
func (p *Pipeline) analysePDB(a *Results) *Results {
	// Create temp PDB on filesystem for analysis with external tools
	a.PDB.LocalFilename = "varq_" + a.PDB.ID
	a.PDB.LocalPath = "/tmp/" + a.PDB.LocalFilename + ".pdb"
	// TODO: code smell?

	err := ioutil.WriteFile(a.PDB.LocalPath, a.PDB.RawPDB, 0644)
	if err != nil {
		a.Error = fmt.Errorf("create tmp PDB: %v", err)
		return a
	}

	defer func() {
		os.Remove(a.PDB.LocalPath)
	}()

	bindingChan := make(chan *binding.Results)
	interactionChan := make(chan *interaction.Results)
	exposureChan := make(chan *exposure.Results)

	if cfg.VarQ.Pipeline.EnableSteps.Binding {
		go binding.Run(a.UniProt, a.PDB, bindingChan)
	}
	if cfg.VarQ.Pipeline.EnableSteps.Interaction {
		go interaction.Run(a.PDB, interactionChan)
	}
	if cfg.VarQ.Pipeline.EnableSteps.Exposure {
		go exposure.Run(a.PDB, exposureChan)
	}

	// TODO: refactor these repeated patterns
	if cfg.VarQ.Pipeline.EnableSteps.Binding {
		bindingRes := <-bindingChan
		if bindingRes.Error != nil {
			a.Error = fmt.Errorf("binding analysis: %v", bindingRes.Error)
			return a
		}
		a.Binding = bindingRes
		p.msg(fmt.Sprintf("PDB %s binding analysis done in %.3f secs", a.PDB.ID, bindingRes.Duration.Seconds()))
	}

	if cfg.VarQ.Pipeline.EnableSteps.Interaction {
		interactionRes := <-interactionChan
		if interactionRes.Error != nil {
			a.Error = fmt.Errorf("interaction analysis: %v", interactionRes.Error)
			return a
		}
		a.Interaction = interactionRes
		p.msg(fmt.Sprintf("PDB %s interaction analysis done in %.3f secs", a.PDB.ID, interactionRes.Duration.Seconds()))
	}

	if cfg.VarQ.Pipeline.EnableSteps.Exposure {
		exposureRes := <-exposureChan
		if exposureRes.Error != nil {
			a.Error = fmt.Errorf("exposure analysis: %v", exposureRes.Error)
			return a
		}
		a.Exposure = exposureRes
		p.msg(fmt.Sprintf("PDB %s exposure analysis done in %.3f secs", a.PDB.ID, exposureRes.Duration.Seconds()))
	}

	if cfg.DebugPrint.Enabled {
		debugPrintChains(a)
	}

	return a
}
