package main

import (
	"fmt"
	"math"
	"time"

	"github.com/tikz/bio"
	"github.com/tikz/bio/fpocket"
	"github.com/tikz/bio/interaction"
	"github.com/tikz/bio/pdb"
	"github.com/tikz/bio/sasa"
	"github.com/tikz/bio/uniprot"
)

type Results struct {
	UniProt      *uniprot.UniProt `json:"uniprot"`
	PDB          *pdb.PDB         `json:"pdb"`
	Variants     []Variant        `json:"variants"`
	Interaction  Interaction      `json:"interaction"`
	Exposure     Exposure         `json:"exposure"`
	Conservation Conservation     `json:"conservation"`
	Fpocket      Fpocket          `json:"fpocket"`
	ActiveSite   ActiveSite       `json:"activeSite"`
}

type Residue struct {
	Residue  *pdb.Residue `json:"residue"`
	Position int64        `json:"position"`
}

type ActiveSite struct {
	Residues []Residue `json:"residues"`
}

type Fpocket struct {
	Pockets []Pocket `json:"pockets"`
}

type Pocket struct {
	Name      string    `json:"name"`
	DrugScore float64   `json:"drugScore"`
	Residues  []Residue `json:"residues"`
}

type Interaction struct {
	Residues []Residue `json:"residues"`
}

type Variant struct {
	// From request
	Residue  *pdb.Residue `json:"-"`
	FromAa   string       `json:"fromAa"`
	ToAa     string       `json:"toAa"`
	Position int64        `json:"position"`
	Change   string       `json:"change"`

	// From UniProt annotations
	Note     string `json:"note"`
	Evidence string `json:"evidence"`
	ID       string `json:"id"`
	DbSNPID  string `json:"dbSNPId"`

	// From ClinVar
	CVName         string `json:"cvName"`
	CVReviewStatus string `json:"cvReviewStatus"`
	CVClinSig      string `json:"cvClinSig"`
	CVPhenotypes   string `json:"cvPhenotypes"`

	// Calculated
	DdG     float64 `json:"ddg"`
	Outcome string  `json:"outcome"`
}

type Exposure struct {
	Residues []ResidueExposure `json:"residues"`
}

type ResidueExposure struct {
	Residue  *pdb.Residue `json:"residue"`
	Position int64        `json:"position"`
	Exposure float64      `json:"exposure"`
}

type Conservation struct {
	Families []Family `json:"families"`
}

type Family struct {
	ID        string                 `json:"id"`
	Name      string                 `json:"name"`
	Desc      string                 `json:"desc"`
	Start     int64                  `json:"start"`
	End       int64                  `json:"end"`
	Positions []PositionConservation `json:"positions"`
}

type PositionConservation struct {
	Position int64   `json:"position"`
	Bitscore float64 `json:"bitscore"`
}

// Pipeline represents a single run of the RespDB pipeline.
type Pipeline struct {
	UniProt  *uniprot.UniProt
	PDBIDs   []string
	Variants []SAS
	Results  map[string]*Results // PDB ID to results
	Duration time.Duration

	Error   error
	msgChan chan string // readable text messages about the status
}

// msg prints and sends a message with added format to the channel.
func (pl *Pipeline) msg(m string) {
	pl.msgChan <- time.Now().Format("15:04:05-0700") + " " + m
}

// NewPipeline constructs a new Pipeline.
func NewPipeline(unp *uniprot.UniProt, pdbIDs []string, variants []SAS, msgChan chan string) (*Pipeline, error) {
	p := Pipeline{
		UniProt:  unp,
		Variants: variants,
		PDBIDs:   pdbIDs,
		Results:  make(map[string]*Results),
		msgChan:  msgChan,
	}

	return &p, nil
}

