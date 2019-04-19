package repos

import "github.com/easy-oj/common/proto/base"

func NewCreateRepoReq() *CreateRepoReq {
	return &CreateRepoReq{
		BaseReq: base.NewBaseReq(),
	}
}

func NewCreateRepoResp() *CreateRepoResp {
	return &CreateRepoResp{
		BaseResp: base.NewBaseResp(),
	}
}

func NewFetchRepoReq() *FetchRepoReq {
	return &FetchRepoReq{
		BaseReq: base.NewBaseReq(),
	}
}

func NewFetchRepoResp() *FetchRepoResp {
	return &FetchRepoResp{
		BaseResp: base.NewBaseResp(),
	}
}

func NewUpdateRepoReq() *UpdateRepoReq {
	return &UpdateRepoReq{
		BaseReq: base.NewBaseReq(),
	}
}

func NewUpdateRepoResp() *UpdateRepoResp {
	return &UpdateRepoResp{
		BaseResp: base.NewBaseResp(),
	}
}
