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
	SingleService    TableLayout = "single_service"
	MultipleServices TableLayout = "multiple_services"
)

func NewTable(layout TableLayout) *Table {
	t := tablewriter.NewWriter(os.Stdout)

	switch layout {
	case SingleService:
		singleServiceLayout(t)
	case MultipleServices:
		multipleServiceLayout(t)
	}

	return &Table{table: t}
}

func singleServiceLayout(t *tablewriter.Table) {
	t.SetHeader([]string{"Location", "Enabled", "HA Enabled", "Reason"})
}

func multipleServiceLayout(t *tablewriter.Table) {
	t.SetHeader([]string{"Service", "Location", "Enabled", "HA Enabled", "Reason"})
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
