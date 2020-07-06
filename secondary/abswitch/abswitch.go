package abswitch

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os/exec"
	"strconv"
	"strings"
)

type Residue struct {
	Gor string  `json:"gor"`
	PaH float64 `json:"pah"`
	PaE float64 `json:"pae"`
	PaC float64 `json:"pac"`
	Amb float64 `json:"amb"`
	Ins float64 `json:"ins"`
	Swi float64 `json:"swi"`
	S5s float64 `json:"s5s"`
}

func Run(name string, seq string) ([]*Residue, error) {
	// Write fasta
	fastaFile := "abswitch_" + name + ".fasta"
	ioutil.WriteFile("bin/"+fastaFile, []byte(">"+name+"\n"+seq), 0644)

	// Write cfg
	cfgFile := "abswitch_" + name + ".cfg"
	outFile := name + ".s5"
	cfg := fmt.Sprintf("command=Switch5\nfasta=%s\noFile=%s", fastaFile, outFile)
	ioutil.WriteFile("bin/"+cfgFile, []byte(cfg), 0644)

	// defer func() {
	// 	os.RemoveAll("bin/" + fastaFile)
	// 	os.RemoveAll("bin/" + cfgFile)
	// 	os.RemoveAll("bin/" + outFile)
	// }()

	// Run
	cmd := exec.Command("./abSwitch", "-f", cfgFile)
	cmd.Dir = "bin/"
	out, err := cmd.CombinedOutput()
	strOut := string(out)
	if err != nil || !strings.Contains(strings.ToLower(strOut), "printed results") {
		fmt.Println(strOut)
		return nil, errors.New(strOut)
	}

	raw, err := ioutil.ReadFile("bin/" + outFile)
	if err != nil {
		return nil, err
	}

	var results []*Residue

	for _, line := range strings.Split(string(raw), "\n")[1:] {
		fields := strings.Fields(line)
		if len(fields) >= 9 {
			toFloat := func(str string) float64 {
				f, _ := strconv.ParseFloat(str, 64)
				return f
			}
			results = append(results, &Residue{
				Gor: fields[1],
				PaH: toFloat(fields[2]),
				PaE: toFloat(fields[3]),
				PaC: toFloat(fields[4]),
				Amb: toFloat(fields[5]),
				Ins: toFloat(fields[6]),
				Swi: toFloat(fields[7]),
				S5s: toFloat(fields[8]),
			})
		}
	}

	return results, nil
}
