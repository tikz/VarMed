package pdb

import (
	"io/ioutil"
	"testing"
)

func LoadTestFile(path string) ([]byte, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func TestChains(t *testing.T) {
	raw, err := LoadTestFile("./testdata/3con.pdb")
	if err != nil {
		t.Errorf("cannot open file: %s", err)
	}

	pdb, err := NewPDBFromRaw(raw)
	if err != nil {
		t.Error(err)
	}

	t.Logf("Testing PDB chains")

	actual := pdb.TotalLength
	expected := int64(156)
	if actual != expected {
		t.Errorf("Expected %d, got %d", expected, actual)
	}

	res := pdb.Chains["A"][160]
	if res.Name != "Valine" {
		t.Errorf("Expected Valine in A-160, got %s", res.Name)
	}

}

func TestSeqRes(t *testing.T) {
	raw, err := LoadTestFile("./testdata/3con.pdb")
	if err != nil {
		t.Errorf("cannot open file: %s", err)
	}

	pdb, err := NewPDBFromRaw(raw)
	if err != nil {
		t.Error(err)
	}

	err = pdb.ExtractSeqRes()
	if err != nil {
		t.Error(err)
	}

	t.Logf("Testing SEQ RES")

	res := pdb.SeqRes["A"][0]
	expected := "Methionine"
	if res.Name != expected {
		t.Errorf("Expected %s in SEQRES A-1, got %s", expected, res.Name)
	}

	res = pdb.SeqRes["A"][189]
	expected = "Asparagine"
	if res.Name != expected {
		t.Errorf("Expected %s in SEQRES A-190, got %s", expected, res.Name)
	}

}

func TestCIF(t *testing.T) {
	raw, err := LoadTestFile("./testdata/3con.pdb")
	if err != nil {
		t.Errorf("cannot open file: %s", err)
	}

	pdb, err := NewPDBFromRaw(raw)
	if err != nil {
		t.Error(err)
	}

	rawCIF, err := LoadTestFile("./testdata/3con.cif")
	if err != nil {
		t.Errorf("cannot open file: %s", err)
	}

	t.Logf("Testing CIF parse")
	pdb.RawCIF = rawCIF

	err = pdb.ExtractCIFData()
	if err != nil {
		t.Errorf("cannot extract CIF: %s", err)
	}

	res := pdb.Title
	expected := "Crystal structure of the human NRAS GTPase bound with GDP"
	if res != expected {
		t.Errorf("Expected %s, got %s", expected, res)
	}

	res = pdb.Method
	expected = "X-RAY DIFFRACTION"
	if res != expected {
		t.Errorf("Expected %s, got %s", expected, res)
	}

	if pdb.Date.Day() != 28 || pdb.Date.Month() != 3 || pdb.Date.Year() != 2008 {
		t.Errorf("Expected date to be 2008-03-28")
	}

	if pdb.Resolution != 1.649 {
		t.Errorf("Expected %f, got %f", 1.649, pdb.Resolution)
	}
}
