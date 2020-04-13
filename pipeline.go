package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
	"varq/binding"
	"varq/exposure"
	"varq/interaction"
	"varq/pdb"
	"varq/uniprot"
)

// Step is an interface representing a group of analyses
// related together, to run on a PDB structure.
type Step interface {
	Run() interface{} // resulting data type varies per case. TODO: see if can be unified
}

// Results contains all steps results for a single PDB entry.
type Results struct {
	PDB         *pdb.PDB
	Binding     *binding.Results
	Interaction *interaction.Results
	Exposure    *exposure.Results
	Error       error `json:"-"`
}

type Pipeline struct {
	UniProt *uniprot.UniProt
	PDBIDs  []string

	PDBs     map[string]*Results
	Duration time.Duration

	msg func(string) // private callback for passing messages
}

func NewPipeline(unpID string, pdbIDs []string, msgChan chan<- string) (*Pipeline, error) {
	uniprot, err := uniprot.NewUniProt(unpID)
	if err != nil {
		return nil, err
	}

	msgHook := func(m string) {
		log.Println(m)
		msgChan <- m
	}

	p := Pipeline{
		UniProt: uniprot,
		PDBIDs:  pdbIDs,
		msg:     msgHook,
		PDBs:    make(map[string]*Results),
	}

	return &p, nil
}

// RunPipeline grabs and analyses all structures from a given UniProt ID.
func (p *Pipeline) RunPipeline() error {
	start := time.Now()

	pdbChan := make(chan *pdb.PDB)
	resultsChan := make(chan *Results)

	for w := 1; w <= cfg.VarQ.Pipeline.StructureWorkers; w++ {
		go p.pdbWorker(pdbChan, resultsChan)
	}

	for _, id := range p.PDBIDs {
		idu := strings.ToUpper(id)
		if _, ok := p.UniProt.PDBs[idu]; !ok {
			return fmt.Errorf("PDB ID %s not found", idu)
		}
		pdbChan <- p.UniProt.PDBs[idu]
	}

	close(pdbChan)

	for a := 1; a <= len(p.PDBIDs); a++ {
		result := <-resultsChan
		if result.Error != nil {
			return fmt.Errorf("step error: %v", result.Error)
		}

		p.PDBs[result.PDB.ID] = result
	}

	p.Duration = time.Since(start)
	p.msg(fmt.Sprintf("Finished UniProt %s in %.3f secs", p.UniProt.ID, p.Duration.Seconds()))
	return nil
}

// pdbWorker fetches and loads a single PDB file.
func (p *Pipeline) pdbWorker(pdbChan <-chan *pdb.PDB, aChan chan<- *Results) {
	for pdb := range pdbChan {
		analysis := Results{PDB: pdb}

		start := time.Now()
		p.msg(fmt.Sprintf("Loading PDB %s...", pdb.ID))
		err := pdb.Load()
		if err != nil {
			analysis.Error = fmt.Errorf("load PDB %s: %v", pdb.ID, err)
			aChan <- &analysis
			continue
		}
		end := time.Since(start)
		p.msg(fmt.Sprintf("PDB %s loaded in %.3f secs", pdb.ID, end.Seconds()))

		aChan <- p.analysePDB(&analysis)
	}
}

// analysePDB runs each available analysis in parallel for a single structure.
func (p *Pipeline) analysePDB(a *Results) *Results {
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

	bindingChan := make(chan *binding.Results)
	interactionChan := make(chan *interaction.Results)
	exposureChan := make(chan *exposure.Results)

	if cfg.VarQ.Pipeline.EnableSteps.Binding {
		go binding.Run(a.PDB, bindingChan)
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
