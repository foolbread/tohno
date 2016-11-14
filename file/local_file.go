/*
@author: foolbread
@time: 2016-11-14
@file:tohno/file/sync_file.go
*/
package file

type LocalFile struct {
	MD5str         string   `json:"md5"`
	FilePath       string   `json:"filepath"`
	FileBackupPath string   `json:"backuppath"`
	BackUpCount    int      `json:"backupcount"`
	SyncType       []string `json:"synctype"`
}

func NewLocalFile(md5str string, filepath string, backupCnt int, stype []string) *LocalFile {
	r := new(LocalFile)
	r.MD5str = md5str
	r.FilePath = filepath
	r.FileBackupPath = filepath
	r.BackUpCount = backupCnt
	r.SyncType = stype

	return r
}

func (f *LocalFile) UpdateContent(md5str string, content []byte) error {
	return nil
}
