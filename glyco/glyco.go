package glyco

import (
	"errors"
	"varq/interaction"
	"varq/pdb"
	"varq/uniprot"
)

type GlycoDist struct {
	MinSeqDist    int64   `json:"minSeqDist"`
	MinStructDist float64 `json:"minStructDist"`
}

func CalculateMinGlycoDist(pos int64, unp *uniprot.UniProt, p *pdb.PDB) *GlycoDist {
	if len(unp.PTMs.Glycosilations) == 0 {
		return nil
	}

	minStructDist, err := MinGlycoStructDist(pos, unp, p)
	if err != nil {
		return nil
	}

	minSeqDist := MinGlycoSeqDist(pos, unp)

	return &GlycoDist{MinSeqDist: minSeqDist, MinStructDist: minStructDist}

}

func MinGlycoStructDist(pos int64, unp *uniprot.UniProt, p *pdb.PDB) (float64, error) {
	var minDist float64
	for _, g := range unp.PTMs.Glycosilations {
		if glycoSites, ok := p.UniProtPositions[unp.ID][g.Position]; ok {
			if residues, ok := p.UniProtPositions[unp.ID][pos]; ok {
				if len(residues) == 0 || len(glycoSites) == 0 {
					return 0, errors.New("not in structure")
				}
				minDist = residuesDistance(residues[0], glycoSites[0])
				for _, glycoRes := range glycoSites {
					for _, res := range residues {
						dist := residuesDistance(res, glycoRes)
						if dist < minDist {
							minDist = dist
						}
					}
				}
			} else {
				return 0, errors.New("position not in structure")
			}
		} else {
			return 0, errors.New("glycosilation site not in structure")
		}
	}

	return minDist, nil
}

func MinGlycoSeqDist(pos int64, unp *uniprot.UniProt) int64 {
	minDist := int64(len(unp.Sequence))

	for _, g := range unp.PTMs.Glycosilations {
		dist := g.Position - pos
		if abs(dist) < abs(minDist) {
			minDist = dist
		}
	}

	return minDist
}

func abs(i int64) int64 {
	if i < 0 {
		return i * -1
	}
	return i
}

func residuesDistance(res1 *pdb.Residue, res2 *pdb.Residue) float64 {
	minDist := interaction.Distance(res1.Atoms[0], res2.Atoms[0])
	for _, a1 := range res1.Atoms {
		for _, a2 := range res2.Atoms {
			dist := interaction.Distance(a1, a2)
			if dist < minDist {
				minDist = dist
			}
		}
	}

	return minDist
}
