package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"testing"

	clienttests "github.com/jfrog/jfrog-client-go/utils/tests"
)

const versionApiResponse = `
{
  "version": "2.2.2"
}
`
const searchApiResponse = `
{
    "results": [
        {
            "repo": "repo1",
            "path": "a",
            "name": "srikrishna",
            "type": "file",
            "size": 11,
            "created": "2020-06-25T10:21:17.488+03:00",
            "modified": "2020-06-25T10:21:17.496+03:00",
            "actual_md5": "e53098d3d8ee1f5eb38c2ec3c783ef3d",
            "actual_sha1": "a",
            "properties": [
                {
                    "key": "ca",
                    "value": "1"
                },
                {
                    "key": "build.name",
                    "value": "myBuild"
                }
            ]
        },
        {
            "repo": "repo1",
            "path": "a",
            "name": "sathish",
            "type": "file",
            "size": 11,
            "created": "2020-06-25T10:21:17.502+03:00",
            "modified": "2020-06-25T10:21:17.510+03:00",
            "actual_md5": "e53098d3d8ee1f5eb38c2ec3c783ef3d",
            "actual_sha1": "a",
            "properties": [
                {
                    "key": "c",
                    "value": "1"
                },
                {
                    "key": "build.name",
                    "value": "myBuild"
                },
                {
                    "key": "build.number",
                    "value": "1"
                }
            ]
        },
        {},
        {
            "repo": "repo1",
            "path": "a",
            "name": "no.in",
            "type": "file",
            "size": 11,
            "created": "2020-06-25T10:21:17.502+03:00",
            "modified": "2020-06-25T10:21:17.510+03:00",
            "actual_md5": "e53098d3d8ee1f5eb38c2ec3c783ef3d",
            "actual_sha1": "a",
            "properties": [
                {
                    "key": "c",
                    "value": "1"
                }
            ]
        },
        {},
        {
            "repo": "repo1",
            "path": "a",
            "name": "yes.in",
            "type": "file",
            "size": 11,
            "created": "2020-06-25T10:21:17.502+03:00",
            "modified": "2020-06-25T10:21:17.510+03:00",
            "actual_md5": "e53098d3d8ee1f5eb38c2ec3c783ef3d",
            "actual_sha1": "b",
            "properties": [
                {
                    "key": "css",
                    "value": "1"
                },
                {
                    "key": "build.name",
                    "value": "myBuild"
                }
            ]
        },
        {
            "repo": "repo1",
            "path": "a",
            "name": "no.in",
            "type": "file",
            "size": 11,
            "created": "2020-06-25T10:21:17.502+03:00",
            "modified": "2020-06-25T10:21:17.510+03:00",
            "actual_md5": "e53098d3d8ee1f5eb38c2ec3c783ef3d",
            "actual_sha1": "b",
            "properties": [
                {
                    "key": "css",
                    "value": "1"
                }
            ]
        },
        {
            "repo": "repo1",
            "path": "a",
            "name": "c1.in",
            "type": "file",
            "size": 11,
            "created": "2020-06-25T10:21:17.502+03:00",
            "modified": "2020-06-25T10:21:17.510+03:00",
            "actual_md5": "e53098d3d8ee1f5eb38c2ec3c783ef3d",
            "actual_sha1": "c",
            "properties": [
                {
                    "key": "c",
                    "value": "1"
                }
            ]
        }
    ],
    "range": {
        "start_pos": 0,
        "end_pos": 5,
        "total": 5
    }
}
`

func getVersionHandler(w http.ResponseWriter, r *http.Request) {
	_, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, _ = fmt.Fprint(w, versionApiResponse)
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	_, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, _ = fmt.Fprint(w, searchApiResponse)

}

func mockServer(t *testing.T) int {
	handlers := clienttests.HttpServerHandlers{}
	handlers["/artifactory/api/system/version"] = getVersionHandler
	handlers["/artifactory/api/search/aql"] = searchHandler
	handlers["/"] = http.NotFound

	port, err := clienttests.StartHttpServer(handlers)
	if err != nil {
		t.Fatal(err)
	}
	return port
}

func Test_searchArtifacts(t *testing.T) {
	port := mockServer(t)
	tests := []struct {
		name    string
		wantErr bool
	}{
		{"mock searchApi",
			false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := searchArtifacts("127.0.0.1", strconv.Itoa(port)); (err != nil) != tt.wantErr {
				t.Errorf("searchArtifacts() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
