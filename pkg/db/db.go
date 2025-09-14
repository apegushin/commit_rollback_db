package db

import "fmt"

const (
	noVal = ""
)

type Database struct {
	Txn     *txn
	Storage map[int]string
}

func NewDatabase() *Database {
	return &Database{
		Storage: make(map[int]string),
	}
}

func (d *Database) Get(id int) string {
	if d.isTxnInProgress() {
		if d.isKeyDeleteRequested(id) {
			return noVal
		}
		if value, ok := d.Txn.keysToUpdate[id]; ok {
			return value
		}
	}
	if value, ok := d.Storage[id]; ok {
		return value
	}
	return noVal
}

func (d *Database) Set(id int, value string) {
	if d.isTxnInProgress() {
		if d.isKeyDeleteRequested(id) {
			d.Txn.keysToDelete.Remove(id)
		}
		d.Txn.keysToUpdate[id] = value
	} else {
		d.Storage[id] = value
	}
}

func (d *Database) DeleteByID(id int) {
	if d.isTxnInProgress() {
		if _, ok := d.Txn.keysToUpdate[id]; ok {
			delete(d.Txn.keysToUpdate, id)
		}
		d.Txn.keysToDelete.Add(id)
	} else {
		if _, ok := d.Storage[id]; ok {
			delete(d.Storage, id)
		}
	}
}

func (d *Database) DeleteByValue(value string) {
	// anything that's in the current txn with the requested value gets removed
	// anything that's in storage with the requested value, if in txn the value
	// is not updated, the key is put on the txn.keysToDelete
	if d.isTxnInProgress() {
		for key, val := range d.Txn.keysToUpdate {
			if val == value {
				delete(d.Txn.keysToUpdate, key)
			}
		}
	}
	for key, val := range d.Storage {
		if val == value {
			if d.isTxnInProgress() {
				if updVal, ok := d.Txn.keysToUpdate[key]; !ok || updVal == value {
					d.Txn.keysToDelete.Add(key)
				}
			} else {
				delete(d.Storage, key)
			}
		}
	}
}

func (d *Database) Begin() error {
	if d.isTxnInProgress() {
		return fmt.Errorf("transaction already in progress, can not begin another one")
	}
	d.Txn = newTxn()
	return nil
}

func (d *Database) Commit() error {
	if !d.isTxnInProgress() {
		return fmt.Errorf("no transaction in progress, nothing to commit")
	}
	// delete keys from keysToDelete set
	for kr := range d.Txn.keysToDelete.Items() {
		delete(d.Storage, kr)
	}
	// update values in Storage with values from keysToUpdate
	for ku, vu := range d.Txn.keysToUpdate {
		d.Storage[ku] = vu
	}
	return nil
}

func (d *Database) Rollback() error {
	if !d.isTxnInProgress() {
		return fmt.Errorf("no transaction in progress, nothing to rollback")
	}
	d.Txn = nil
	return nil
}

func (d *Database) isTxnInProgress() bool {
	return d.Txn != nil
}

func (d *Database) isKeyDeleteRequested(id int) bool {
	return d.Txn.keysToDelete.Contains(id)
}

type txn struct {
	keysToDelete *Set
	keysToUpdate map[int]string
}

func newTxn() *txn {
	return &txn{
		keysToDelete: NewSet(),
		keysToUpdate: make(map[int]string),
	}
}
