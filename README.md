# report/odf

Minimalistic Open Document Format (ODF) report generation tool for Go.

## About

This tool aims to be a minimalistic aproach to render ODF files
using the `text/template` package from stdlib.

Using this library, one can create a Open Document or Open Graphic,
mark the file with the standard Go templating utilities and
render the resulting document to a file.

The file can later be converted by any format that LibreOffice or
Unoconv exporting filters allow.

# Quick start

Import the package

    import "github.com/ronoaldo/report/odf"

Open a template, call Execute() and then, WriteFile:

    doc, _ := odf.Open("template.odt")
    doc.Execute(values)
    doc.WriteFile("output.odt")

# Tips

Being minimalistic, we assume the input file is not outrageously complex.
Hence, the bare minimum changes to the XML structure is done.

## Errors while calling .Execute()

If you call odf.Execute() and you see an error about invalid characters
on content.xml or styles.xml, make sure try the following:

*Remove extra formatting using CTRL+M* while selecting the template
instructions between `{{ }}`.

*Make sure that, between `{{ }}`, all characters as ASCII*. You may need
to disable the convertion of quotes in the local auto-correction preferences.