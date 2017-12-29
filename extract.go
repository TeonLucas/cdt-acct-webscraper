package cdtacctwebscraper

import (
	"fmt"
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"strconv"
)

// This function is used to extract a list as a table
func (page *ScrapePage) GetList(acctLog *log.Logger, rowTag, rowAttr string) (err error) {
	var doc *goquery.Document

	//fmt.Println("INNER HTML:", html)
	doc, err = goquery.NewDocumentFromReader(strings.NewReader(page.TableHtml))
	if err != nil {
		err = fmt.Errorf("Issue parsing Table: %v", err)
		return
	}

	// Find each list item
	doc.Find(rowTag).Each(func(i int, row *goquery.Selection) {
		// Get item attribute
		rowText := []string{}
		if len(rowAttr) > 0 {
			attText, _ := row.Attr(rowAttr)
			rowText = append(rowText, strings.Replace(strings.TrimSpace(attText), "\u00a0", " ", -1))
		}
		// Get item text
		rowText = append(rowText, strings.Replace(strings.TrimSpace(row.Text()), "\u00a0", " ", -1))
		//fmt.Printf("Row %.2d: %q\n", i+1,rowText)
		if len(rowText) > 0 {
			page.Table = append(page.Table, rowText)
		}
	})

	// Does not print if only header row
	numRows := len(page.Table)
	if numRows > 1 {
		acctLog.Printf("Extracted table, %d rows\n", numRows)
	}
	return
}

// This function can get called twice, once for header row, then for data rows
func (page *ScrapePage) GetTable(acctLog *log.Logger, rowTag, colTag string) (err error) {
	var doc *goquery.Document

	//fmt.Println("INNER HTML:", html)
	doc, err = goquery.NewDocumentFromReader(strings.NewReader(page.TableHtml))
	if err != nil {
		err = fmt.Errorf("Issue parsing Table: %v", err)
		return
	}

	// Find each row
	doc.Find(rowTag).Each(func(i int, row *goquery.Selection) {
		// For each col, get text
		rowText := []string{}
		row.Find(colTag).Each(func(j int, col *goquery.Selection) {
			rowText = append(rowText, strings.Replace(strings.TrimSpace(col.Text()), "\u00a0", " ", -1))
		})
		//fmt.Printf("Row %.2d: %q\n", i+1,rowText)
		if len(rowText) > 0 {
			page.Table = append(page.Table, rowText)
		}
	})

	// Does not print if only header row
	numRows := len(page.Table)
	if numRows > 1 {
		acctLog.Printf("Extracted table, %d rows\n", numRows)
	}
	return
}

// This function works only after GetTable
func (browser *Browser) WalkDownTable(page *ScrapePage, rowTag, colTag, linkText string, commands Commands, indent int) (err error) {
	var doc *goquery.Document
	var text string
	var didHeader bool

	numRows := len(page.Table)
	if numRows < 2 {
		err = fmt.Errorf("Table too short: %d rows", numRows)
		return
	}

	doc, err = goquery.NewDocumentFromReader(strings.NewReader(page.TableHtml))
	if err != nil {
		err = fmt.Errorf("Issue parsing Table: %v", err)
		return
	}

	linkText = strings.ToLower(linkText)

	// Visit each row
	var rowNum int
	rowNum = 1
	doc.Find(rowTag).Each(func(i int, row *goquery.Selection) {
		// For each col, get text
		row.Find(colTag).Each(func(j int, col *goquery.Selection) {
			text = strings.ToLower(strings.TrimSpace(col.Text()))
			if strings.Contains(text, linkText) {
				col.Find("a").Each(func(k int, link *goquery.Selection) {
					href, hasHref := link.Attr("href")
					i := rowNum
					if hasHref {
						// Substitute tags
						var new Commands
						new, err = HrefSubstitute(commands, href)
						if err != nil {
							return
						}

						// Execute commands
						err = browser.ExecuteCommands(page, new, &Tasks{Name: "WalkDownTable", Indent: indent + 1})
						if err != nil {
							return
						}

						// Append results to table
						if !didHeader {
							page.Table[0] = append(page.Table[0], page.Heading)
							fmt.Printf("Appending heading %q to Table[0]\n", page.Heading)
							didHeader = true
						}
						fmt.Printf("Appending cell %q to Table[%d]\n", page.Text, i)
						page.Table[i] = append(page.Table[i], page.Text)
					} else {
						fmt.Println("Href not found")
					}
				})
				rowNum++
			}
		})
	})

	return
}

func (page *ScrapePage) GetCellWithHeading(rowTag, headTag, colTag, tHeading string) (err error) {
	var doc *goquery.Document
	var heading string
	var rPos, cPos int

	tHeading = strings.ToLower(tHeading)

	// fmt.Println("INNER HTML:", html)
	doc, err = goquery.NewDocumentFromReader(strings.NewReader(page.TextHtml))
	if err != nil {
		err = fmt.Errorf("Issue parsing Cell: %v", err)
		return
	}

	doc.Find(rowTag).Each(func(r int, row *goquery.Selection) {
		row.Find(headTag).Each(func(c int, hCol *goquery.Selection) {
			// Locate heading with matching text
			heading = hCol.Text()
			if strings.Contains(strings.ToLower(heading), tHeading) {
				rPos = r
				cPos = c
				if !page.Flag {
					page.Heading = strings.Replace(strings.TrimSpace(heading), "\u00a0", " ", -1)
					page.Flag = true
				}
			}
		})
	})
	doc.Find(rowTag).Each(func(r int, row *goquery.Selection) {
		if r == rPos {
			row.Find(colTag).Each(func(c int, col *goquery.Selection) {
				if c == cPos {
					page.Text = strings.Replace(strings.TrimSpace(col.Text()), "\u00a0", " ", -1)
				}
			})
		}
	})
	return
}

func (page *ScrapePage) GetCell(tag, line string) (err error) {
	var doc *goquery.Document
	var n int

	// fmt.Println("INNER HTML:", html)
	doc, err = goquery.NewDocumentFromReader(strings.NewReader(page.TextHtml))
	if err != nil {
		err = fmt.Errorf("Issue parsing Table: %v", err)
		return
	}

	// Optional "line=" parameter selects a specific line of text from a cell
	if len(line) > 0 {
		n, err = strconv.Atoi(strings.Replace(strings.ToLower(line), "line=", "", 1))
		if err != nil {
			err = fmt.Errorf("Issue with 'line=' param: %v", err)
			return
		}
	}

	page.Text = ""
	doc.Find(tag).Each(func(i int, elem *goquery.Selection) {
		page.Text += strings.Replace(strings.TrimSpace(elem.Text()), "\u00a0", " ", -1)
	})

	if len(line) > 0 {
		text := strings.Split(page.Text, "\n")
		if len(text) > n-1 {
			page.Text = strings.TrimSpace(text[n-1])
		}
	}
	return
}
