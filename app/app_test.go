package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/franela/goblin"
	"github.com/gin-gonic/gin"
)

const Cookie = "Set-Cookie"

func client(method, path, cookie string) *httptest.ResponseRecorder {
	return request(method, path, cookie, "")
}

func request(method, path, cookie string, body string) *httptest.ResponseRecorder {
	gin.SetMode("test")
	app := NewApp()
	payload := bytes.NewBufferString(body)
	req, _ := http.NewRequest(method, path, payload)
	if len(cookie) != 0 {
		req.Header.Set("Cookie", cookie)
	}
	w := httptest.NewRecorder()
	app.ServeHTTP(w, req)
	return w
}

func Test(t *testing.T) {
	g := Goblin(t)
	g.Describe("App api", func() {
		var cookie string

		g.It("Should return 302 on / to redirect to file name ", func() {
			w := client("GET", "/", "")

			g.Assert(w.Code).Equal(302)
			cookie = w.HeaderMap.Get(Cookie)
		})

		g.It("Should return 200 on /slides.md ", func() {
			w := client("GET", "/slides.md", cookie)
			g.Assert(w.Code).Equal(200)
		})

		g.It("Should return 200 on PUT /slides.md ", func() {
			w := request("PUT", "/slides.md", cookie, "whatever")
			g.Assert(w.Code).Equal(200)
		})

		g.It("Should return list of slides ", func() {
			w := client("GET", "/stash", cookie)
			g.Assert(w.Code).Equal(200)
		})

		g.It("Should return specific slide in preview", func() {
			w := client("GET", "/published/slides/shy-cell.md", cookie)
			g.Assert(w.Code).Equal(200)
		})

		g.It("Should return specific slide in edit mode", func() {
			w := client("GET", "/stash/edit/shy-cell.md", cookie)
			g.Assert(w.Code).Equal(200)
		})

		g.It("Should return specific slide in preview without session", func() {
			w := client("GET", "/published/slides/shy-cell.md", "")
			g.Assert(w.Code).Equal(200)
		})

		g.It("Should return specific slide in edit mode without session", func() {
			w := client("GET", "/stash/edit/shy-cell.md", "")
			g.Assert(w.Code).Equal(200)
		})

		g.It("Should works")

	})
}
