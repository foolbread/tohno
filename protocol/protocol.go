/*
@author: foolbread
@time: 2016-11-14
@file:tohno/protocol/protocol.go
*/
package protocol

type FileCreateReq struct {
	File        string       `json:"file"`
	Content     string       `json:"content"`
	BackupCount int          `json:"backupcount"`
	SyncFile    SyncFileType `json:"syncfile"`
}

type FileUpdateReq struct {
	File        string       `json:"file"`
	Content     string       `json:"content"`
	BackupCount int          `json:"backupcount"`
	SyncFile    SyncFileType `json:"syncfile"`
}

type FileGetReq struct {
	File string `json:"file"`
}

type FileGetRes struct {
	File        string       `json:"file"`
	Content     string       `json:"content"`
	BackupCount int          `json:"backupcount"`
	SyncFile    SyncFileType `json:"syncfile"`
}

type FileDeleteReq struct {
	File string `json:"file"`
}

type FileRenameReq struct {
	ParentPath string `json:"parentdir"`
	NewName    string `json:"newname"`
	OldName    string `json:"oldname"`
}

type DirCreateReq struct {
	Dir string `json:"dir"`
}

type DirDelReq struct {
	Dir string `json:"dir"`
}

type DirRenameReq struct {
	ParentPath string `json:"parentdir"`
	NewName    string `json:"newname"`
	OldName    string `json:"oldname"`
}

type DirScanReq struct {
	Dir string `json:"dir"`
}

type SyncFileType struct {
	Infos []*SyncFilePair `json:"infos"`
}

type CommonRes struct {
	Info string `json:"info"`
}

type DirInfoRes struct {
	Dir   string `json:"dir"`
	Infos []*DirInfoPair
}
