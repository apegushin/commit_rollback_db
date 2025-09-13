package db

const (
	noVal = ""
)

type Database struct {
	Name    string
	Txn     *txn
	Storage map[int]string
}

func NewDatabase(name string) *Database {
	return &Database{
		Name:    name,
		Storage: make(map[int]string),
	}
}

func (d *Database) Get(id int) string {
	if d.Txn != nil {
		if d.Txn.keysToRemove.Contains(id) {
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
	if d.Txn != nil {
		if d.Txn.keysToRemove.Contains(id) {
			d.Txn.keysToRemove.Remove(id)
		}
		d.Txn.keysToUpdate[id] = value
	} else {
		d.Storage[id] = value
	}
}

func (d *Database) DeleteByID(id int) {
	if d.Txn != nil {
		if _, ok := d.Txn.keysToUpdate[id]; ok {
			delete(d.Txn.keysToUpdate, id)
		}
		d.Txn.keysToRemove.Add(id)
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

}

type txn struct {
	keysToRemove *Set
	keysToUpdate map[int]string
}

func newTxn() *txn {
	return &txn{
		keysToRemove: NewSet(),
	}
}
