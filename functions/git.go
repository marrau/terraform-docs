package functions

import (
	"log"

	git "gopkg.in/src-d/go-git.v4"
)

func gitURL(path string) string {
	repo, err := git.PlainOpenWithOptions(path, &git.PlainOpenOptions{DetectDotGit: true})
	if err != nil {
		log.Fatal(err)
	}

	remote, err := repo.Remote("origin")
	if err != nil {
		log.Fatal(err)
	}

	return remote.Config().URLs[0]
}
