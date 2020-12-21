package track

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/mitchellh/go-homedir"
)

// Local file system time tracker.
type Local struct {
	LogLocation string
}

// Start adds a new entry to the log.
func (tracker *Local) Start(entry Entry) Entry {
	entries := tracker.LoadEntries()
	recent := MostRecentEntry(entries)
	if (recent != Entry{} && recent.InProgress()) {
		log.Fatal(fmt.Sprintf("The last item in the log is missing a finish time:\n %v", recent))
	}

	tracker.SaveEntries([]Entry{entry})
	return entry
}

// Finish adds an end time to the most recent entry in the log.
func (tracker *Local) Finish(entry Entry) Entry {
	entries := tracker.LoadEntries()
	if len(entries) == 0 {
		log.Fatal("There are no entries to update.")
	}

	recent := MostRecentEntry(entries)
	if !recent.FinishedAt.IsZero() {
		log.Fatal("This would overwrite the most recent entry: \n", recent)
	}

	if entry.Description != "" {
		recent.Description = entry.Description
	}

	duration, err := entry.GetDuration()
	if err != nil {
		log.Fatal(err)
	}

	recent.End(duration, entry.FinishedAt)

	fmt.Println("Updated entry in log at position: ", len(entries)-1)

	tracker.SaveEntries([]Entry{recent})
	return entry
}

// LoadEntries loads all Entries from a local file.
func (tracker *Local) LoadEntries() []Entry {
	logLocation := tracker.LogLocation
	if strings.Contains(logLocation, "~") {
		expanded, err := homedir.Expand(logLocation)
		if err != nil {
			log.Fatal(err)
		}
		logLocation = expanded
	}

	var entries []Entry

	exists, _ := Exists(logLocation)
	if exists {
		content, err := ioutil.ReadFile(logLocation)
		if err != nil {
			log.Fatal(err)
		}
		if err := json.Unmarshal(content, &entries); err != nil {
			log.Fatal(err)
		}
	}
	return entries
}

// SaveEntries saves a list of entries to a local file.
func (tracker *Local) SaveEntries(entries []Entry) []Entry {
	logLocation := tracker.LogLocation
	if strings.Contains(logLocation, "~") {
		expanded, err := homedir.Expand(logLocation)
		if err != nil {
			log.Fatal(err)
		}
		logLocation = expanded
	}
	current := tracker.LoadEntries()
	updated, err := UpdateEntries(current, entries, "ID")
	nextID := NextID(updated)
	for ix := range updated {
		if updated[ix].ID == 0 {
			updated[ix].ID = nextID
			nextID++
		}
	}
	if err != nil {
		log.Fatal(err)
	}
	SortEntries(updated)

	data, err := json.MarshalIndent(updated, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	if err := ioutil.WriteFile(logLocation, data, 0644); err != nil {
		log.Fatal(err)
	}

	return entries
}

// Exists checks if the file 'name' exists.
func Exists(name string) (bool, error) {
	_, err := os.Stat(name)
	if os.IsNotExist(err) {
		return false, nil
	}
	return err == nil, err
}
