/*
@author: foolbread
@time: 2016-11-14
@file:tohno/storage/common.go
*/
package storage

import (
	"github.com/foolbread/fbcommon/golog"
	"github.com/syndtr/goleveldb/leveldb/util"
)

func InitStorage() {
	golog.Info("tohno storage initing......")
}

var g_leveldb *tohnoLevelDB

////////////////////////////////////////////////////////////////
type DataPair struct {
	Key  []byte
	Data []byte
}

////////////////////////////////////////////////////////////////
func GetDataByPrefix(prefix []byte) ([]*DataPair, error) {
	iter := g_leveldb.db.NewIterator(util.BytesPrefix(prefix), nil)
	var ret []*DataPair
	for iter.Next() {
		ret = append(ret, &DataPair{iter.Key(), iter.Value()})
	}

	iter.Release()

	return ret, iter.Error()
}

func GetData(key []byte) ([]byte, error) {
	return g_leveldb.db.Get(key, nil)
}

func PutData(key []byte, data []byte) error {
	return g_leveldb.db.Put(key, data, nil)
}

func DeleteData(key []byte) error {
	return g_leveldb.db.Delete(key, nil)
}

func ExistData(key []byte) bool {
	ret, _ := g_leveldb.db.Has(key, nil)
	return ret
}
