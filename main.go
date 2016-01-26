package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"text/template"

	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/facebook"
	"github.com/stretchr/gomniauth/providers/github"
	"github.com/stretchr/gomniauth/providers/google"
	"github.com/stretchr/objx"
	"github.com/stretchr/signature"

	"github.com/upamune/ws-chat/trace"
)

const (
	securityKey           = "SECURITY_KEY"
	githubAuthKey         = "GITHUB_AUTH_KEY"
	githubAuthSecretKey   = "GITHUB_AUTH_SECRET_KEY"
	facebookAuthKey       = "FACEBOOK_AUTH_KEY"
	facebookAuthSecretKey = "FACEBOOK_AUTH_SECRET_KEY"
	googleAuthKey         = "GOOGLE_AUTH_KEY"
	googleAuthSecretKey   = "GOOGLE_AUTH_SECRET_KEY"
)

type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	data := map[string]interface{}{
		"Host": r.Host,
	}
	if authCookie, err := r.Cookie("auth"); err == nil {
		data["UserData"] = objx.MustFromBase64(authCookie.Value)
	}
	err := t.templ.Execute(w, data)
	if err != nil {
		log.Fatal("Template Excute:", err)
	}
}

// 現在アクティブなAvatarの実装
var avatars = TryAvatars{
	UseFileSystemAvatar,
	UseAuthAvatar,
	UseGravatar,
}

func main() {
	var addr = flag.String("addr", ":8080", "アプリケーションのアドレス")
	flag.Parse()
	gomniauth.SetSecurityKey(signature.RandomKey(64))
	gomniauth.WithProviders(
		facebook.New(os.Getenv(facebookAuthKey), os.Getenv(facebookAuthSecretKey), "http://localhost:8080/auth/callback/facebook"),
		github.New(os.Getenv(githubAuthKey), os.Getenv(githubAuthSecretKey), "http://localhost:8080/auth/callback/github"),
		google.New(os.Getenv(googleAuthKey), os.Getenv(googleAuthSecretKey), "http://localhost:8080/auth/callback/google"),
	)

	r := newRoom()
	r.tracer = trace.New(os.Stdout)
	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.Handle("/upload", &templateHandler{filename: "upload.html"})
	http.Handle("/avatars/", http.StripPrefix("/avatars/", http.FileServer(http.Dir("./avatars"))))
	http.HandleFunc("/uploader", uploaderHandler)
	http.HandleFunc("/auth/", loginHandler)
	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:   "auth",
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		})
		w.Header()["Location"] = []string{"/chat"}
		w.WriteHeader(http.StatusTemporaryRedirect)
	})
	http.Handle("/room", r)
	go r.run()

	log.Println("Webサーバーを開始します．ポート: ", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
