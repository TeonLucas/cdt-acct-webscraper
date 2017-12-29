package cdtacctwebscraper

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
)

func (data *AccountScrapeData) ReadBusinessConfigs() (err error) {
	var matches, matches_shortname []string
	var filename, shortname string
	var i int
	var b []byte
	var ok bool

	if matches, err = filepath.Glob(data.CmdLine.BusinessConfigDir + "/*.json"); err != nil {
		return
	}
	if len(matches) == 0 {
		return fmt.Errorf("No business config files found")
	}
	for _, filename = range matches {
		matches_shortname = append(matches_shortname, filepath.Base(filename))
	}

	fmt.Printf("Reading %d business config files %v\n", len(matches), matches_shortname)
	for i, filename = range matches {
		shortname = matches_shortname[i]
		if b, err = ioutil.ReadFile(filename); err != nil {
			err = fmt.Errorf("Issue reading business config: %v", err)
			return
		}

		// Parse business config file
		bMap := make(map[string]BusinessConfig)
		if err = json.Unmarshal(b, &bMap); err != nil {
			err = fmt.Errorf("Issue reading business config %s: %v", shortname, err)
			return
		}

		for bid, config := range bMap {
			_, ok = data.BusinessMap[bid]
			if ok {
				// Duplicate Business Id
				err = fmt.Errorf("Duplicate Business Id %s specified in %s", bid, shortname)
				return
			}

			// New Business Id, insert
			data.BusinessMap[bid] = config
		}
	}
	return
}

func (data *AccountScrapeData) ReadAccountConfigs() (err error) {
	var matches, matches_shortname []string
	var filename, shortname string
	var i int
	var b []byte
	var bc BusinessConfig
	var ok bool

	if matches, err = filepath.Glob(data.CmdLine.AccountConfigDir + "/*.json"); err != nil {
		return
	}
	if len(matches) == 0 {
		return fmt.Errorf("No account config files found")
	}
	for _, filename = range matches {
		matches_shortname = append(matches_shortname, filepath.Base(filename))
	}

	fmt.Printf("Reading %d account config files %v\n", len(matches), matches_shortname)
	for i, filename = range matches {
		shortname = matches_shortname[i]
		if b, err = ioutil.ReadFile(filename); err != nil {
			err = fmt.Errorf("Issue reading account config: %v", err)
			return
		}

		// Parse account config file
		configs := []AccountConfig{}
		if err = json.Unmarshal(b, &configs); err != nil {
			err = fmt.Errorf("Issue reading account config: %v", err)
			return
		}

		for _, config := range configs {
			_, ok = data.AccountMap[config.AccountId]
			if ok {
				// Duplicate Account
				err = fmt.Errorf("Duplicate Account Number %s (%q in %s)",
					config.AccountId, config.AccountName, shortname)
				return
			}

			bc, ok = data.BusinessMap[config.BusinessId]
			if !ok {
				// Business Id not configured
				err = fmt.Errorf("Business Id %s not configured in BusinessMap (%q in %s)",
					config.BusinessId, config.AccountName, shortname)
				// Print as warning for now
				fmt.Println("Warn:", err)
				err = nil
				continue
			}

			// If no override for account, use business defaults
			if len(config.Url) == 0 {
				config.Url = bc.Url
			}
			if len(config.DateFormat) == 0 {
				config.DateFormat = bc.DateFormat
			}
			if len(config.CurrencyFormat) == 0 {
				config.CurrencyFormat = bc.CurrencyFormat
			}
			if len(config.AmountColumn) == 0 {
				config.AmountColumn = bc.AmountColumn
			}
			if len(config.DateColumn) == 0 {
				config.DateColumn = bc.DateColumn
			}

			// Remove spaces
			config.CompactName = strings.Replace(config.AccountName, " ", "", -1)

			// New Business Id, insert
			data.AccountMap[config.AccountId] = config
		}
	}
	return
}
