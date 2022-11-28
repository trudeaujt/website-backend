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
	w.Header().Set("Content-Type", jsonContentType)
	js, err := json.Marshal(b.posts)
	if err != nil {
		http.Error(w, "The server encountered a problem and could not process your request.", http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(js)
}
