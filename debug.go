package main

import (
	"fmt"
	"strconv"
	"varq/pdb"

	"github.com/logrusorgru/aurora"
)

func debugPrintChains(a *Analysis) {
	var pocketResidues []*pdb.Residue
	for _, pocket := range a.Binding.Pockets {
		pocketResidues = append(pocketResidues, pocket.Residues...)
	}

	if len(a.Interaction.Residues) > 0 {
		debugPrintChainsMarkedResidues("Interface residues by distance", a.PDB, a.Interaction.Residues, nil)
	}

	if a.Exposure != nil {
		if len(a.Exposure.ExposedResidues) > 0 {
			debugPrintChainsMarkedResidues("Exposed residues", a.PDB, a.Exposure.ExposedResidues, nil)
		}
	}

	if len(pocketResidues) > 0 {
		debugPrintChainsMarkedResidues("Fpocket", a.PDB, pocketResidues, nil)
	}

	if a.Binding.Catalytic != nil {
		debugPrintChainsMarkedResidues("M-CSA", a.PDB, a.Binding.Catalytic.Residues, nil)
	}

	if len(a.PDB.BindingSite) > 0 {
		var residues []*pdb.Residue

		for _, rs := range a.PDB.BindingSite {
			residues = append(residues, rs...)
		}
		e := func() {
			for site, desc := range a.PDB.BindingSiteDesc {
				fmt.Print(aurora.BrightGreen(site), ": ", desc, " | ")
			}
			fmt.Println()
		}
		debugPrintChainsMarkedResidues("PDB SITE records", a.PDB, residues, e)
	}

	var famRes []*pdb.Residue
	for _, fam := range a.PDB.SIFTS.Pfam {
		for _, m := range fam.Mappings {
			for i := m.PDBStart.ResidueNumber; i < m.PDBEnd.ResidueNumber; i++ {
				famRes = append(famRes, a.PDB.SeqResChains[m.ChainID][i])
			}
		}
	}
	e := func() {
		for id, fam := range a.PDB.SIFTS.Pfam {
			fmt.Println(aurora.BrightGreen(id), aurora.BrightGreen(fam.Name), "-", fam.Description)
		}
	}
	debugPrintChainsMarkedResidues("Pfam", a.PDB, famRes, e)

	fmt.Println()
}

func residueExists(res *pdb.Residue, resList []*pdb.Residue) bool {
	for _, r := range resList {
		if r == res {
			return true
		}
	}
	return false
}

func debugPrintChainsMarkedResidues(analysisName string, pdb *pdb.PDB, aRes []*pdb.Residue, extra func()) {
	fmt.Println("==============================================================================")
	fmt.Println(aurora.BgBlack(aurora.Bold(aurora.Cyan(analysisName))))

	if extra != nil {
		fmt.Println("-----------------------------------------")
		extra()
	}

	for _, mapping := range pdb.SIFTS.UniProt[pdb.UniProtID].Mappings {
		residues := pdb.SeqRes[mapping.ChainID]
		unpStart := int(mapping.UnpStart)
		pdbStart := int(mapping.PDBStart.ResidueNumber)
		fmt.Println("---------", pdb.ID, "Chain", mapping.ChainID, "-", pdb.UniProtID, "---------")

		// Ruler
		fmt.Print("             ")
		for i := 0; i < pdbStart; i++ {
			fmt.Print(" ")
		}
		fmt.Print(aurora.Underline("1"), "        ")
		for i := 10; i < len(pdb.UniProtSequence)-20; i = i + 10 {
			n := strconv.Itoa(i)
			fmt.Print(aurora.Bold(aurora.Underline(n[:1])), n[1:], "         ")
		}
		fmt.Println()

		fmt.Print(">UNIPROT     ")
		for i := 0; i < pdbStart; i++ {
			fmt.Print(" ")
		}
		fmt.Print(pdb.UniProtSequence)
		fmt.Println()
		// fmt.Print(">SEQRES      ")
		// for i := 0; i < unpStart; i++ {
		// 	fmt.Print(" ")
		// }
		// for _, res := range residues {
		// 	fmt.Printf(res.Abbrv1)
		// }
		// fmt.Println()
		fmt.Print(">PDB         ")
		for i := 1; i < unpStart; i++ {
			fmt.Print(" ")
		}
		for i := range residues {
			res, ok := pdb.SeqResChains[mapping.ChainID][int64(i)]
			if ok {
				if residueExists(res, aRes) {
					fmt.Print(aurora.BgRed(aurora.Bold(aurora.Yellow(res.Abbrv1))))
					fmt.Print()
				} else {
					fmt.Print(res.Abbrv1)
				}
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Println()
	}
	// fmt.Println("==============================================================================")
}
