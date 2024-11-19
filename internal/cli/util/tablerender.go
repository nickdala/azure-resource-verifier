package util

import (
	"os"

	"github.com/olekukonko/tablewriter"
)

type Table struct {
	table *tablewriter.Table
}

type TableLayout string

const (
	MultipleServices TableLayout = "multiple_services"
)

func NewTable(layout TableLayout) *Table {
	t := tablewriter.NewWriter(os.Stdout)

	switch layout {
	case MultipleServices:
		multipleServiceLayout(t)
	}

	return &Table{table: t}
}

func multipleServiceLayout(t *tablewriter.Table) {
	t.SetHeader([]string{"Service", "Location", "Enabled"})
	t.SetAutoMergeCellsByColumnIndex([]int{0})
}

func (t *Table) SetHeader(header []string) {
	t.table.SetHeader(header)
}

func (t *Table) AppendRow(row []string) {
	t.table.Append(row)
}

func (t *Table) AppendBulk(rows [][]string) {
	t.table.AppendBulk(rows)
}

func (t *Table) Render() {
	t.table.Render()
}
