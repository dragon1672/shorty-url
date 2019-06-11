package db


type Status struct {
	Success bool
	ErrorText string
}

type Database interface  {
	Init()
	Get(entry string) (result string, exists bool)
	AddMapping(from, to string) error
	RemoveEntry(entry string) error
}

