package url

import (
	"log"
	"net/url"
)

func Join(base string, paths ...string) string {
	p, err := url.JoinPath(base, paths...)
	if err != nil {
		log.Println(err)
		return p
	}

	return p
}
