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
	UniProt       *uniprot.UniProt `json:"uniprot"`
	PDB           *pdb.PDB         `json:"pdb"`
	Variants      []Variant        `json:"variants"`
	Interaction   Interaction      `json:"interaction"`
	Exposure      Exposure         `json:"exposure"`
	Conservation  Conservation     `json:"conservation"`
	Fpocket       Fpocket          `json:"fpocket"`
	ActiveSite    ActiveSite       `json:"activeSite"`
	Switchability Switchability    `json:"switchability"`
	Aggregability Aggregability    `json:"aggregability"`
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

type Switchability struct {
	Positions []PositionValue `json:"positions"`
}

type Aggregability struct {
	Positions []PositionValue `json:"positions"`
}

type PositionValue struct {
	Position int64   `json:"position"`
	Value    float64 `json:"value"`
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
	start := time.Now()
	pl.msg("Job started")

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
		activeSiteChan := pl.activeSiteRunner(u, p)
		exposureChan := pl.exposureRunner(p)
		conservationChan := pl.conservationRunner(pl.UniProt)
		fpocketChan := pl.fpocketRunner(p)
		switchabilityChan := pl.switchabilityRunner(u, p)
		aggregabilityChan := pl.aggregabilityRunner(u, p)

		results.Interaction = <-interactionChan
		results.ActiveSite = <-activeSiteChan
		results.Exposure = <-exposureChan
		results.Conservation = <-conservationChan
		results.Fpocket = <-fpocketChan
		results.Switchability = <-switchabilityChan
		results.Aggregability = <-aggregabilityChan

		pl.Results[pdbID] = &results
	}

	pl.Duration = time.Now().Sub(start)
	pl.msg(fmt.Sprintf("Pipeline finished in %s", pl.Duration.String()))

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

func (pl *Pipeline) activeSiteRunner(u *uniprot.UniProt, p *pdb.PDB) chan ActiveSite {
	rchan := make(chan ActiveSite)
	go func() {
		results := ActiveSite{}

		pl.msg(fmt.Sprintf("Compute active site by distance to catalytic residues with PDB %s", p.ID))
		for _, site := range u.Sites {
			if site.Type == "active" {
				if residues, ok := p.UniProtPositions[u.ID][site.Position]; ok {
					for _, catRes := range residues {
						for _, res := range pdb.CloseResidues(p, catRes, 5) {
							results.Residues = append(results.Residues, Residue{
								Residue:  res,
								Position: res.UnpPosition,
							})
						}
					}
				}
			}
		}

		pl.msg(fmt.Sprintf("Found %d active site residues in PDB %s", len(results.Residues), p.ID))

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

func (pl *Pipeline) switchabilityRunner(u *uniprot.UniProt, p *pdb.PDB) chan Switchability {
	rchan := make(chan Switchability)
	go func() {
		results := Switchability{}
		pl.msg(fmt.Sprintf("Running abSwitch for chain seqs in PDB %s", p.ID))
		res, err := instances.AbSwitch.Switchability(u, p)
		if err != nil {
			pl.Error = err
			rchan <- results
			return
		}

		for pos, r := range res {
			if r.S5s > 5 {
				results.Positions = append(results.Positions, PositionValue{Position: pos, Value: r.S5s})
			}
		}

		pl.msg(fmt.Sprintf("%d high switchability residues found for PDB %s",
			len(results.Positions), p.ID))

		rchan <- results
	}()
	return rchan
}

func (pl *Pipeline) aggregabilityRunner(u *uniprot.UniProt, p *pdb.PDB) chan Aggregability {
	rchan := make(chan Aggregability)
	go func() {
		results := Aggregability{}
		pl.msg(fmt.Sprintf("Running Tango for chain seqs in PDB %s", p.ID))
		res, err := instances.Tango.Aggregability(u, p)
		if err != nil {
			pl.Error = err
			rchan <- results
			return
		}

		for pos, r := range res {
			if r.Aggregation > 5 {
				results.Positions = append(results.Positions, PositionValue{Position: pos, Value: r.Aggregation})
			}
		}

		pl.msg(fmt.Sprintf("%d high aggregability residues found for PDB %s",
			len(results.Positions), p.ID))

		rchan <- results
	}()
	return rchan
}
