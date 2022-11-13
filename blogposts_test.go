package websitebackend_test

import (
	"errors"
	"io/fs"
	"reflect"
	"testing"
	"testing/fstest"
)

func TestNewBlogPosts(t *testing.T) {
	t.Run("it returns posts equal to the number of files", func(t *testing.T) {
		fs := fstest.MapFS{
			"hello world.md":  {Data: []byte("Title: hi there")},
			"hello-world2.md": {Data: []byte("Title: ohayou gozaimasu")},
		}

		posts, err := websitebackend.NewPostsFromFS(fs)

		if err != nil {
			t.Fatal(err)
		}

		if len(posts) != len(fs) {
			t.Errorf("got %d posts want %d posts", len(posts), len(fs))
		}
	})
	t.Run("it returns errors when something goes wrong", func(t *testing.T) {
		fs := StubFailingFS{}

		posts, err := blogposts.NewPostsFromFS(fs)

		if err == nil {
			t.Error("expected an error, didn't get one")
		}
		if posts != nil {
			t.Errorf("didn't expect any posts, got %v", posts)
		}
	})
	t.Run("it returns the post title", func(t *testing.T) {
		fs := fstest.MapFS{
			"hello world.md":  {Data: []byte("Title: Post 1")},
			"hello-world2.md": {Data: []byte("Title: Post 2")},
		}

		posts, err := blogposts.NewPostsFromFS(fs)
		if err != nil {
			t.Fatal(err)
		}

		assertPost(t, posts[0], blogposts.Post{
			Title: "Post 1",
		})
	})
	t.Run("it returns the post description", func(t *testing.T) {
		const (
			firstBody = `Title: Post 1
Description: Description 1`
			secondBody = `Title: Post 2
Description: Description 2`
		)

		fs := fstest.MapFS{
			"hello world.md":  {Data: []byte(firstBody)},
			"hello-world2.md": {Data: []byte(secondBody)},
		}

		posts, err := blogposts.NewPostsFromFS(fs)
		if err != nil {
			t.Fatal(err)
		}

		assertPost(t, posts[0], blogposts.Post{
			Title:       "Post 1",
			Description: "Description 1",
		})
	})
	t.Run("it returns the post tags as a slice", func(t *testing.T) {
		const (
			firstBody = `Title: Post 1
Description: Description 1
Tags: tdd, go`
			secondBody = `Title: Post 2
Description: Description 2
Tags: tdd2, go2`
		)

		fs := fstest.MapFS{
			"hello world.md":  {Data: []byte(firstBody)},
			"hello-world2.md": {Data: []byte(secondBody)},
		}

		posts, err := blogposts.NewPostsFromFS(fs)
		if err != nil {
			t.Fatal(err)
		}

		assertPost(t, posts[0], blogposts.Post{
			Title:       "Post 1",
			Description: "Description 1",
			Tags:        []string{"tdd", "go"},
		})
	})
	t.Run("it returns the post body", func(t *testing.T) {
		const (
			firstBody = `Title: Post 1
Description: Description 1
Tags: tdd, go
---
First line
Second line`
			secondBody = `Title: Post 2
Description: Description 2
Tags: tdd2, go2
---
A
B
C`
		)

		fs := fstest.MapFS{
			"hello world.md":  {Data: []byte(firstBody)},
			"hello-world2.md": {Data: []byte(secondBody)},
		}

		posts, err := blogposts.NewPostsFromFS(fs)
		if err != nil {
			t.Fatal(err)
		}
		assertPost(t, posts[0], blogposts.Post{
			Title:       "Post 1",
			Description: "Description 1",
			Tags:        []string{"tdd", "go"},
			Body: `First line
Second line`,
		})
	})
	t.Run("it reads the metadata in any order", func(t *testing.T) {
		const (
			firstBody = `Description: Description 1
Title: Post 1
Tags: tdd, go
---
First line
Second line`
			secondBody = `Tags: tdd2, go2
Title: Post 2
Description: Description 2
---
A
B
C`
		)

		fs := fstest.MapFS{
			"hello world.md":  {Data: []byte(firstBody)},
			"hello-world2.md": {Data: []byte(secondBody)},
		}

		posts, err := blogposts.NewPostsFromFS(fs)
		if err != nil {
			t.Fatal(err)
		}
		assertPost(t, posts[0], blogposts.Post{
			Title:       "Post 1",
			Description: "Description 1",
			Tags:        []string{"tdd", "go"},
			Body: `First line
Second line`,
		})
	})
	t.Run("it returns an error when the metadata is wrong", func(t *testing.T) {
		const (
			firstBody = `Description: Description 1
Title: Post 1
this-should-throw-an-error: yes
Tags: tdd, go
---
First line
Second line`
		)

		fs := fstest.MapFS{
			"hello world.md": {Data: []byte(firstBody)},
		}

		_, err := blogposts.NewPostsFromFS(fs)
		if err.Error() != "invalid parameter: this-should-throw-an-error: yes" {
			t.Errorf("expected an invalid parameter error, got %v", err)
		}
	})
}

func assertPost(t *testing.T, got blogposts.Post, want blogposts.Post) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %+v want %+v", got, want)
	}
}

type StubFailingFS struct{}

func (s StubFailingFS) Open(name string) (fs.File, error) {
	return nil, errors.New("always fails")
}
