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
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"

	"github.com/foolbread/fbcommon/golog"
	"github.com/foolbread/tohno/protocol"
	"github.com/foolbread/tohno/storage"
	"github.com/foolbread/tohno/sync"
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
func marshalCommonRes(status int, str string) []byte {
	var res protocol.CommonRes
	res.Info = str
	res.Status = status

	data, _ := json.Marshal(&res)
	return data
}

///////////////////////////////////////////////////////////////////////////////////////////
func getTimeMD5KeyStr() string {
	da, _ := time.Now().MarshalBinary()
	bymd5 := md5.Sum(da)

	return hex.EncodeToString(bymd5[:])
}

func getMD5KeyStrByData(data []byte) string {
	bymd5 := md5.Sum(data)

	return hex.EncodeToString(bymd5[:])
}

///////////////////////////////////////////////////////////////////////////////////////////
func HandlerFile(action string, data []byte) {
	switch action {
	case protocol.FILE_CREATE:
	case protocol.FILE_GET:
	case protocol.FILE_DELETE:
	case protocol.FILE_UPDATE:
	case protocol.BACKUP_FILE_GET:
	}
}

func HandlerDir(action string, data []byte) {
	switch action {
	case protocol.DIR_CREATE:
	case protocol.DIR_SCAN:
	case protocol.DIR_DEL:
	case protocol.DIR_RENAME:
	}
}

/////////////////////////////////////////////////////////////////////////////////////////////////
func handlerFileCreate(data []byte) ([]byte, error) {
	var req protocol.FileCreateReq
	err := json.Unmarshal(data, &req)
	if err != nil {
		return nil, err
	}

	if storage.ExistData([]byte(req.File)) {
		return nil, errors.New("file is already exist!")
	}

	lf := NewLocalFile(getMD5KeyStrByData([]byte(req.Content)), req.File, req.BackupCount)
	err = lf.CreateFile([]byte(req.Content))
	if err != nil {
		return nil, err
	}

	for _, v := range req.SyncFile.Infos {
		lf.SyncFile.Put(v.IP, v.File)
	}

	da, err := json.Marshal(lf)
	if err != nil {
		return nil, err
	}

	err = storage.PutData([]byte(req.File), da)
	if err != nil {
		return nil, err
	}

	//sync file to remote
	sync.SyncFileByFerry(path.Join(file_dir, req.File), lf.SyncFile)
	return nil, nil
}

func handlerFileGet(data []byte) ([]byte, error) {
	var req protocol.FileGetReq
	err := json.Unmarshal(data, &req)
	if err != nil {
		return nil, err
	}

	da, err := storage.GetData([]byte(req.File))
	if err != nil {
		return nil, err
	}

	var lf LocalFile
	err = json.Unmarshal(da, &lf)
	if err != nil {
		return nil, err
	}

	var res protocol.FileGetRes
	res.File = lf.FilePath
	da, err = lf.GetFileContent()
	if err != nil {
		return nil, err
	}
	res.Content = string(da)
	for _, v := range res.SyncFile.Infos {
		res.SyncFile.Infos = append(res.SyncFile.Infos, &protocol.SyncFilePair{v.IP, v.File})
	}
	res.BackupFiles, err = lf.GetBackupFileList()
	if err != nil {
		return nil, err
	}
	res.BackupCount = lf.BackUpCount

	rdata, err := json.Marshal(&res)
	if err != nil {
		return nil, err
	}

	return rdata, nil
}

func handlerFileDel(data []byte) ([]byte, error) {
	var req protocol.FileDeleteReq
	err := json.Unmarshal(data, &req)
	if err != nil {
		return nil, err
	}

	da, err := storage.GetData([]byte(req.File))
	if err != nil {
		return nil, err
	}

	var lf LocalFile
	err = json.Unmarshal(da, &lf)
	if err != nil {
		return nil, err
	}

	err = lf.DeleteFile()
	if err != nil {
		return nil, err
	}

	err = storage.DeleteData([]byte(req.File))
	if err != nil {
		return nil, err
	}

	return marshalCommonRes(protocol.OK_STATUS, "file delete ok!"), nil
}

