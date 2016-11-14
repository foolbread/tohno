/*
@author: foolbread
@time: 2016-11-14
@file:tohno/protocol/protocol.go
*/
package protocol

type FileContentCreate struct {
	File        string   `json:"file"`
	Content     string   `json:"content"`
	BackupCount int      `json:"backupcount"`
	SyncTypes   []string `json:"synctypes"`
	SyncInfos   []string `json:"syncinfos"`
}

type FileContentUpdate struct {
	File    string `json:"file"`
	Content string `json:"content"`
}

type FileContentDelete struct {
	File string `json:"file"`
}

type FileNameChange struct {
	NewName string `json:"newname"`
	OldName string `json:"oldname"`
}

type FileBackupChange struct {
	File        string `json:"file"`
	BackupCount int    `json:"backupcount"`
}

type SyncFileType struct {
	Infos []*SyncFilePair `json:"infos"`
}
