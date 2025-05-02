package gcp

import "fmt"

type errUnexpectedStatusCode struct {
	statusCode int
}

func (e errUnexpectedStatusCode) Error() string {
	return fmt.Sprintf("unexpected status code: %d", e.statusCode)
}
