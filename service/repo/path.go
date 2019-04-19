package repo

import (
	"path"

	"github.com/easy-oj/common/settings"
)

func GenBareRepoDirPath(repoName string) string {
	return path.Join(settings.Repos.Path, "bare", repoName)
}

func GenTmpRepoDirPath(repoName string) string {
	return path.Join(settings.Repos.Path, "tmp", repoName)
}
