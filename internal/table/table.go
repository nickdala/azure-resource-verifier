package table

import (
	"os"

	"github.com/olekukonko/tablewriter"
)

type Table struct {
	table *tablewriter.Table
}

type TableLayout string

const (
	Locations         TableLayout = "locations"
	PostgreSqlService TableLayout = "postgresql_service"
	RedisService      TableLayout = "redis_service"
	WebApp            TableLayout = "web_app"
	MultipleServices  TableLayout = "multiple_services"
)

func NewTable(layout TableLayout) *Table {
	t := tablewriter.NewWriter(os.Stdout)

	switch layout {
	case Locations:
		locationsLayout(t)
	case PostgreSqlService:
		singleServiceLayout(t)
	case WebApp:
		webAppLayout(t)
	case RedisService:
		redisLayout(t)
	case MultipleServices:
		multipleServiceLayout(t)
	}

	return &Table{table: t}
}

func redisLayout(t *tablewriter.Table) {
	t.SetHeader([]string{"Location", "Display Name", "Enabled"})
}

func webAppLayout(t *tablewriter.Table) {
	t.SetHeader([]string{"Location", "Display Name", "Enabled"})
}

func locationsLayout(t *tablewriter.Table) {
	t.SetHeader([]string{"Name", "Display Name"})
}

func singleServiceLayout(t *tablewriter.Table) {
	t.SetHeader([]string{"Location", "Display Name", "Enabled", "HA Enabled", "Reason"})
	t.SetAutoWrapText(true)
}

func multipleServiceLayout(t *tablewriter.Table) {
	t.SetHeader([]string{"Service", "Location", "Enabled", "HA Enabled", "Reason"})
	t.SetAutoMergeCellsByColumnIndex([]int{0})
	t.SetAutoWrapText(true)
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
