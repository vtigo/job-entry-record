package main

import (
	"encoding/csv"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type Storager interface {
	Add(Entry) error
	Update(id int) error
	Delete(id int) error
	GetAll() ([]Entry, error)
	GetOne(id int) (Entry, error)
}

type CsvStorage struct {
	Path string
	Filename string
}

func NewCsvStorage(path, filename string) *CsvStorage {
	return &CsvStorage{
		Path: path,
		Filename: filename,
	}
}

func (s *CsvStorage) Add(entry Entry) error {
	filePath := filepath.Join(s.Path, s.Filename + ".csv")
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

	//TODO: add new entry to the records and save the file
	log.Println(records)
	
	return nil
}

func(s *CsvStorage) Update(id int) error {
	return nil
}

func (s *CsvStorage) Delete(id int) error {
	return nil
}

func (s *CsvStorage) GetAll() ([]Entry, error) {
	filePath := filepath.Join(s.Path, s.Filename + ".csv")
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	entries := []Entry{}
	
	// Skip header row
	for i := 1; i < len(records); i++ {
		record := records[i]
		
		// Skip weird records
		if len(record) != 6 {
			log.Printf("Skipped record #%d for unexpected format", i)
			continue
		}
		
		applyDate, err := time.Parse(time.RFC3339Nano, record[4])
		if err != nil {
			log.Printf("Skipped record #%d. Failed to parse apply date", i)
			continue
		}

		contactReplied, err := strconv.ParseBool(record[5])
		if err != nil {
			log.Printf("Skipped record #%d. Failed to parse contact replied", i)
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
	
	return entries, nil
}
