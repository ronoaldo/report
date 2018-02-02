/*
Package odf implements a minimalistic templating engine for Open Document files.

The pourpose of this package is to allow the rendering of Open Document
files using the Go standard library text/template package.

Since Open Document are packaged XML files, one can make edit an Open Document
file and add pipelines to it in order to render an output interpolated with data.
This package essentially provides the minimal tooling to unpack, merge and render
the resulting file.

To edit templates, you can use any Office Suite that is capable to output
ODF files (.odt, odg and .ods).

Only simple instructions are supported, such as rendering a single value.
Minimal support for {{range}} ... {{end}} inside annotations delimiting a table-row
allow for a ODT file to be rendered and output a table with several items.
See the testdata directory for samples.
*/
package odf
