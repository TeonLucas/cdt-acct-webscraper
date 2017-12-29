package cdtacctwebscraper

import (
	"fmt"
	"strconv"
	"strings"
)

// This function pads for an optional column
func (page *ScrapePage) PadTableColumn(padspec string) (err error) {
	var col, width int

	pspec := strings.ToLower(padspec)
	params := strings.Split(pspec, "of")
	if !strings.Contains(pspec, "of") || len(params) != 2 {
		err = fmt.Errorf("PadTableColumn expects '# of #' (column of width): %q provided", padspec)
		return
	}

	col, err = strconv.Atoi(strings.TrimSpace(params[0]))
	if err != nil {
		err = fmt.Errorf("PadTableColumn expects '# of #' (column of width): non integer %v", padspec, err)
		return
	}
	width, err = strconv.Atoi(strings.TrimSpace(params[1]))
	if err != nil {
		err = fmt.Errorf("PadTableColumn expects '# of #' (column of width): non integer %v", padspec, err)
		return
	}
	if col > width {
		err = fmt.Errorf("PadTableColumn expects '# of #' (column of width): invalid %d > %d", col, width)
		return
	}

	for i, row := range page.Table {
		if len(row) == width-1 {
			page.Table[i] = append(page.Table[i][:col-1], append([]string{""}, page.Table[i][col-1:]...)...)
		}
	}
	return
}

// This function splits one column into two
func (page *ScrapePage) SplitTableColumn(sep, heading, newHeading string) (err error) {
	var i, j, splitCol int
	var row []string
	var col string

	if len(page.Table) == 0 {
		err = fmt.Errorf("SplitTableColumn on empty table")
		return
	}

	splitColHeading := strings.ToLower(heading)
	for i, row = range page.Table {
		if i == 0 {
			// Locate column to split
			splitCol = -1
			for j, col = range row {
				heading := strings.ToLower(col)
				if strings.Contains(heading, splitColHeading) {
					splitCol = j
				}
			}
			if splitCol == -1 {
				err = fmt.Errorf("Failed to locate column for Account %s", page.AccountId)
				return
			}

			// Insert new heading
			page.Table[i] = append(page.Table[i], "")
			if splitCol > len(row)-1 {
				copy(page.Table[i][splitCol+2:], page.Table[i][splitCol+1:])
			}
			page.Table[i][splitCol+1] = newHeading
		} else {
			// Insert split column
			page.Table[i] = append(page.Table[i], "")
			if splitCol > len(row)-1 {
				copy(page.Table[i][splitCol+2:], page.Table[i][splitCol+1:])
			}
			newCols := strings.Split(row[splitCol], sep)
			page.Table[i][splitCol] = newCols[0]
			if len(newCols) > 1 {
				page.Table[i][splitCol+1] = newCols[1]
			}
		}
	}
	return
}

func (data *AccountScrapeData) StoreTable(id string, table [][]string) {
	fmt.Printf("• Storing extracted table in HistoryData[%q]\n", id)
	pd := HistoryData{}

	// Copy table rows into HistoryRows
	for _, row := range table {
		pd.HistoryTable = append(pd.HistoryTable, HistoryRow{Extracted: row})
	}
	pd.EmptyLen = len(pd.HistoryTable[0].Extracted)

	data.Save.HistoryMap[id] = pd
}
