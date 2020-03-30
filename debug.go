package main

import (
	"fmt"
	"varq/pdb"

	"github.com/logrusorgru/aurora"
)

func debugPrintChains(a *Analysis) {
	var pocketResidues []*pdb.Residue
	for _, pocket := range a.Binding.Pockets {
		pocketResidues = append(pocketResidues, pocket.Residues...)
	}

	if len(pocketResidues) > 0 {
		debugPrintChainsMarkedResidues("Fpocket", a.PDB, pocketResidues)
	}

	if a.Binding.Catalytic != nil {
		debugPrintChainsMarkedResidues("M-CSA", a.PDB, a.Binding.Catalytic.Residues)
	}

	if len(a.Interaction.Residues) > 0 {
		debugPrintChainsMarkedResidues("Interface residues by distance", a.PDB, a.Interaction.Residues)
	}

	if len(a.Exposure.ExposedResidues) > 0 {
		debugPrintChainsMarkedResidues("Exposed residues", a.PDB, a.Exposure.ExposedResidues)
	}

}

func residueExists(res *pdb.Residue, resList []*pdb.Residue) bool {
	for _, r := range resList {
		if r == res {
			return true
		}
	}
	return false
}

func debugPrintChainsMarkedResidues(analysisName string, pdb *pdb.PDB, aRes []*pdb.Residue) {
	fmt.Println("==============================================================================")
	fmt.Println(aurora.BgBlack(aurora.Bold(aurora.Cyan(analysisName))))
	for chain, mapping := range pdb.SIFTS.UniProtIDs[pdb.UniProtID].Chains {
		residues := pdb.SeqRes[chain]
		unpStart := int(mapping.UniProtStart)
		pdbStart := int(mapping.PDBStart)
		fmt.Println("---------", pdb.ID, "Chain", chain, "-", pdb.UniProtID, "---------")
		fmt.Print(">UNIPROT     ")
		for i := 0; i < pdbStart; i++ {
			fmt.Print(" ")
		}
		fmt.Print(pdb.UniProtSequence)
		fmt.Println()
		fmt.Print(">SEQRES      ")
		for i := 0; i < unpStart; i++ {
			fmt.Print(" ")
		}
		for _, res := range residues {
			fmt.Printf(res.Abbrv1)
		}
		fmt.Println()
		fmt.Print(">PDB         ")
		for i := 1; i < unpStart; i++ {
			fmt.Print(" ")
		}
		for i := range residues {
			res, ok := pdb.SeqResChains[chain][int64(i)]
			if ok {
				if residueExists(res, aRes) {
					fmt.Print(aurora.BgRed(aurora.Bold(aurora.Yellow(res.Abbrv1))))
				} else {
					fmt.Print(res.Abbrv1)
				}
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Println()
	}
	fmt.Println("==============================================================================")
}
