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
	/**
	/posts - all posts
	/posts/1 - first 10?
	/posts/2 - next 10?
	/post/this-is-a-slug - specific post
	*/
	srv.Handler = router

	return srv
}

func (b *BlogServer) handlePosts(w http.ResponseWriter, r *http.Request) {
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
	//TODO implement a PostStore and create a getPost function on it to return a single post.
	//Don't want to put this logic into the Post class.
}
