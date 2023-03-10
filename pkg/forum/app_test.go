/*
*
Package forum_test provides API test for thumuht.
*/
package forum_test

import (
	"backend/pkg/forum"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"testing"
	"unicode"
)

const (
	TestPath = "127.0.0.1:9998"
	TestGQL  = "http://127.0.0.1:9998/query"
)

var utoken string

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

/**
utility. make gql requests with variable

TODO
*/
// func makeGQLRequestV(gs string, hdr map[string]string) *strings.Reader{
// 	return nil
// }

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

	newuserResp, err := http.DefaultClient.Do(gqlreq)
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
utility. send gql request & receive gql resp, and compare 'em
*/
func SendAndCompareGQL(req string, resp string, hdr map[string]string) (bool, error) {
	resp = KillWhitespaces(resp)

	tresp, err := SendAndGetGQL(req, hdr)
	if err != nil {
		return false, err
	}

	*tresp = KillWhitespaces(*tresp)

	println("Get Response\n", *tresp)

	return strings.Compare(*tresp, resp) == 0, nil
}

// set up test environment.
// for now, run forum server is ok..
func TestMain(m *testing.M) {
	app := forum.NewForum()

	// launch app. use goroutine, because Run() will block.
	go func() {
		app.Run(TestPath)
	}()

	m.Run()
}

func TestNewUser(t *testing.T) {
	newuser := `mutation {
		createUser(input: {
			loginName: "thumuht"
			password: "harmful"
		}) {
			id
			loginName
		}
	}`

	newuserResp := `{"data":{"createUser":{"id":1,"loginName":"thumuht"}}}`
	if ok, err := SendAndCompareGQL(newuser, newuserResp, nil); ok == false {
		t.Error(fmt.Errorf("cannot new user: %w", err))
	}
}

func TestLogin(t *testing.T) {
	login := `mutation {
		login(input: {
			loginName: "thumuht"
			password: "harmful"
		})
	}`

	s, err := SendAndGetGQL(login, nil)
	if err != nil {
		t.Error(fmt.Errorf("cannot login: %w", err))
	}
	re := regexp.MustCompile(`"login":"(.+)"`)
	sm := re.FindAllSubmatch([]byte(*s), -1)
	fmt.Printf("match: %q\n", sm)
	if sm == nil {
		t.Error("no token")
	}
	utoken = string(sm[0][1])
	fmt.Printf("token: %s\n", utoken)
}

func TestNewPost(t *testing.T) {
	newpost := `mutation {
		createPost(input: {
			userId: 1
			title: "go"
			content: "too good"
		}) {
			id
			title
			content
		}
	}`

	hdrs := map[string]string{"Token": utoken}
	newpostResp := `{"data":{
		"createPost":
		{"id": 1,
		"title": "go",
		"content": "too good"
		}}
	}`

	if ok, err := SendAndCompareGQL(newpost, newpostResp, hdrs); ok == false {
		t.Error(fmt.Errorf("cannot new post: %w", err))
	}
}

func TestNewComment(t *testing.T) {
	newcomment := `mutation {
		createComment(input: {
			userId: 1
			postId: 1
			content: "I agree with u!!"
		}) {
			id
			content
		}
	}`

	hdrs := map[string]string{"Token": utoken}
	newcommentResp := `{"data":{
		"createComment": {
			"id": 1,
			"content": "I agree with u!!"
			}
		}
	}`
	if ok, err := SendAndCompareGQL(newcomment, newcommentResp, hdrs); ok == false {
		t.Error(fmt.Errorf("cannot new comment: %w", err))
	}
}

func NewPost(title string, content string) error {
	newpost := `mutation {
		createPost(input: {
			userId: 1
			title: "$1$"
			content: "$2$"
		}) {
			id
			title
			content
		}
	}`
	newpost = strings.Replace(
		strings.Replace(newpost, "$1$", title, 1),
		"$2$",
		content,
		1,
	)

	hdrs := map[string]string{"Token": utoken}
	_, err := SendAndGetGQL(newpost, hdrs)
	return err
}

func TestPostPaging(t *testing.T) {
	alp := string("zxcvbnmlkjhgfdsapoiuytrewq")
	for idx, cha := range alp {
		err := NewPost(string(cha), fmt.Sprintf("%d", idx))
		if err != nil {
			t.Errorf("insert post failed")
		}
	}

	queryPosts := `query{
		posts(input: {
			limit: 10
			offset: 0
			orderBy: title
			order: ASC
		}) {
			title
			content
		}
	}`

	queryresp := `
	{"data":{"posts":[{"title":"a","content":"15"},{"title":"b","content":"4"},
	{"title":"c","content":"2"},{"title":"d","content":"13"},{"title":"e","content":"23"},
	{"title":"f","content":"12"},{"title":"g","content":"11"},{"title":"go","content":"too good"},
	{"title":"h","content":"10"},{"title":"i","content":"18"}]}}
	`

	resp, err := SendAndCompareGQL(queryPosts, queryresp, nil)
	if err != nil {
		t.Errorf("can not query")
	}
	if resp == false {
		t.Errorf("incorrect order")
	}

	queryPosts = `query{
		posts(input: {
			limit: 10
			offset: 0
			orderBy: created_at
			order: ASC
		}) {
			title
			content
		}
	}`

	queryresp = `
	{"data":{"posts":[{"title":"go","content":"too good"},{"title":"z","content":"0"},
	{"title":"x","content":"1"},{"title":"c","content":"2"},{"title":"v","content":"3"},
	{"title":"b","content":"4"},{"title":"n","content":"5"},{"title":"m","content":"6"},
	{"title":"l","content":"7"},{"title":"k","content":"8"}]}}
	`

	resp, err = SendAndCompareGQL(queryPosts, queryresp, nil)
	if err != nil {
		t.Errorf("can not query")
	} else if resp == false {
		t.Error("incorrect time result")
	}
}

// TODO(wj): filesystem test