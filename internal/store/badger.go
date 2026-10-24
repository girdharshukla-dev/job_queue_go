package store

import (
	"encoding/json"
	"fmt"
	"github.com/dgraph-io/badger/v4"
	"go_job_queue/internal/job"
)

type BadgerStore struct {
	db *badger.DB
}

func NewBadgerStore(path string) (*BadgerStore, error) {
	badgerdb, err := badger.Open(badger.DefaultOptions(path))
	if err != nil {
		fmt.Println("Error in opening badger db")
		return nil, err
	}
	return &BadgerStore{
		db: badgerdb,
	}, nil
}

func (b *BadgerStore) Get(id string) (job.Job, error) {

	var j job.Job

	err := b.db.View(func(txn *badger.Txn) error {
		item, err2 := txn.Get([]byte(id))
		if err2 != nil {
			return err2
		}

		val, err2 := item.ValueCopy(nil)
		if err2 != nil {
			return err2
		}

		err2 = json.Unmarshal(val, &j)
		if err2 != nil {
			return err2
		}
		return nil
	})

	if err != nil {
		return job.Job{}, err
	}

	return j, nil
}

func (b *BadgerStore) ListByStatus(status job.Status) ([]job.Job, error) {
	var jobs []job.Job

	err := b.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10

		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			err2 := item.Value(func(val []byte) error {
				j, err3 := job.Decode(val)
				if err3 != nil {
					return err3
				}
				if j.Status == status {
					jobs = append(jobs, j)
				}
				return nil
			})

			if err2 != nil {
				return err2
			}
		}
		return nil
	})

	if err != nil {
		return []job.Job{}, err
	}
	return jobs, nil
}

func (b *BadgerStore) Save(j job.Job) error {

	j_bytes, err := job.Encode(j)
	if err != nil {
		return err
	}

	err = b.db.Update(func(txn *badger.Txn) error {
		err2 := txn.Set([]byte(j.ID), j_bytes)
		if err2 != nil {
			return err2
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (b *BadgerStore) Close() error {
	return b.db.Close()
}
