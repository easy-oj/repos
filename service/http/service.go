package http

import (
	"compress/gzip"
	"fmt"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/easy-oj/common/logs"
	"github.com/easy-oj/repos/service/repo"
	"github.com/easy-oj/repos/service/submit"
	"github.com/easy-oj/repos/tools/cmd"
)

func (c *httpContext) packService(serviceName string) bool {
	if c.Req.Header.Get("Content-Type") != fmt.Sprintf("application/x-git-%s-request", serviceName) {
		c.Status(http.StatusUnauthorized)
		return false
	}
	c.Resp.Header().Set("Content-Type", fmt.Sprintf("application/x-git-%s-result", serviceName))
	reqBody := c.Req.Request.Body
	if c.Req.Header.Get("Content-Encoding") == "gzip" {
		var err error
		if reqBody, err = gzip.NewReader(reqBody); err != nil {
			http.Error(c.Resp, err.Error(), http.StatusInternalServerError)
			return false
		}
	}

	if err := cmd.Git(serviceName, "--stateless-rpc", c.DirPath).WorkDir(c.DirPath).Stdin(reqBody).Stdout(c.Resp).Run(); err != nil {
		http.Error(c.Resp, err.Error(), http.StatusInternalServerError)
		return false
	}
	return true
}

func (c *httpContext) uploadPackService() {
	c.packService("upload-pack")
}

func (c *httpContext) receivePackService() {
	if !c.packService("receive-pack") {
		return
	}
	if err := c.Archive(); err != nil {
		http.Error(c.Resp, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpRepo := repo.NewTmpRepo(c.BareRepo)
	defer tmpRepo.Clean()
	if err := tmpRepo.Clone(); err != nil {
		return
	}
	if content, err := tmpRepo.Read(); err == nil {
		if _, err := submit.Submit(c.UUID, content); err != nil {
			logs.Error("[HTTP] uuid = %s submit error: %s", c.UUID, err.Error())
		}
	}
}

func (c *httpContext) refsService() {
	c.setHeaderNoCache()
	service := c.Query("service")
	if strings.HasPrefix(service, "git-") {
		service = strings.TrimPrefix(service, "git-")
	} else {
		service = ""
	}
	if service != "upload-pack" && service != "receive-pack" {
		c.fileService("text/plain; charset=utf-8")
		return
	}
	str := "# service=git-" + service + "\n"
	header := strconv.FormatInt(int64(len(str)+4), 16)
	if len(header)%4 != 0 {
		header = strings.Repeat("0", 4-len(header)%4) + header
	}
	c.Resp.Header().Set("Content-Type", fmt.Sprintf("application/x-git-%s-advertisement", service))
	c.Resp.Write([]byte(header + str + "0000"))
	cmd.Git(service, "--stateless-rpc", "--advertise-refs", ".").WorkDir(c.DirPath).Stdout(c.Resp).JustRun()
}

func (c *httpContext) fileService(contentType string) {
	filePath := path.Join(c.DirPath, c.fileName)
	info, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		c.Status(http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(c.Resp, err.Error(), http.StatusInternalServerError)
		return
	}
	c.Resp.Header().Set("Content-Type", contentType)
	c.Resp.Header().Set("Content-Length", fmt.Sprintf("%d", info.Size()))
	c.Resp.Header().Set("Last-Modified", info.ModTime().Format(http.TimeFormat))
	http.ServeFile(c.Resp, c.Req.Request, filePath)
}

func (c *httpContext) textFileService() {
	c.setHeaderNoCache()
	c.fileService("text/plain")
}

func (c *httpContext) infoPacksService() {
	c.setHeaderCacheForever()
	c.fileService("text/plain; charset=utf-8")
}

func (c *httpContext) looseObjectService() {
	c.setHeaderCacheForever()
	c.fileService("application/x-git-loose-object")
}

func (c *httpContext) packFileService() {
	c.setHeaderCacheForever()
	c.fileService("application/x-git-packed-objects")
}

func (c *httpContext) idxFileService() {
	c.setHeaderCacheForever()
	c.fileService("application/x-git-packed-objects-toc")
}
