package main

import (
	"github.com/braintree/manners"
	"github.com/googollee/go-middleware/bind"
	"github.com/googollee/go-middleware/routermiddle"
	"github.com/julienschmidt/httprouter"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	conf, err := Init()
	if err != nil {
		log.Fatal("init error:", err)
	}
	templ, err := template.ParseGlob(conf.Templates + "/*")
	if err != nil {
		log.Fatal("template error:", err)
	}

	videos, err := NewVideos(conf.mdb, conf.Ak, conf.Sk, conf)
	if err != nil {
		log.Fatal("create videos error:", err)
	}

	router := httprouter.New()
	m := routermiddle.New()
	m.Compose(conf.logger.Serve)

	getId := func(p httprouter.Params) string {
		return p.ByName("id")
	}

	bindError := func(err error) error {
		return Error{http.StatusBadRequest, err.Error(), 900001}
	}

	router.GET("/", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		if err := templ.ExecuteTemplate(w, "index.html", nil); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	router.GET("/player/:id", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		id := p.ByName("id")
		video, err := videos.Get(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		data := struct {
			Video *Video
		}{}
		data.Video = video
		if err := templ.ExecuteTemplate(w, "player.html", data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	router.ServeFiles("/assets/*filepath", conf.Assets)

	router.POST("/api/v1/videos",
		m.Handle(bind.Bind(Video{}, bindError), videos.Create))

	router.GET("/api/v1/videos",
		m.Handle(videos.List))

	router.GET("/api/v1/videos/:id",
		m.Handle(getId, videos.Get))

	router.DELETE("/api/v1/videos/:id",
		m.Handle(getId, videos.Delete))

	router.GET("/api/v1/videos/:id/push_url",
		m.Handle(getId, videos.GetPushUrl))

	conf.logger.Printf("start server: %s", conf.ListenAddr)

	var handler http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE")
			w.Header().Set("Access-Control-Max-Age", "600")
			w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
			return
		}
		router.ServeHTTP(w, r)
	}
	server := manners.NewServer()

	termSig := make(chan os.Signal)
	go func() {
		signal.Notify(termSig, syscall.SIGTERM)
		<-termSig
		conf.logger.Printf("received SIGTERM")
		server.Shutdown <- true
	}()

	err = server.ListenAndServe(conf.ListenAddr, handler)
	if err != nil {
		conf.logger.Printf("server launch error: %s", err)
	}
	conf.logger.Printf("quitting")
}
