package odf

import "testing"

func TestOpen(t *testing.T) {
	odf, err := Open("testdata/report.odt")
	if err != nil {
		t.Fatalf("Error opening file: %v", err)
	}

	// Check if buffers where read
	t.Logf("DEBUG: odf.cache: %#v", odf.listFiles())
	for _, f := range []string{"content.xml", "styles.xml"} {
		if odf.fileBuffer(f).Len() == 0 {
			t.Errorf("Unexpected nil zipfile inside ODF struct: %v", f)
		}
	}
}
