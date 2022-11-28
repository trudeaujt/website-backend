package blogposts

import "strings"
import "strconv"


type Slug string

func (s Slug) MarshalJSON() ([]byte, error) {
	jsonSlug := strings.ReplaceAll(string(s), " ", "_")
    quotedJsonSlug := strconv.Quote(jsonSlug)
	return []byte(quotedJsonSlug), nil
}