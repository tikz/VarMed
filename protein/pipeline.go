package protein

import (
	"fmt"
	"log"
	"time"
	"varq/pdb"
)

func pipelineCrystalWorker(crystalChan <-chan *pdb.PDB, errChan chan<- error) {
	for c := range crystalChan {
		err := c.Fetch()
		if err != nil {
			errChan <- err
			continue
		}
		errChan <- nil
	}
}

// RunPipeline manages the workers for parallel fetching and processing of protein data
func RunPipeline(uniprotID string) (*Protein, error) {
	start := time.Now()

	p, err := NewProtein(uniprotID)
	if err != nil {
		return nil, fmt.Errorf("run pipeline: %v", err)
	}

	length := len(p.Crystals)
	crystalChan := make(chan *pdb.PDB, length)
	errChan := make(chan error, length)

	// Fetch all crystals in parallel
	for w := 1; w <= 20; w++ {
		go pipelineCrystalWorker(crystalChan, errChan)
	}

	for _, crystal := range p.Crystals {
		crystalChan <- crystal
	}
	close(crystalChan)

	for a := 1; a <= length; a++ {
		err := <-errChan
		if err != nil {
			return nil, err
		}
	}

	end := time.Since(start)
	log.Printf("Finished UniProt %s in %f secs", p.UniProt.ID, end.Seconds())
	return p, nil
}
