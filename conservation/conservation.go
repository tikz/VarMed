package conservation

import (
	"fmt"
	"sync"
	"time"
	"varq/conservation/pfam"
	"varq/uniprot"
)

var fileMux sync.Mutex

// Results holds the collected data in the conservation analysis step
type Results struct {
	Families []*pfam.Family `json:"families"`
	Duration time.Duration  `json:"duration"`
	Error    error          `json:"error"`
}

// Run starts the conservation analysis step
func Run(unp *uniprot.UniProt, results chan<- *Results, msg func(string)) {
	start := time.Now()

	fams, err := pfam.LoadFamilies(unp, &fileMux)
	if err != nil {
		results <- &Results{Error: fmt.Errorf("Pfam: %v", err)}
	}

	results <- &Results{
		Families: fams,
		Duration: time.Since(start),
	}
}