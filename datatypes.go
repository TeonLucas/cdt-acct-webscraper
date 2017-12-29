package cdtacctwebscraper

import (
	"context"
	"log"
	"time"

	"github.com/chromedp/chromedp"
)

func NewAccountScrapeData() (data *AccountScrapeData) {
	data = &AccountScrapeData{}
	data.AccountMap = make(AccountMap)
	data.BusinessMap = make(BusinessMap)
	data.Save.HistoryMap = make(HistoryMap)
	return
}

// Command line arguments
type CmdLine struct {
	Debug             bool
	Chromium          bool
	ShowTable         bool
	AccountIds        string
	RedoAccountIds    string
	AccountConfigDir  string
	BusinessConfigDir string
	SnapShotFile      string
	Logfile           string
}

// AccountScrape data - main functions are methods of this
type AccountScrapeData struct {
	CmdLine     CmdLine
	Save        SaveSettings
	AccountMap  AccountMap
	BusinessMap BusinessMap
}

// Settings that will be saved in / loaded from snapshot
type SaveSettings struct {
	AccountIds []string
	HistoryMap HistoryMap
}

// Business configuration
type BusinessConfig struct {
	BusinessName string     `json:"business_name"`
	Login        [][]string `json:"login" scrape:"Login"`
	History      [][]string `json:"history" scrape:"History"`
	Logout       [][]string `json:"logout" scrape:"Logout"`
	// Default values for business
	Url            string `json:"url"`
	DateFormat     string `json:"date_format"`
	CurrencyFormat string `json:"currency_format"`
	AmountColumn   string `json:"amount_column"`
	DateColumn     string `json:"date_column"`
}

type BusinessMap map[string]BusinessConfig

// Account configuration
type AccountConfig struct {
	AccountName string `json:"account_name"`
	AccountId   string `json:"account_id"`
	BusinessId  string `json:"business_id"`
	CompactName string
	Username    string `json:"username"`
	Password    string `json:"password"`
	// Account-specific Overrides (optional)
	Url            string `json:"url"`
	DateFormat     string `json:"date_format"`
	CurrencyFormat string `json:"currency_format"`
	AmountColumn   string `json:"amount_column"`
	DateColumn     string `json:"date_column"`
}

type AccountMap map[string]AccountConfig

// Scrape commands and data for a particular page (login, history, logout)
type ScrapePage struct {
	Name      string
	AccountId string
	Commands  Commands
	Tasks     chromedp.Tasks
	TableHtml string
	Table     [][]string
	Heading   string
	TextHtml  string
	Text      string
	Flag      bool
}

// Result of a scrape on a given account, returns via channel from go routine
type ScrapeResult struct {
	AccountId string
	Err       error
	Account   string
	Table     [][]string
}

// Browser configuration
type Browser struct {
	Cdp    *chromedp.CDP
	Ctxt   context.Context
	Cancel context.CancelFunc
	Log    *log.Logger
	Id     string
	Port   int
}

type Command struct {
	Instruction string
	Params      []string
	Commands    Commands
}

type Commands []Command

type Tasks struct {
	Name   string
	Tasks  chromedp.Tasks
	Indent int
}

// Extracted results
type HistoryData struct {
	MinDateT       time.Time
	MaxDateRow     int
	EmptyLen       int
	AmountCol      int
	DateCol        int
	HistoryToMatch int
	HistoryTable   []HistoryRow
}

type HistoryRow struct {
	Skip      bool
	HistoryId string
	Amount    float64
	Date      string
	DateT     time.Time
	Extracted []string
}

type HistoryMap map[string]HistoryData

// Struct to accumulate a list of errors
type ErrList struct {
	Prefix string
	Text   string
	Count  int
}
