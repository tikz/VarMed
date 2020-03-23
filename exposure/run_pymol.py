# Usage: python3 run_pymol.py <PDB file path> <start> <end>
# Launches PyMOL in headless mode and for each residue outputs [chain residue_number b_factor sasa/area] to stdout.
import sys

from pymol import cmd, stored
import pymol
import __main__
__main__.pymol_argv = ['pymol', '-qc']
pymol.finish_launching()

path, start, end, *_ = sys.argv[1:]

cmd.load(path)
stored.residues = []
cmd.iterate('name ca', 'stored.residues.append([chain, resi, b])')

residues = stored.residues[int(start):int(end)]

cmd.set('dot_density', 1)

cmd.set('dot_solvent', 0)
res_area = [cmd.get_area('chain %s and resi %s' % (res[0], res[1])) for res in residues]
cmd.set('dot_solvent', 1)
res_sasa = [cmd.get_area('chain %s and resi %s' % (res[0], res[1])) for res in residues]


for l in [res + [sasa / area] for res, sasa, area in zip(residues, res_sasa, res_area)]:
    print(" ".join(map(str, l)))
