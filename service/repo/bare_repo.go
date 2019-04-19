package repo

import (
	"bytes"
	"context"
	"os"
	"path"

	"github.com/easy-oj/common/logs"
	"github.com/easy-oj/common/proto/oss"
	"github.com/easy-oj/repos/common/caller"
	"github.com/easy-oj/repos/tools/cmd"
)

type BareRepo struct {
	*Repo

	UUID string
}

func NewBareRepo(UUID string) *BareRepo {
	return &BareRepo{
		Repo: NewRepo(GenBareRepoDirPath(UUID)),
		UUID: UUID,
	}
}

func (r *BareRepo) Init() error {
	if err := os.MkdirAll(r.DirPath, 0777); err != nil {
		logs.Info("[BareRepo] mkdir error: %s", err.Error())
		return err
	}
	if err := cmd.Git("init", "--bare").WorkDir(r.DirPath).Run(); err != nil {
		return err
	}
	return nil
}

func (r *BareRepo) Release() (bool, error) {
	req := oss.NewGetObjectReq()
	req.Path = path.Join("eoj/repos", r.UUID)
	resp, err := caller.OSSClient.GetObject(context.Background(), req)
	if err != nil {
		logs.Error("[BareRepo] call OSSClient.GetObject error: %s", err.Error())
		return false, err
	}
	if resp.Object == nil {
		logs.Warn("[BareRepo] repository %s not found", r.UUID)
		return false, nil
	}
	if err := os.MkdirAll(r.DirPath, 0777); err != nil {
		logs.Error("[BareRepo] mkdir error: %s", err.Error())
		return false, err
	}
	if err := cmd.TarExtract(r.DirPath, bytes.NewReader(resp.Object)).Run(); err != nil {
		logs.Error("[BareRepo] release repository %s error: %s", r.UUID, err.Error())
		return false, err
	}
	return true, nil
}

func (r *BareRepo) Archive() error {
	if err := cmd.Git("update-server-info").WorkDir(r.DirPath).Run(); err != nil {
		logs.Error("[BareRepo] update repository %s error: %s", r.UUID, err.Error())
		return err
	}
	buffer := bytes.NewBuffer(nil)
	if err := cmd.TarCreate(r.DirPath, buffer).Run(); err != nil {
		logs.Error("[BareRepo] archive repository %s error: %s", r.UUID, err.Error())
		return err
	}
	req := oss.NewPutObjectReq()
	req.Path = path.Join("eoj/repos", r.UUID)
	req.Object = buffer.Bytes()
	_, err := caller.OSSClient.PutObject(context.Background(), req)
	if err != nil {
		logs.Error("[BareRepo] call OSSClient.PutObject error: %s", err.Error())
		return err
	}
	return nil
}
