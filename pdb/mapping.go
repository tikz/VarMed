package pdb

import "fmt"

func (pdb *PDB) makeMappings() {
	// SEQRES chain and pos to residues
	pdb.SeqResChains = make(map[string]map[int64]*Aminoacid)

	for chain, offset := range pdb.SeqResOffsets {
		pdb.SeqResChains[chain] = make(map[int64]*Aminoacid)
		minPos := pdb.minChainPos(chain)
		for pos, aa := range pdb.Chains[chain] {
			pdb.SeqResChains[chain][pos+offset-minPos] = aa
		}
	}

	// UniProt canonical sequence to residues
	pdb.UniProtPositions = make(map[int64][]*Aminoacid)
	fmt.Println(pdb.UniProtID, pdb.ID)
	chainMapping := pdb.SIFTS.UniProtIDs[pdb.UniProtID].Chains
	for chain, mapping := range chainMapping {
		for i := mapping.UniProtStart; i <= mapping.UniProtEnd; i++ {
			if res, ok := pdb.SeqResChains[chain][mapping.PDBStart+i-2]; ok {
				pdb.UniProtPositions[i] = append(pdb.UniProtPositions[i], res)
			}
		}
	}

	pdb.debugAlignment()
}

// This alignment needs to be done since the residue numbers in ATOM tags doesn't always coincide with SEQRES positions.
// TODO: see if there is a value available somewhere to skip this.
func (pdb *PDB) calculateChainsOffset() {
	pdb.SeqResOffsets = make(map[string]int64)
	for chain := range pdb.Chains {
		var bestOffset, bestScore int

		minPos := pdb.minChainPos(chain)
		chainLength := pdb.maxChainPos(chain) - minPos
		steps := len(pdb.SeqRes[chain]) - int(chainLength)

		for offset := 0; offset < steps; offset++ {
			score := 0
			for pos, aa := range pdb.Chains[chain] {
				seqResPos := pos + int64(offset) - minPos
				if aa.Abbrv1 == pdb.SeqRes[chain][seqResPos].Abbrv1 {
					score++
				}
			}
			if score > bestScore {
				bestScore = score
				bestOffset = offset
			}
		}

		pdb.SeqResOffsets[chain] = int64(bestOffset)
	}
}

func (pdb *PDB) chainKeys(chain string) (k []int64) {
	for pos := range pdb.Chains[chain] {
		k = append(k, pos)
	}
	return k
}

func (pdb *PDB) minChainPos(chain string) int64 {
	ck := pdb.chainKeys(chain)
	min := ck[0]
	for _, pos := range ck {
		if pos < min {
			min = pos
		}
	}
	return min
}

func (pdb *PDB) maxChainPos(chain string) int64 {
	ck := pdb.chainKeys(chain)
	max := ck[0]
	for _, pos := range ck {
		if pos > max {
			max = pos
		}
	}
	return max
}

func (pdb *PDB) debugAlignment() {
	for chain, mapping := range pdb.SIFTS.UniProtIDs[pdb.UniProtID].Chains {
		aas := pdb.SeqRes[chain]
		unpStart := int(mapping.UniProtStart)
		pdbStart := int(mapping.PDBStart)
		fmt.Println("-----")
		fmt.Println(pdb.ID, "Chain", chain, "-", pdb.UniProtID)
		fmt.Print(">UNIPROT     ")
		for i := 1; i < pdbStart; i++ {
			fmt.Print(" ")
		}
		fmt.Print(pdb.UniProtSequence)
		fmt.Println()
		fmt.Print(">SEQRES      ")
		for i := 1; i < unpStart; i++ {
			fmt.Print(" ")
		}
		for _, aa := range aas {
			fmt.Printf(aa.Abbrv1)
		}
		fmt.Println()
		fmt.Print(">CRYSTAL     ")
		for i := 1; i < unpStart; i++ {
			fmt.Print(" ")
		}
		for i := range aas {
			aa, ok := pdb.SeqResChains[chain][int64(i)]
			if ok {
				fmt.Print(aa.Abbrv1)
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Println()
	}
	fmt.Println("================")
}
