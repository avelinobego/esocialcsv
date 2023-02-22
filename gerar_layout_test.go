package layout

import (
	"bufio"
	"fmt"
	"os"
	"testing"

	"github.com/gocolly/colly"
)

func TestLayout(t *testing.T) {
	c := colly.NewCollector()
	c.OnHTML("html", makeFunc(t, c))
	c.Visit("https://www.gov.br/esocial/pt-br/documentacao-tecnica/leiautes-esocial-v-s1.1-nt-01-2023/index.html")
}

func makeFunc(t *testing.T, c *colly.Collector) func(e *colly.HTMLElement) {
	return func(e *colly.HTMLElement) {

		eventos_id := e.ChildAttrs("h3", "id")

		file_desc, err := os.OpenFile("files/descritivo.csv", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
		if err != nil {
			t.Fatal(err)
		}

		defer func() {
			if err := file_desc.Close(); err != nil {
				t.Fatal(err)
			}
		}()

		desc := bufio.NewWriter(file_desc)
		desc.WriteString("Descricao;Arquivo\n")

		for _, evt := range eventos_id {

			descricao := e.ChildText(fmt.Sprintf(`h3[id="%s"]`, evt))

			file, err := os.OpenFile(fmt.Sprintf("files/%s.csv", evt), os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
			if err != nil {
				t.Fatal(err)
			}

			defer func() {
				if err := file.Close(); err != nil {
					t.Fatal(err)
				}
			}()

			columns := bufio.NewWriter(file)

			e.ForEachWithBreak(fmt.Sprintf(`h3[id=%s]~h4~table`, evt), func(index_table int, table *colly.HTMLElement) bool {

				desc.WriteString(fmt.Sprintf(`"%s";%s.csv`, descricao, evt))
				desc.WriteString("\n")
				desc.Flush()

				table.ForEach("tr", func(i int, row *colly.HTMLElement) {
					ExcreverColunas(row, columns)
				})

				return false
			})
		}

	}
}

func ExcreverColunas(row *colly.HTMLElement, columns *bufio.Writer) {
	row.ForEach("th,td", func(i int, col *colly.HTMLElement) {
		if i > 0 {
			columns.WriteString(";")
		}
		columns.WriteString(fmt.Sprintf(`"%s"`, col.Text))
	})

	columns.WriteString("\n")
	columns.Flush()

}
