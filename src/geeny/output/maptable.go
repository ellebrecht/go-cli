package output

import (
	"bytes"
	"fmt"

	"geeny/util"
	log "geeny/log"
)

type MapTable struct {
	tableTitle  []string
	tableField  []string
	value		[]interface{}
}

func NewMapTable(tableTitle []string, tableField []string, value []interface{}) *MapTable {
	return &MapTable{tableTitle, tableField, value}
}

func (t *MapTable) String() (*string, error) {
	tableWidth := []int{}
	var buf bytes.Buffer

	// dynamically determine column width based on cell contents
	for f, field := range t.tableField {
		tableWidth = append(tableWidth, len(t.tableTitle[f]))
		for _, value := range t.value {
			v, ok := value.(map[string]interface{})
			if ok {
				tableWidth[f] = util.Max(len(fmt.Sprintf("%v", v[field])), tableWidth[f])
			} else {
				log.Warnf("Table row is not a map: %v", value)
			}
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
	for _, value := range t.value {
		util.AppendLine(&buf, tableWidth, 1)
		for i, field := range t.tableField {
			buf.WriteString("\x1b[39;1m│\x1b[0m \x1b[34;1m")
			v, ok := value.(map[string]interface{})
			if ok {
				buf.WriteString(util.Pad(fmt.Sprintf("%v", v[field]), tableWidth[i]))
			}
			buf.WriteString("\x1b[0m ")
		}
		buf.WriteString("\x1b[39;1m│\x1b[0m\n")
	}
	util.AppendLine(&buf, tableWidth, 2)

	s := buf.String()
	return &s, nil
}
