package filedb

import (
	"io/ioutil"
	"fmt"
	"encoding/json"
	"crypto/md5"
	"time"
)

type DB struct {
	Host string
	Records map[string]map[string]string
}

func MakeHash() string {
	return fmt.Sprintf("%x", md5.Sum([]byte(time.Now().String())))
}

func (d *DB) ReadAll() []byte {
	f, err := ioutil.ReadFile(d.Host)
	json.Unmarshal(f, &d.Records)
	if err != nil {
		panic(err)
	}
	return f
}

func (d *DB) WriteDB() bool {
	b, _ := json.Marshal(d.Records)
	err := ioutil.WriteFile(d.Host, b, 0644)
	if err != nil {
		return false
	}
	return true
}

func (d *DB) InsertRecord(record map[string]string) string {
	f := d.ReadAll()
	json.Unmarshal(f, &d.Records)
 	newId := MakeHash()
	record["id"] = newId
	d.Records[newId] = record
	b, _ := json.Marshal(d.Records)
	err := ioutil.WriteFile(d.Host, b, 0644)
	if err != nil {
		return ""
	}
	d.WriteDB()
	return newId
}

func (d *DB) FindID(id string) map[string]string {
	f := d.ReadAll()
	json.Unmarshal(f, &d.Records)
	if val, ok := d.Records[id]; ok {
		return val
	}
	return map[string]string{}
}

func (d *DB) RemoveID(id string) bool {
	f := d.ReadAll()
	json.Unmarshal(f, &d.Records)
	if _, ok := d.Records[id]; ok {
		delete(d.Records, id)
		d.WriteDB()
		return true
	}
	return false
}

func (d *DB) UpdateID(id string, value map[string]string) bool {
	f := d.ReadAll()
	json.Unmarshal(f, &d.Records)
	if _, ok := d.Records[id]; ok {
		d.Records[id] = value
		d.WriteDB()
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
