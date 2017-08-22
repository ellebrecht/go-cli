package output

import (
	"bytes"
	"geeny/util"
)

type JsonTable struct {
	tableTitle  []string
	tableField  []func(interface{}) string
	valuesSlice interface{}
}

func NewTable(tableTitle []string, tableField []func(interface{}) string, valuesSlice interface{}) *JsonTable {
	return &JsonTable{tableTitle, tableField, valuesSlice}
}

func (t *JsonTable) String() (*string, error) {
	tableWidth := []int{}
	var buf bytes.Buffer
	//TODO: using reflection here to achieve a generic array argument is horrible - surely there's a better way?
	values := util.TakeSliceArg(t.valuesSlice)

	// dynamically determine column width based on cell contents
	for f, field := range t.tableField {
		tableWidth = append(tableWidth, len(t.tableTitle[f]))
		for _, value := range values {
			tableWidth[f] = util.Max(len(field(value)), tableWidth[f])
		}
	}

	// append header
	util.AppendLine(&buf, tableWidth, 0)
	for i, title := range t.tableTitle {
		buf.WriteString("\x1b[39;1m│\x1b[0m \x1b[31m")
		buf.WriteString(util.Pad(title, tableWidth[i]))
		buf.WriteString("\x1b[0m ")
	}
	buf.WriteString("\x1b[39;1m│\x1b[0m\n")

	// append rows
	for _, value := range values {
		util.AppendLine(&buf, tableWidth, 1)
		for i, field := range t.tableField {
			buf.WriteString("\x1b[39;1m│\x1b[0m \x1b[34;1m")
			buf.WriteString(util.Pad(field(value), tableWidth[i]))
			buf.WriteString("\x1b[0m ")
		}
		buf.WriteString("\x1b[39;1m│\x1b[0m\n")
	}
	util.AppendLine(&buf, tableWidth, 2)

	s := buf.String()
	return &s, nil
}
