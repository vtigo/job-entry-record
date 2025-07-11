package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func main() {
	state := NewState("data", "csv")
	fmt.Println(state)
	err := state.LoadEntries("test")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(state.String())

	// e1 := Entry{
	// 	Company: "emp",
	// 	Role: "dev",
	// 	Status: "paz",
	// 	Platform: "web",
	// 	ApplyDate: time.Now(),
	// 	ContactReplied: false,
	// }
	//
	// SaveEntries("test", []Entry{e1})
}

type Entry struct {
	Company        string
	Role       	   string
	Status         string
	Platform       string
	ApplyDate 	   time.Time
	ContactReplied bool
}

type State struct {
	Entries    []Entry
	DataDir    string
	DataFormat string
}

func NewState(dataDir, dataFormat string) *State {
	return & State{
		Entries: []Entry{},
		DataDir: dataDir,
		DataFormat: fmt.Sprintf(".%s", dataFormat),
	}
}

func (s *State) String() string {
	return fmt.Sprintf(
		"State\nEntries: %d\nDataDir: %s\nDataFormat: %s",
		len(s.Entries),
		s.DataDir,
		s.DataFormat,
	)
}

// This method reads the content of a CSV file in the Entry format and saves it to the state
func (s *State) LoadEntries(filename string) error {
	filePath := filepath.Join(s.DataDir, filename + s.DataFormat)
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
		
		log.Printf("reading record %d from %s\n", i, filePath)
		// Skip weird records
		if len(record) != 6 {
			log.Println("skipping weird row...")
			continue
		}
		
		applyDate, err := time.Parse(time.RFC3339Nano, record[4])
		if err != nil {
			log.Println("time parsing error, skipping record...")
			continue
		}

		contactReplied, err := strconv.ParseBool(record[5])
		if err != nil {
			log.Println("bool parsing error, skipping record...")
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
		log.Println("record appended")
	}

	s.Entries = entries
	log.Println("finished loading entries to state")
	
	return nil
}

// This method saves all the entries of the state into a CSV file with the filename specified.
func (s *State) SaveEntries(filename string) error {
	os.MkdirAll(s.DataDir, 0755)
	filePath := filepath.Join(s.DataDir, filename + s.DataFormat)
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

func (s *State) AddEntry(entry Entry) {
	s.Entries = append(s.Entries, entry)
}

