package blogposts_test

import (
	"errors"
	"github.com/trudeaujt/blogposts"
	"io/fs"
	"reflect"
	"testing"
	"testing/fstest"
)

func TestBlogposts(t *testing.T) {
	t.Run("it returns posts correctly", func(t *testing.T) {
		const (
			firstBody = `Title: Post 1
Description: Description 1
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

		testFs := fstest.MapFS{
			"hello world.md":  {Data: []byte(firstBody)},
			"hello-world2.md": {Data: []byte(secondBody)},
		}

		posts, err := blogposts.NewPostsFromFS(testFs)
		if err != nil {
			t.Fatal(err)
		}
		assertPosts(t, posts, map[string]blogposts.Post{
			"post_1": {
				Title:       "Post 1",
				Slug:        "post_1",
				Description: "Description 1",
				Tags:        []string{"tdd", "go"},
				Body: `First line
Second line`,
			},
			"post_2": {
				Title:       "Post 2",
				Slug:        "post_2",
				Description: "Description 2",
				Tags:        []string{"tdd2", "go2"},
				Body: `A
B
C`,
			}})
	})
	t.Run("it returns an error when the filesystem is bad", func(t *testing.T) {
		testFs := StubFailingFS{}

		posts, err := blogposts.NewPostsFromFS(testFs)

		if err == nil {
			t.Error("expected an error, didn't get one")
		}
		if posts != nil {
			t.Errorf("didn't expect any posts, got %v", posts)
		}
	})
	t.Run("it returns an error when the metadata is bad", func(t *testing.T) {
		const (
			firstBody = `Description: Description 1
Title: Post 1
this-should-throw-an-error: yes
Tags: tdd, go
---
First line
Second line`
		)

		testFs := fstest.MapFS{
			"hello world.md": {Data: []byte(firstBody)},
		}

		posts, err := blogposts.NewPostsFromFS(testFs)
		if err.Error() != "invalid parameter: this-should-throw-an-error: yes" {
			t.Errorf("expected an invalid parameter error, got %v", err)
		}
		if posts != nil {
			t.Errorf("didn't expect any posts, got %v", posts)
		}
	})
}

func assertPosts(t *testing.T, got, want map[string]blogposts.Post) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %+v want %+v", got, want)
	}
}

type StubFailingFS struct{}

func (s StubFailingFS) Open(string) (fs.File, error) {
	return nil, errors.New("always fails")
}
