package tango

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type TangoResidue struct {
	// http://tango.crg.es/Tango_Handbook.pdf
	// Those files will have the following columns:
	Beta             float64 // Percentage of beta-strand conformation
	BetaTurn         float64 // Percentage of beta-turn conformation
	Helix            float64 // Percentage of alpha-helical conformation
	Aggregation      float64 // Percentage of Aggregation
	HelixAggregation float64 // Percentage of Helical Aggregation
}

func Run(name string, seq string) ([]*TangoResidue, error) {
	// http://tango.crg.es/Tango_Handbook.pdf
	// The format of the sequences to be run is as follows:
	// Name Cter Nter pH Temp Ionic Sequence
	// Name = name of the sequence (less than 25 characters)
	// Cter = status of the C-terminus of the peptide (amidated Y, free N)
	// Nter = status of the N-terminus of the peptide (acetylated A, succinilated S and free N)
	// pH = pH
	// Temp = temperature in Kelvin
	// Ionic = ionic strength in M
	// sequence = sequence of the peptide in one letter code.
	cmd := exec.Command("./tango", name, "ct=N", "nt=N", "ph=7.4", "te=303", "io=0.05", "seq="+seq)
	cmd.Dir = "bin/"
	out, err := cmd.CombinedOutput()
	strOut := string(out)
	if err != nil || strings.Contains(strings.ToLower(strOut), "error") {
		fmt.Println(strOut)
		return nil, errors.New(strOut)
	}

	outPath := "bin/" + name + ".txt"
	defer func() {
		os.RemoveAll(outPath)
	}()

	raw, err := ioutil.ReadFile(outPath)
	if err != nil {
		return nil, err
	}

	var results []*TangoResidue

	for _, line := range strings.Split(string(raw), "\n")[1:] {
		fields := strings.Fields(line)
		if len(fields) >= 7 {
			toFloat := func(str string) float64 {
				f, _ := strconv.ParseFloat(str, 64)
				return f
			}
			results = append(results, &TangoResidue{
				Beta:             toFloat(fields[2]),
				BetaTurn:         toFloat(fields[3]),
				Helix:            toFloat(fields[4]),
				Aggregation:      toFloat(fields[5]),
				HelixAggregation: toFloat(fields[6]),
			})
		}
	}

	return results, nil
}
