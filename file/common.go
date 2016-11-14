/*
@author: foolbread
@time: 2016-11-14
@file:tohno/file/common.go
*/
package file

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"

	"github.com/foolbread/fbcommon/golog"
	"github.com/foolbread/tohno/protocol"
	"github.com/foolbread/tohno/storage"
)

func InitFile() {
	golog.Info("tohno file initing......")
}

const (
	SYNC_TYPE_FILE = "file"
)

var file_dir string
var backup_dir string

///////////////////////////////////////////////////////////////////////////////////////////
func HandlerFileContent(action string, data []byte) {
	switch action {
	case protocol.FILE_CONTENT_CREATE:
		handlerCreateFileContent(data)
	case protocol.FILE_CONTENT_UPDATE:
		handlerUpdateFileContent(data)
	case protocol.FILE_CONTENT_DELETE:
	}
}

func HandlerFileInfo(action string, data []byte) {
	switch action {
	case protocol.FILE_NAME_CHANGE:
		handlerChangeFileName(data)
	case protocol.FILE_BACKUP_CHANGE:
		handlerChangeBackupCnt(data)
	}
}

func HandlerSyncInfo(action string, data []byte) {
	switch action {
	case protocol.SYNC_FILE_TYPE:
		handlerUpdateSyncFile(data)
	case protocol.SYNC_ZOOKEEPER_TYPE:
	}
}

/////////////////////////////////////////////////////////////////////////////////////////////////
func handlerUpdateFileContent(data []byte) error {
	var req protocol.FileContentUpdate
	err := json.Unmarshal(data, &req)
	if err != nil {
		return err
	}

	//get fileinfo
	da, err := storage.GetFileInfo(req.File)
	if err != nil {
		return err
	}
	var file LocalFile
	err = json.Unmarshal(da, &file)
	if err != nil {
		return err
	}
	//compare md5
	bytemd5 := md5.Sum([]byte(req.Content))
	newmd5 := hex.EncodeToString(bytemd5[:])
	if file.MD5str == newmd5 {
		return nil
	}
	//update file content
	err = file.UpdateContent(newmd5, []byte(req.Content))
	if err != nil {
		return err
	}
	//sync file
	for _, v := range file.SyncType {
		switch v {
		case protocol.SYNC_FILE_TYPE:
			da, err := storage.GetSyncFileInfo(file.FilePath)
			if err != nil {
				return err
			}
			var info protocol.SyncFileType
			err = json.Unmarshal(da, &info)
			if err != nil {
				return err
			}
			//to sync

		case protocol.SYNC_ZOOKEEPER_TYPE:
		}
	}

	return nil
}

func handlerCreateFileContent(data []byte) error {
	var req protocol.FileContentCreate
	err := json.Unmarshal(data, &req)
	if err != nil {
		return err
	}

	if storage.ExistFileInfo(req.File) {
		return errors.New("file is already exsit!")
	}

	if len(req.SyncInfos) != len(req.SyncTypes) {
		return errors.New("syncinfo is not equal synctypes!")
	}

	byteCont := []byte(req.Content)

	bytemd5 := md5.Sum(byteCont)
	md5str := hex.EncodeToString(bytemd5[:])
	file := NewLocalFile(md5str, req.File, req.BackupCount, req.SyncTypes)
	//store file info and sync info
	da, err := json.Marshal(&file)
	if err != nil {
		return err
	}

	err = storage.PutFileInfo(req.File, da)
	if err != nil {
		return err
	}

	//update file content
	err = file.UpdateContent(md5str, byteCont)
	if err != nil {
		return err
	}

	for k, v := range req.SyncTypes {
		switch v {
		case protocol.SYNC_FILE_TYPE:
			var info protocol.SyncFileType
			err = json.Unmarshal([]byte(req.SyncInfos[k]), &info)
			if err != nil {
				return err
			}
			err = storage.PutSyncFileInfo(req.File, []byte(req.SyncInfos[k]))
			if err != nil {
				return err
			}
			//to sync
		case protocol.SYNC_ZOOKEEPER_TYPE:
		}
	}

	return nil
}

func handlerDeleteFileContent(data []byte) {

}

func handlerChangeFileName(data []byte) {

}

func handlerChangeBackupCnt(data []byte) {

}

func handlerUpdateSyncFile(data []byte) {

}
