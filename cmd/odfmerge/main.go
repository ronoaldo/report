package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"

	"github.com/ronoaldo/report/odf"
)

var (
	dataFile, templateFile, outFile string
)

func init() {
	flag.StringVar(&dataFile, "i", "", "The `json data` file to read from.")
	flag.StringVar(&templateFile, "t", "", "The `template` file to be merged.")
	flag.StringVar(&outFile, "o", "", "The `output` file to be written.")
}

func main() {
	flag.Parse()

	tpl, err := odf.Open(templateFile)
	if err != nil {
		log.Fatalf("odfmerge: error opening template %v", err)
	}
	data, err := loadJSONData(dataFile)
	if err != nil {
		log.Fatalf("odfmerge: error loading data file: %v", err)
	}

	err = tpl.Execute(data)
	if err != nil {
		log.Fatalf("odfmerge: error merging data with template: %v", err)
	}

	err = tpl.WriteFile(outFile)
	if err != nil {
		log.Fatalf("odfmerge: error writing output: %v", err)
	}

	log.Printf("odfmerge: output file written: %v", outFile)
}
func loadJSONData(fname string) (map[string]interface{}, error) {
	f, err := os.Open(fname)
	if err != nil {
		return nil, err
	}

	res := make(map[string]interface{})
	if err := json.NewDecoder(f).Decode(&res); err != nil {
		return nil, err
	}

	return res, nil
}
