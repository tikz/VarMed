package main

import (
	"fmt"
	"os"
	"time"
	"varq/binding"
	"varq/energy"
	"varq/exposure"
	"varq/interaction"
	"varq/pdb"
	"varq/uniprot"
)

// Results represents a group of results from each one of the available analysis steps.
type Results struct {
	UniProt     *uniprot.UniProt     `json:"uniprot"`
	PDB         *pdb.PDB             `json:"pdb"`
	Binding     *binding.Results     `json:"binding"`
	Interaction *interaction.Results `json:"interaction"`
	Exposure    *exposure.Results    `json:"exposure"`
	Energy      *energy.Results      `json:"energy"`
	Error       error                `json:"-"`
}

// Pipeline represents a single run of the VarQ pipeline.
type Pipeline struct {
	UniProt  *uniprot.UniProt
	Variants map[int]string
	Results  map[string]*Results // PDB ID to results
	Duration time.Duration

	msgChan chan string // readable text messages about the status
	pdbIDs  []string
}

// msg prints and sends a message with added format to the channel.
func (p *Pipeline) msg(m string) {
	p.msgChan <- time.Now().Format("15:04:05-0700") + " " + m
}

// NewPipeline constructs a new Pipeline.
func NewPipeline(unpID string, pdbIDs []string, variants map[int]string, msgChan chan string) (*Pipeline, error) {
	uniprot, err := loadUniProt(unpID)
	if err != nil {
		return nil, err
	}

	p := Pipeline{
		UniProt:  uniprot,
		Variants: variants,
		msgChan:  msgChan,
		pdbIDs:   pdbIDs,
		Results:  make(map[string]*Results),
	}

	return &p, nil
}

// Run starts the process of analyzing given PDB IDs corresponding to an UniProt ID.
func (p *Pipeline) Run() error {
	start := time.Now()
	p.msg("Job started")

	pdbIDChan := make(chan string, len(p.pdbIDs))
	resChan := make(chan *Results, len(p.pdbIDs))

	for w := 1; w <= cfg.VarQ.Pipeline.StructureWorkers; w++ {
		go p.pdbWorker(pdbIDChan, resChan)
	}

	for _, id := range p.pdbIDs {
		if !p.UniProt.PDBIDExists(id) {
			return fmt.Errorf("PDB ID %s not found", id)
		}
		pdbIDChan <- id
	}

	for a := 1; a <= len(p.pdbIDs); a++ {
		result := <-resChan
		if result.Error != nil {
			return fmt.Errorf("step error: %v", result.Error)
		}

		p.Results[result.PDB.ID] = result
	}

	p.Duration = time.Since(start)
	p.msg(fmt.Sprintf("Pipeline finished in %.3f secs", p.Duration.Seconds()))
	return nil
}

// pdbWorker fetches and loads a single PDB file.
func (p *Pipeline) pdbWorker(pdbIDChan <-chan string, resChan chan<- *Results) {
	for pdbID := range pdbIDChan {
		results := Results{}

		start := time.Now()
		p.msg(fmt.Sprintf("Loading PDB %s...", pdbID))
		pdb, err := loadPDB(pdbID)
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
func (p *Pipeline) analysePDB(r *Results) *Results {
	// Create temp PDB on filesystem for analysis with external tools
	path := "bin/" + r.PDB.ID + ".pdb"
	r.PDB.WriteFile(path)

	defer func() {
		os.Remove(path)
	}()

	bindingChan := make(chan *binding.Results)
	interactionChan := make(chan *interaction.Results)
	exposureChan := make(chan *exposure.Results)
	energyChan := make(chan *energy.Results)

	idStr := fmt.Sprintf("PDB %s ", r.PDB.ID)
	msgPDB := func(msg string) {
		p.msg(idStr + msg)
	}

	if cfg.VarQ.Pipeline.EnableSteps.Binding {
		go binding.Run(r.UniProt, r.PDB, bindingChan, msgPDB)
		msgPDB("started binding analysis")
	}
	if cfg.VarQ.Pipeline.EnableSteps.Interaction {
		go interaction.Run(r.PDB, interactionChan, msgPDB)
		msgPDB("started interaction analysis")
	}
	if cfg.VarQ.Pipeline.EnableSteps.Exposure {
		go exposure.Run(r.PDB, exposureChan, msgPDB)
		msgPDB("started exposure analysis")
	}
	if cfg.VarQ.Pipeline.EnableSteps.Energy {
		go energy.Run(p.Variants, r.UniProt, r.PDB, foldxDir, energyChan, msgPDB)
		msgPDB("started energy analysis")
	}

	// TODO: refactor these repeated patterns
	if cfg.VarQ.Pipeline.EnableSteps.Binding {
		bindingRes := <-bindingChan
		if bindingRes.Error != nil {
			r.Error = fmt.Errorf("binding analysis: %v", bindingRes.Error)
			return r
		}
		r.Binding = bindingRes
		msgPDB(fmt.Sprintf("binding analysis done in %.3f secs", bindingRes.Duration.Seconds()))
	}

	if cfg.VarQ.Pipeline.EnableSteps.Interaction {
		interactionRes := <-interactionChan
		if interactionRes.Error != nil {
			r.Error = fmt.Errorf("interaction analysis: %v", interactionRes.Error)
			return r
		}
		r.Interaction = interactionRes
		msgPDB(fmt.Sprintf("interaction analysis done in %.3f secs", interactionRes.Duration.Seconds()))
	}

	if cfg.VarQ.Pipeline.EnableSteps.Exposure {
		exposureRes := <-exposureChan
		if exposureRes.Error != nil {
			r.Error = fmt.Errorf("exposure analysis: %v", exposureRes.Error)
			return r
		}
		r.Exposure = exposureRes
		msgPDB(fmt.Sprintf("exposure analysis done in %.3f secs", exposureRes.Duration.Seconds()))
	}

	if cfg.VarQ.Pipeline.EnableSteps.Energy {
		energyRes := <-energyChan
		if energyRes.Error != nil {
			r.Error = fmt.Errorf("energy analysis: %v", energyRes.Error)
			return r
		}
		r.Energy = energyRes
		msgPDB(fmt.Sprintf("energy analysis done in %.3f secs", energyRes.Duration.Seconds()))
	}

	if cfg.DebugPrint.Enabled {
		printResults(r)
	}

	return r
}