func handlerFileUpdate(data []byte) ([]byte, error) {
	var req protocol.FileUpdateReq
	err := json.Unmarshal(data, &req)
	if err != nil {
		return nil, err
	}

	var contentsync bool = false
	var infosync bool = false
	//get file info
	da, err := storage.GetData([]byte(req.File))
	if err != nil {
		return nil, err
	}

	var lf LocalFile
	err = json.Unmarshal(da, &lf)
	if err != nil {
		return nil, err
	}
	//check content
	md5str := getMD5KeyStrByData([]byte(req.Content))
	if md5str != lf.MD5str {
		contentsync = true
		err = lf.UpdateContent([]byte(req.Content))
		if err != nil {
			return nil, err
		}
		lf.MD5str = md5str
	}

	var syncinfos *sync.SyncFileInfoSet
	//check sync info
	if !contentsync {
		syncinfos = sync.NewSyncFileInfoSet()
		lf.SyncFile.Parse()
		for _, v := range req.SyncFile.Infos {
			if !lf.SyncFile.Exist(v.IP, v.File) {
				syncinfos.Put(v.IP, v.File)
			}
		}
	} else {
		syncinfos = lf.SyncFile
	}

	//sync
	if contentsync || infosync {

	}
	return nil, nil
}

func handlerDirCreate(data []byte) ([]byte, error) {
	var req protocol.DirCreateReq
	err := json.Unmarshal(data, &req)
	if err != nil {
		return nil, err
	}

	err = os.Mkdir(req.Dir, 0777)
	if err != nil {
		return nil, err
	}

	return marshalCommonRes(protocol.OK_STATUS, "dir create ok!"), nil
}

func handlerDirScan(data []byte) ([]byte, error) {
	var req protocol.DirScanReq
	err := json.Unmarshal(data, &req)
	if err != nil {
		return nil, err
	}

	info, err := ioutil.ReadDir(req.Dir)
	if err != nil {
		return nil, err
	}

	var res protocol.DirInfoRes
	res.Dir = req.Dir
	for _, v := range info {
		res.Infos = append(res.Infos, &protocol.DirInfoPair{v.Name(), v.IsDir()})
	}

	rdata, err := json.Marshal(&res)
	if err != nil {
		return nil, err
	}

	return rdata, nil
}

func handlerDirDel(data []byte) ([]byte, error) {
	var req protocol.DirDelReq
	err := json.Unmarshal(data, &req)
	if err != nil {
		return nil, err
	}

	datas, err := storage.GetDataByPrefix([]byte(req.Dir))
	if err != nil {
		return nil, err
	}

	err = os.RemoveAll(path.Join(file_dir, req.Dir))
	if err != nil {
		return nil, err
	}

	err = os.RemoveAll(path.Join(backup_dir, req.Dir))
	if err != nil {
		return nil, err
	}

	for _, v := range datas {
		storage.DeleteData(v.Key)
	}

	return marshalCommonRes(protocol.OK_STATUS, "dir delete ok!"), nil
}

func handlerDirRename(data []byte) ([]byte, error) {
	var req protocol.DirRenameReq
	err := json.Unmarshal(data, &req)
	if err != nil {
		return nil, err
	}

	trueNewDir := path.Join(file_dir, req.ParentPath, req.NewName)
	trueOldDir := path.Join(file_dir, req.ParentPath, req.OldName)
	err = os.Rename(trueOldDir, trueNewDir)
	if err != nil {
		return nil, err
	}

	fullNewDir := path.Join(req.ParentPath, req.NewName)
	fullOldDir := path.Join(req.ParentPath, req.OldName)
	datas, err := storage.GetDataByPrefix([]byte(fullOldDir))
	if err != nil {
		return nil, err
	}

	for _, v := range datas {
		nk := strings.Replace(string(v.Key), fullOldDir, fullNewDir, 1)
		storage.PutData([]byte(nk), v.Data)
		storage.DeleteData(v.Key)
	}

	return marshalCommonRes(protocol.OK_STATUS, "dir rename ok!"), nil
}
