package cdtacctwebscraper

import (
	"errors"
	"fmt"
)

// Make new error list
func NewErrList() (err_list *ErrList) {
	err_list = &ErrList{Prefix: "", Text: "", Count: 0}
	return
}

func NewErrListPrefix(prefix string) (err_list *ErrList) {
	err_list = &ErrList{Prefix: prefix, Text: "", Count: 0}
	return
}

// Appends err_text to list
func (err_list *ErrList) Add(err_text string) {
	err_list.Count++

	// return the err_text, plain and simple
	if err_list.Count == 1 {
		err_list.Text = err_text
		return
	}

	// use the count to display errors, one per line
	if err_list.Count == 2 {
		err_list.Text = fmt.Sprintf("multiple errors\n(#1) %s\n(#2) %s", err_list.Text, err_text)
	} else {
		err_list.Text = fmt.Sprintf("%s\n(#%d) %s", err_list.Text, err_list.Count, err_text)
	}
}

// Formats list into error type
func (err_list ErrList) Get() (err error) {
	var prefix string
	if err_list.Prefix != "" {
		prefix = err_list.Prefix + ": "
	}
	if err_list.Count != 0 {
		err = errors.New(prefix + err_list.Text)
	}
	return
}
