package blogposts

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

const jsonContentType = "application/json"

type BlogServer struct {
	posts map[string]Post
	http.Handler
}

func NewBlogServer(posts map[string]Post) *BlogServer {
	srv := new(BlogServer)

	srv.posts = posts
	router := http.NewServeMux()
	router.Handle("/posts", http.HandlerFunc(srv.handleAllPosts))
	router.Handle("/post/", http.HandlerFunc(srv.handleSinglePost))
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
	_, _ = w.Write(js)
}

func (b *BlogServer) handleSinglePost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", jsonContentType)
	parsed := strings.Split(r.URL.Path, "/post/")

	if len(parsed) != 2 {
		//handle URL incorrectly parsed case
	}

	post, err := getPostBySlug(b.posts, parsed[1])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	js, _ := json.MarshalIndent(post, "\t", "\t")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(js)
}

func getPostBySlug(posts map[string]Post, slug string) (Post, error) {
	post, ok := posts[slug]
	if ok {
		return post, nil
	}
	return Post{}, fmt.Errorf("post not found: %s", slug)
}
