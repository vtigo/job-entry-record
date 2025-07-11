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

type Entry struct {
	Company        string
	Role       	   string
	Status         string
	Platform       string
	ApplyDate 	   time.Time
	ContactReplied bool
}

func LoadEntries(filename string) ([]Entry, error) {
	f, err := os.Open(filepath.Join("data", filename + ".csv"))
	if err != nil {
		return nil, err
	}
	defer f.Close()
	
	reader := csv.NewReader(f)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	entries := []Entry{}
	
	// Skip header row
	for i := 1; i < len(records); i++ {
		r := records[i]
		
		log.Printf("reading record %d from file %s...\n", i, filename)
		// Skip weird records
		if len(r) != 6 {
			log.Println("skipping weird row")
			continue
		}
		
		applyDate, err := time.Parse(time.RFC3339Nano, r[4])
		if err != nil {
			log.Println("time parsing error, skipping record")
			continue
		}

		contactReplied, err := strconv.ParseBool(r[5])
		if err != nil {
			log.Println("bool parsing error, skipping record")
			continue
		}

		entry := Entry{
			Company: r[0],
			Role: r[1],
			Status: r[2],
			Platform: r[3],
			ApplyDate: applyDate,
			ContactReplied: contactReplied,
		}
		entries = append(entries, entry)
		log.Println("record appended")
	}

	return entries, nil
}

// This method saves all the entries passed to it in a file with the filename specified.
func SaveEntries(filename string, entries []Entry) error {
	os.MkdirAll("data", 0755)
	f, err := os.Create(filepath.Join("data", filename + ".csv"))
	if err != nil {
		return err
	}
	defer f.Close()

	writer := csv.NewWriter(f)
	defer writer.Flush()

	header := []string{"company", "position", "status", "platform", "apply_date", "contact_replied"}
	if err = writer.Write(header); err != nil {
		return err
	}

	for _, e := range entries {
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

func main() {
	e1 := Entry{
		Company: "emp",
		Role: "dev",
		Status: "paz",
		Platform: "web",
		ApplyDate: time.Now(),
		ContactReplied: false,
	}

	SaveEntries("test", []Entry{e1})
	fmt.Println(LoadEntries("test"))
}
