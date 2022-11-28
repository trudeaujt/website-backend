package blogposts

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

type Post struct {
	Title       string
	Slug        string
	Description string
	Tags        []string
	Body        string
}

func newPost(postFile io.Reader) (Post, error) {
	scanner := bufio.NewScanner(postFile)
	post := Post{}

	for scanner.Scan() {
		line := scanner.Text()
		if line == "---" {
			break
		}

		parameter := strings.Split(line, ": ")
		switch parameter[0] {
		case "Title":
			post.Title = parameter[1]
			post.Slug = strings.ToLower(strings.ReplaceAll(post.Title, " ", "_"))
		case "Description":
			post.Description = parameter[1]
		case "Tags":
			post.Tags = strings.Split(parameter[1], ", ")
		default:
			return Post{}, errors.New(fmt.Sprintf("invalid parameter: %v", line))
		}
	}

	post.Body = readBody(scanner)
	return post, nil
}

func readBody(scanner *bufio.Scanner) string {
	body := bytes.Buffer{}

	for scanner.Scan() {
		_, err := fmt.Fprintln(&body, scanner.Text())
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Fprintln: %v\n", err)
		}

	}
	return strings.TrimSuffix(body.String(), "\n")
}
