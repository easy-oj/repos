package service

import (
	"context"
	"errors"

	"github.com/easy-oj/common/proto/repos"
	"github.com/easy-oj/common/tools"
	"github.com/easy-oj/repos/common/database"
	"github.com/easy-oj/repos/service/repo"
	"github.com/easy-oj/repos/service/submit"
)

type reposHandler struct{}

var (
	repoDataNotFoundError = errors.New("repo data not found")
)

func NewReposHandler() *reposHandler {
	return &reposHandler{}
}

func (h *reposHandler) CreateRepo(ctx context.Context, req *repos.CreateRepoReq) (*repos.CreateRepoResp, error) {
	resp := repos.NewCreateRepoResp()
	uuid := tools.GenUUID()
	bareRepo := repo.NewBareRepo(uuid)
	defer bareRepo.Clean()
	if err := bareRepo.Init(); err != nil {
		return resp, err
	}
	if err := h.updateRepo(bareRepo, req.Content); err != nil {
		return resp, err
	}
	_, err := database.DB.Exec(
		"INSERT INTO tb_repo (uid, pid, lid, uuid) VALUES (?, ?, ?, ?)", req.Uid, req.Pid, req.Lid, uuid)
	if err != nil {
		return resp, err
	}
	resp.Sid, err = submit.Submit(uuid, req.Content)
	return resp, err
}

func (*reposHandler) FetchRepo(ctx context.Context, req *repos.FetchRepoReq) (*repos.FetchRepoResp, error) {
	resp := repos.NewFetchRepoResp()
	bareRepo := repo.NewBareRepo(req.Uuid)
	defer bareRepo.Clean()
	if ok, err := bareRepo.Release(); err != nil {
		return resp, err
	} else if !ok {
		return resp, repoDataNotFoundError
	}
	tmpRepo := repo.NewTmpRepo(bareRepo)
	defer tmpRepo.Clean()
	if err := tmpRepo.Clone(); err != nil {
		return resp, err
	}
	if content, err := tmpRepo.Read(); err != nil {
		return resp, err
	} else {
		resp.Content = content
	}
	return resp, nil
}

func (h *reposHandler) UpdateRepo(ctx context.Context, req *repos.UpdateRepoReq) (*repos.UpdateRepoResp, error) {
	resp := repos.NewUpdateRepoResp()
	bareRepo := repo.NewBareRepo(req.Uuid)
	defer bareRepo.Clean()
	ok, err := bareRepo.Release()
	if err != nil {
		return resp, err
	} else if !ok {
		return resp, repoDataNotFoundError
	}
	if err := h.updateRepo(bareRepo, req.Content); err != nil {
		return resp, err
	}
	resp.Sid, err = submit.Submit(req.Uuid, req.Content)
	return resp, nil
}

func (*reposHandler) updateRepo(bareRepo *repo.BareRepo, content map[string]string) error {
	tmpRepo := repo.NewTmpRepo(bareRepo)
	defer tmpRepo.Clean()
	if err := tmpRepo.Clone(); err != nil {
		return err
	}
	if err := tmpRepo.Write(content); err != nil {
		return err
	}
	if err := tmpRepo.Push(); err != nil {
		return err
	}
	if err := bareRepo.Archive(); err != nil {
		return err
	}
	return nil
}
