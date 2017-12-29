package cdtacctwebscraper

import (
	"context"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
)

// Create list of tasks to execute from config file command list
func (browser *Browser) ExecuteCommands(page *ScrapePage, commands Commands, tasks *Tasks) (err error) {
	pref := strings.Repeat("  ", tasks.Indent)
	browser.Log.Printf("%sExecute %s tasks\n", pref, tasks.Name)

	// Loop through commands, make a list of tasks
	for _, command := range commands {
		browser.Log.Printf("%s%s%q\n", pref, command.Instruction, command.Params)

		err = browser.PushTask(page, command, tasks)
		if err != nil {
			return
		}
	}
	// Then run list of tasks
	err = browser.Cdp.Run(browser.Ctxt, tasks.Tasks)
	return
}

func (browser *Browser) PushTask(page *ScrapePage, command Command, tasks *Tasks) (err error) {
	var buf []byte

	l := len(command.Params)
	switch strings.ToLower(command.Instruction) {
	case "accountname":
		if l < 2 {
			err = fmt.Errorf("Too few parameters for %q", command)
			return
		}
		url := command.Params[0]
		tag := command.Params[1]
		tasks.Tasks = append(tasks.Tasks, chromedp.InnerHTML(url, &page.TextHtml, chromedp.ByQuery))
		line := ""
		if l == 3 {
			line = command.Params[2]
		}
		tasks.Tasks = append(tasks.Tasks, chromedp.ActionFunc(func(context.Context, cdp.Executor) error {
			err = page.GetCell(tag, line)
			browser.Log.Printf("Logged in: %q\n", page.Text)
			return err
		}))
	case "clear":
		sel := command.Params[0]
		tasks.Tasks = append(tasks.Tasks, chromedp.Clear(sel, chromedp.ByQuery))
	case "click":
		sel := command.Params[0]
		tasks.Tasks = append(tasks.Tasks, chromedp.Click(sel, chromedp.ByQuery))
	case "getcell":
		if l < 2 {
			err = fmt.Errorf("Too few parameters for %q", command)
			return
		}
		sel := command.Params[0]
		tasks.Tasks = append(tasks.Tasks, chromedp.InnerHTML(sel, &page.TextHtml, chromedp.ByQuery))
		ctag := command.Params[1]
		line := ""
		if l == 3 {
			line = command.Params[2]
		}
		tasks.Tasks = append(tasks.Tasks, chromedp.ActionFunc(func(context.Context, cdp.Executor) error {
			return page.GetCell(ctag, line)
		}))
	case "getcellwithheading":
		if l < 5 {
			err = fmt.Errorf("Too few parameters for %q", command)
			return
		}
		sel := command.Params[0]
		tasks.Tasks = append(tasks.Tasks, chromedp.InnerHTML(sel, &page.TextHtml, chromedp.ByQuery))
		rtag := command.Params[1]
		htag := command.Params[2]
		ctag := command.Params[3]
		hdg := command.Params[4]
		page.Flag = false
		tasks.Tasks = append(tasks.Tasks, chromedp.ActionFunc(func(context.Context, cdp.Executor) error {
			return page.GetCellWithHeading(rtag, htag, ctag, hdg)
		}))
	case "setcellheading":
		page.Heading = command.Params[0]
	case "getlist":
		if l < 3 {
			err = fmt.Errorf("Too few parameters for %q", command)
			return
		}
		sel := command.Params[0]
		tasks.Tasks = append(tasks.Tasks, chromedp.InnerHTML(sel, &page.TableHtml, chromedp.ByQuery))
		rtag := command.Params[1]
		ctag := command.Params[2]
		tasks.Tasks = append(tasks.Tasks, chromedp.ActionFunc(func(context.Context, cdp.Executor) error {
			return page.GetList(browser.Log, rtag, ctag)
		}))
	case "gettable":
		if l < 3 {
			err = fmt.Errorf("Too few parameters for %q", command)
			return
		}
		sel := command.Params[0]
		tasks.Tasks = append(tasks.Tasks, chromedp.InnerHTML(sel, &page.TableHtml, chromedp.ByQuery))
		rtag := command.Params[1]
		rattr := command.Params[2]
		tasks.Tasks = append(tasks.Tasks, chromedp.ActionFunc(func(context.Context, cdp.Executor) error {
			return page.GetTable(browser.Log, rtag, rattr)
		}))
	case "gettablewithheader":
		if l < 5 {
			err = fmt.Errorf("Too few parameters for %q", command)
			return
		}
		sel := command.Params[0]
		tasks.Tasks = append(tasks.Tasks, chromedp.InnerHTML(sel, &page.TableHtml, chromedp.ByQuery))
		rtagH := command.Params[1]
		ctagH := command.Params[2]
		tasks.Tasks = append(tasks.Tasks, chromedp.ActionFunc(func(context.Context, cdp.Executor) error {
			return page.GetTable(browser.Log, rtagH, ctagH)
		}))
		rtag := command.Params[3]
		ctag := command.Params[4]
		tasks.Tasks = append(tasks.Tasks, chromedp.ActionFunc(func(context.Context, cdp.Executor) error {
			return page.GetTable(browser.Log, rtag, ctag)
		}))
	case "settableheading":
		var params []string
		params = append(params, command.Params...)
		tasks.Tasks = append(tasks.Tasks, chromedp.ActionFunc(func(context.Context, cdp.Executor) error {
			if len(page.Table) > 0 {
				page.Table[0] = append(page.Table[0], params...)
			} else {
				page.Table = append(page.Table, params)
			}
			return nil
		}))
	case "splittablecolumn":
		if l < 3 {
			err = fmt.Errorf("Too few parameters for %q", command)
			return
		}
		sep := command.Params[0]
		heading := command.Params[1]
		newHeading := command.Params[2]
		tasks.Tasks = append(tasks.Tasks, chromedp.ActionFunc(func(context.Context, cdp.Executor) error {
			return page.SplitTableColumn(sep, heading, newHeading)
		}))
	case "padtablecolumn":
		padspec := command.Params[0]
		tasks.Tasks = append(tasks.Tasks, chromedp.ActionFunc(func(context.Context, cdp.Executor) error {
			return page.PadTableColumn(padspec)
		}))
	case "walkdowntable":
		if l < 3 {
			err = fmt.Errorf("Too few parameters for %q", command)
			return
		}
		rtag := command.Params[0]
		ctag := command.Params[1]
		linkt := command.Params[2]
		commands := command.Commands
		tasks.Tasks = append(tasks.Tasks, chromedp.ActionFunc(func(context.Context, cdp.Executor) error {
			return browser.WalkDownTable(page, rtag, ctag, linkt, commands, tasks.Indent)
		}))
	case "input":
		if l < 2 {
			err = fmt.Errorf("Too few parameters for %q", command)
			return
		}
		tasks.Tasks = append(tasks.Tasks, chromedp.SendKeys(command.Params[0], command.Params[1], chromedp.ByQuery))
	case "navigate":
		url := command.Params[0]
		tasks.Tasks = append(tasks.Tasks, chromedp.Navigate(url))
	case "setvalue":
		if l < 2 {
			err = fmt.Errorf("Too few parameters for %q", command)
			return
		}
		sel := command.Params[0]
		val := command.Params[1]
		tasks.Tasks = append(tasks.Tasks, chromedp.SetValue(sel, val, chromedp.ByQuery))
	case "sleep":
		var seconds float64
		seconds, err = strconv.ParseFloat(command.Params[0], 64)
		if err != nil {
			err = fmt.Errorf("Invalid seconds %v", err)
			return
		}
		tasks.Tasks = append(tasks.Tasks, chromedp.Sleep(time.Duration(seconds*1000)*time.Millisecond))
	case "snapshot":
		tasks.Tasks = append(tasks.Tasks, chromedp.CaptureScreenshot(&buf))
		ssfile := command.Params[0] + ".png"
		tasks.Tasks = append(tasks.Tasks, chromedp.Sleep(5*time.Second))
		tasks.Tasks = append(tasks.Tasks, chromedp.ActionFunc(func(context.Context, cdp.Executor) error {
			browser.Log.Printf("Saving snapshot %q\n", ssfile)
			return ioutil.WriteFile(ssfile, buf, 0644)
		}))
	case "submit":
		sel := command.Params[0]
		tasks.Tasks = append(tasks.Tasks, chromedp.Submit(sel, chromedp.ByQuery))
	case "waitnotpresent":
		sel := command.Params[0]
		tasks.Tasks = append(tasks.Tasks, chromedp.WaitNotPresent(sel, chromedp.ByQuery))
	case "waitnotvisible":
		sel := command.Params[0]
		tasks.Tasks = append(tasks.Tasks, chromedp.WaitNotVisible(sel, chromedp.ByQuery))
	case "waitready":
		sel := command.Params[0]
		tasks.Tasks = append(tasks.Tasks, chromedp.WaitReady(sel, chromedp.ByQuery))
	default:
		err = fmt.Errorf("Unrecognized command %q", command)
		return
	}
	return
}