// Run starts the process of analyzing given PDB IDs corresponding to an UniProt ID.
func (pl *Pipeline) Run() error {
	// start := time.Now()
	// p.msg("Job started")

	for _, pdbID := range pl.PDBIDs {
		u := pl.UniProt
		results := Results{UniProt: u}

		pl.msg(fmt.Sprintf("Loading PDB %s", pdbID))
		p, err := bio.LoadPDB(pdbID)
		if err != nil {
			pl.Error = err
			continue
		}
		results.PDB = p

		// In coverage
		var coveredVariants []SAS
		for _, v := range pl.Variants {
			inStructure := len(p.UniProtPositions[u.ID][v.Position]) > 0
			if inStructure {
				coveredVariants = append(coveredVariants, v)
			} else {
				pl.msg(fmt.Sprintf("Variant %s position not covered by PDB %s", v.Change, pdbID))
			}
		}

		// FoldX
		if len(pl.Variants) > 0 {
			pl.msg(fmt.Sprintf("Running FoldX RepairPDB %s", pdbID))
			rp, err := instances.FoldX.Repair(p)
			if err != nil {
				return err
			}
			pl.msg(fmt.Sprintf("RepairPDB %s done", pdbID))

			n := len(pl.Variants)
			varJobs := make(chan SAS, n)
			varRes := make(chan Variant, n)

			for w := 1; w <= 16; w++ {
				go pl.variantWorker(rp, u, p, varJobs, varRes)
			}

			for _, v := range coveredVariants {
				varJobs <- v
			}
			close(varJobs)

			for range coveredVariants {
				v := <-varRes
				results.Variants = append(results.Variants, v)
				pl.msg(fmt.Sprintf("BuildModel for variant %s with PDB %s done", v.Change, pdbID))
			}
		}

		// Start other runners in parallel
		interactionChan := pl.interactionRunner(p)
		exposureChan := pl.exposureRunner(p)
		conservationChan := pl.conservationRunner(pl.UniProt)
		fpocketChan := pl.fpocketRunner(p)

		results.Interaction = <-interactionChan
		results.Exposure = <-exposureChan
		results.Conservation = <-conservationChan
		results.Fpocket = <-fpocketChan

		pl.Results[pdbID] = &results
	}

	return pl.Error
}

func (pl *Pipeline) variantWorker(repairPDB string, u *uniprot.UniProt, p *pdb.PDB, sas <-chan SAS, rchan chan<- Variant) {
	for v := range sas {
		results := Variant{}
		ddg, err := instances.FoldX.BuildModelUniProt(repairPDB, p, u.ID, v.Position, v.ToAa)
		if err != nil {
			pl.Error = err
			rchan <- results
			continue
		}
		results.DdG = ddg

		results.Residue = p.UniProtPositions[u.ID][v.Position][0]
		results.FromAa = v.FromAa
		results.ToAa = v.ToAa
		results.Position = v.Position
		results.Change = v.Change

		for _, av := range u.Variants {
			if av.Change == v.Change {
				results.Note = av.Note
				results.Evidence = av.Evidence
				results.ID = av.ID
				results.DbSNPID = av.DbSNP
				if av.DbSNP != "" {
					allele := instances.ClinVar.GetVariation(av.DbSNP, av.Change)
					if allele != nil {
						results.CVName = allele.Name
						results.CVReviewStatus = allele.ReviewStatus
						results.CVClinSig = allele.ClinSig
						results.CVPhenotypes = allele.Phenotypes
					}
				}
				break
			}
		}

		rchan <- results
	}
}

func (pl *Pipeline) interactionRunner(p *pdb.PDB) chan Interaction {
	rchan := make(chan Interaction)
	go func() {
		results := Interaction{}
		interacts := interaction.Chains(p, 5)

		pl.msg(fmt.Sprintf("Compute interface by distance with PDB %s", p.ID))
		for res, interRes := range interacts {
			if len(interRes) > 0 {
				results.Residues = append(results.Residues, Residue{res, res.UnpPosition})
			}
		}

		pl.msg(fmt.Sprintf("Found %d interface residues in PDB %s", len(interacts), p.ID))

		rchan <- results
	}()
	return rchan
}

func (pl *Pipeline) exposureRunner(p *pdb.PDB) chan Exposure {
	rchan := make(chan Exposure)
	go func() {
		results := Exposure{}
		sr, err := sasa.SASA(p)
		if err != nil {
			pl.Error = err
			rchan <- results
			return
		}

		pl.msg(fmt.Sprintf("Compute solvent accesible surface area for PDB %s", p.ID))
		for res, sasa := range sr.Residues {
			if sasa.RelSide < 50 {
				re := ResidueExposure{
					Residue:  res,
					Position: res.UnpPosition,
					Exposure: sasa.RelSide,
				}
				if math.IsNaN(re.Exposure) {
					re.Exposure = 0
				}
				results.Residues = append(results.Residues, re)
			}
		}
		pl.msg(fmt.Sprintf("Done SASA for %d residues, %d buried, in PDB %s",
			len(sr.Residues), len(results.Residues), p.ID))

		rchan <- results
	}()
	return rchan
}

