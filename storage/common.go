/*
@author: foolbread
@time: 2016-11-14
@file:tohno/storage/common.go
*/
package storage

import (
	"github.com/foolbread/fbcommon/golog"
)

func InitStorage() {
	golog.Info("tohno storage initing......")
}

var g_leveldb *tohnoLevelDB

const (
	sync_file_tail      = "-file"
	sync_zookeeper_tail = "-zk"
)

func ExistFileInfo(file string) bool {
	ret, _ := g_leveldb.db.Has([]byte(file), nil)

	return ret
}

func ExistSyncFileInfo(file string) bool {
	ret, _ := g_leveldb.db.Has([]byte(file+sync_file_tail), nil)

	return ret
}

////////////////////////////////////////////////////////////////
func GetFileInfo(file string) ([]byte, error) {
	return g_leveldb.db.Get([]byte(file), nil)
}

func GetSyncFileInfo(file string) ([]byte, error) {
	return g_leveldb.db.Get([]byte(file+sync_file_tail), nil)
}

////////////////////////////////////////////////////////////////

func PutFileInfo(file string, data []byte) error {
	return g_leveldb.db.Put([]byte(file), data, nil)
}

func PutSyncFileInfo(file string, data []byte) error {
	return g_leveldb.db.Put([]byte(file+sync_zookeeper_tail), data, nil)
}
