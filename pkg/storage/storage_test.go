package storage

import (
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"
	"testing/fstest"
)

func TestLoadBackend(t *testing.T) {
	// init the folder storage
	fs := NewFolderStorage()

	// setup the in memory filesystem
	testdata := "This is the data"
	mfs := &fstest.MapFS{"Nokia/SRLinux/bar.txt": {Data: []byte(testdata)}}

	err := fs.LoadBackend(mfs)
	if err != nil {
		t.Errorf("error loading storage backend: %v", err)
	}

	// ###### request existing file
	// init a httptest.responserecorder
	rr := httptest.NewRecorder()
	fs.Handle(rr, "Nokia/SRLinux/bar.txt")
	if rr.Code != 200 || rr.Body.String() != testdata {
		t.Errorf("error on handling filedelivery via storage. HTTP Code: {expected: %d, actual %d}, body: {expected: %s, actual: %s }", 200, rr.Code, testdata, rr.Body.String())
	}

	// ###### request non existing file
	// init a new httptest.responserecorder
	rr = httptest.NewRecorder()
	fs.Handle(rr, "Nokia/SRLinux/foo.txt")
	if rr.Code != 500 || !strings.Contains(rr.Body.String(), "unable to retrieve file") {
		t.Errorf("error on handling file delivery via storage. HTTP Code: {expected: %d, actual %d}", 500, rr.Code)
	}

	// ###### write throws error
	// init a new httptest.responserecorder
	mrr := NewMyResponseRecorder()

	fs.Handle(mrr, "Nokia/SRLinux/bar.txt")
	if mrr.Code != 500 {
		t.Errorf("error on handling file delivery via storage. HTTP Code: {expected: %d, actual %d}", 500, mrr.Code)
	}

}

// MyResponseRecorder modified httptest.ResponseRecorder, where the Write method can throw erros
type MyResponseRecorder struct {
	httptest.ResponseRecorder
}

func (mrr *MyResponseRecorder) Write(buf []byte) (int, error) {
	return 0, fmt.Errorf("Error")
}

// NewRecorder returns an initialized ResponseRecorder.
func NewMyResponseRecorder() *MyResponseRecorder {
	return &MyResponseRecorder{
		ResponseRecorder: *httptest.NewRecorder(),
	}
}
