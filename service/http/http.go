package http

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"gopkg.in/macaron.v1"

	"github.com/easy-oj/common/logs"
	"github.com/easy-oj/common/settings"
)

type httpService struct {
	routes []*httpRoute
}

type httpRoute struct {
	regex   *regexp.Regexp
	method  string
	handler func(*httpContext)
}

func StartHTTPService() {
	m := macaron.New()
	if settings.Repos.HTTP.Log {
		m.Use(macaron.Logger())
	}
	m.Use(macaron.Recovery())
	m.Use(macaron.Renderer())
	m.Route("/repos/:uuid/*", "GET,POST", newHTTPService().handle)
	go func() {
		address := fmt.Sprintf("0.0.0.0:%d", settings.Repos.HTTP.Port)
		if err := http.ListenAndServe(address, m); err != nil {
			panic(err)
		}
		logs.Info("[HTTP] service served on %s", address)
	}()
}

func newHTTPService() *httpService {
	s := &httpService{}
	s.routes = []*httpRoute{
		{regexp.MustCompile("(.*?)/git-upload-pack$"), "POST", s.uploadPackService},
		{regexp.MustCompile("(.*?)/git-receive-pack$"), "POST", s.receivePackService},
		{regexp.MustCompile("(.*?)/info/refs$"), "GET", s.refsService},
		{regexp.MustCompile("(.*?)/HEAD$"), "GET", s.textFileService},
		{regexp.MustCompile("(.*?)/objects/info/alternates$"), "GET", s.textFileService},
		{regexp.MustCompile("(.*?)/objects/info/http-alternates$"), "GET", s.textFileService},
		{regexp.MustCompile("(.*?)/objects/info/packs$"), "GET", s.infoPacksService},
		{regexp.MustCompile("(.*?)/objects/info/[^/]*$"), "GET", s.textFileService},
		{regexp.MustCompile("(.*?)/objects/[0-9a-f]{2}/[0-9a-f]{38}$"), "GET", s.looseObjectService},
		{regexp.MustCompile("(.*?)/objects/pack/pack-[0-9a-f]{40}\\.pack$"), "GET", s.packFileService},
		{regexp.MustCompile("(.*?)/objects/pack/pack-[0-9a-f]{40}\\.idx$"), "GET", s.idxFileService},
	}
	return s
}

func (s *httpService) handle(ctx *macaron.Context) {
	httpCtx := newHTTPContext(ctx)
	if httpCtx.prepare() {
		s.handleContext(httpCtx)
	}
}

func (s *httpService) handleContext(c *httpContext) {
	reqPath := c.Req.URL.Path
	for _, route := range s.routes {
		ss := route.regex.FindStringSubmatch(reqPath)
		if ss == nil {
			continue
		}
		c.fileName = strings.TrimPrefix(reqPath, ss[1]+"/")
		s.handleRoute(c, route)
		return
	}
	c.Status(http.StatusNotFound)
}

func (s *httpService) handleRoute(c *httpContext, r *httpRoute) {
	if c.Req.Method != r.method {
		c.Status(http.StatusMethodNotAllowed)
		return
	}
	defer c.BareRepo.Clean()
	if ok, err := c.Release(); err != nil {
		http.Error(c.Resp, err.Error(), http.StatusInternalServerError)
	} else if !ok {
		c.Status(http.StatusNotFound)
	} else {
		r.handler(c)
	}
}

func (s *httpService) uploadPackService(c *httpContext) {
	c.uploadPackService()
}

func (s *httpService) receivePackService(c *httpContext) {
	c.receivePackService()
}

func (s *httpService) refsService(c *httpContext) {
	c.refsService()
}

func (s *httpService) textFileService(c *httpContext) {
	c.textFileService()
}

func (s *httpService) infoPacksService(c *httpContext) {
	c.infoPacksService()
}

func (s *httpService) looseObjectService(c *httpContext) {
	c.looseObjectService()
}

func (s *httpService) packFileService(c *httpContext) {
	c.packFileService()
}

func (s *httpService) idxFileService(c *httpContext) {
	c.idxFileService()
}
