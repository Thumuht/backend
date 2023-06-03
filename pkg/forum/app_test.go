/*
*
Package forum_test provides API test for thumuht.
*/
package forum_test

import (
	"backend/pkg/forum"
	"backend/pkg/utils"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"regexp"
	"strings"
	"sync"
	"testing"

	"github.com/spf13/viper"
)

// set up test environment.
// for now, run forum server is ok..
func TestMain(m *testing.M) {
	app = forum.NewForum()

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
		}){
			token
			userId
		}
	}`

	s, err := SendAndGetGQL(login, nil)
	fmt.Print(*s)
	if err != nil {
		t.Error(fmt.Errorf("cannot login: %w", err))
	}
	re := regexp.MustCompile(`"login":{"token":"(.+)","u`)
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
		// abort all tests
		t.FailNow()
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
	{"data":{"posts":[{"title":"z","content":"0"},{"title":"x","content":"1"},
	{"title":"c","content":"2"},{"title":"v","content":"3"},{"title":"b","content":"4"},
	{"title":"n","content":"5"},{"title":"m","content":"6"},{"title":"l","content":"7"},
	{"title":"k","content":"8"},{"title":"j","content":"9"}]}}
	`

	resp, err = SendAndCompareGQL(queryPosts, queryresp, nil)
	if err != nil {
		t.Errorf("can not query")
	} else if resp == false {
		t.Error("incorrect time result")
	}
}

func TestCache(t *testing.T) {
	browse := `mutation {
		likePost(input: 1)
	}`

	getlike := `query {
		postDetail(input: 1) {
			like
		}
	}`

	hdrs := map[string]string{"Token": utoken}

	var wg sync.WaitGroup
	for i := 0; i < 20000; i++ {
		wg.Add(1)
		// i must be in parameter
		go func() {
			defer wg.Done()
			_, _ = SendAndGetGQL(browse, hdrs)
		}()
	}

	wg.Wait()

	app.Cache.PostLike.Flush()
	ok, err := SendAndCompareGQL(getlike, `{"data":{"postDetail":{"like":20000}}}`, nil)
	if err != nil || !ok {
		t.Error("testcache failed")
	}
}

func TestCacheCSP(t *testing.T) {
	browse := `mutation {
		likePost(input: 1)
	}`

	getlike := `query {
		postDetail(input: 1) {
			like
		}
	}`

	hdrs := map[string]string{"Token": utoken}

	var wg sync.WaitGroup
	for i := 0; i < 20000; i++ {
		wg.Add(1)
		// i must be in parameter
		go func() {
			defer wg.Done()
			_, _ = SendAndGetGQL(browse, hdrs)
		}()
	}

	wg.Wait()

	app.Cache.PostLike.Flush()
	ok, err := SendAndCompareGQL(getlike, `{"data":{"postDetail":{"like":40000}}}`, nil)
	if err != nil || !ok {
		t.Error("testcache failed")
	}
}

func TestFS(t *testing.T) {
	s := utils.GenRandStr(2000)
	err := os.WriteFile("testthisfile", []byte(s), 0777)
	if err != nil {
		t.Error("cannot test fs.")
	}
	viper.Set("fs_route", "./testfs")
	t.Cleanup(func() {
		os.Remove("testthisfile")
	})

	uploadfieldsK := []string{
		"operations",
		"map",
	}
	uploadfieldsV := []string{
		`{"query": "mutation($req: Upload!) { fileUpload(input: {parentId: 1,parentType: post,upload: $req}) }","variables": { "req": null } }`,
		`{ "0": ["variables.req"] }`,
	}

	files := map[string]string{
		"0": "testthisfile",
	}

	req := makeMultiPartRequest(uploadfieldsK, uploadfieldsV, files)
	req.Header.Add("Token", utoken)
	str, _ := ReqToRespStr(req)
	if strings.Compare(*str, `{"data":{"fileUpload":true}}`) != 0 {
		t.Error("upload failed")
	}

	// now query post, should have attachment
	qpost := `query {
		postDetail(input: 1) {
			attachment {
				parentId
				parentType
				fileName
			}
		}
	}`
	ok, _ := SendAndCompareGQL(qpost, `{"data":{"postDetail":{"attachment":[{"parentId":1,"parentType":"post","fileName":"testthisfile"}]}}}`, nil)
	if !ok {
		t.Error("attachment did not goto database")
	}

	// check file sanity
	resp, err := http.Get("http://" + path.Join(TestPath, "fs/testthisfile"))
	if err != nil {
		t.Error("cannot get fs")
	}
	t.Cleanup(func() {
		os.RemoveAll("./testfs")
		os.Remove("./testfs")
	})
	defer resp.Body.Close()

	results, _ := io.ReadAll(resp.Body)
	if strings.Compare(string(results), s) != 0 {
		t.Error("file do not identical")
	}

}

// Test @login API before THIS
func TestLogout(t *testing.T) {
	logout := `mutation {
		logout
	}`

	hdrs := map[string]string{"Token": utoken}
	ok, err := SendAndCompareGQL(logout, `{"data":{"logout":true}}`, hdrs)
	if err != nil || !ok {
		t.Error("cannot send gql")
	}

	np := `mutation {
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
	ok, err = SendAndCompareGQL(np, fmt.Sprintf(`{"errors":[{"message":"no token %s access denied","path":["createPost"]}],"data":null}`, utoken), hdrs)
	if err != nil || !ok {
		t.Error("cannot get gql")
	}

}
