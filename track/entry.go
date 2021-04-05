package track

import (
	"encoding/json"
	"fmt"
	"sort"
	"time"
)

// Entry contains information for a period of time
type Entry struct {
	ID          int       `json:"id"`
	StartedAt   time.Time `json:"started_at,omitempty"`
	FinishedAt  time.Time `json:"finished_at,omitempty"`
	Duration    int       `json:"duration,omitempty"` // Seconds
	Description string    `json:"description"`
	ClientID    int       `json:"client_id"`
	ProjectID   int       `json:"project_id,omitempty"`
	ExternalID  int       `json:"external_id,omitempty"`
}

// GetDuration converts Duration into a time.Duration object.
func (entry *Entry) GetDuration() (time.Duration, error) {
	return time.ParseDuration(fmt.Sprintf("%d", entry.Duration) + "s")
}

// End derives the duration and finished at times.
func (entry *Entry) End(duration time.Duration, finishedAt time.Time) {
	if duration.Seconds() > 0 {
		entry.FinishedAt = entry.StartedAt.Add(duration)
		entry.Duration = int(duration.Seconds())
	} else if !finishedAt.IsZero() {
		entry.FinishedAt = finishedAt
		duration := entry.FinishedAt.Sub(entry.StartedAt)
		entry.Duration = int(duration.Seconds())
	}
}

// InProgress returns true if the entry has not been completed.
func (entry *Entry) InProgress() bool {
	return entry.FinishedAt == time.Time{}
}

// JSON returns Entry as JSON object with indent.
func (entry *Entry) JSON() ([]byte, error) {
	return json.MarshalIndent(entry, "", "  ")
}

// String returns a string representation of the Entry object.
func (entry *Entry) String() string {
	s := entry.StartedAt.Local().Format(time.UnixDate)
	f := "In progress"
	if !entry.FinishedAt.IsZero() {
		f = entry.FinishedAt.Local().Format(time.UnixDate)
	}
	var id string
	if entry.ExternalID > 0 {
		id = fmt.Sprintf("ID: %v (External ID: %v)", entry.ID, entry.ExternalID)
	} else {
		id = fmt.Sprintf("ID: %v", entry.ID)
	}
	d, _ := entry.GetDuration()
	return fmt.Sprintf("Description: %s\nStarted At: %s\nFinished At: %s\nDuration: %v\n%s\nClient ID: %d", entry.Description, s, f, d.Round(time.Minute), id, entry.ClientID)
}

// MostRecentEntry returns the most recent entry if the entries slice
// is not empty.
func MostRecentEntry(entries []Entry) Entry {
	recent := Entry{}
	for _, entry := range entries {
		if entry.StartedAt.After(recent.StartedAt) {
			recent = entry
		}
	}
	return recent
}

// NextID determines the next ID from a list of entries.
func NextID(entries []Entry) int {
	nextID := 0
	for _, entry := range entries {
		if entry.ID > nextID {
			nextID = entry.ID
		}
	}
	return nextID + 1
}

// FilterParameters declares the parameters for FilterEntries.
type FilterParameters struct {
	Since       time.Time
	Until       time.Time
	Description string
	Limit       int
}

// FilterEntries searches by start, end, and description.
func FilterEntries(entries []Entry, params FilterParameters) []Entry {
	res := make([]Entry, len(entries))
	copy(res, entries)
	SortEntries(res)
	if !params.Since.IsZero() {
		i := sort.Search(len(res), func(i int) bool { return res[i].StartedAt.Sub(params.Since).Seconds() >= 0 })
		if i < len(entries) {
			res = res[i:]
		} else {
			// No entries with started at less than params.Since
			return []Entry{}
		}
	}
	if !params.Until.IsZero() {
		i := sort.Search(len(res), func(i int) bool { return res[i].StartedAt.Sub(params.Until).Seconds() >= 0 })
		if i > 0 {
			res = res[0:i]
		} else {
			// No entries with date greater than params.Until.
			return []Entry{}
		}
	}

	// if (params.Description != "") {
	// 	TODO
	// }
	if len(res) > params.Limit && params.Limit > 0 {
		res = res[len(res)-params.Limit:]
	}
	return res
}

// SortEntries using the StartedAt attribute.
func SortEntries(entries []Entry) {
	sort.SliceStable(entries, func(i int, j int) bool { return entries[i].StartedAt.Before(entries[j].StartedAt) })
}

// UpdateEntries merges entries in left with entries in right or add
// new entries.
func UpdateEntries(left []Entry, right []Entry, on string) ([]Entry, error) {
	if on != "ID" && on != "ExternalID" {
		return []Entry{}, fmt.Errorf("Update on must be ID or ExternalID. Got %s", on)
	}

	type lookup struct {
		Index int
		Entry Entry
	}
	index := map[int]lookup{}
	result := []Entry{}
	for ix, entry := range left {
		var key int
		if on == "ID" {
			key = entry.ID
		} else {
			key = entry.ExternalID
		}
		if _, exists := index[key]; exists {
			return []Entry{}, fmt.Errorf("Entries has duplicate %s: %v", on, key)
		}

		index[key] = lookup{Index: ix, Entry: entry}
		result = append(result, entry)
	}

	newEntries := []Entry{}
	for _, entry := range right {
		var key int
		if on == "ID" {
			key = entry.ID
		} else {
			key = entry.ExternalID
		}

		if val, exists := index[key]; exists {
			// ensure ID isn't lost if matching on ExternalID
			entry.ID = val.Entry.ID
			result[val.Index] = entry
		} else {
			newEntries = append(newEntries, entry)
		}
	}

	return append(result, newEntries...), nil
}
