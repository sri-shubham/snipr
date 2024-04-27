package shorten

import (
	"net/url"
)

type Shortener interface {
	Shorten(url url.URL) url.URL
	ShortenCustom(url url.URL, customString string) url.URL
}
