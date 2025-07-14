package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type State struct {
	Entries    []Entry
	DataDir    string
}

func NewState(dataDir string) *State {
	return & State{
		Entries: []Entry{},
		DataDir: dataDir,
	}
}

// Returns the state in a formated string
func (s *State) String() string {
	return fmt.Sprintf(
		"State\n-----\nEntries: %d\nDataDir: %s\n-----\n",
		len(s.Entries),
		s.DataDir,
	)
}

// This method reads the content of a CSV file in the Entry format and saves it to the state
func (s *State) LoadEntries(filename string) error {
	filePath := filepath.Join(s.DataDir, filename + ".csv")
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	entries := []Entry{}
	
	// Skip header row
	for i := 1; i < len(records); i++ {
		record := records[i]
		
		// Skip weird records
		if len(record) != 6 {
			continue
		}
		
		applyDate, err := time.Parse(time.RFC3339Nano, record[4])
		if err != nil {
			continue
		}

		contactReplied, err := strconv.ParseBool(record[5])
		if err != nil {
			continue
		}

		entry := Entry{
			Company: record[0],
			Role: record[1],
			Status: record[2],
			Platform: record[3],
			ApplyDate: applyDate,
			ContactReplied: contactReplied,
		}
		entries = append(entries, entry)
	}

	s.Entries = entries
	
	return nil
}

// This method saves all the entries of the state into a CSV file with the filename specified.
func (s *State) SaveEntries(filename string) error {
	os.MkdirAll(s.DataDir, 0755)
	filePath := filepath.Join(s.DataDir, filename + ".csv")
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	header := []string{"company", "position", "status", "platform", "apply_date", "contact_replied"}
	if err = writer.Write(header); err != nil {
		return err
	}

	for _, e := range s.Entries {
		record := []string{
			e.Company,
			e.Role,
			e.Status,
			e.Platform,
			e.ApplyDate.Format(time.RFC3339Nano),
			strconv.FormatBool(e.ContactReplied),
		}
		if err := writer.Write(record); err != nil {
			return err
		}
	}
	
	return nil
}

func (s *State) ListEntries() {
	for i, e :=  range s.Entries {
		y, m, d := e.ApplyDate.Date() 

		var replied string
		if e.ContactReplied == true {
			replied = "Yes"
		} else {
			replied = "No"
		}

		fmt.Printf("%d. %s\n", i + 1, e.Company)
		fmt.Printf("%d/%d/%d\n", d, m, y)
		fmt.Printf("%s -> replied: %s\n", e.Status, replied)
		fmt.Println("----------")
	}
}
