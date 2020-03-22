package main

import (
	"fmt"
	"log"
	"time"
	"varq/binding"
	"varq/interaction"
	"varq/pdb"
	"varq/uniprot"
)

// Analysis contains all pipeline steps results for a single PDB entry.
type Analysis struct {
	PDB         *pdb.PDB
	Binding     *binding.BindingAnalysis
	Interaction *interaction.InteractionAnalysis
	Error       error `json:"-"`
}

// pipelineStructureWorker fetches a single PDB crystal file.
func pipelineStructureWorker(crystalChan <-chan *pdb.PDB, analysisChan chan<- *Analysis) {
	for crystal := range crystalChan {
		err := crystal.Fetch()
		if err != nil {
			analysisChan <- &Analysis{PDB: crystal, Error: fmt.Errorf("PDB %s: %v", crystal.ID, err)}
			continue
		}

		analysisChan <- analyseStructure(&Analysis{PDB: crystal})
	}
}

// analyseStructure runs each available analysis in parallel for a single structure.
func analyseStructure(analysis *Analysis) *Analysis {
	bindingChan := make(chan *binding.BindingAnalysis)
	interactionChan := make(chan *interaction.InteractionAnalysis)

	go binding.RunBindingAnalysis(analysis.PDB, bindingChan)
	go interaction.RunInteractionAnalysis(analysis.PDB, interactionChan)

	bindingRes := <-bindingChan
	if bindingRes.Error != nil {
		analysis.Error = fmt.Errorf("binding analysis: %v", bindingRes.Error)
		return analysis
	}
	analysis.Binding = bindingRes
	log.Printf("PDB %s binding analysis done in %d ms", analysis.PDB.ID, bindingRes.Duration.Milliseconds())

	interactionRes := <-interactionChan
	if interactionRes.Error != nil {
		analysis.Error = fmt.Errorf("interaction analysis: %v", interactionRes.Error)
		return analysis
	}
	analysis.Interaction = interactionRes
	log.Printf("PDB %s interaction analysis done in %d ms", analysis.PDB.ID, interactionRes.Duration.Milliseconds())

	return analysis
}

func RunPipeline(crystals []*pdb.PDB) (analyses []*Analysis, err error) {
	length := len(crystals)
	crystalChan := make(chan *pdb.PDB, length)
	analysisChan := make(chan *Analysis, length)

	// Launch crystal workers
	for w := 1; w <= 20; w++ {
		go pipelineStructureWorker(crystalChan, analysisChan)
	}

	for _, crystal := range crystals {
		crystalChan <- crystal
	}
	close(crystalChan)

	for a := 1; a <= length; a++ {
		analysis := <-analysisChan
		if analysis.Error != nil {
			return nil, analysis.Error
		}
		analyses = append(analyses, analysis)
	}

	return analyses, nil
}

// RunPipelineForUniProt grabs and analyses all structures from a given UniProt ID.
func RunPipelineForUniProt(uniprotID string) ([]*Analysis, error) {
	start := time.Now()

	u, err := uniprot.NewUniProt(uniprotID)
	if err != nil {
		return nil, fmt.Errorf("run pipeline: %v", err)
	}

	analyses, err := RunPipeline(u.Crystals)
	if err != nil {
		return nil, fmt.Errorf("analyzing crystals: %v", err)
	}

	end := time.Since(start)
	log.Printf("Finished UniProt %s in %f secs", u.ID, end.Seconds())
	return analyses, nil
}

// RunPipelineForPDBs grabs and analyses structures from a slice of given PDB IDs.
func RunPipelineForPDBs(PDBIDs []string) ([]*Analysis, error) {
	start := time.Now()

	var crystals []*pdb.PDB
	for _, ID := range PDBIDs {
		crystals = append(crystals, &pdb.PDB{ID: ID})
	}

	analyses, err := RunPipeline(crystals)
	if err != nil {
		return nil, fmt.Errorf("analyzing crystals: %v", err)
	}

	end := time.Since(start)
	log.Printf("Finished PDBs %s in %f secs", PDBIDs, end.Seconds())
	return analyses, nil
}
