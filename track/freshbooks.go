package track

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/hdoupe/ttrack/oauth"
)

// TimeEntry represents the FreshBooks Time Entry object.
// This is analagous to the Track object.
type TimeEntry struct {
	Note      string    `json:"note"`
	Duration  int       `json:"duration"`
	ClientID  int       `json:"client_id,omitempty"`
	ProjectID int       `json:"project_id,omitempty"`
	IsLogged  bool      `json:"is_logged"`
	StartedAt time.Time `json:"started_at"`
	Active    bool      `json:"active"`
	ID        int       `json:"id"`
}

// TimeEntryPayload is the data structure for data posted to freshbooks.com.
type TimeEntryPayload struct {
	TimeEntry TimeEntry `json:"time_entry"`
}

// ToEntry converts a TimeEntry to an Entry.
func (timeEntry *TimeEntry) ToEntry() Entry {
	finishedAt := time.Time{}
	duration, err := time.ParseDuration(fmt.Sprintf("%d", timeEntry.Duration) + "s")
	if err != nil {
		log.Fatal("Unable to convert timeEntry second to duration", timeEntry.Duration)
	}
	if !timeEntry.StartedAt.IsZero() && duration.Seconds() > 0 {
		finishedAt = timeEntry.StartedAt.Add(duration)
	}
	return Entry{
		StartedAt:   timeEntry.StartedAt,
		FinishedAt:  finishedAt,
		Duration:    int(duration.Seconds()),
		Description: timeEntry.Note,
		ClientID:    timeEntry.ClientID,
		ProjectID:   timeEntry.ProjectID,
		ExternalID:  timeEntry.ID,
	}
}

func (entry *Entry) toTimeEntry() TimeEntry {
	duration, err := entry.GetDuration()
	if err != nil {
		log.Fatal(err)
	}
	// a couple defaults until this is handled better.
	// if entry.ClientID == 0 {
	// 	entry.ClientID = 67837
	// }
	// if entry.ProjectID == 0 {
	// 	entry.ProjectID = 5476277
	// }
	return TimeEntry{
		Active:    true,
		StartedAt: entry.StartedAt,
		Duration:  int(duration.Seconds()),
		Note:      entry.Description,
		ClientID:  entry.ClientID,
		ProjectID: entry.ProjectID,
		ID:        entry.ExternalID,
		IsLogged:  true,
	}
}

// FreshBooks integrates ttrack and freshbooks.com.
type FreshBooks struct {
	LogLocation string
	Credentials oauth.Credentials
}

// Start entry on FreshBooks.
func (tracker *FreshBooks) Start(entry Entry) Entry {
	entries := tracker.LoadEntries()

	recent := MostRecentEntry(entries)
	if (recent != Entry{} && recent.InProgress()) {
		log.Fatal(fmt.Sprintf("The last item in the log is missing a finish time:\n %v", recent))
	}

	entry = tracker.CreateEntry(entry)

	local := Local{LogLocation: tracker.LogLocation}
	local.SaveEntries([]Entry{entry})
	return entry
}

// Finish entry on FreshBooks.
func (tracker *FreshBooks) Finish(entry Entry) Entry {
	entries := tracker.LoadEntries()
	if len(entries) == 0 {
		log.Fatal("There are no entries to update.")
	}

	recent := MostRecentEntry(entries)

	if (recent.FinishedAt != time.Time{}) {
		data, err := json.MarshalIndent(recent, "", "  ")
		if err != nil {
			log.Fatal(err)
		}
		log.Fatal("This would overwrite the most recent entry", string(data))
	}

	if entry.Description != "" {
		recent.Description = entry.Description
	}

	duration, err := entry.GetDuration()
	if err != nil {
		log.Fatal(err)
	}

	recent.End(duration, entry.FinishedAt)

	recent = tracker.UpdateEntry(recent)
	local := Local{LogLocation: tracker.LogLocation}
	local.SaveEntries([]Entry{recent})

	return recent
}

// LoadEntries loads all entries from freshbooks. The entries are synced
// with the local entries using their ExternalID.
func (tracker *FreshBooks) LoadEntries() []Entry {
	businessID := RetrieveBusinessID(tracker.Credentials)
	timeEntries := RetrieveTimeEntries(businessID, tracker.Credentials)

	entries := []Entry{}
	for _, timeEntry := range timeEntries {
		entries = append(entries, timeEntry.ToEntry())
	}

	local := Local{LogLocation: tracker.LogLocation}
	locEntries := local.LoadEntries()

	entries, err := UpdateEntries(locEntries, entries, "ExternalID")
	if err != nil {
		log.Fatal(err)
	}

	SortEntries(entries)

	return entries
}

// SaveEntries creates new entries and updates existing entries to FreshBooks.
func (tracker *FreshBooks) SaveEntries(entries []Entry) []Entry {
	currEntries := tracker.LoadEntries()

	res := []Entry{}
	for _, entry := range entries {
		if entry.ExternalID == 0 {
			res = append(res, tracker.CreateEntry(entry))
		} else {
			for _, curr := range currEntries {
				if curr.ExternalID == entry.ExternalID && curr != entry {
					res = append(res, tracker.UpdateEntry(entry))
				}
			}
		}
	}

	return res
}

