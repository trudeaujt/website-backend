package blogposts_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/trudeaujt/blogposts"
)

func TestServer(t *testing.T) {
	posts := map[string]blogposts.Post{
		"a_title": {
			Title:       "A Title",
			Slug:        "a_title",
			Description: "Description",
			Tags:        []string{"one", "two"},
			Body:        "Body",
			Published:   true,
			Date:        "2023-01-17",
		},
		"a_title2": {
			Title:       "A Title2",
			Slug:        "a_title2",
			Description: "Description2",
			Tags:        []string{"two", "three"},
			Body:        "Body2",
			Published:   false,
			Date:        "2020-09-17",
		},
	}
	server := blogposts.NewBlogServer(posts)

	t.Run("it returns all the posts as JSON on GET", func(t *testing.T) {
		req := newAllPostsRequest()
		res := httptest.NewRecorder()
		server.ServeHTTP(res, req)

		assertStatus(t, res.Code, http.StatusOK)
		assertBody(t, res.Body.Bytes(), []byte(`{
	"a_title": {
		"title": "A Title",
		"slug": "a_title",
		"description": "Description",
		"tags": [
			"one",
			"two"
		],
		"published": true,
		"date": "2023-01-17"
	},
	"a_title2": {
		"title": "A Title2",
		"slug": "a_title2",
		"description": "Description2",
		"tags": [
			"two",
			"three"
		],
		"published": false,
		"date": "2020-09-17"
	}
}`,
		))
	})
	t.Run("it returns a single post as JSON on GET", func(t *testing.T) {
		req := newSinglePostRequest("a_title")
		res := httptest.NewRecorder()
		server.ServeHTTP(res, req)

		assertStatus(t, res.Code, http.StatusOK)
		assertBody(t, res.Body.Bytes(), []byte(`{
		"title": "A Title",
		"slug": "a_title",
		"description": "Description",
		"tags": [
			"one",
			"two"
		],
		"published": true,
		"date": "2023-01-17"
	}`))
		errReq := newSinglePostRequest("bad-slug")
		errRes := httptest.NewRecorder()
		server.ServeHTTP(errRes, errReq)

		assertStatus(t, errRes.Code, http.StatusNotFound)
	})
}

func newSinglePostRequest(slug string) *http.Request {
	url := fmt.Sprintf("/post/%s", slug)
	print(url)
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	return req
}

func assertBody(t *testing.T, got, want []byte) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", string(got), string(want))
	}
}

func assertStatus(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("incorrect status, got %d want %d", got, want)
	}
}

func newAllPostsRequest() *http.Request {
	req, _ := http.NewRequest(http.MethodGet, "/posts", nil)
	return req
}
