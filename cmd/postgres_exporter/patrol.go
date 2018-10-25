package main

import (
    "database/sql"
    _ "strconv"
    "syscall"
	"github.com/op/go-logging"
)

type DbInfo struct {
    Maxage uint64
    DiskUsage float32
}

var mlog = logging.MustGetLogger("patrol")

func HandleError(err error) {
    if err !=nil {
        mlog.Error(err)
    }
}

func diskUsage(volumePath string) float32 {

    var stat syscall.Statfs_t
    syscall.Statfs(volumePath, &stat)

    msize := stat.Blocks * uint64(stat.Bsize)
    // available = free - reserved filesystem blocks(for root)
    mavail := stat.Bavail* uint64(stat.Bsize)

    return float32(msize - mavail) / float32(msize)
}

func getDbAge(mdb *sql.DB) uint64 {
    msql := `SELECT age(datfrozenxid)
                FROM pg_database
                WHERE datname <> 'template1'
                     AND datname <> 'template0'
                     AND datname <> 'postgres'
                      AND datname <> 'monitordb'
                ORDER BY age(datfrozenxid) DESC LIMIT 1;`

    var age uint64
    err := mdb.QueryRow(msql).Scan(&age)
    if err != nil && err != sql.ErrNoRows {
        mlog.Error(err)
    }

    return age
}

func getDbDiskUsage(mdb *sql.DB) float32 {
    msql := `show data_directory;`

    var dataDir string
    err := mdb.QueryRow(msql).Scan(&dataDir)
    if err != nil && err != sql.ErrNoRows {
        mlog.Error(err)
    }

    return diskUsage(dataDir)
}
