package main

import (
	"flag"
	"fmt"
	"github.com/trudeaujt/blogposts"
	"log"
	"net/http"
	"os"
)

func main() {
	dirFlag := flag.String("directory", "/var/posts", "directory containing blogposts")
	portFlag := flag.String("port", "5001", "port to listen on")
	flag.Parse()

	dir := os.DirFS(*dirFlag)
	posts, _ := blogposts.NewPostsFromFS(dir)
	server := blogposts.NewBlogServer(posts)

	port := fmt.Sprintf(":%v", *portFlag)

	if err := http.ListenAndServe(port, server); err != nil {
		log.Fatalf("could not listen on port 5001 %v", err)
	}
}
