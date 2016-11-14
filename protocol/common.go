/*
@author: foolbread
@time: 2016-11-14
@file:tohno/protocol/common.go
*/
package protocol

const (
	FILE_CONTENT_CREATE = "create"
	FILE_CONTENT_UPDATE = "update"
	FILE_CONTENT_DELETE = "delete"
)

const (
	FILE_NAME_CHANGE   = "filename"
	FILE_BACKUP_CHANGE = "backupcount"
)

const (
	SYNC_FILE_TYPE      = "syncfile"
	SYNC_ZOOKEEPER_TYPE = "synczk"
)

type SyncFilePair struct {
	IP   string `json:"ip"`
	File string `json:"file"`
}
