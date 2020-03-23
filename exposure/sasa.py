# Usage: python3 exposed.py <PDB file path> <first residue number> <...> <n residue number>
# Launches PyMOL in headless mode and outputs the
# residue [residue number, sasa area, total area, sasa/total, exposed as 0 1 flag] to stdout.

import sys

from pymol import cmd
import pymol

import __main__
__main__.pymol_argv = ["pymol", "-qc"]

pymol.finish_launching()

path, *res_numbers = sys.argv[1:]

cmd.load(path)

cmd.set("solvent_radius", 1.4)  # default. TODO: default?
for n in res_numbers:
    # Residue SASA
    cmd.set("dot_solvent", 1)
    sasa = cmd.get_area("resi " + n)

    # Residue total area
    cmd.set("dot_solvent", 0)
    total = cmd.get_area("resi " + n)

    area_prop = sasa / total
    exposed = 1 if area_prop > 0.5 else 0

    print(n, sasa, total, area_prop, exposed)
