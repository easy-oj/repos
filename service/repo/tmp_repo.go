package repo

import (
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/easy-oj/repos/tools/cmd"
)

type TmpRepo struct {
	*Repo

	BareRepo *BareRepo
}

func NewTmpRepo(bareRepo *BareRepo) *TmpRepo {
	return &TmpRepo{
		Repo:     NewRepo(GenTmpRepoDirPath(bareRepo.UUID)),
		BareRepo: bareRepo,
	}
}

func (r *TmpRepo) Clone() error {
	cloneDirPath := path.Dir(r.DirPath)
	if err := os.MkdirAll(cloneDirPath, 0777); err != nil {
		return err
	}
	if err := cmd.Git("clone", r.BareRepo.DirPath).WorkDir(cloneDirPath).Run(); err != nil {
		return err
	}
	return nil
}

func (r *TmpRepo) Read() (map[string]string, error) {
	content := make(map[string]string)
	var fn func(dirPath string) error
	fn = func(dirPath string) error {
		fs, err := ioutil.ReadDir(dirPath)
		if err != nil {
			return err
		}
		for _, f := range fs {
			if dirPath == r.DirPath && f.Name() == ".git" {
				continue
			}
			p := path.Join(dirPath, f.Name())
			if f.IsDir() {
				if err := fn(p); err != nil {
					return err
				}
			} else {
				if bs, err := ioutil.ReadFile(p); err != nil {
					return err
				} else {
					content[strings.TrimPrefix(p, r.DirPath)] = string(bs)
				}
			}
		}
		return nil
	}
	if err := fn(r.DirPath); err != nil {
		return nil, err
	}
	return content, nil
}

func (r *TmpRepo) Write(content map[string]string) error {
	fs, err := ioutil.ReadDir(r.DirPath)
	if err != nil {
		return err
	}
	for _, f := range fs {
		if f.Name() == ".git" {
			continue
		}
		if err := os.RemoveAll(path.Join(r.DirPath, f.Name())); err != nil {
			return err
		}
	}
	for p, c := range content {
		p = path.Join(r.DirPath, p)
		if err := ioutil.WriteFile(p, []byte(c), 0666); err != nil {
			return err
		}
	}
	return nil
}

func (r *TmpRepo) Push() error {
	if err := cmd.Git("add", "-A").WorkDir(r.DirPath).Run(); err != nil {
		return err
	}
	if err := cmd.Git("commit", "--allow-empty", "-m", "Auto committed").WorkDir(r.DirPath).Run(); err != nil {
		return err
	}
	if err := cmd.Git("push").WorkDir(r.DirPath).Run(); err != nil {
		return err
	}
	return nil
}
