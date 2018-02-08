package odf

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

var (
	odfFiles, dataFile, outFilePrefix string
)

func init() {
	flag.StringVar(&odfFiles, "odf", "testdata/report.*", "The `glob pattern` to run the tests against")
	flag.StringVar(&dataFile, "data", "testdata/data.json", "The `json file` to use as input")
	flag.StringVar(&outFilePrefix, "out", "testdata/out", "The `destination prefix` to use as output. WARNING: will be overwritten.")
}

func TestRenderODF(t *testing.T) {
	m, err := filepath.Glob(odfFiles)
	if err != nil {
		t.Fatalf("Unable to parse input glob: %v", err)
	}
	for _, odfFile := range m {
		ext := filepath.Ext(odfFile)
		outFile := outFilePrefix + ext
		os.Remove(outFile)
		t.Logf("Testing template source %v (saving as %s)", odfFile, outFile)
		odf, err := Open(odfFile)
		if err != nil {
			t.Fatalf("Error opening file: %v", err)
		}

		data := loadTestData(t)

		t.Logf("Executing template")
		if err = odf.Execute(data); err != nil {
			t.Fatalf("Error executing odf template: %v", err)
		}

		t.Logf("Saving file")
		if err = odf.WriteFile(outFile); err != nil {
			t.Fatalf("Error saving file: %v", err)
		}

		t.Logf("Zip file saved. Checking a PDF convertion")

		// Check a PDF conversion to prove the output is correctly rendered
		if err = libreofficeConvert(t, outFile, "pdf", "./testdata/"); err != nil {
			t.Fatalf("Error reading resulting template: %v", err)
		}
	}
}

func loadTestData(t *testing.T) map[string]interface{} {
	// Execute with data
	var data map[string]interface{}
	data, err := ReadJSON(dataFile)
	if err != nil {
		t.Fatalf("Failed to load data file %v", err)
	}
	items := make([]interface{}, 0)
	gt := 0.0

	randomSuffix := func(i int) string {
		switch i % 5 {
		case 1:
			return "- Small Piece"
		case 2:
			return "- Optimized for Reading"
		default:
			return ""
		}
	}

	for i := 1; i <= 100; i++ {
		v := float64((i * 3) % 4)
		q := (i * 2) % 10
		item := map[string]interface{}{
			"ItemNo": i,
			"Name":   fmt.Sprintf("Random item #%d %s", i, randomSuffix(i)),
			"Quant":  q,
			"Price":  fmt.Sprintf("$ %.02f", v),
			"Total":  fmt.Sprintf("$ %.02f", v*float64(q)),
		}
		items = append(items, item)
		gt += v
	}
	// Replaces data from JSON with the generated values
	data["Items"] = items
	data["GrandTotal"] = fmt.Sprintf("$ %.2f", gt)

	return data
}

func libreofficeConvert(t *testing.T, src, format, outdir string) error {
	t.Logf("Converting: src=%v, format=%v, outdir=%v", src, format, outdir)

	libreoffice := exec.Command("libreoffice",
		"-env:UserInstallation=file:///tmp/odftest",
		"--headless",
		"--convert-to", format,
		"--outdir", outdir,
		src,
	)
	t.Logf("Running %v", libreoffice.Args)
	out, err := libreoffice.CombinedOutput()
	if err != nil {
		return err
	}

	t.Logf("Result: %v", string(out))
	return nil
}
