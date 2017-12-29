package cdtacctwebscraper

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func (ac AccountConfig) ParseFormats(pd *HistoryData) (err error) {
	var col, dateColHeading, amountColHeading string
	var dateGreedy, amountGreedy bool
	var i, j int
	var zeroT, maxDateT time.Time

	errlist := NewErrListPrefix("ParseFormats for " + ac.AccountName)

	if pd.HistoryToMatch == 0 {
		err = fmt.Errorf("No history to process for %s", ac.AccountName)
		return
	}

	// Initialize column match tags
	if len(ac.DateColumn) == 0 {
		dateColHeading, dateGreedy = "date", false
	} else {
		dateColHeading, dateGreedy = CheckGreedy(ac.DateColumn)
	}
	if len(ac.AmountColumn) == 0 {
		amountColHeading, amountGreedy = "amount", false
	} else {
		amountColHeading, amountGreedy = CheckGreedy(ac.AmountColumn)
	}

	for i = range pd.HistoryTable {
		if i == 0 {
			// Locate date and amount columns
			pd.DateCol = -1
			pd.AmountCol = -1
			for j, col = range pd.HistoryTable[i].Extracted {
				heading := strings.ToLower(col)
				if strings.Contains(heading, dateColHeading) {
					if pd.DateCol == -1 || pd.DateCol != -1 && dateGreedy {
						pd.DateCol = j
					}
				}
				if strings.Contains(heading, amountColHeading) {
					if pd.AmountCol == -1 || pd.AmountCol != -1 && amountGreedy {
						pd.AmountCol = j
					}
				}
			}
			if pd.DateCol == -1 {
				err = fmt.Errorf("Failed to locate date column for %s", ac.AccountName)
				return
			}
			if pd.AmountCol == -1 {
				err = fmt.Errorf("Failed to locate amount column for %s", ac.AccountName)
				return
			}
		} else if !pd.HistoryTable[i].Skip {

			l := len(pd.HistoryTable[i].Extracted)

			// Process dates
			if l <= pd.DateCol {
				col = ""
			} else {
				col = pd.HistoryTable[i].Extracted[pd.DateCol]
			}
			pd.HistoryTable[i].DateT, pd.HistoryTable[i].Date, err = ParseDate(ac.DateFormat, col)
			if err != nil {
				errlist.Add(fmt.Sprintf("Could not parse date %s in row %d", col, i))
			}

			// Update min date
			if pd.MinDateT.After(pd.HistoryTable[i].DateT) && pd.HistoryTable[i].DateT != zeroT ||
				pd.MinDateT == zeroT {
				pd.MinDateT = pd.HistoryTable[i].DateT
			}

			// Update max date
			if maxDateT.Before(pd.HistoryTable[i].DateT) && pd.HistoryTable[i].DateT != zeroT ||
				maxDateT == zeroT {
				maxDateT = pd.HistoryTable[i].DateT
				pd.MaxDateRow = i
			}

			// Process amounts
			if l <= pd.AmountCol {
				col = ""
			} else {
				col = pd.HistoryTable[i].Extracted[pd.AmountCol]
			}
			pd.HistoryTable[i].Amount, err = ParseCurrency(ac.CurrencyFormat, col)
			if err != nil {
				errlist.Add(fmt.Sprintf("Could not parse amount %s in row %d", col, i))
			}
		}
	}

	err = errlist.Get()
	return
}

func ParseCurrency(format, amount string) (result float64, err error) {
	var prefix, suffix, comma, decimal string
	var i, j, k, l, s int

	if len(amount) == 0 {
		result = 0
		return
	}

	// Example format: $1,234.56
	l = len(format)

	// locate digit sequence start and end
	i = strings.Index(format, "1")
	j = strings.Index(format, "4")
	k = strings.Index(format, "6")

	if l != 0 {
		if i < 0 {
			err = fmt.Errorf("Invalid format %q: missing '1' in digit sequence\n", format)
			return
		}
		if j < 0 {
			err = fmt.Errorf("Invalid format %q: missing '4' in digit sequence\n", format)
			return
		}

		// Strip off prefix currency symbol
		if i > 0 {
			prefix = format[:i]
			amount = strings.TrimPrefix(amount, prefix)
		}

		// Strip off suffix currency symbol
		if k > 0 {
			s = k
		} else {
			s = j
		}
		if s < l-1 {
			suffix = format[s+1:]
			amount = strings.TrimSuffix(amount, suffix)
		}

		// Update decimal if needed
		if k > 0 {
			decimal = format[j+1 : j+2]
			if decimal != "." {
				amount = strings.Replace(amount, decimal, ".", 1)
			}
		}

		// comma is 2 means there is no comma
		comma = format[i+1 : i+2]
		if comma != "2" {
			amount = strings.Replace(amount, comma, "", -1)
		}
	}

	result, err = strconv.ParseFloat(amount, 64)
	return
}

func ParseDate(format, date string) (resultT time.Time, result string, err error) {
	var zeroT time.Time

	if len(date) == 0 {
		resultT = zeroT
	} else {
		if len(format) == 0 {
			format = "2006-01-02"
		}

		resultT, err = time.Parse(format, date)
		if err != nil {
			resultT = zeroT
		}

		// Add current year if date format only had Month/Day
		if resultT.Year() == 0 {
			resultT = resultT.AddDate(time.Now().Year(), 0, 0)
		}
	}

	result = resultT.Format("2006-01-02")
	return
}

func CheckGreedy(tag string) (finalTag string, greedy bool) {
	if tag[len(tag)-1:] == "?" {
		finalTag = strings.ToLower(tag[:len(tag)-1])
		greedy = false
	} else {
		finalTag = strings.ToLower(tag)
		greedy = true
	}

	return
}
