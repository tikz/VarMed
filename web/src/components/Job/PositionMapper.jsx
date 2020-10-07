export default class PositionMapper {
  constructor(res) {
    this.seqResOffsets = res.pdb.seqResOffsets;
    this.chainStartResN = res.pdb.chainStartResNumber;
    this.mappings = res.pdb.SIFTS.UniProt[res.uniprot.id].mappings;
    this.loadOffsets();
    this.loadChains();
    this.loadChainRanges();
    this.makeMappings();
  }

  makeMappings() {
    this.mapUnpPDB = {};
    this.mapPDBUnp = {};
    this.mappings.forEach((chain) => {
      for (let i = chain.unp_start; i < chain.unp_end; i++) {
        this.mapUnpPDB[i] = this.mapUnpPDB[i] || [];
        this.mapUnpPDB[i].push({
          chain: chain.chain_id,
          position: i - chain.unp_start + chain.start.residue_number,
        });

        this.mapPDBUnp[chain.chain_id] = this.mapPDBUnp[chain.chain_id] || {};
        this.mapPDBUnp[chain.chain_id][
          i + chain.unp_start - chain.start.residue_number
        ] = i;
      }
    });
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

  loadChainRanges() {
    this.chainRanges = {};
    this.mappings.forEach((chain) => {
      this.chainRanges[chain.chain_id] = this.chainRanges[chain.chain_id] || [];
      this.chainRanges[chain.chain_id].push({
        start: chain.unp_start,
        end: chain.unp_end,
      });
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
    return this.mapPDBUnp[chain][pos];
    return pos + this.unpOffsets[chain];
  }

  unpToPDB(pos) {
    return this.mapUnpPDB[pos];
    let residues = [];
    this.mappings.forEach((chain) => {
      residues.push({
        chain: chain.chain_id,
        position: pos + this.pdbOffsets[chain.chain_id],
      });
    });
    return residues;
  }
}
