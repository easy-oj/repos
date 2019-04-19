package repo

import (
	"os"

	"github.com/easy-oj/common/logs"
)

type Repo struct {
	DirPath string
}

func NewRepo(dirPath string) *Repo {
	return &Repo{
		DirPath: dirPath,
	}
}

func (r *Repo) Clean() {
	if err := os.RemoveAll(r.DirPath); err != nil {
		logs.Warn("[Repo] clean error: %s", err.Error())
	}
}
