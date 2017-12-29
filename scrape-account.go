package cdtacctwebscraper

import (
	"fmt"
	"reflect"
)

// Scrape a Business Id
//   use browser.Log for status messages
//   return extracted data and error status via channel

func (data *AccountScrapeData) ScrapeAccount(browser Browser, id string, result chan ScrapeResult) {
	var ac AccountConfig
	var bc BusinessConfig
	var ok bool
	var tag string
	var i int
	var pages []ScrapePage
	var err error

	scrapeResult := ScrapeResult{AccountId: id}

	// See if already scraped (when Snapshot was loaded)
	_, ok = data.Save.HistoryMap[id]
	if ok {
		scrapeResult.Err = fmt.Errorf("Account Id %s already scraped\n", id)
		result <- scrapeResult
		return
	}

	ac, ok = data.AccountMap[id]
	if !ok {
		scrapeResult.Err = fmt.Errorf("No account configured for Id %s", id)
		result <- scrapeResult
		return
	}

	bc, ok = data.BusinessMap[ac.BusinessId]
	if !ok {
		scrapeResult.Err = fmt.Errorf("No business configured for Id %s", ac.BusinessId)
		result <- scrapeResult
		return
	}

	err = browser.LaunchBrowser(data.CmdLine.Chromium)
	if err != nil {
		scrapeResult.Err = err
		result <- scrapeResult
		return
	}
	defer browser.Cancel()

	st := reflect.TypeOf(bc)
	val := reflect.ValueOf(&bc).Elem()

	// Iterate through struct fields, access scrape tags to compile commands
	for i = 0; i < st.NumField(); i++ {
		field := st.Field(i)
		if tag, ok = field.Tag.Lookup("scrape"); ok {
			if tag != "" {
				// Create new scrape-page from list
				page := ScrapePage{Name: tag, AccountId: id}
				page.Commands, err = GenerateCommands(val.Field(i).Interface().([][]string))
				if err != nil {
					scrapeResult.Err = fmt.Errorf("%s page error %v", page.Name, err)
					result <- scrapeResult
					return
				}
				pages = append(pages, page)
			}
		}
	}

	browser.Log.Printf("Starting scrape for %q\n", ac.AccountName)
	for _, page := range pages {
		// Substitute tags
		page.Commands, err = ac.VariableSubstitute(page.Commands)
		if err != nil {
			scrapeResult.Err = fmt.Errorf("%s page error %v", page.Name, err)
			result <- scrapeResult
			return
		}
		//browser.PrintCommands(page.Commands, 0)

		// Execute commands
		tasks := Tasks{Name: page.Name, Indent: 0}
		err = browser.ExecuteCommands(&page, page.Commands, &tasks)
		if err != nil {
			scrapeResult.Err = fmt.Errorf("%s config error %v", page.Name, err)
			result <- scrapeResult
			return
		}

		// Save extracted account name if any
		if len(page.Text) > 0 && len(scrapeResult.Account) == 0 {
			scrapeResult.Account = page.Text
		}

		// Save extracted table if any
		if len(page.Table) > 0 && len(scrapeResult.Table) == 0 {
			scrapeResult.Table = page.Table
		}
	}

	err = browser.Shutdown()
	if err != nil {
		scrapeResult.Err = err
	}

	result <- scrapeResult
	return
}
