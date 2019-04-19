package http

import (
	"fmt"
	"net/http"
	"time"

	"gopkg.in/macaron.v1"

	"github.com/easy-oj/repos/service/auth"
	"github.com/easy-oj/repos/service/repo"
)

type httpContext struct {
	*macaron.Context
	*repo.BareRepo

	fileName string
}

func newHTTPContext(ctx *macaron.Context) *httpContext {
	return &httpContext{
		Context: ctx,
	}
}

func (c *httpContext) prepare() bool {
	username, password, ok := c.Req.BasicAuth()
	if !ok {
		c.askCredentials(http.StatusUnauthorized)
		return false
	}
	uuid := c.Params(":uuid")
	if ok, err := auth.Authenticate(username, password, uuid); err != nil {
		http.Error(c.Resp, err.Error(), http.StatusInternalServerError)
		return false
	} else if !ok {
		c.askCredentials(http.StatusForbidden)
		return false
	}
	c.BareRepo = repo.NewBareRepo(uuid)
	return true
}

func (c *httpContext) askCredentials(status int) {
	c.Resp.Header().Set("WWW-Authenticate", "Basic realm=\".\"")
	c.Status(status)
}

func (c *httpContext) setHeaderNoCache() {
	c.Resp.Header().Set("Expires", "Fri, 01 Jan 1980 00:00:00 GMT")
	c.Resp.Header().Set("Pragma", "no-cache")
	c.Resp.Header().Set("Cache-Control", "no-cache, max-age=0, must-revalidate")
}

func (c *httpContext) setHeaderCacheForever() {
	now := time.Now().Unix()
	expires := now + 31536000
	c.Resp.Header().Set("Date", fmt.Sprintf("%d", now))
	c.Resp.Header().Set("Expires", fmt.Sprintf("%d", expires))
	c.Resp.Header().Set("Cache-Control", "public, max-age=31536000")
}
