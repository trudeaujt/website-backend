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
	Title       string   `json:"title"`
	Slug        string   `json:"slug"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
	Body        string   `json:"-"`
	Published   bool     `json:"published"`
	Date        string   `json:"date"`
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

//func (p Post) MarshalJSON() ([]byte, error) {
//	type PostAlias Post
//
//	aux := struct {
//		PostAlias
//		Body string `json:"-"`
//	}{
//		PostAlias: PostAlias(p),
//		Body:      p.Body,
//	}
//
//	return json.Marshal(aux)
//}
