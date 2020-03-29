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
	Binding     *binding.BindingAnalysis
	Interaction *interaction.InteractionAnalysis
	Exposure    *exposure.ExposureAnalysis
	Error       error `json:"-"`
}

// pipelinePDBWorker fetches a single PDB crystal file.
func pipelinePDBWorker(pdbChan <-chan *pdb.PDB, analysisChan chan<- *Analysis) {
	for crystal := range pdbChan {
		crystal.Fetch()
		if crystal.Error != nil {
			analysisChan <- &Analysis{PDB: crystal, Error: fmt.Errorf("PDB %s: %v", crystal.ID, crystal.Error)}
			continue
		}

		analysisChan <- analysePDB(&Analysis{PDB: crystal})
	}
}

// analysePDB runs each available analysis in parallel for a single structure.
func analysePDB(analysis *Analysis) *Analysis {
	return analysis // TODO: remove
	// Create temp PDB on filesystem for analysis with external tools
	analysis.PDB.LocalFilename = "varq_" + analysis.PDB.ID
	analysis.PDB.LocalPath = "/tmp/" + analysis.PDB.LocalFilename + ".pdb"

	err := ioutil.WriteFile(analysis.PDB.LocalPath, analysis.PDB.RawPDB, 0644)
	if err != nil {
		analysis.Error = fmt.Errorf("create tmp PDB: %v", err)
		return analysis
	}

	defer func() {
		os.Remove(analysis.PDB.LocalPath)
	}()

	bindingChan := make(chan *binding.BindingAnalysis)
	interactionChan := make(chan *interaction.InteractionAnalysis)
	exposureChan := make(chan *exposure.ExposureAnalysis)

	go binding.RunBindingAnalysis(analysis.PDB, bindingChan)
	go interaction.RunInteractionAnalysis(analysis.PDB, interactionChan)
	go exposure.RunExposureAnalysis(analysis.PDB, exposureChan)

	// TODO: Maybe refactor these repeated patterns
	bindingRes := <-bindingChan
	if bindingRes.Error != nil {
		analysis.Error = fmt.Errorf("binding analysis: %v", bindingRes.Error)
		return analysis
	}
	analysis.Binding = bindingRes
	log.Printf("PDB %s binding analysis done in %.3f secs", analysis.PDB.ID, bindingRes.Duration.Seconds())

	interactionRes := <-interactionChan
	if interactionRes.Error != nil {
		analysis.Error = fmt.Errorf("interaction analysis: %v", interactionRes.Error)
		return analysis
	}
	analysis.Interaction = interactionRes
	log.Printf("PDB %s interaction analysis done in %.3f secs", analysis.PDB.ID, interactionRes.Duration.Seconds())

	exposureRes := <-exposureChan
	if exposureRes.Error != nil {
		analysis.Error = fmt.Errorf("exposure analysis: %v", exposureRes.Error)
		return analysis
	}
	analysis.Exposure = exposureRes
	log.Printf("PDB %s exposure analysis done in %.3f secs", analysis.PDB.ID, exposureRes.Duration.Seconds())

	return analysis
}

func RunPipeline(pdbs []*pdb.PDB) (analyses []*Analysis, err error) {
	length := len(pdbs)
	pdbChan := make(chan *pdb.PDB, length)
	analysisChan := make(chan *Analysis, length)

	// Launch crystal workers
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
		}
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

	u.CleanCrystals()

	// for _, crystal := range u.Crystals {
	// 	if _, ok := crystal.SIFTS.UniProtIDs[u.ID].Chains["A"]; ok {
	// 		fmt.Println(u.ID, crystal.ID)
	// 		fmt.Println(crystal.SIFTS.UniProtIDs)
	// 		fmt.Println("UNIPROT", u.Sequence[crystal.SIFTS.UniProtIDs[u.ID].Chains["A"].UniProtStart-1:crystal.SIFTS.UniProtIDs[u.ID].Chains["A"].UniProtEnd])
	// 		fmt.Print("SEQRES  ")
	// 		start := crystal.SIFTS.UniProtIDs[u.ID].Chains["A"].PDBStart - 1
	// 		for i := start; i < int64(len(crystal.SeqRes["A"])); i++ {
	// 			fmt.Print(crystal.SeqRes["A"][i].Abbrv1)
	// 		}
	// 		fmt.Println()
	// 		fmt.Print("PDB     ")
	// 		for i := 1; i < len(crystal.Chains["A"]); i++ {
	// 			if pos, ok := crystal.Chains["A"][int64(i)-crystal.ChainsOffsets["A"]-1]; ok {
	// 				fmt.Print(pos.Abbrv1)
	// 			} else {
	// 				fmt.Print(" ")
	// 			}
	// 		}
	// 		fmt.Println()
	// 	}

	// }

	end := time.Since(start)
	log.Printf("Finished UniProt %s in %.3f secs", u.ID, end.Seconds())
	return analyses, nil
}

// RunPipelineForPDBs grabs and analyses structures from a given slice of PDB IDs.
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
	log.Printf("Finished PDBs %s in %.3f secs", PDBIDs, end.Seconds())
	return analyses, nil
}
