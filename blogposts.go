package blogposts

import (
	"io/fs"
)

func NewPostsFromFS(fileSystem fs.FS) (map[string]Post, error) {
	dir, err := fs.ReadDir(fileSystem, ".")

	if err != nil {
		return nil, err
	}

	posts := make(map[string]Post)

	for _, file := range dir {
		post, err := getPost(fileSystem, file.Name())
		if err != nil {
			return nil, err
		}
		posts[post.Slug] = post
	}
	return posts, nil
}

func getPost(fileSystem fs.FS, fileName string) (Post, error) {
	postFile, err := fileSystem.Open(fileName)
	if err != nil {
		return Post{}, err
	}
	defer func(postFile fs.File) {
		e := postFile.Close()
		if e != nil {
			err = e
		}
	}(postFile)

	post, err := newPost(postFile)
	return post, err
}
