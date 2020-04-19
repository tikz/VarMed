package main

import (
	"fmt"
	"strconv"
	"strings"
	"varq/pdb"
	"varq/uniprot"

	"github.com/logrusorgru/aurora"
)

func printResults(r *Results) {
	if r.Interaction != nil && len(r.Interaction.Residues) > 0 {
		printInterface(r)
	}

	if r.Exposure != nil && len(r.Exposure.Residues) > 0 {
		printBuried(r)
	}

	if r.Binding != nil && len(r.Binding.Pockets) > 0 {
		printFpocket(r)
	}

	if r.Binding != nil && r.Binding.Catalytic != nil {
		printMCSA(r)
	}

	if r.Binding != nil && len(r.Binding.Ligands) > 0 {
		printNearLigands(r)
	}

	if len(r.PDB.BindingSite) > 0 {
		printSiteRecords(r)
	}

	fmt.Println()
}

func printInterface(a *Results) {
	printResultsBlock("Interface residues by distance", a.UniProt, a.PDB, a.Interaction.Residues, nil)
}

func printBuried(a *Results) {
	printResultsBlock("Buried residues", a.UniProt, a.PDB, a.Exposure.Residues, nil)
}

func printFpocket(a *Results) {
	var pocketResidues []*pdb.Residue
	for _, pocket := range a.Binding.Pockets {
		pocketResidues = append(pocketResidues, pocket.Residues...)
	}
	printResultsBlock("Fpocket", a.UniProt, a.PDB, pocketResidues, nil)
}

func printMCSA(a *Results) {
	printResultsBlock("M-CSA", a.UniProt, a.PDB, a.Binding.Catalytic.Residues, nil)
}

func printNearLigands(r *Results) {
	e := func() {
		for name, res := range r.Binding.Ligands {
			var residues []string
			for _, r := range res {
				residues = append(residues, r.Chain+"-"+r.Abbrv3+strconv.FormatInt(r.Position, 10))
			}
			fmt.Println("Ligand", aurora.BrightGreen(name), "-", aurora.Red(strings.Join(residues, " ")))
		}
	}
	var res []*pdb.Residue
	for _, ligand := range r.Binding.Ligands {
		res = append(res, ligand...)
	}
	printResultsBlock("Residues near ligands", r.UniProt, r.PDB, res, e)
}

func printSiteRecords(r *Results) {
	var residues []*pdb.Residue

	for _, rs := range r.PDB.BindingSite {
		residues = append(residues, rs...)
	}
	e := func() {
		for site, desc := range r.PDB.BindingSiteDesc {
			var residues []string
			for _, res := range r.PDB.BindingSite[site] {
				residues = append(residues, res.Chain+"-"+res.Abbrv3+strconv.FormatInt(res.Position, 10))
			}
			fmt.Print(aurora.BrightGreen(site), " (", aurora.Red(strings.Join(residues, " ")), "): ", desc, " | ")
		}
		fmt.Println()
	}
	printResultsBlock("PDB SITE records", r.UniProt, r.PDB, residues, e)
}
func printPfam(r *Results) {
	var famRes []*pdb.Residue
	for _, fam := range r.PDB.SIFTS.Pfam {
		for _, m := range fam.Mappings {
			for i := m.PDBStart.ResidueNumber; i <= m.PDBEnd.ResidueNumber; i++ {
				famRes = append(famRes, r.PDB.SeqResChains[m.ChainID][i])
			}
		}
	}
	e := func() {
		for id, fam := range r.PDB.SIFTS.Pfam {
			var chains []string
			for _, m := range fam.Mappings {
				chains = append(chains, m.ChainID)
			}
			fmt.Println(aurora.BrightGreen(id), aurora.BrightGreen(fam.Name), "("+strings.Join(chains, ", ")+")", "-", fam.Description)
		}
	}
	printResultsBlock("Pfam", r.UniProt, r.PDB, famRes, e)
}

func residueExists(res *pdb.Residue, resList []*pdb.Residue) bool {
	for _, r := range resList {
		if r == res {
			return true
		}
	}
	return false
}

func printResultsBlock(name string, unp *uniprot.UniProt, pdb *pdb.PDB, res []*pdb.Residue, extra func()) {
	fmt.Println("==============================================================================")
	fmt.Println(aurora.BgBlack(aurora.Bold(aurora.Cyan(name))))
	if extra != nil {
		fmt.Println("-----------------------------------------")
		extra()
	}

	for _, mapping := range pdb.SIFTS.UniProt[unp.ID].Mappings {
		residues := pdb.SeqRes[mapping.ChainID]
		unpStart := int(mapping.UnpStart)
		pdbStart := int(mapping.PDBStart.ResidueNumber)
		fmt.Println("---------", pdb.ID, "Chain", mapping.ChainID, "-", unp.ID, "---------")

		if cfg.DebugPrint.Rulers.UniProt {
			fmt.Print("             ")
			for i := 0; i < pdbStart; i++ {
				fmt.Print(" ")
			}
			fmt.Print(aurora.Underline("1"), "        ")
			for i := 10; i < len(unp.Sequence); i = i + 10 {
				n := strconv.Itoa(i)
				fmt.Print(aurora.Bold(aurora.Underline(n[:1])), n[1:])
				for s := 0; s < 10-len(n); s++ {
					fmt.Print(" ")
				}
			}
			fmt.Println()
		}

		fmt.Print(">UNIPROT     ")
		for i := 0; i < pdbStart; i++ {
			fmt.Print(" ")
		}
		fmt.Print(unp.Sequence)
		fmt.Println()

		fmt.Print(">SEQRES      ")
		for i := 0; i < unpStart; i++ {
			fmt.Print(" ")
		}
		for _, r := range residues {
			fmt.Printf(r.Abbrv1)
		}
		fmt.Println()

		if cfg.DebugPrint.Rulers.PDB {
			fmt.Print("             ")
			for i := 0; i < unpStart; i++ {
				fmt.Print(" ")
			}
			fmt.Print(aurora.Underline("1"), "        ")
			for i := 10; i < len(pdb.Chains[mapping.ChainID]); i = i + 10 {
				n := strconv.Itoa(i)
				fmt.Print(aurora.Bold(aurora.Underline(n[:1])), n[1:])
				for s := 0; s < 10-len(n); s++ {
					fmt.Print(" ")
				}
			}
			fmt.Println()
		}

		fmt.Print(">PDB         ")
		for i := 1; i < unpStart; i++ {
			fmt.Print(" ")
		}
		for i := range residues {
			r, ok := pdb.SeqResChains[mapping.ChainID][int64(i)]
			if ok {
				if residueExists(r, res) {
					fmt.Print(aurora.BgRed(aurora.Bold(aurora.Yellow(r.Abbrv1))))
				} else {
					fmt.Print(r.Abbrv1)
				}
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Println()
	}
}
