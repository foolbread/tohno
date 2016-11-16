/*
@author: foolbread
@time: 2016-11-14
@file:tohno/protocol/common.go
*/
package protocol

const (
	FILE_CREATE     = "create"
	FILE_UPDATE     = "update"
	FILE_DELETE     = "delete"
	FILE_GET        = "get"
	BACKUP_FILE_GET = "backupget"
)

const (
	DIR_CREATE = "create"
	DIR_DEL    = "delete"
	DIR_RENAME = "rename"
	DIR_SCAN   = "scan"
)

const (
	ERROR_STATUS = -1
	OK_STATUS    = 0
)

const (
	SYNC_FILE_TYPE      = "syncfile"
	SYNC_ZOOKEEPER_TYPE = "synczk"
)

type SyncFilePair struct {
	IP   string `json:"ip"`
	File string `json:"file"`
}

type DirInfoPair struct {
	Name  string `json:"name"`
	IsDir bool   `json:"isdir"`
}
