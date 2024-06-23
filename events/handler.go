package events

var Events = make(map[string]*Event)

type Event struct {
	Name    string
	Handler func(string)
}
