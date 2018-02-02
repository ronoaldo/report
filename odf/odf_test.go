package odf

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

var (
	odfFile, dataFile, outFile string
)

func init() {
	flag.StringVar(&odfFile, "odf", "testdata/report.odt", "`Source odf file` to run the test against")
	flag.StringVar(&dataFile, "data", "testdata/data.json", "`Data file` to use as input")
	flag.StringVar(&outFile, "out", "testdata/out.odt", "`Destination file` to use as output. WARNING: will be overwritten")
}

func TestExecuteTemplate(t *testing.T) {
	odf, err := Open(odfFile)
	if err != nil {
		t.Fatalf("Error opening file: %v", err)
	}

	// Check if buffers where read
	t.Logf("DEBUG: odf.cache: %#v", odf.listFiles())
	for _, f := range []string{"content.xml", "styles.xml"} {
		fd := odf.fileBuffer(f)
		if fd.Len() == 0 {
			t.Errorf("Unexpected nil zipfile inside ODF struct: %v", f)
		}
		// Check if prepare will work as expected
		rawXML := fd.String()
		cleanXML, err := prepareXMLForTemplate(rawXML)
		if err != nil {
			t.Errorf("Unable to prepare XML for template: %v", err)
		}
		if err := ioutil.WriteFile("/tmp/"+f, []byte(cleanXML), 0644); err != nil {
			t.Errorf("Error generating temp file %v", err)
		}
	}

	// Execute with data
	var data map[string]interface{}
	if data, err = loadJSONFile(dataFile); err != nil {
		t.Fatalf("Failed to load data file %v", err)
	}
	items := make([]interface{}, 0)
	gt := 0.0
	for i := 0; i < 100; i++ {
		v := float64(i * 3 % 4)
		q := i * 2 % 10
		item := map[string]interface{}{
			"ItemNo": i,
			"Name":   fmt.Sprintf("Random item #%d", i),
			"Quant":  q,
			"Value":  fmt.Sprintf("$ %.02f", v),
			"Total":  fmt.Sprintf("$ %.02f", v*float64(q)),
		}
		items = append(items, item)
		gt += v
	}
	data["Items"] = items
	data["GrandTotal"] = fmt.Sprintf("$ %.2f", gt)

	if err = odf.Execute(data); err != nil {
		t.Fatalf("Error executing odf template: %v", err)
	}

	if err = odf.WriteFile(outFile); err != nil {
		t.Fatalf("Error saving file")
	}
}

