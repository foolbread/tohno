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
func marshalCommonRes(str string) []byte {
	var res protocol.CommonRes
	res.Info = str

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
func handlerFileCreate(data []byte) []byte {
	var req protocol.FileCreateReq
	err := json.Unmarshal(data, &req)
	if err != nil {
		return marshalCommonRes(err.Error())
	}

	if storage.ExistData([]byte(req.File)) {
		return marshalCommonRes("file is already exist!")
	}

	lf := NewLocalFile(getMD5KeyStrByData([]byte(req.Content)), req.File, req.BackupCount)
	err = lf.CreateFile([]byte(req.Content))
	if err != nil {
		return marshalCommonRes(err.Error())
	}

	for _, v := range req.SyncFile.Infos {
		lf.SyncFile.Put(v.IP, v.File)
	}

	da, err := json.Marshal(lf)
	if err != nil {
		return marshalCommonRes(err.Error())
	}

	err = storage.PutData([]byte(req.File), da)
	if err != nil {
		return marshalCommonRes(err.Error())
	}

	//sync file to remote
	return nil
}

func handlerFileGet(data []byte) []byte {
	var req protocol.FileGetReq
	err := json.Unmarshal(data, &req)
	if err != nil {
		return marshalCommonRes(err.Error())
	}

	da, err := storage.GetData([]byte(req.File))
	if err != nil {
		return marshalCommonRes(err.Error())
	}

	var lf LocalFile
	err = json.Unmarshal(da, &lf)
	if err != nil {
		return marshalCommonRes(err.Error())
	}

	var res protocol.FileGetRes
	res.File = lf.FilePath
	da, err = lf.GetFileContent(0)
	if err != nil {
		return marshalCommonRes(err.Error())
	}
	res.Content = string(da)
	for _, v := range res.SyncFile.Infos {
		res.SyncFile.Infos = append(res.SyncFile.Infos, &protocol.SyncFilePair{v.IP, v.File})
	}

	return nil
}

func handlerFileDel(data []byte) []byte {
	var req protocol.FileDeleteReq
	err := json.Unmarshal(data, &req)
	if err != nil {
		return marshalCommonRes(err.Error())
	}

	da, err := storage.GetData([]byte(req.File))
	if err != nil {
		return marshalCommonRes(err.Error())
	}

	var lf LocalFile
	err = json.Unmarshal(da, &lf)
	if err != nil {
		return marshalCommonRes(err.Error())
	}

	err = lf.DeleteFile()
	if err != nil {
		return marshalCommonRes(err.Error())
	}

	err = storage.DeleteData([]byte(req.File))
	if err != nil {
		return marshalCommonRes(err.Error())
	}
	return nil
}

func handlerFileUpdate(data []byte) []byte {
	var req protocol.FileUpdateReq
	err := json.Unmarshal(data, &req)
	if err != nil {
		return marshalCommonRes(err.Error())
	}

	var contentsync bool = false
	var infosync bool = false
	//get file info
	da, err := storage.GetData([]byte(req.File))
	if err != nil {
		return marshalCommonRes(err.Error())
	}

	var lf LocalFile
	err = json.Unmarshal(da, &lf)
	if err != nil {
		return marshalCommonRes(err.Error())
	}
	//check content
	md5str := getMD5KeyStrByData([]byte(req.Content))
	if md5str != lf.MD5str {
		contentsync = true
		err = lf.UpdateContent([]byte(req.Content))
		if err != nil {
			return marshalCommonRes(err.Error())
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
	return nil
}

func handlerDirCreate(data []byte) []byte {
	var req protocol.DirCreateReq
	err := json.Unmarshal(data, &req)
	if err != nil {
		return marshalCommonRes(err.Error())
	}

	err = os.Mkdir(req.Dir, 0777)
	if err != nil {
		return marshalCommonRes(err.Error())
	}
	return nil
}

func handlerDirScan(data []byte) []byte {
	var req protocol.DirScanReq
	err := json.Unmarshal(data, &req)
	if err != nil {
		return marshalCommonRes(err.Error())
	}

	info, err := ioutil.ReadDir(req.Dir)
	if err != nil {
		return marshalCommonRes(err.Error())
	}

	var res protocol.DirInfoRes
	res.Dir = req.Dir
	for _, v := range info {
		res.Infos = append(res.Infos, &protocol.DirInfoPair{v.Name(), v.IsDir()})
	}

	return nil
}

func handlerDirDel(data []byte) []byte {
	var req protocol.DirDelReq
	err := json.Unmarshal(data, &req)
	if err != nil {
		return marshalCommonRes(err.Error())
	}

	datas, err := storage.GetDataByPrefix([]byte(req.Dir))
	if err != nil {
		return marshalCommonRes(err.Error())
	}

	err = os.RemoveAll(path.Join(file_dir, req.Dir))
	if err != nil {
		return marshalCommonRes(err.Error())
	}

	err = os.RemoveAll(path.Join(backup_dir, req.Dir))
	if err != nil {
		return marshalCommonRes(err.Error())
	}

	for _, v := range datas {
		storage.DeleteData(v.Key)
	}

	return nil
}

func handlerDirRename(data []byte) []byte {
	var req protocol.DirRenameReq
	err := json.Unmarshal(data, &req)
	if err != nil {
		return marshalCommonRes(err.Error())
	}

	trueNewDir := path.Join(file_dir, req.ParentPath, req.NewName)
	trueOldDir := path.Join(file_dir, req.ParentPath, req.OldName)
	err = os.Rename(trueOldDir, trueNewDir)
	if err != nil {
		return marshalCommonRes(err.Error())
	}

	fullNewDir := path.Join(req.ParentPath, req.NewName)
	fullOldDir := path.Join(req.ParentPath, req.OldName)
	datas, err := storage.GetDataByPrefix([]byte(fullOldDir))
	if err != nil {
		return marshalCommonRes(err.Error())
	}

	for _, v := range datas {
		nk := strings.Replace(string(v.Key), fullOldDir, fullNewDir, 1)
		storage.PutData([]byte(nk), v.Data)
		storage.DeleteData(v.Key)
	}

	return nil
}