func (pl *Pipeline) conservationRunner(u *uniprot.UniProt) chan Conservation {
	rchan := make(chan Conservation)
	go func() {
		results := Conservation{}
		pl.msg(fmt.Sprintf("Loading Pfam families for %s sequence", u.ID))
		fams, err := instances.Pfam.Families(u)
		if err != nil {
			pl.Error = err
			rchan <- results
			return
		}

		for _, fam := range fams {
			pl.msg(fmt.Sprintf("Conservation for Pfam family %s %s", fam.ID, fam.HMM.Desc))
			rf := Family{
				ID:   fam.ID,
				Name: fam.HMM.Name,
				Desc: fam.HMM.Desc,
			}

			for _, mp := range fam.Mappings {
				rf.Positions = append(rf.Positions, PositionConservation{
					Position: int64(mp.Position),
					Bitscore: mp.Bitscore,
				})
			}

			pl.msg(fmt.Sprintf("%d aligned residues to HMM model for Pfam family %s", len(fam.Mappings), fam.ID))

			if len(rf.Positions) > 1 {
				rf.Start = rf.Positions[0].Position
				rf.End = rf.Positions[len(rf.Positions)-1].Position
			}
			pl.msg(fmt.Sprintf("%s %s ranges: %d-%d", fam.ID, fam.HMM.Desc, rf.Start, rf.End))

			results.Families = append(results.Families, rf)
		}

		rchan <- results
	}()
	return rchan
}

func (pl *Pipeline) fpocketRunner(p *pdb.PDB) chan Fpocket {
	rchan := make(chan Fpocket)
	go func() {
		results := Fpocket{}
		pl.msg(fmt.Sprintf("Searching pockets for PDB %s", p.ID))
		fp, err := fpocket.Run(cfg.Paths.Fpocket, p)
		if err != nil {
			pl.Error = err
			rchan <- results
			return
		}

		for i, pocket := range fp.Pockets {
			if pocket.DrugScore > 0.5 {
				p := Pocket{
					Name:      string(i),
					DrugScore: pocket.DrugScore,
				}

				for _, res := range pocket.Residues {
					p.Residues = append(p.Residues, Residue{res, res.UnpPosition})
				}
				results.Pockets = append(results.Pockets, p)
			}
		}

		pl.msg(fmt.Sprintf("%d suitable pockets found, %d total for PDB %s", len(results.Pockets), len(fp.Pockets), p.ID))

		rchan <- results
	}()
	return rchan
}

// 	for _, id := range p.pdbIDs {
// 		if !p.UniProt.PDBIDExists(id) {
// 			return fmt.Errorf("PDB ID %s not found", id)
// 		}
// 		pdbIDChan <- id
// 	}

// 	for a := 1; a <= len(p.pdbIDs); a++ {
// 		result := <-resChan
// 		if result.Error != nil {
// 			return fmt.Errorf("step error: %v", result.Error)
// 		}

// 		p.Results[result.PDB.ID] = result
// 	}

// 	p.Duration = time.Since(start)
// 	p.msg(fmt.Sprintf("Pipeline finished in %.3f secs", p.Duration.Seconds()))
// 	return nil
// }

// // pdbWorker fetches and loads a single PDB file.
// func (p *Pipeline) pdbWorker(pdbIDChan <-chan string, resChan chan<- *Results) {
// 	for pdbID := range pdbIDChan {
// 		results := Results{}

// 		start := time.Now()
// 		p.msg(fmt.Sprintf("Loading PDB %s...", pdbID))
// 		pdb, err := loadPDB(pdbID)
// 		if err != nil {
// 			results.Error = fmt.Errorf("load PDB %s: %v", pdbID, err)
// 			resChan <- &results
// 			continue
// 		}
// 		results.PDB = pdb
// 		results.UniProt = p.UniProt

// 		if _, ok := pdb.SIFTS.UniProt[p.UniProt.ID]; !ok {
// 			results.Error = errors.New("UniProt ID not in SIFTS data of PDB " + pdbID)
// 			resChan <- &results
// 			continue
// 		}

// 		end := time.Since(start)
// 		p.msg(fmt.Sprintf("PDB %s loaded in %.3f secs", pdbID, end.Seconds()))

// 		resChan <- p.analysePDB(&results)
// 	}
// }

// // analysePDB runs each available analysis in parallel for a single structure.
// func (p *Pipeline) analysePDB(r *Results) *Results {
// 	// Create temp PDB on filesystem for analysis with external tools
// 	path := "bin/" + r.PDB.ID + ".pdb"
// 	r.PDB.WriteFile(path)

// 	defer func() {
// 		os.Remove(path)
// 	}()

// 	bindingChan := make(chan *binding.Results)
// 	interactionChan := make(chan *interaction.Results)
// 	secondaryChan := make(chan *secondary.Results)
// 	conservationChan := make(chan *conservation.Results)
// 	exposureChan := make(chan *exposure.Results)
// 	stabilityChan := make(chan *stability.Results)