// CreateEntry saves just one Entry to FreshBooks.
func (tracker *FreshBooks) CreateEntry(entry Entry) Entry {
	businessID := RetrieveBusinessID(tracker.Credentials)
	url := fmt.Sprintf("https://api.freshbooks.com/timetracking/business/%s/time_entries", fmt.Sprint(businessID))

	timeEntry := TimeEntryPayload{TimeEntry: entry.toTimeEntry()}
	payload, err := json.Marshal(timeEntry)
	if err != nil {
		log.Fatal(err)
	}
	req, err := http.NewRequest("POST", url, bytes.NewReader(payload))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Authorization", "Bearer "+tracker.Credentials.AccessToken)
	req.Header.Add("API-Version", "alpha")
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		log.Fatal("Unexpected error when creating time entry (", resp.StatusCode, ")", string(body[:]))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	type TimeEntryResponse struct {
		TimeEntry struct {
			ID int `json:"ID"`
		} `json:"time_entry"`
	}

	var data TimeEntryResponse
	if err := json.Unmarshal(body, &data); err != nil {
		fmt.Println("Unable to parse response from FreshBooks:", string(body[:]))
	}
	entry.ExternalID = data.TimeEntry.ID
	return entry
}

// UpdateEntry updates an entry on freshbooks.com
func (tracker *FreshBooks) UpdateEntry(entry Entry) Entry {
	if entry.ExternalID == 0 {
		log.Fatal("Unable to update entry because the external id is not defined.", entry)
	}
	businessID := RetrieveBusinessID(tracker.Credentials)
	url := fmt.Sprintf(
		"https://api.freshbooks.com/timetracking/business/%s/time_entries/%s",
		fmt.Sprint(businessID),
		fmt.Sprint(entry.ExternalID),
	)

	timeEntry := TimeEntryPayload{TimeEntry: entry.toTimeEntry()}
	payload, err := json.Marshal(timeEntry)
	if err != nil {
		log.Fatal(err)
	}
	req, err := http.NewRequest("PUT", url, bytes.NewReader(payload))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Authorization", "Bearer "+tracker.Credentials.AccessToken)
	req.Header.Add("API-Version", "alpha")
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		log.Fatal("Unexpected error when updating time entry (", resp.StatusCode, ")", string(body[:]))
	}
	return entry
}

// SyncEntries on FreshBooks with a local file.
func (tracker *FreshBooks) SyncEntries() {
	entries := tracker.LoadEntries()
	local := Local{LogLocation: tracker.LogLocation}
	local.SaveEntries(entries)
}

// RetrieveTimeEntries returns a list of time entries from FreshBooks.
func RetrieveTimeEntries(businessID int, credentials oauth.Credentials) []TimeEntry {
	var result []TimeEntry
	timeEntries, hasMore := RetrieveTimeEntriesPage(businessID, credentials, 0)
	result = append(result, timeEntries...)
	page := 1
	for hasMore {
		timeEntries, hasMore = RetrieveTimeEntriesPage(businessID, credentials, page)
		result = append(result, timeEntries...)
		page++
	}
	return result
}

// RetrieveTimeEntriesPage returns a page of time entries from FreshBooks.
func RetrieveTimeEntriesPage(businessID int, credentials oauth.Credentials, page int) ([]TimeEntry, bool) {
	url := fmt.Sprintf("https://api.freshbooks.com/timetracking/business/%s/time_entries", fmt.Sprint(businessID))
	req, err := http.NewRequest("GET", url, bytes.NewReader([]byte{}))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Authorization", "Bearer "+credentials.AccessToken)
	req.Header.Add("API-Version", "alpha")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		log.Fatal("Unexpected error when retrieving page of time entries (", resp.StatusCode, ")", string(body[:]))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	var timeEntries struct {
		TimeEntries []TimeEntry `json:"time_entries"`
		Meta        struct {
			Pages int `json:"pages"`
			Page  int `json:"page"`
		}
	}
	if err := json.Unmarshal(body, &timeEntries); err != nil {
		log.Fatal(err)
	}
	return timeEntries.TimeEntries, timeEntries.Meta.Page < timeEntries.Meta.Pages
}

// Me is the data from the Me response that is necessary to use ttrack.
type Me struct {
	Response struct {
		ID                  int `json:"id"`
		BusinessMemberships []struct {
			Business struct {
				ID int `json:"id"`
			} `json:"business"`
		} `json:"business_memberships"`
	} `json:"response"`
}

// RetrieveBusinessID gets the user's business ID to be used for the
// time tracking API calls.
func RetrieveBusinessID(credentials oauth.Credentials) int {
	url := "https://api.freshbooks.com/auth/api/v1/users/me"

	req, err := http.NewRequest("GET", url, bytes.NewReader([]byte{}))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Authorization", "Bearer "+credentials.AccessToken)
	req.Header.Add("API-Version", "alpha")
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Println("Unexpected error when getting user identity (", resp.StatusCode, ")")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var data Me
	if err := json.Unmarshal(body, &data); err != nil {
		panic(err)
	}

	if len(data.Response.BusinessMemberships) == 0 {
		log.Fatal("One business membership is required.")
	} else if len(data.Response.BusinessMemberships) > 1 {
		log.Fatal("TODO: Let user select business membership.")
	}

	return data.Response.BusinessMemberships[0].Business.ID
}
