package track

// Tracker defines the time tracker interface.
type Tracker interface {
	Start(entry Entry) Entry
	Finish(entry Entry) Entry
	LoadEntries() []Entry
	SaveEntries(entries []Entry) []Entry
}
