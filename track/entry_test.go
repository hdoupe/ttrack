package track

import (
	"log"
	"testing"
	"time"
)

func timePair(ts string, ds string) (time.Time, time.Duration) {
	utc, _ := time.LoadLocation("")
	t, err := time.ParseInLocation("2006-01-02 3:04:05 PM", ts, utc)
	if err != nil {
		log.Fatal(err)
	}
	d, err := time.ParseDuration("2h")
	if err != nil {
		log.Fatal(err)
	}
	return t, d
}

func mockEntries() []Entry {
	time1, dur1 := timePair("2020-11-21 10:00:00 AM", "2h")
	time2, dur2 := timePair("2020-11-21 1:00:00 PM", "30m")
	time3, dur3 := timePair("2020-11-21 3:00:00 PM", "45m")
	time4, dur4 := timePair("2020-11-21 5:00:00 PM", "15m")
	return []Entry{
		{
			ID:          0,
			StartedAt:   time1,
			FinishedAt:  time1.Add(dur1),
			Duration:    int(dur1.Seconds()),
			Description: "Write some code",
			ExternalID:  123,
		},
		{
			ID:          1,
			StartedAt:   time2,
			FinishedAt:  time2.Add(dur2),
			Duration:    int(dur2.Seconds()),
			Description: "Write some tests",
			ExternalID:  456,
		},
		{
			ID:          2,
			StartedAt:   time3,
			FinishedAt:  time3.Add(dur3),
			Duration:    int(dur3.Seconds()),
			Description: "Write some docs",
			ExternalID:  789,
		},
		{
			ID:          3,
			StartedAt:   time4,
			FinishedAt:  time4.Add(dur4),
			Duration:    int(dur3.Seconds()),
			Description: "Oops fix some bugs whilst drinking a beer",
			ExternalID:  257,
		},
	}
}

func TestAddEntries(t *testing.T) {
	entries := mockEntries()
	left := entries[0:2]
	right := entries[2:4]

	for _, on := range []string{"ID", "ExternalID"} {
		result, err := UpdateEntries(left, right, on)

		if err != nil {
			t.Error(err)
		}

		if len(result) != 4 {
			t.Errorf("(%s) Length of entries is not 4, got %v", on, len(result))
		}
		if len(left) != 2 {
			t.Errorf("(%s) Length of left entries is not 2, got %v", on, len(left))
		}
		if len(right) != 2 {
			t.Errorf("(%s) Length of right entries is not 2, got %v", on, len(right))
		}
	}
}

func TestUpdateEntries(t *testing.T) {
	entries := mockEntries()

	newEntry := entries[2]
	newEntry.Description = "Ugh, write some docs"

	if newEntry.Description == entries[2].Description {
		t.Errorf("These shouldn't be equal")
	}

	result, err := UpdateEntries(entries, []Entry{newEntry}, "ID")

	if err != nil {
		t.Error(err)
	}

	if len(result) != 4 {
		t.Errorf("Length of entries is not 4, got %v", len(result))
	}

	if result[2].Description != "Ugh, write some docs" {
		t.Errorf("Description is: %s", result[2].Description)
	}
}
