package cdtacctwebscraper

import (
	"fmt"

	"github.com/tealeg/xlsx"
)

func (data *AccountScrapeData) MakeHistoryReport() (err error) {
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var cell *xlsx.Cell
	var row *xlsx.Row
	var filename string
	var ac AccountConfig
	var ok bool

	filename = "HistoryReport.xlsx"
	file = xlsx.NewFile()

	fmt.Println("Writing History Report to XLSX")
	for id, pd := range data.Save.HistoryMap {
		ac, ok = data.AccountMap[id]
		if !ok {
			err = fmt.Errorf("Could not find account for Id %s", id)
			return
		}

		sheet, err = file.AddSheet(ac.CompactName)
		if err != nil {
			err = fmt.Errorf("Error Adding Sheet: %v", err)
			return
		}

		if len(pd.HistoryTable) == 0 {
			fmt.Printf("Warning: No data for %s Id %s\n", ac.CompactName, id)
			row = sheet.AddRow()
			cell = row.AddCell()
			cell.SetString("No data")
			continue
		}

		// History
		for i, hrow := range pd.HistoryTable {
			row = sheet.AddRow()

			for j, field := range hrow.Extracted {
				cell = row.AddCell()
				cell.SetString(field)
				if i == 0 {
					sheet.Col(j).Width = 15.0
					cell.SetStyle(getStyle("th"))
				}
			}

			if i == 0 {
				cell = row.AddCell()
				sheet.Col(len(hrow.Extracted)).Width = 15.0
				cell.SetString("Parsed Amount")
				cell.SetStyle(getStyle("th2"))
				cell = row.AddCell()
				sheet.Col(len(hrow.Extracted) + 1).Width = 15.0
				cell.SetString("Parsed Date")
				cell.SetStyle(getStyle("th2"))
			} else {
				cell = row.AddCell()
				cell.SetString("Parsed Amount")
				cell.SetFloat(hrow.Amount)
				cell.NumFmt = styleFmt_ACCT
				cell = row.AddCell()
				cell.SetString(hrow.Date)
			}
		}
	}

	err = file.Save(filename)
	if err != nil {
		return fmt.Errorf("Error writing XLSX file: %v", err)
	}
	return
}
