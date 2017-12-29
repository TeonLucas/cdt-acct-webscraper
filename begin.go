package cdtacctwebscraper

import (
	"log"
	"os"
	"strings"
)

var debugLog *log.Logger

func Begin() (data *AccountScrapeData, err error) {
	var fLog *os.File

	data = NewAccountScrapeData()
	err = data.GetCmdLine()
	if err != nil {
		return
	}

	// Log debug detail
	fLog, err = os.OpenFile(data.CmdLine.Logfile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return
	}
	debugLog = log.New(fLog, "CDB: ", log.LstdFlags)

	// Get list of Account Ids
	if len(data.CmdLine.SnapShotFile) > 0 {
		// Ids loaded directly from Snapshot
		err = data.LoadSnapshot()
		if err != nil {
			return
		}

		// Delete all Redo Id's from History Map
		if len(data.CmdLine.RedoAccountIds) > 0 {
			array := strings.Split(data.CmdLine.RedoAccountIds, ",")
			for _, id := range array {
				delete(data.Save.HistoryMap, id)
			}
		}
	} else {
		// Get Ids corresponding to command-line options
		data.Save.AccountIds = strings.Split(data.CmdLine.AccountIds, ",")
	}

	// Read config files
	err = data.ReadBusinessConfigs()
	if err != nil {
		return
	}
	err = data.ReadAccountConfigs()
	return
}
