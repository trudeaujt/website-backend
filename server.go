package blogposts

import (
	"encoding/json"
	"net/http"
	"strings"
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
	router.Handle("/posts", http.HandlerFunc(srv.handleAllPosts))
	router.Handle("/post/", http.HandlerFunc(srv.handleSinglePost))
	/**
	/posts - all posts
	/posts/1 - first 10?
	/posts/2 - next 10?
	/post/this-is-a-slug - specific post
	*/
	srv.Handler = router

	return srv
}

func (b *BlogServer) handleAllPosts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", jsonContentType)

	js, err := json.MarshalIndent(b.posts, "", "\t")
	if err != nil {
		http.Error(w, "The server encountered a problem and could not process your request.", http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(js)
}

func (b *BlogServer) handleSinglePost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", jsonContentType)
	parsed := strings.Split(r.URL.Path, "/post/")

	if len(parsed) != 2 {
		//handle URL incorrectly parsed case
	}

	post := getPostBySlug(parsed[1])
	if post == nil {
		//handle post not found
	}
	js, _ := json.MarshalIndent(post, "\t", "\t")
	w.WriteHeader(http.StatusOK)
	w.Write(js)
}
