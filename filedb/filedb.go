package filedb

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sync"
	"time"
)

type Records map[string]Record
type Record map[string]string

type DB struct {
	Host    string
	records Records
	sync.Mutex
}

func MakeHash() string {
	return fmt.Sprintf("%x", md5.Sum([]byte(time.Now().String())))
}

func (d *DB) readAll() {
	f, _ := ioutil.ReadFile(d.Host)
	json.Unmarshal(f, &d.records)
	if d.records == nil {
		d.records = make(Records)
	}
}

func (d *DB) writeDB() error {
	b, err := json.Marshal(d.records)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(d.Host, b, 0644)
}

func (d *DB) InsertRecord(record Record) string {
	d.Lock()
	defer d.Unlock()
	d.readAll()
	newID := MakeHash()
	record["id"] = newID
	d.records[newID] = record
	d.writeDB()
	return newID
}

// All returns underlying records.
// NOTE: A copy of d.records should be returned to prevent caller from modifying
// underlying map
func (d *DB) All() Records {
	return d.records
}

func (d *DB) FindID(id string) Record {
	d.Lock()
	defer d.Unlock()
	d.readAll()
	return d.records[id]
}

func (d *DB) RemoveID(id string) bool {
	d.Lock()
	defer d.Unlock()
	d.readAll()
	if _, ok := d.records[id]; ok {
		delete(d.records, id)
		d.writeDB()
		return true
	}
	return false
}

func (d *DB) UpdateID(id string, value map[string]string) bool {
	d.Lock()
	defer d.Unlock()
	d.readAll()
	if _, ok := d.records[id]; ok {
		d.records[id] = value
		d.writeDB()
		return true
	}
	return false
}

//
//func main() {
//	d := DB{Host: "data/db.json"}
//	m := map[string]string{"name": "Steven Andrew"}
//	newId := d.InsertRecord(m)
//	e := d.FindID(newId)
//	fmt.Println(e)
//	e["name"] = "for reals thought"
//	d.UpdateID(newId, e)
//	d.RemoveID(newId)
//}
