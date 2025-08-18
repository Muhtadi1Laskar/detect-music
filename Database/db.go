package database

import (
	"encoding/json"
	"log"

	"github.com/dgraph-io/badger/v4"
)

type DBEntry struct {
	SongID     string
	TimeOffset int
}


func OpenBadger(path string) *badger.DB {
	opts := badger.DefaultOptions(path)
	db, err := badger.Open(opts)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func AddFingerprint(db *badger.DB, hash string, entry DBEntry) error {
	return db.Update(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(hash))
		if err == badger.ErrKeyNotFound {
			buf, _ := json.Marshal([]DBEntry{entry})
			return txn.Set([]byte(hash), buf)
		}
		if err != nil {
			return err
		}

		var list []DBEntry
		err = item.Value(func(v []byte) error {
			return json.Unmarshal(v, &list)
		})
		if err != nil {
			return err
		}

		list = append(list, entry)
		buf, _ := json.Marshal(list)
		return txn.Set([]byte(hash), buf)
	})
}

func LookupFingerprint(db *badger.DB, hash string) ([]DBEntry, error) {
	var list []DBEntry
	err := db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(hash))
		if err != nil {
			return err
		}
		return item.Value(func(v []byte) error {
			return json.Unmarshal(v, &list)
		})
	})
	return list, err
}