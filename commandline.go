package cdtacctwebscraper

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func (cmdline *CmdLine) init() {
	// Define command-line arguments
	flag.StringVar(&cmdline.AccountConfigDir, "ac", "", "Directory of JSON account config files")
	flag.StringVar(&cmdline.BusinessConfigDir, "bc", "", "Directory of JSON business config files")
	flag.StringVar(&cmdline.SnapShotFile, "fromss", "", "Snapshot file to load and continue from previous run")
	flag.StringVar(&cmdline.RedoAccountIds, "redo", "", "Redo scrape on account_id[,account_id]")
	flag.StringVar(&cmdline.AccountIds, "id", "", "One or more account_id[,account_id]")
	flag.StringVar(&cmdline.Logfile, "log", "scrape.log", "Specify logging filename")
	flag.BoolVar(&cmdline.Debug, "debug", false, "Debug logging")
	flag.BoolVar(&cmdline.Chromium, "chromium", false, "Use Chromium browser")
	flag.BoolVar(&cmdline.ShowTable, "show", false, "Show extracted tables")
}

// Make sure required arguments are specified
func (cmdline *CmdLine) validateFlags() (err error) {
	var absolutePath, relPath string
	var dirInfo os.FileInfo

	// Validate Logfile
	if len(cmdline.Logfile) != 0 {
		relPath = filepath.Dir(cmdline.Logfile)
		if absolutePath, err = filepath.Abs(cmdline.Logfile); err != nil {
			return fmt.Errorf("-logfile %s %v", cmdline.Logfile, err)
		}
		dirInfo, err = os.Stat(relPath)
		if os.IsNotExist(err) {
			return fmt.Errorf("-logfile %s: directory %s %v", cmdline.Logfile, relPath, err)
		} else if !dirInfo.IsDir() {
			return fmt.Errorf("-logfile %s: %s is not a directory", cmdline.Logfile, relPath)
		}
		cmdline.Logfile = absolutePath
	}

	// Validate BusinessId
	if len(cmdline.AccountIds) > 0 && len(cmdline.SnapShotFile) > 0 {
		return fmt.Errorf("Only one of -id or -fromss allowed")
	}

	if len(cmdline.AccountIds) > 0 {
		cmdline.AccountIds = strings.Replace(cmdline.AccountIds, " ", "", -1)
		if err = validateNumList(cmdline.AccountIds); err != nil {
			return fmt.Errorf("-id %v", err)

		}
	} else if len(cmdline.SnapShotFile) == 0 {
		return fmt.Errorf("Must select -id, or -fromss")
	}

	// Validate redo option
	if len(cmdline.RedoAccountIds) > 0 {
		if len(cmdline.SnapShotFile) == 0 {
			return fmt.Errorf("Must specify -fromss with -redo")
		}
		cmdline.RedoAccountIds = strings.Replace(cmdline.RedoAccountIds, " ", "", -1)
		if err = validateNumList(cmdline.RedoAccountIds); err != nil {
			return fmt.Errorf("-redo %v", err)
		}
	}

	// Validate Account Config
	if len(cmdline.AccountConfigDir) == 0 {
		return fmt.Errorf("flag -ac <dir> required")
	}
	if absolutePath, err = filepath.Abs(cmdline.AccountConfigDir); err != nil {
		return fmt.Errorf("-ac %s %v", cmdline.AccountConfigDir, err)
	}
	dirInfo, err = os.Stat(absolutePath)
	if os.IsNotExist(err) {
		return fmt.Errorf("-ac %s %v", cmdline.AccountConfigDir, err)
	} else if !dirInfo.IsDir() {
		return fmt.Errorf("-ac %s is not a directory", cmdline.AccountConfigDir)
	}
	cmdline.AccountConfigDir = absolutePath

	// Validate Business Config
	if len(cmdline.BusinessConfigDir) == 0 {
		return fmt.Errorf("flag -bc <dir> required")
	}
	if absolutePath, err = filepath.Abs(cmdline.BusinessConfigDir); err != nil {
		return fmt.Errorf("-bc %s %v", cmdline.BusinessConfigDir, err)
	}
	dirInfo, err = os.Stat(absolutePath)
	if os.IsNotExist(err) {
		return fmt.Errorf("-bc %s %v", cmdline.BusinessConfigDir, err)
	} else if !dirInfo.IsDir() {
		return fmt.Errorf("-bc %s is not a directory", cmdline.BusinessConfigDir)
	}
	cmdline.BusinessConfigDir = absolutePath

	// Validate Snapshot file
	if len(cmdline.SnapShotFile) > 0 {
		if absolutePath, err = filepath.Abs(cmdline.SnapShotFile); err != nil {
			return fmt.Errorf("-fromss %s %v", cmdline.SnapShotFile, err)
		}
		dirInfo, err = os.Stat(absolutePath)
		if os.IsNotExist(err) {
			return fmt.Errorf("-fromss %s %v", cmdline.SnapShotFile, err)
		} else if dirInfo.IsDir() {
			return fmt.Errorf("-fromss %s is a directory", cmdline.SnapShotFile)
		}
		cmdline.SnapShotFile = absolutePath
	}
	return
}

func validateNumList(list string) (err error) {
	array := strings.Split(list, ",")
	for i, num := range array {
		if _, err = strconv.Atoi(num); err != nil {
			err = fmt.Errorf("item %d: %q invalid number", i+1, num)
			return
		}
	}
	return
}

func (data *AccountScrapeData) GetCmdLine() (err error) {

	data.CmdLine.init()

	// Parse commandline flag arguments
	flag.Parse()

	// Validate
	err = data.CmdLine.validateFlags()
	return
}
