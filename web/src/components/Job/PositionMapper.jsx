export default class PositionMapper {
  constructor(res) {
    this.seqResOffsets = res.PDB.SeqResOffsets;
    this.chainStartResN = res.PDB.ChainStartResNumber;
    this.mappings = res.PDB.SIFTS.UniProt[res.UniProt.ID].mappings;
    this.loadOffsets();
    this.loadChains();
  }

  loadOffsets() {
    this.unpOffsets = {};
    this.mappings.forEach((chain) => {
      this.unpOffsets[chain.chain_id] =
        chain.unp_start +
        this.seqResOffsets[chain.chain_id] -
        chain.start.residue_number -
        this.chainStartResN[chain.chain_id] +
        1;
    });

    this.pdbOffsets = {};
    this.mappings.forEach((chain) => {
      this.pdbOffsets[chain.chain_id] =
        -chain.unp_start + chain.start.residue_number;
    });
  }

  loadChains() {
    this.chains = this.mappings.map((chain) => {
      return {
        id: chain.chain_id,
        start:
          chain.unp_start +
          this.seqResOffsets[chain.chain_id] -
          chain.start.residue_number +
          1,
        end:
          chain.unp_end +
          this.seqResOffsets[chain.chain_id] -
          chain.start.residue_number +
          1,
      };
    });
  }

  pdbToUnp(chain, pos) {
    return pos + this.unpOffsets[chain];
  }

  unpToPDB(pos) {
    let residues = [];
    this.mappings.forEach((chain) => {
      residues.push({
        Chain: chain.chain_id,
        Position: pos + this.pdbOffsets[chain.chain_id],
      });
    });
    return residues;
  }
}
