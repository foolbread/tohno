/*
@author: foolbread
@time: 2016-11-14
@file:tohno/storage/leveldb.go
*/
package storage

import (
	"github.com/syndtr/goleveldb/leveldb"
)

type tohnoLevelDB struct {
	dbFile string
	db     *leveldb.DB
}

func newTohnoLevelDB(dbFile string) (*tohnoLevelDB, error) {
	r := new(tohnoLevelDB)
	r.dbFile = dbFile
	var err error
	r.db, err = leveldb.OpenFile(dbFile, nil)
	if err != nil {
		return nil, err
	}

	return r, nil
}
