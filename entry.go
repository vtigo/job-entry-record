package main

import (
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

func NewEntry(company, role, status, platform string, applyDate time.Time, replied bool) *Entry {
	return &Entry{
		Company: company,
		Role: role,
		Status: status,
		Platform: platform,
		ApplyDate: applyDate,
		ContactReplied: replied,
	}
}

func (e Entry) AddToState(s *State) {
	s.Entries = append(s.Entries, e)
}
