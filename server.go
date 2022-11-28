package blogposts

import (
	"encoding/json"
	"net/http"
)

const jsonContentType = "application/json"

type BlogServer struct {
	posts []Post
	http.Handler
}

func NewBlogServer(posts []Post) *BlogServer {
	srv := new(BlogServer)

	srv.posts = posts
	router := http.NewServeMux()
	router.Handle("/posts", http.HandlerFunc(srv.handlePosts))
	srv.Handler = router

	return srv
}

func (b *BlogServer) handlePosts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", jsonContentType)
	js, _ := json.Marshal(b.posts)
	w.WriteHeader(http.StatusOK)
	w.Write(js)
}