// 	idStr := fmt.Sprintf("PDB %s ", r.PDB.ID)
// 	msgPDB := func(msg string) {
// 		p.msg(idStr + msg)
// 	}

// 	if cfg.RespDB.Pipeline.EnableSteps.Binding {
// 		go binding.Run(r.UniProt, r.PDB, bindingChan, msgPDB)
// 		msgPDB("started binding analysis")
// 	}
// 	if cfg.RespDB.Pipeline.EnableSteps.Interaction {
// 		go interaction.Run(r.PDB, interactionChan, msgPDB)
// 		msgPDB("started interaction analysis")
// 	}
// 	if cfg.RespDB.Pipeline.EnableSteps.Secondary {
// 		go secondary.Run(r.UniProt, r.PDB, secondaryChan, msgPDB)
// 		msgPDB("started secondary structure analysis")
// 	}
// 	if cfg.RespDB.Pipeline.EnableSteps.Conservation {
// 		go conservation.Run(r.UniProt, conservationChan, msgPDB)
// 		msgPDB("started conservation analysis")
// 	}
// 	if cfg.RespDB.Pipeline.EnableSteps.Exposure {
// 		go exposure.Run(r.PDB, exposureChan, msgPDB)
// 		msgPDB("started exposure analysis")
// 	}
// 	if cfg.RespDB.Pipeline.EnableSteps.Stability {
// 		go stability.Run(p.SAS, r.UniProt, r.PDB, stabilityChan, msgPDB)
// 		msgPDB("started stability analysis")
// 	}

// 	// TODO: refactor these repeated patterns
// 	if cfg.RespDB.Pipeline.EnableSteps.Binding {
// 		bindingRes := <-bindingChan
// 		if bindingRes.Error != nil {
// 			r.Error = fmt.Errorf("binding analysis: %v", bindingRes.Error)
// 			return r
// 		}
// 		r.Binding = bindingRes
// 		msgPDB(fmt.Sprintf("binding analysis done in %.3f secs", bindingRes.Duration.Seconds()))
// 	}

// 	if cfg.RespDB.Pipeline.EnableSteps.Interaction {
// 		interactionRes := <-interactionChan
// 		if interactionRes.Error != nil {
// 			r.Error = fmt.Errorf("interaction analysis: %v", interactionRes.Error)
// 			return r
// 		}
// 		r.Interaction = interactionRes
// 		msgPDB(fmt.Sprintf("interaction analysis done in %.3f secs", interactionRes.Duration.Seconds()))
// 	}

// 	if cfg.RespDB.Pipeline.EnableSteps.Secondary {
// 		secondaryRes := <-secondaryChan
// 		if secondaryRes.Error != nil {
// 			r.Error = fmt.Errorf("secondary structure analysis: %v", secondaryRes.Error)
// 			return r
// 		}
// 		r.Secondary = secondaryRes
// 		msgPDB(fmt.Sprintf("secondary structure analysis done in %.3f secs", secondaryRes.Duration.Seconds()))
// 	}

// 	if cfg.RespDB.Pipeline.EnableSteps.Conservation {
// 		conservationRes := <-conservationChan
// 		if conservationRes.Error != nil {
// 			r.Error = fmt.Errorf("conservation analysis: %v", conservationRes.Error)
// 			return r
// 		}
// 		r.Conservation = conservationRes
// 		msgPDB(fmt.Sprintf("conservation analysis done in %.3f secs", conservationRes.Duration.Seconds()))
// 	}

// 	if cfg.RespDB.Pipeline.EnableSteps.Exposure {
// 		exposureRes := <-exposureChan
// 		if exposureRes.Error != nil {
// 			r.Error = fmt.Errorf("exposure analysis: %v", exposureRes.Error)
// 			return r
// 		}
// 		r.Exposure = exposureRes
// 		msgPDB(fmt.Sprintf("exposure analysis done in %.3f secs", exposureRes.Duration.Seconds()))
// 	}

// 	if cfg.RespDB.Pipeline.EnableSteps.Stability {
// 		stabilityRes := <-stabilityChan
// 		if stabilityRes.Error != nil {
// 			r.Error = fmt.Errorf("stability analysis: %v", stabilityRes.Error)
// 			return r
// 		}
// 		r.Stability = stabilityRes
// 		msgPDB(fmt.Sprintf("stability analysis done in %.3f secs", stabilityRes.Duration.Seconds()))
// 	}

// 	if cfg.DebugPrint.Enabled {
// 		printResults(r)
// 	}

// 	return r
// }
