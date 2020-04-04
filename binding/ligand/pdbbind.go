package ligand

import (
	"fmt"
	"io/ioutil"
	"regexp"
)

// PDBBind represents added data from the index files.
type PDBBind struct {
	Ligands map[string][]*PDBBindEntry // unique 3-letter ligand code to entries containing it
}

// PDBBindEntry represents a single PDB-ligand entry, either nucleic acid or protein.
type PDBBindEntry struct {
	PDBID       string
	Resolution  string
	Year        string
	BindingData string
	LigandID    string
	Desc        string
}

// LoadPDBBind loads the index files to a map.
func LoadPDBBind() (map[string][]*PDBBindEntry, error) {
	codes := make(map[string][]*PDBBindEntry)

	err := ParseIndexFile("static/pdbbind/INDEX_general_NL.2019", codes)
	if err != nil {
		return nil, err
	}

	err = ParseIndexFile("static/pdbbind/INDEX_general_PL.2019", codes)
	if err != nil {
		return nil, err
	}

	return codes, nil
}

// ParseIndexFile parses a single PDBBind index file.
func ParseIndexFile(path string, m map[string][]*PDBBindEntry) error {
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("loading file: %v", err)
	}

	// https://regex101.com/r/nmqWQj/3
	r, err := regexp.Compile("(?m)^([a-z0-9]{4})[ ]+(.*?)[ ]+([0-9]{4})[ ]+(.*?)[ ]+.*?\\((.*?)\\)[ ]?(.*?)$")
	lines := r.FindAllStringSubmatch(string(raw), -1)

	for _, l := range lines {
		m[l[5]] = append(m[l[5]], &PDBBindEntry{
			PDBID:       l[1],
			Resolution:  l[2],
			Year:        l[3],
			BindingData: l[4],
			LigandID:    l[5],
			Desc:        l[6],
		})

	}

	return nil
}
