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

func GetPathFromURL(urlStr string) string {
	u, err := url.Parse(urlStr)
	if err != nil {
		log.Println(err)
		return ""
	}

	return u.Path
}
