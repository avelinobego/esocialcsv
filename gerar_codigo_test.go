package layout

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"text/template"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func TestCriarCodigo(t *testing.T) {

	files, err := filepath.Glob("files/evt*.csv")
	if err != nil {
		t.Fatal(err)
	}

	tpl := template.New("struct.tpl")

	tpl.Funcs(template.FuncMap{
		"Title": Title,
	})

	parser, err := tpl.ParseFiles("struct.tpl")
	if err != nil {
		t.Fatal(err)
	}

	var base_name string
	var clean_name string

	for _, n := range files {

		base_name = filepath.Base(n)
		clean_name = strings.TrimSuffix(base_name, ".csv")

		d := CriarData(n, clean_name)

		iow := bytes.NewBufferString("")
		if err := parser.Execute(iow, d); err != nil {
			t.Fatal(err)
		}

		os.WriteFile(fmt.Sprintf("eventos/%s.go", clean_name), iow.Bytes(), fs.ModePerm)

		t.Log(clean_name)
	}

}

func CriarData(file_name string, clean_name string) (result *Data) {
	result = new(Data)

	result.Name = clean_name

	f, e := os.Open(file_name)
	if e != nil {
		log.Fatal(e)
	}

	defer func() {
		f.Close()
	}()

	read := csv.NewReader(f)
	read.Comma = ';'

	var record []string
	_, eof := read.Read()
	for eof != io.EOF {
		record, eof = read.Read()
		for value := range record {

			log.Print(record[value])
		}
	}
	return
}

func Title(s string) string {
	c := cases.Title(language.BrazilianPortuguese, cases.NoLower)
	return c.String(s)
}

type Data struct {
	Name   string
	Fields []string
}
