package pdb

func (pdb *PDB) makeMappings() {
	// SEQRES chain and position to crystal residues.
	pdb.SeqResChains = make(map[string]map[int64]*Residue)

	for chain, offset := range pdb.SeqResOffsets {
		pdb.SeqResChains[chain] = make(map[int64]*Residue)
		minPos := pdb.minChainPos(chain)
		for pos, res := range pdb.Chains[chain] {
			pdb.SeqResChains[chain][pos+offset-minPos] = res
		}
	}

	// UniProt canonical sequence position to crystal residues.
	pdb.UniProtPositions = make(map[int64][]*Residue)
	chainMapping := pdb.SIFTS.UniProtIDs[pdb.UniProtID].Chains
	for chain, mapping := range chainMapping {
		var i int64
		for i = 0; i < mapping.PDBEnd-mapping.PDBStart; i++ {
			seqResPos := i + pdb.SeqResOffsets[chain]
			unpPos := seqResPos + mapping.UniProtStart
			if res, ok := pdb.SeqResChains[chain][seqResPos]; ok {
				pdb.UniProtPositions[unpPos] = append(pdb.UniProtPositions[unpPos], res)
			}
		}
	}
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
			for pos, res := range pdb.Chains[chain] {
				seqResPos := pos + int64(offset) - minPos
				if res.Abbrv1 == pdb.SeqRes[chain][seqResPos].Abbrv1 {
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
