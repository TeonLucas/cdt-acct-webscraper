package cdtacctwebscraper

import (
	"context"

	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/runner"
)

const Chromium = "/Applications/Chromium.app/Contents/MacOS/Chromium"

func (browser *Browser) LaunchBrowser(chromium bool) (err error) {
	// Create context
	browser.Ctxt, browser.Cancel = context.WithCancel(context.Background())
	name := "Google Chrome"

	options := []runner.CommandLineOption{
		runner.Flag("disable-save-password-bubble", true),
		runner.Flag("window-size", "1200,800"),
	}

	// Specify Chromium instead of default Chrome
	if chromium {
		name = "Chromium"
		options = append(options, []runner.CommandLineOption{
			runner.Flag("no-first-run", true),
			runner.Flag("no-default-browser-check", true),
			runner.Flag("remote-debugging-port", browser.Port),
			runner.ExecPath(Chromium),
		}...)
	}

	// Create browser instance
	browser.Cdp, err = chromedp.New(browser.Ctxt, chromedp.WithLog(debugLog.Printf),
		chromedp.WithRunnerOptions(options...))

	// Create a target for each port
	browser.Cdp.NewTarget(&browser.Id)
	browser.Log.Printf("Launch %s, Target Id %q\n", name, browser.Cdp.ListTargets()[0])

	return
}

func (browser *Browser) Shutdown() (err error) {
	// Shutdown browser
	err = browser.Cdp.Shutdown(browser.Ctxt)
	if err != nil {
		return
	}

	// Wait to finish
	err = browser.Cdp.Wait()
	browser.Cancel()
	return
}
