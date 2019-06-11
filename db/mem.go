package db

import (
	"github.com/pkg/errors"
)

type inMemDb struct {
	m map[string]string
}

func GetInMemDb() Database {
	return &inMemDb{}
}

func (db *inMemDb) Init() {
	db.m = make(map[string]string)
}

func (db *inMemDb) Get(entry string) (result string, exists bool) {
	result, exists = db.m[entry]
	return result, exists

}

func (db *inMemDb) AddMapping(from, to string) error {
	if _, exists := db.m[from]; exists {
		return errors.Errorf("%s already has a mapping", from)
	}
	db.m[from] = to
	return nil
}
func (db *inMemDb) RemoveEntry(entry string) error {
	delete(db.m, entry)
	return nil
}
