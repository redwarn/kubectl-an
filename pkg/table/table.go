package table

import (
	"github.com/jedib0t/go-pretty/v6/table"
)

var title = table.Row{
	"NAMESPACE",
	"TYPE",
	"NAME",
	"ANNOTATION_KEY",
	"ANNOTATION_VALUE",
}

type Resource struct {
	Namespace   string
	Type        string
	Name        string
	Annotations map[string]string
}

func GenTable(rs []Resource) *table.Table {
	t := &table.Table{}
	tr := make([]table.Row, 0)
	t.AppendHeader(title)
	for i := range rs {
		tr = append(tr, table.Row{rs[i].Namespace, rs[i].Type, rs[i].Name, "----------------------------------------", "----------------------------------------"})
		for k, v := range rs[i].Annotations {
			tr = append(tr, table.Row{"", "", "", k, v})
		}
	}
	t.AppendRows(tr)
	return t
}
