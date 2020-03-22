package main

import (
	"fmt"
	"log"
	"time"
	"varq/binding"
	"varq/interaction"
	"varq/pdb"
)

// PDBAnalysis contains all pipeline steps results for an associated PDB entry
type PDBAnalysis struct {
	PDB         *pdb.PDB
	Binding     *binding.BindingAnalysis
	Interaction *interaction.InteractionAnalysis
	Error       error `json:"-"`
}

// pipelineCrystalWorker fetches a single PDB crystal file, then fires more goroutines to do each analysis in parallel
func pipelineCrystalWorker(crystalChan <-chan *pdb.PDB, analysisChan chan<- *PDBAnalysis) {
	for crystal := range crystalChan {
		err := crystal.Fetch()
		if err != nil {
			analysisChan <- &PDBAnalysis{PDB: crystal, Error: fmt.Errorf("PDB %s: %v", crystal.ID, err)}
			continue
		}

		analysisChan <- analyseCrystal(&PDBAnalysis{PDB: crystal})
	}
}

func analyseCrystal(analysis *PDBAnalysis) *PDBAnalysis {
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
	log.Printf("PDB %s interaction analysis done in %d ns", analysis.PDB.ID, interactionRes.Duration.Nanoseconds())

	return analysis
}

// RunPipeline manages the workers for parallel fetching and processing of protein data
func RunPipeline(uniprotID string, pdbIDsFilter []string) (*Protein, error) {
	start := time.Now()

	p, err := NewProtein(uniprotID)
	if err != nil {
		return nil, fmt.Errorf("run pipeline: %v", err)
	}

	length := len(p.Crystals)
	if len(pdbIDsFilter) != 0 {
		length = len(pdbIDsFilter)
	}

	crystalChan := make(chan *pdb.PDB, length)
	analysisChan := make(chan *PDBAnalysis, length)

	// Fetch all crystals in parallel
	for w := 1; w <= 20; w++ {
		go pipelineCrystalWorker(crystalChan, analysisChan)
	}

	if len(pdbIDsFilter) == 0 {
		// No PDB IDs specified, grab all crystals in the UniProt entry
		for _, crystal := range p.Crystals {
			crystalChan <- crystal
		}
	} else {
		for _, pdbID := range pdbIDsFilter {
			var exists bool
			for _, crystal := range p.Crystals {
				if pdbID == crystal.ID {
					exists = true
					crystalChan <- crystal
				}
			}
			// User has specified a PDB ID that's not in the UniProt entry. Fail loudly.
			if !exists {
				return nil, fmt.Errorf("specified PDB ID %s not found inside UniProt entry %s", pdbID, p.UniProt.ID)
			}
		}
	}
	close(crystalChan)

	for a := 1; a <= length; a++ {
		analysis := <-analysisChan
		if analysis.Error != nil {
			return nil, analysis.Error
		}
		p.PDBAnalysis = append(p.PDBAnalysis, analysis)
	}

	end := time.Since(start)
	log.Printf("Finished UniProt %s in %f secs", p.UniProt.ID, end.Seconds())
	return p, nil
}
