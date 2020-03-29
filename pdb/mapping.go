package pdb

import "fmt"

func (pdb *PDB) makeMappings() {
	// SEQRES chain and pos to residues
	pdb.SeqResChains = make(map[string]map[int64]*Aminoacid)

	for chain, offset := range pdb.ChainsOffsets {
		pdb.SeqResChains[chain] = make(map[int64]*Aminoacid)
		minPos := pdb.minChainPos(chain)
		for pos, aa := range pdb.Chains[chain] {
			pdb.SeqResChains[chain][pos+offset-minPos] = aa
		}
	}

	pdb.debugAlignment()
}

// This alignment needs to be done since the residue numbers in ATOM tags doesn't always coincide with SEQRES positions.
// TODO: see if there is a value available somewhere to skip this.
func (pdb *PDB) calculateChainsOffset() {
	pdb.ChainsOffsets = make(map[string]int64)
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

		pdb.ChainsOffsets[chain] = int64(bestOffset)
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

// 6D5J 6axg 6BVK
func (pdb *PDB) debugAlignment() {
	for chain, aas := range pdb.SeqRes {
		fmt.Println("-----")
		fmt.Println(pdb.ID, "Chain", chain)
		fmt.Print("SEQRES      ")
		for _, aa := range aas {
			fmt.Printf(aa.Abbrv1)
		}
		fmt.Println()
		fmt.Print("CRYSTAL     ")
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
