package forum_test

import (
	"backend/pkg/forum"
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"unicode"
)

const (
	TestPath = "127.0.0.1:9998"
	TestGQL  = "http://127.0.0.1:9998/query"
)

var utoken string
var app *forum.App

/*
utility. make graphql requests

gql request looks like:
{"query":"mutation {\n  createUser(input:{\n    loginName: \"thumuht\"\n    password: \"harmful\"\n  }) {\n    id\n    loginName\n  }\n}"}

but we only need inside mutation..
*/
func makeGQLRequest(gs string) *strings.Reader {
	// DO NOT escape strings, RFC 7159 section #7
	//
	// https://www.rfc-editor.org/rfc/rfc7159#section-7
	//
	// send \n, rather than lf(0a)
	//
	// https://stackoverflow.com/questions/50054666/golang-not-escape-a-string-variable
	return strings.NewReader(fmt.Sprintf(`{"query": %#v}`, gs))
}

/*
*
utility. make a multipart/form-data request
*/
func makeMultiPartRequest(fieldsK, fieldsV []string, files map[string]string) *http.Request {
	var req bytes.Buffer
	multiw := multipart.NewWriter(&req)
	for k := range fieldsK {
		err := multiw.WriteField(fieldsK[k], fieldsV[k])
		if err != nil {
			panic("cannot new form")
		}
	}

	for k, v := range files {
		filewriter, err := multiw.CreateFormFile(k, v)
		if err != nil {
			panic("cannot new filewriter")
		}
		file, err := os.Open(v)
		if err != nil {
			panic("file cannot open")
		}
		fsize, err := io.Copy(filewriter, file)
		if fstat, _ := file.Stat(); fstat.Size() != fsize || err != nil {
			panic("file cannot upload")
		}
	}
	multiw.Close()

	request, err := http.NewRequest("POST", TestGQL, &req)
	if err != nil {
		panic("cannot new post")
	}

	request.Header.Set("Content-Type", multiw.FormDataContentType())

	return request
}

func ReqToRespStr(req *http.Request) (*string, error) {
	newuserResp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer newuserResp.Body.Close()

	newuserBuf, err := io.ReadAll(newuserResp.Body)
	if err != nil {
		return nil, err
	}
	newuserS := string(newuserBuf)
	return &newuserS, nil
}

/*
*
utility. remove whitespaces for compare

https://stackoverflow.com/questions/32081808/strip-all-whitespace-from-a-string
*/
func KillWhitespaces(s string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, s)
}

func SendAndGetGQL(req string, hdr map[string]string) (*string, error) {
	gqlreq, err := http.NewRequest("POST", TestGQL, makeGQLRequest(req))
	if err != nil {
		return nil, err
	}

	// header
	gqlreq.Header.Set("content-type", "application/json")
	for k, v := range hdr {
		gqlreq.Header.Set(k, v)
	}

	return ReqToRespStr(gqlreq)
}

/*
*
utility. send gql request & receive gql resp, and compare 'em
return true if ok
*/
func SendAndCompareGQL(req string, resp string, hdr map[string]string) (bool, error) {
	resp = KillWhitespaces(resp)

	tresp, err := SendAndGetGQL(req, hdr)
	if err != nil {
		return false, err
	}

	*tresp = KillWhitespaces(*tresp)

	println("Get Response\n", *tresp, resp, strings.Compare(*tresp, resp))

	return strings.Compare(*tresp, resp) == 0, nil
}
