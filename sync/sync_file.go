/*
@author: foolbread
@time: 2016-11-14
@file:tohno/sync/sync_file.go
*/
package sync

type SyncFileInfo struct {
	IP   string `json:"ip"`
	File string `json:"file"`
}

type SyncFileInfoSet struct {
	Infos []*SyncFileInfo `json:"syncfileinfos"`
	set   map[string]string
}

func NewSyncFileInfoSet() *SyncFileInfoSet {
	r := new(SyncFileInfoSet)
	r.set = make(map[string]string)

	return r
}

func (s *SyncFileInfoSet) Put(ip string, file string) {
	s.Infos = append(s.Infos, &SyncFileInfo{ip, file})
}

func (s *SyncFileInfoSet) Exist(ip string, file string) bool {
	_, ok := s.set[ip+file]

	return ok
}

func (s *SyncFileInfoSet) Parse() {
	for _, v := range s.Infos {
		s.set[v.IP+v.File] = v.File
	}
}
