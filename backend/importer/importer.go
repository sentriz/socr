package importer

import (
	"fmt"
	"io/ioutil"
	"sync/atomic"
	"time"

	"go.senan.xyz/socr/db"
)

type Importer struct {
	isRunning *int32
	DB        db.DB
	Dirs      []string
}

type StatusError struct {
	Time  time.Time `json:"time"`
	Error string    `json:"error"`
}

type Status struct {
	Errors         []StatusError `json:"errors,omitempty"`
	LastID         string        `json:"last_id,omitempty"`
	CountProcessed int           `json:"count_processed"`
	CountTotal     int           `json:"count_total"`
}

func (s *Status) AddError(err error) {
	s.Errors = append(s.Errors, StatusError{
		Time:  time.Now(),
		Error: err.Error(),
	})
}

func (i *Importer) IsRunning() bool { return atomic.LoadInt32(i.isRunning) == 1 }
func (i *Importer) SetRunning()     { atomic.StoreInt32(i.isRunning, 1) }
func (i *Importer) SetFinished()    { atomic.StoreInt32(i.isRunning, 0) }

func (i *Importer) Scan() error {
	for _, dir := range i.Dirs {
		files, err := ioutil.ReadDir(dir)
		if err != nil {
			return fmt.Errorf("listing dir: %w", err)
		}

		for _, file := range files {
			timeLast, err := i.DB.GetModTime(dir, file.Name())
			if err != nil {
				return fmt.Errorf("get last mod time: %v", err)
			}

			if timeLast == nil {

			}

			// bytes, err := ioutil.ReadFile(file.Name())
			// if err != nil {
			// 	return fmt.Errorf("reading from disk: %v", err)
			// }
		}

	}
	return nil
}