func loadJSONFile(fname string) (map[string]interface{}, error) {
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

func xTestMergeTags(t *testing.T) {
	input := `<?xml version="1.0" encoding="UTF-8"?>
<office:document-content xmlns:office="urn:oasis:names:tc:opendocument:xmlns:office:1.0" xmlns:style="urn:oasis:names:tc:opendocument:xmlns:style:1.0" xmlns:text="urn:oasis:names:tc:opendocument:xmlns:text:1.0" xmlns:table="urn:oasis:names:tc:opendocument:xmlns:table:1.0" xmlns:draw="urn:oasis:names:tc:opendocument:xmlns:drawing:1.0" xmlns:fo="urn:oasis:names:tc:opendocument:xmlns:xsl-fo-compatible:1.0" xmlns:xlink="http://www.w3.org/1999/xlink" xmlns:dc="http://purl.org/dc/elements/1.1/" xmlns:meta="urn:oasis:names:tc:opendocument:xmlns:meta:1.0" xmlns:number="urn:oasis:names:tc:opendocument:xmlns:datastyle:1.0" xmlns:presentation="urn:oasis:names:tc:opendocument:xmlns:presentation:1.0" xmlns:svg="urn:oasis:names:tc:opendocument:xmlns:svg-compatible:1.0" xmlns:chart="urn:oasis:names:tc:opendocument:xmlns:chart:1.0" xmlns:dr3d="urn:oasis:names:tc:opendocument:xmlns:dr3d:1.0" xmlns:math="http://www.w3.org/1998/Math/MathML" xmlns:form="urn:oasis:names:tc:opendocument:xmlns:form:1.0" xmlns:script="urn:oasis:names:tc:opendocument:xmlns:script:1.0" xmlns:ooo="http://openoffice.org/2004/office" xmlns:ooow="http://openoffice.org/2004/writer" xmlns:oooc="http://openoffice.org/2004/calc" xmlns:dom="http://www.w3.org/2001/xml-events" xmlns:xforms="http://www.w3.org/2002/xforms" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:smil="urn:oasis:names:tc:opendocument:xmlns:smil-compatible:1.0" xmlns:anim="urn:oasis:names:tc:opendocument:xmlns:animation:1.0" xmlns:rpt="http://openoffice.org/2005/report" xmlns:of="urn:oasis:names:tc:opendocument:xmlns:of:1.2" xmlns:xhtml="http://www.w3.org/1999/xhtml" xmlns:grddl="http://www.w3.org/2003/g/data-view#" xmlns:officeooo="http://openoffice.org/2009/office" xmlns:tableooo="http://openoffice.org/2009/table" xmlns:drawooo="http://openoffice.org/2010/draw" xmlns:calcext="urn:org:documentfoundation:names:experimental:calc:xmlns:calcext:1.0" xmlns:loext="urn:org:documentfoundation:names:experimental:office:xmlns:loext:1.0" xmlns:field="urn:openoffice:names:experimental:ooo-ms-interop:xmlns:field:1.0" xmlns:formx="urn:openoffice:names:experimental:ooxml-odf-interop:xmlns:form:1.0" xmlns:css3t="http://www.w3.org/TR/css3-text/" office:version="1.2">
<draw:frame draw:style-name="gr25" draw:text-style-name="P4" draw:layer="layout" svg:width="6.852cm" svg:height="0.31cm" svg:x="6.102cm" svg:y="9.568cm">
  <draw:text-box>
	<text:p>
	  <text:span text:style-name="T7">{{or .</text:span>
	  <text:span text:style-name="T7">Mar }</text:span>
	  <text:span text:style-name="T7">}</text:span>
	</text:p>
  </draw:text-box>
</draw:frame>`
	t.Logf("input=%s", input)

	output, err := prepareXMLForTemplate(input)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	t.Logf("*** output= %v", output)
}

func xTestParseAnnotations(t *testing.T) {
	input := ` <office:body>
	<office:text>
	 <text:sequence-decls>
	  <text:sequence-decl text:display-outline-level="0" text:name="Illustration"/>
	  <text:sequence-decl text:display-outline-level="0" text:name="Table"/>
	  <text:sequence-decl text:display-outline-level="0" text:name="Text"/>
	  <text:sequence-decl text:display-outline-level="0" text:name="Drawing"/>
	 </text:sequence-decls>
	 <table:table table:name="Tabela1" table:style-name="Tabela1">
	  <table:table-column table:style-name="Tabela1.A"/>
	  <table:table-column table:style-name="Tabela1.B"/>
	  <table:table-column table:style-name="Tabela1.C"/>
	  <table:table-column table:style-name="Tabela1.D"/>
	  <table:table-column table:style-name="Tabela1.E"/>
	  <table:table-row table:style-name="Tabela1.1">
	   <table:table-cell table:style-name="Tabela1.A1" office:value-type="string">
		<text:p text:style-name="P7"><office:annotation>
		  <dc:creator>Ronoaldo JLP</dc:creator>
		  <dc:date>2018-02-01T20:57:59.666690844</dc:date>
		  <text:p text:style-name="P9"><text:span text:style-name="T7">{{range .Items}}</text:span></text:p>
		 </office:annotation>{{.ItemNo}}</text:p>
	   </table:table-cell>
	   <table:table-cell table:style-name="Tabela1.A1" office:value-type="string">
		<text:p text:style-name="Table_20_Contents">{{.Name}}</text:p>
	   </table:table-cell>
	   <table:table-cell table:style-name="Tabela1.A1" office:value-type="string">
		<text:p text:style-name="Table_20_Contents">{{.Quant}}</text:p>
	   </table:table-cell>
	   <table:table-cell table:style-name="Tabela1.A1" office:value-type="string">
		<text:p text:style-name="Table_20_Contents">{{.Value}}</text:p>
	   </table:table-cell>
	   <table:table-cell table:style-name="Tabela1.E1" office:value-type="string">
		<text:p text:style-name="P8">{{.Total}}<office:annotation>
		  <dc:creator>Ronoaldo JLP</dc:creator>
		  <dc:date>2018-02-01T20:58:14.379188531</dc:date>
		  <text:p text:style-name="P9"><text:span text:style-name="T7">{{end}}</text:span></text:p>
		 </office:annotation></text:p>
	   </table:table-cell>
	  </table:table-row>
	 </table:table>
	 <text:p text:style-name="P6"/>
	</office:text>
   </office:body>`

	output, err := prepareXMLForTemplate(input)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("output=%v", output)
}
