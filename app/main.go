package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"

	haikunator "github.com/atrox/haikunatorgo"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	cache "github.com/hashicorp/golang-lru"
	"github.com/msoedov/hacker-slides/auth"
	"github.com/msoedov/hacker-slides/files"
	log "github.com/sirupsen/logrus"
)

const sessionHeader = "slide-session"

func SlidePath(name string) string {
	return fmt.Sprintf("slides/%s.md", name)
}

func NewApp() *gin.Engine {

	r := gin.Default()

	store := cookie.NewStore([]byte("secret"))
	arc, err := cache.NewARC(10)
	if err != nil {
		log.Fatalf("Failied to allocate cache %#v", err)
	}
	r.Use(sessions.Sessions(sessionHeader, store))
	r.Use(auth.BasicAuth())

	r.LoadHTMLGlob("templates/*.tmpl")
	r.Static("/static", "./static")
	r.Static("/images", "./slides/images")

	r.GET("/", func(c *gin.Context) {
		isNew := c.Query("new")
		latest := files.LatestFileIn("slides")
		log.WithFields(log.Fields{
			"name":  latest,
			"isNew": isNew,
		}).Info("Restoring latest point")

		var path, name string
		if latest == "" || isNew != "" {
			haikunator := haikunator.New()
			haikunator.TokenLength = 0
			name = haikunator.Haikunate()
		} else {
			name = strings.Replace(latest, ".md", "", 1)
		}
		path = SlidePath(name)

		log.WithFields(log.Fields{
			"path": path,
		}).Info("A new session")
		session := sessions.Default(c)
		session.Set("name", path)
		session.Save()

		c.Writer.Header().Set("Location", fmt.Sprintf("/stash/edit/%s", name))
		c.HTML(302, "index.tmpl", gin.H{
			"pubTo": path,
		})
	})

	mustHaveSession := func(c *gin.Context) (string, error) {
		session := sessions.Default(c)
		val := session.Get("name")
		emptySession := errors.New("Empty session")
		if val == nil {
			c.String(400, "No context")
			return "", emptySession
		}
		log.WithFields(log.Fields{
			"path": val,
		}).Info("Got session")
		path, ok := val.(string)
		if !ok {
			c.String(400, "No context")
			return "", emptySession
		}
		return path, nil
	}

	r.GET("/slides.md", func(c *gin.Context) {
		path, err := mustHaveSession(c)
		if err != nil {
			return
		}
		if _, err := os.Stat(path); err != nil {
			// copy sample markdown file to the path
			body, err := ioutil.ReadFile("initial-slides.md")
			if err != nil {
				panic(err)
			}
			ioutil.WriteFile(path, body, 0644)
			c.String(200, string(body))
			return
		}

		var slide string
		cached, ok := arc.Get(path)
		if ok {
			slide = string(cached.([]byte))
		} else {
			body, err := ioutil.ReadFile(path)
			if err != nil {
				log.Errorf("Failied to read file %#v", err)
				c.Abort()
				return
			}
			slide = string(body)
		}
		c.String(200, slide)
	})

	r.PUT("/slides.md", func(c *gin.Context) {
		path, err := mustHaveSession(c)
		if err != nil {
			return
		}
		body, _ := ioutil.ReadAll(c.Request.Body)
		arc.Add(path, body)
		go ioutil.WriteFile(path, body, 0644)
		log.WithFields(log.Fields{
			"size": len(body),
			"file": path,
		}).Info("Async wrote to file")
		c.String(200, "")
	})

	r.GET("/stash", func(c *gin.Context) {
		files, err := ioutil.ReadDir("slides")
		if err != nil {
			log.Fatal(err)
		}

		sort.Slice(files, func(i, j int) bool {
			return files[i].ModTime().Unix() > files[j].ModTime().Unix()
		})

		var stash []string
		for _, file := range files {
			if file.IsDir() {
				continue
			}
			stash = append(stash, file.Name())
		}
		c.HTML(200, "stash.tmpl", gin.H{
			"stash": stash,
		})
	})

	r.GET("/stash/edit/:name", func(c *gin.Context) {

		name := c.Param("name")
		log.WithFields(log.Fields{
			"name": name,
		}).Info("Restore session?")

		if strings.HasSuffix(name, ".md") {
			name = name[0 : len(name)-3]
		}
		path := SlidePath(name)
		session := sessions.Default(c)
		session.Set("name", path)
		session.Save()

		c.HTML(200, "index.tmpl", gin.H{
			"pubTo": path,
		})
	})

	r.GET("/published/slides/:name", func(c *gin.Context) {

		name := c.Param("name")
		log.WithFields(log.Fields{
			"name": name,
		}).Info("Published")

		if strings.HasSuffix(name, ".md") {
			name = name[0 : len(name)-3]
		}
		path := SlidePath(name)
		session := sessions.Default(c)
		session.Set("name", path)
		session.Save()
		c.HTML(200, "slides.tmpl", gin.H{
			"pubTo": path,
		})
	})

	return r

}

func main() {
	r := NewApp()
	port := "8080"
	if len(os.Args) > 1 {
		port = os.Args[1]
	} else {
		envPort := os.Getenv("PORT")
		if len(envPort) > 0 {
			port = envPort
		}
	}
	log.Info("Started http://0.0.0.0:8080")
	r.Run(fmt.Sprintf(":%s", port))
}
