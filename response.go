package main

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

// a response is a wrapper around an HTTP response;
// it contains the request value for context.
type response struct {
	request request

	status     string
	statusCode int
	headers    []string
	body       []byte
	err        error
}

// String returns a string representation of the request and response
func (r response) String() string {
	b := &bytes.Buffer{}

	b.WriteString(r.request.URL())
	b.WriteString("\n\n")

	b.WriteString(fmt.Sprintf("> %s %s HTTP/1.1\n", r.request.method, r.request.path))

	// request headers
	for _, h := range r.request.headers {
		b.WriteString(fmt.Sprintf("> %s\n", h))
	}
	b.WriteString("\n")

	// status line
	b.WriteString(fmt.Sprintf("< HTTP/1.1 %s\n", r.status))

	// response headers
	for _, h := range r.headers {
		b.WriteString(fmt.Sprintf("< %s\n", h))
	}
	b.WriteString("\n")

	// body
	b.Write(r.body)

	return b.String()
}

func (r response) StringNoHeaders() string {
	b := &bytes.Buffer{}

	b.Write(r.body)

	return b.String()
}

// save write a request and response output to disk
func (r response) save(pathPrefix string, staticOutput, noHeaders bool) (string, error) {
	content := []byte(r.String())
	if noHeaders {
		content = []byte(r.StringNoHeaders())
	}

	parts := []string{pathPrefix}
	parts = append(parts, r.request.Hostname())

	if staticOutput {
		fName := strings.ReplaceAll(r.request.Hostname(), ".", "_")
		parts = append(parts, fName)
	} else {
		checksum := sha1.Sum(content)
		parts = append(parts, fmt.Sprintf("%x", checksum))
	}
	p := path.Join(parts...)

	if _, err := os.Stat(path.Dir(p)); os.IsNotExist(err) {
		err = os.MkdirAll(path.Dir(p), 0750)
		if err != nil {
			return p, err
		}
	}

	err := ioutil.WriteFile(p, content, 0640)
	if err != nil {
		return p, err
	}

	return p, nil
}
