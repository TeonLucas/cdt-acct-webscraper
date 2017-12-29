package cdtacctwebscraper

import (
	"fmt"
	"strings"
)

func (data *AccountScrapeData) ParseHistory() (err error) {
	var ac AccountConfig
	var ok bool

	fmt.Println("Begin History Parsing")

	for id, pd := range data.Save.HistoryMap {
		ac, ok = data.AccountMap[id]
		if !ok {
			err = fmt.Errorf("Could not find account for Id %s", id)
			return
		}

		// Identify rows to skip
		pd.IdentifyRowsToSkip()

		fmt.Printf("• Parsing %q history: %d rows\n", ac.AccountName, pd.HistoryToMatch)

		// Parse out dates and amounts from History table
		err = ac.ParseFormats(&pd)
		if err != nil {
			fmt.Printf("• Skipping %s: %v\n", ac.AccountName, err)
			err = nil
			continue
		}

		//pd.Dump()

		// Save updates in map
		data.Save.HistoryMap[id] = pd
	}

	return
}

func (pd *HistoryData) IdentifyRowsToSkip() {
	var col string

	for i, historyRow := range pd.HistoryTable {
		// Skip header row
		if i == 0 {
			pd.HistoryTable[i].Skip = true
			continue
		}

		// Skip any total rows
		for _, col = range historyRow.Extracted {
			if strings.Contains(strings.ToLower(col), "total") {
				pd.HistoryTable[i].Skip = true
				continue
			}
		}

		// Skip short rows
		if len(historyRow.Extracted) < pd.EmptyLen {
			pd.HistoryTable[i].Skip = true
			continue
		}

		// Count of rows not skipped
		pd.HistoryToMatch++
	}

	return
}

// For debugging
func (pd *HistoryData) Dump() {
	for i, historyRow := range pd.HistoryTable {
		fmt.Printf("# HistoryRow[%d]: Skip=%t\n", i, historyRow.Skip)
		fmt.Printf("    %q\n", historyRow.Extracted)
		fmt.Printf("    Amount=%.2f Date=%s\n", historyRow.Amount, historyRow.Date)
	}
}
