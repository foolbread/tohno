/*
@author: foolbread
@time: 2016-11-14
@file:tohno/file/sync_file.go
*/
package file

import (
	"io/ioutil"
	"os"
	"path"
	"time"

	"github.com/foolbread/tohno/sync"
)

type LocalFile struct {
	MD5str         string                `json:"md5"`
	FilePath       string                `json:"filepath"`
	FileBackupPath string                `json:"backuppath"`
	BackUpCount    int                   `json:"backupcount"`
	CurBackCount   int                   `json:"-"`
	SyncFile       *sync.SyncFileInfoSet `json:"syncfile"`
}

func NewLocalFile(md5str string, filepath string, backupCnt int) *LocalFile {
	r := new(LocalFile)
	r.MD5str = md5str
	r.FilePath = filepath
	r.FileBackupPath = filepath
	r.BackUpCount = backupCnt
	r.CurBackCount = 0
	r.SyncFile = new(sync.SyncFileInfoSet)

	return r
}

//idx = 0 => file content idx > 0 => backup content
func (f *LocalFile) GetFileContent() ([]byte, error) {
	return ioutil.ReadFile(path.Join(file_dir, f.FilePath))
}

func (f *LocalFile) GetBackupContent(filename string) ([]byte, error) {
	return ioutil.ReadFile(path.Join(backup_dir, f.FileBackupPath, filename))
}

func (f *LocalFile) CreateFile(content []byte) error {
	return ioutil.WriteFile(path.Join(file_dir, f.FilePath), content, 0777)
}

func (f *LocalFile) UpdateContent(content []byte) error {
	old, err := ioutil.ReadFile(path.Join(file_dir, f.FilePath))
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(path.Join(file_dir, f.FilePath), content, 0777)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(path.Join(backup_dir, f.FileBackupPath, time.Now().Format("2006-01-02-15:04:05")), old, 0777)
	if err != nil {
		return err
	}

	f.CurBackCount++
	if f.CurBackCount > f.BackUpCount {
		err = f.ChangeBackUp(f.BackUpCount)
		if err != nil {
			return err
		}

		f.CurBackCount = f.BackUpCount
	}
	return nil
}

func (f *LocalFile) GetBackupFileList() ([]string, error) {
	infos, err := ioutil.ReadDir(path.Join(backup_dir, f.FileBackupPath))
	if err != nil {
		return nil, err
	}

	var ret []string
	for _, v := range infos {
		ret = append(ret, v.Name())
	}
	return ret, nil
}

func (f *LocalFile) DeleteFile() error {
	err := os.Remove(path.Join(file_dir, f.FilePath))
	if err != nil {
		return err
	}

	return os.RemoveAll(path.Join(backup_dir, f.FileBackupPath))
}

func (f *LocalFile) ChangeName(newpath string) error {
	return os.Rename(path.Join(file_dir, f.FilePath), newpath)
}

func (f *LocalFile) ChangeBackUp(cnt int) error {
	if cnt > f.BackUpCount {
		return nil
	}

	infos, err := ioutil.ReadDir(path.Join(backup_dir, f.FileBackupPath))
	if err != nil {
		return err
	}

	for i := 0; i < len(infos)-1; i++ {
		for j := i + 1; j < len(infos); j++ {
			if infos[i].ModTime().Before(infos[j].ModTime()) {
				infos[i], infos[j] = infos[j], infos[i]
			}
		}
	}

	for i := cnt; i < len(infos); i++ {
		err = os.Remove(path.Join(backup_dir, f.FileBackupPath, infos[i].Name()))
		if err != nil {
			return err
		}
	}

	return nil
}
