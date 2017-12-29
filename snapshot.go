package cdtacctwebscraper

import (
	"encoding/gob"
	"fmt"
	"os"
)

func (data *AccountScrapeData) WriteSnapshot() (err error) {

	filename := "Scrape-Snapshot.gob"

	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("Cannot open snapshot file %s: %v", filename, err)
	}
	e := gob.NewEncoder(f)
	err = e.Encode(data.Save)
	if err != nil {
		return fmt.Errorf("Cannot write %s: %v", filename, err)
	}
	f.Close()

	fmt.Println("Wrote snapshot", filename)
	return
}

func (data *AccountScrapeData) LoadSnapshot() (err error) {
	var f *os.File

	f, err = os.Open(data.CmdLine.SnapShotFile)
	if err != nil {
		return fmt.Errorf("Cannot read snapshot %s: %v", data.CmdLine.SnapShotFile, err)
	}

	d := gob.NewDecoder(f)
	err = d.Decode(&data.Save)
	if err != nil {
		return fmt.Errorf("Cannot decode snapshot %s: %v", data.CmdLine.SnapShotFile, err)
	}
	f.Close()

	fmt.Println("Loaded snapshot", data.CmdLine.SnapShotFile)
	return
}
