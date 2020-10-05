export default class PositionMapper {
  constructor(res) {
    this.seqResOffsets = res.pdb.seqResOffsets;
    this.chainStartResN = res.pdb.chainStartResNumber;
    this.mappings = res.pdb.SIFTS.UniProt[res.uniprot.id].mappings;
    this.loadOffsets();
    this.loadChains();
  }

  loadOffsets() {
    this.unpOffsets = {};
    this.mappings.forEach((chain) => {
      this.unpOffsets[chain.chain_id] =
        chain.unp_start - chain.start.residue_number;
    });

    this.pdbOffsets = {};
    this.mappings.forEach((chain) => {
      this.pdbOffsets[chain.chain_id] =
        -chain.unp_start + chain.start.residue_number;
    });
  }

  loadChains() {
    this.chainsStruct = {};
    this.chains = this.mappings.map((chain) => {
      this.chainsStruct[chain.chain_id] = chain.struct_asym_id;
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

  structChain(chain) {
    return this.chainsStruct[chain];
  }

  pdbToUnp(chain, pos) {
    return pos + this.unpOffsets[chain];
  }

  unpToPDB(pos) {
    let residues = [];
    this.mappings.forEach((chain) => {
      residues.push({
        chain: chain.chain_id,
        position: pos + this.pdbOffsets[chain.chain_id],
      });
    });
    return residues;
  }

  // unpResiduesToPDB(unpResidues) {
  //   let residues = [];
  //   unpResidues.forEach((res) => {
  //     this.mappings.forEach((chain) => {
  //       residues.push({
  //         chain: chain.chain_id,
  //         position: res.position + this.pdbOffsets[chain.chain_id],
  //       });
  //     });
  //   });

  //   return residues;
  // }
}
