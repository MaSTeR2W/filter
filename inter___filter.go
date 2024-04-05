package filter

import "net/url"

type Filter interface {
	validate_and_construct(v url.Values, language string) (string, error)
}
