package blogposts_test

import (
	"fmt"
	"github.com/trudeaujt/blogposts"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestServer(t *testing.T) {
	posts := []blogposts.Post{
		{
			Title:       "A Title",
			Slug:        "a_title",
			Description: "Description",
			Tags:        []string{"one", "two"},
			Body:        "Body",
		},
		{
			Title:       "A Title2",
			Slug:        "a_title2",
			Description: "Description2",
			Tags:        []string{"two", "three"},
			Body:        "Body2",
		},
	}
	server := blogposts.NewBlogServer(posts)

	t.Run("It returns all the posts as JSON on GET", func(t *testing.T) {
		req := newAllPostsRequest()
		res := httptest.NewRecorder()
		server.ServeHTTP(res, req)

		fmt.Println(res.Body)
		assertStatus(t, res.Code, http.StatusOK)
		assertBody(t, res.Body.Bytes(), []byte(`[
	{
		"title": "A Title",
		"slug": "a_title",
		"description": "Description",
		"tags": [
			"one",
			"two"
		]
	},
	{
		"title": "A Title2",
		"slug": "a_title2",
		"description": "Description2",
		"tags": [
			"two",
			"three"
		]
	}
]`,
		))
	})
}

func assertBody(t *testing.T, got, want []byte) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
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
