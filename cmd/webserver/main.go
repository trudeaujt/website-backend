package main

import (
	"flag"
	"github.com/trudeaujt/blogposts"
	"log"
	"net/http"
	"os"
)

func main() {
	dirFlag := flag.String("directory", "/var/posts", "directory containing blogposts")
	flag.Parse()

	dir := os.DirFS(*dirFlag)
	posts, _ := blogposts.NewPostsFromFS(dir)
	server := blogposts.NewBlogServer(posts)

	if err := http.ListenAndServe(":5001", server); err != nil {
		log.Fatalf("could not listen on port 5001 %v", err)
	}
}
