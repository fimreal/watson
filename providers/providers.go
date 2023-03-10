package providers

import (
	"time"
)

type DNSProvider interface {
	GetRecordList(zone string) ([]Record, error)
	GetRecord(zone string, record Record) (Record, error)
	SetRecord(zone string, record Record) error
	DeleteRecord(zone string, record Record) error
}

type Record struct {
	ID    string
	Type  string
	Name  string
	Value string
	TTL   time.Duration
}
