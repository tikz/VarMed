package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"strconv"
	"strings"
)

// ResultsCSV returns the CSV for a given PDB in a job.
func ResultsCSV(results *Results) string {
	buf := new(bytes.Buffer)
	writer := csv.NewWriter(buf)
	writer.Write([]string{"UniProt ID", "PDB ID", "PDB Position", "Position", "From Aa", "To Aa",
		"Family", "Conservation Bitscore", "Active Site", "Interface", "Buried", "High Aggregability",
		"High Switchability", "DDG", "Outcome", "PubMed IDs", "dbSNP ID", "ClinVar Sig",
		"ClinVar Phenotypes"})

	uniprotID := results.UniProt.ID
	pdbID := results.PDB.ID
	for _, v := range results.Variants {
		position := v.Position
		fromAa := v.FromAa
		toAa := v.ToAa

		pdbPosition := results.PDB.UniProtPositions[uniprotID][position][0].Position

		// Conservation
		var consBitscore float64
		var family string
		for _, fam := range results.Conservation.Families {
			family = fam.ID
			for _, p := range fam.Positions {
				if p.Position == position {
					consBitscore = p.Bitscore
				}
			}
		}

		posExistsResidues := func(pos int64, residues []Residue) bool {
			for _, res := range residues {
				if res.Position == position {
					return true
				}
			}
			return false
		}

		posExistsPosVal := func(pos int64, positions []PositionValue) bool {
			for _, res := range positions {
				if res.Position == position {
					return true
				}
			}
			return false
		}

		posExistsResExp := func(pos int64, residues []ResidueExposure) bool {
			for _, res := range residues {
				if res.Position == position {
					return true
				}
			}
			return false
		}

		// Active site
		activeSite := posExistsResidues(position, results.ActiveSite.Residues)
		interaction := posExistsResidues(position, results.Interaction.Residues)
		buried := posExistsResExp(position, results.Exposure.Residues)
		highSwi := posExistsPosVal(position, results.Switchability.Positions)
		highAgg := posExistsPosVal(position, results.Aggregability.Positions)

		ddg := v.DdG
		outcome := v.Outcome
		pubmedIDs := v.PubMedIDs
		dbSNPID := v.DbSNPID
		cvSig := v.CVClinSig
		cvPhenotypes := v.CVPhenotypes

		writer.Write([]string{uniprotID,
			pdbID,
			fmt.Sprintf("%d", pdbPosition),
			fmt.Sprintf("%d", position),
			fromAa,
			toAa,
			family,
			fmt.Sprintf("%f", consBitscore),
			strconv.FormatBool(activeSite),
			strconv.FormatBool(interaction),
			strconv.FormatBool(buried),
			strconv.FormatBool(highAgg),
			strconv.FormatBool(highSwi),
			fmt.Sprintf("%f", ddg),
			outcome,
			strings.Join(pubmedIDs, ", "),
			dbSNPID,
			cvSig,
			cvPhenotypes})

	}

	writer.Flush()
	return buf.String()
}
