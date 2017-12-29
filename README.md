# cdt-acct-webscraper
Web scraper that uses Chrome DevTools to login and access configured accounts via Chrome DevTools.

The [Chrome DevTools Protocol](https://developer.chrome.com/devtools/docs/debugger-protocol) enables a faster, simpler way to drive browsers from Go-lang (Chrome, Edge, Safari, etc.) without external dependencies such as Selenium or PhantomJS.

### Key Features
* Configure general scrape method for each business
* Configure specific credentials for each account
* Parse date and amount using configured formats
* Report scraped table with parsed amounts in XLSX

### Key Packages used
* "**reflect**" for configuration
* "**github.com/chromedp/chromedp**" for Chrome DevTools
* "**github.com/PuerkitoBio/goquery**" to select HTML elements and extract tables
* "**encoding/gob**" for snapshot

### Data types
*AccountConfig* is the credentials for each account to access.
```go
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
```
*BusinessConfig* is the DOM tasks to log in, navigate to history table, and log out for each business.
```go
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
```

## How Scraping is Configured
For each business, identify:
* The URLs to login, and navigate
* The HTML element tags where input is required, or content is to be extracted

Assign a unique id (example: 1) as shown
```json
  "1": {
    "business_name": "Starbucks",
      :
  }
```

Then configure each section of scrape commands:
```json
    "login": [ ... ],
    "history": [ ... ],
    "logout": [ ... ],
````

### Login Section
Here is an example showing a basic set of DOM tasks to login. 
```json
      [
        "Navigate",
        "{Url}/account/signin"
      ],
      [
        "WaitReady",
        "main#content"
      ],
```
* First we **Navigate** to the member login page
* Then we wait until the content displays on the page

```json
      [
        "Input",
        "input#username",
        "{Username}"
      ],
      [
        "Input",
        "input#password",
        "{Password}"
      ],
      [
        "Click",
        "button.sb-frap"
      ],
```
* Then we input the *Username* and *Password*
* And submit the form

#### Login Success
Once logged in, the account name will appear.
```json
      [
        "AccountName",
        "div.profileBanner__container",
        "span"
      ]
```
The **AccountName** command waits for the specified element, and displays its text like so:
```
• Logged in: 
Acct 123: Logged in: "Hi, Santia."

```

### History Section
Next we navigate to the history page.
```json
      [
        "Navigate",
        "{Url}/account/history"
      ],
      [
        "WaitReady",
        "div.historyWrapper"
      ],
      [
        "Sleep",
        "2"
      ],
      [
        "Snapshot",
        "{CompactName}"
      ],
      [
        "GetTableWithHeader",
        "div.historyWrapper",
        "div.column",
        "h2",
        "div.column li",
        "h3,span.historyItemMessage"
      ],
```
We extract first a screenshot.
* The **Snapshot** will be named *AccountName*.png

And then the history table
* The **GetTableWithHeader** command takes 5 parameters:
  1. identifier of the div containing the table
  1. the class on the header row
  1. the class on each header cell
  1. the class on each data row
  1. the class on each data cell

## How Commands Select Nodes
Commands use DOM.querySelector().  You can select elements by id="tag", element type, and/or by css styles.

For example, the following command selects the element with id="login":
```json
      [
        "WaitVisible",
        "#login"
      ],
```
For a syntax reference, see: [CSS selectors](https://developer.mozilla.org/en-US/docs/Web/CSS/CSS_Selectors)

## How to build

Prerequisites:  Go-lang and Chrome (or Chromium)

For example, to log in to Starbucks, edit account_config/MyAccounts.json to provide username and password

#### Build
```sh
cd scrape
go build
./scrape -help
```

#### Run
Scrape with Chromium (automatically saves snapshot, writes *HistoryReport.xlsx*)
```sh
./scrape -ac ../account_config -bc ../business_config -chromium -id 123
```

#### Run from snapshot
Loads scraped data from .gob file, re-writes report (-show option prints dump of extracted table).
```sh
./scrape -ac ../account_config -bc ../business_config -fromss Scrape-Snapshot.gob -show
```
The snapshot is used so you can tweak parsing separately from scraping.


## Chromium (use to fix the version)

Chromium is designed to be installed in parallel to Google Chrome, so you can set a stable version.

* For example, there was a bug for sending keys to an input in Chrome 62: [Chromedp issue 130](https://github.com/chromedp/chromedp/issues/130)
* You can install version 61 of Chromium to work around this issue: [Ex: Mac Snapshot 488534](https://commondatastorage.googleapis.com/chromium-browser-snapshots/index.html?prefix=Mac/488533/)

