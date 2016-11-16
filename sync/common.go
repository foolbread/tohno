/*
@author: foolbread
@time: 2016-11-14
@file:tohno/sync/common.go
*/
package sync

import (
	"github.com/foolbread/fbcommon/golog"
)

func InitSync() {
	golog.Info("tohno sync initing......")
}

func SyncFileByFerry(filepath string, info *SyncFileInfoSet) ([]byte, error) {
	return nil, nil
}
