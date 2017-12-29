package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	ws "github.com/DavidSantia/cdt-acct-webscraper"
)

func main() {
	// Initialize Affiliate data and read command line
	data, err := ws.Begin()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Check accounts specified via command line
	if len(data.Save.AccountIds) == 0 {
		fmt.Println("No accounts specified")
		os.Exit(1)
	}
	fmt.Printf("%d accounts specified\n", len(data.Save.AccountIds))

	// Channel for scrape results
	resultChan := make(chan ws.ScrapeResult)
	workers := 0

	// Start scrape workers
	var skiplist []string
	for _, id := range data.Save.AccountIds {
		_, ok := data.Save.HistoryMap[id]
		if ok {
			skiplist = append(skiplist, id)
			continue
		}

		// Create browser
		browser := ws.Browser{
			Log:  log.New(os.Stderr, "Acct "+id+": ", 0),
			Id:   "Acct" + id,
			Port: 9222 + workers*2,
		}
		go data.ScrapeAccount(browser, id, resultChan)
		workers++
	}

	if len(skiplist) > 0 {
		fmt.Printf("Skipped Account Ids %s - already scraped\n", strings.Join(skiplist, ","))
	}
	fmt.Printf("Waiting for %d workers to finish\n", workers)

	// Wait for workers to finish
	for i := 0; i < workers; i++ {
		// Get response from each
		result := <-resultChan
		fmt.Printf("• %d workers done\n", i+1)

		if result.Err != nil {
			fmt.Printf("Error Account Id %s: %v\n", result.AccountId, result.Err)
		}

		// Store extracted table in HistoryData
		if len(result.Table) > 0 {
			data.StoreTable(result.AccountId, result.Table)
		}
	}

	fmt.Printf("Scraping complete: Have %d out of %d History tables\n",
		len(data.Save.HistoryMap), len(data.Save.AccountIds))

	if len(data.Save.HistoryMap) == 0 {
		fmt.Println("No History data, exiting")
		os.Exit(0)
	}

	// If scraping was done, save in snapshot
	if workers > 0 {
		err = data.WriteSnapshot()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	// Show extracted data
	if data.CmdLine.ShowTable {
		for id, pd := range data.Save.HistoryMap {
			ac, _ := data.AccountMap[id]
			fmt.Printf("== %s [%s] Extracted Table ==\n", ac.AccountName, id)
			for i, hrow := range pd.HistoryTable {
				fmt.Printf("%.2d %q\n", i, hrow.Extracted)
			}
		}
	}

	err = data.ParseHistory()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = data.MakeHistoryReport()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	os.Exit(0)
}
