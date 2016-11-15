/*
@author: foolbread
@time: 2016-11-14
@file:tohno/file/sync_file.go
*/
package file

import (
	"github.com/foolbread/tohno/sync"
)

type LocalFile struct {
	MD5str         string                `json:"md5"`
	FilePath       string                `json:"filepath"`
	FileBackupPath string                `json:"backuppath"`
	BackUpCount    int                   `json:"backupcount"`
	SyncFile       *sync.SyncFileInfoSet `json:"syncfile"`
}

func NewLocalFile(md5str string, filepath string, backupCnt int) *LocalFile {
	r := new(LocalFile)
	r.MD5str = md5str
	r.FilePath = filepath
	r.FileBackupPath = filepath
	r.BackUpCount = backupCnt
	r.SyncFile = new(sync.SyncFileInfoSet)

	return r
}

//idx = 0 => file content idx > 0 => backup content
func (f *LocalFile) GetFileContent(idx int) ([]byte, error) {
	return nil, nil
}

func (f *LocalFile) CreateFile(content []byte) error {
	return nil
}

func (f *LocalFile) UpdateContent(content []byte) error {
	return nil
}

func (f *LocalFile) DeleteFile() error {
	return nil
}

func (f *LocalFile) ChangeName(name string) error {
	return nil
}

func (f *LocalFile) ChangeBackUp(cnt int) error {
	return nil
}
