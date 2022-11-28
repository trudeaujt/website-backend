package blogposts_test

import (
	"encoding/json"
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
			Description: "Description",
			Tags:        []string{"one", "two"},
			Body:        "Body",
		},
		{
			Title:       "A Title2",
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
		assertBody(t, res.Body.Bytes(), posts)
	})
}

func assertBody(t *testing.T, got []byte, want []blogposts.Post) {
	t.Helper()
	j, _ := json.Marshal(want)
	//j = append(j, 10) // add a newline. json.Marshal doesn't add one, while json.Encode does.
	if !reflect.DeepEqual(got, j) {
		t.Errorf("got %v want %v", got, j)
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
