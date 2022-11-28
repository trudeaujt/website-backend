package blogposts


type Slug string

func (s Slug) MarshalJSON() ([]byte, error) {
	return nil, nil
}