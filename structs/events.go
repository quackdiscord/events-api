package structs

type Event struct {
	Name    string
	Handler func(string)
}
