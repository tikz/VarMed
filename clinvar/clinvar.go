package clinvar

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"respdb/pdb"
	"strconv"
	"strings"
)

const (
	clinvarDB   = "https://ftp.ncbi.nlm.nih.gov/pub/clinvar/tab_delimited/variant_summary.txt.gz"
	summaryPath = "data/clinvar/variant_summary.txt"
)

// DbSNP holds all relevant entries from dbSNP
type DbSNP struct {
	snps map[string][]Allele
}

// Allele represents a variation from ClinVar
type Allele struct {
	VariationID   string `json:"variationId"`
	Name          string `json:"name"`
	ClinSig       string `json:"clinSig"`
	ClinSigSimple int    `json:"clinSigSimple"`
	ProteinChange string `json:"proteinChange"`
	ReviewStatus  string `json:"reviewStatus"`
	Phenotypes    string `json:"phenotypes"`
}

func NewDbSNP() *DbSNP {
	GetClinVar()

	var db DbSNP
	db.snps = make(map[string][]Allele)
	db.Load()

	return &db
}

// Load parses and loads all alleles from the summary file.
func (d *DbSNP) Load() {
	fmt.Println("loading ClinVar variants...")

	f, err := os.Open(summaryPath)
	if err != nil {
		panic(err)
	}

	s := bufio.NewScanner(f)
	alleles := 0
	for s.Scan() {
		line := strings.Split(s.Text(), "\t")
		variantType := line[1]
		name := line[2]
		synonymous := strings.Index(name, "=") != -1
		r, _ := regexp.Compile(`\(p.([A-z]{3})([0-9]*)([A-z]{3})\)`)
		m := r.FindAllStringSubmatch(name, -1)
		coding := len(m) > 0
		assembly := line[16] == "GRCh38"
		if variantType == "single nucleotide variant" && coding && !synonymous && assembly {
			dbSNPID := line[9]

			_, _, fromAa := pdb.AminoacidNames(m[0][1])
			_, _, toAa := pdb.AminoacidNames(m[0][3])
			pos := m[0][2]
			change := fromAa + pos + toAa
			clinSigSimple, _ := strconv.Atoi(line[7])
			allele := Allele{
				VariationID:   line[30],
				Name:          name,
				ClinSig:       line[6],
				ClinSigSimple: clinSigSimple,
				ProteinChange: change,
				ReviewStatus:  line[24],
				Phenotypes:    line[13],
			}
			d.snps["rs"+dbSNPID] = append(d.snps["rs"+dbSNPID], allele)
			alleles++
		}
	}

	fmt.Println(fmt.Sprintf("%d SNPs, %d alleles loaded.", len(d.snps), alleles))
}

func (d *DbSNP) GetVariation(dbSNPID string, proteinChange string) *Allele {
	if alleles, ok := d.snps[dbSNPID]; ok {
		for _, allele := range alleles {
			if allele.ProteinChange == proteinChange {
				return &allele
			}
		}
	}
	return nil
}

// GetClinVar downloads and decompresses the variant_summary.txt.gz file from the NCBI FTP.
func GetClinVar() error {
	_, err := os.Stat(summaryPath)
	if os.IsNotExist(err) {
		resp, err := http.Get(clinvarDB)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		out, err := os.Create(summaryPath)
		if err != nil {
			return err
		}
		defer out.Close()

		gr, err := gzip.NewReader(resp.Body)
		defer gr.Close()

		_, err = io.Copy(out, gr)
		return err
	}

	return nil
}
