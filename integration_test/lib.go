package test

import (
	"fmt"
	"github.com/SasukeBo/configer"
	"github.com/SasukeBo/pmes-device-monitor/orm"
	"github.com/SasukeBo/pmes-device-monitor/router"
	"github.com/gavv/httpexpect"
	"net/http"
)

var host string
var AccessTokenCookie string

type Object map[string]interface{}

type Tester struct {
	E       *httpexpect.Expect
	Headers map[string]string
}

type Request struct {
	*httpexpect.Request
}

func (r *Request) GQLObject() *httpexpect.Object {
	return r.Expect().Status(http.StatusOK).JSON().Object()
}

// set Tester header
func (t *Tester) SetHeader(key string, value string) {
	t.Headers[key] = value
}

// send a POST Request with form Data
func (t *Tester) POST(path string, variables interface{}, pathargs ...interface{}) *Request {
	rr := t.E.POST(path, pathargs...).WithHeaders(t.Headers).WithHeader("Cookie", AccessTokenCookie).WithForm(variables)
	return &Request{rr}
}

// send a GET Request with query
func (t *Tester) GET(path string, variables interface{}, pathargs ...interface{}) *Request {
	rr := t.E.GET(path, pathargs...).WithHeaders(t.Headers).WithHeader("Cookie", AccessTokenCookie).WithQueryObject(variables)
	return &Request{rr}
}

// API1Admin post a api/v1/admin Request
func (t *Tester) API1Admin(query string, variables interface{}) *Request {
	return t.api("/api/v1/admin", query, variables)
}

// API1Admin post a api/v1 Request
func (t *Tester) API1(query string, variables interface{}) *Request {
	return t.api("/api/v1", query, variables)
}

func (t *Tester) api(path, query string, variables interface{}) *Request {
	payload := map[string]interface{}{
		"operationName": "",
		"query":         query,
		"variables":     variables,
	}

	rr := t.E.POST(path).WithHeaders(t.Headers).WithHeader("Cookie", AccessTokenCookie).WithJSON(payload)
	return &Request{rr}
}

// send a POST Request with form Data
func (t *Tester) Upload(path string, pathargs ...interface{}) *Request {
	rr := t.E.POST(path, pathargs...).WithHeaders(t.Headers).WithHeader("Cookie", AccessTokenCookie)
	return &Request{rr}
}

// new a Tester
func NewTester(t httpexpect.LoggerReporter) *Tester {
	tst := &Tester{}
	tst.E = httpexpect.New(t, host)
	tst.Headers = make(map[string]string)
	return tst
}

func init() {
	orm.DB.LogMode(false)
	tearDown()
	host = fmt.Sprintf("http://localhost:%v", configer.GetEnv("port"))
	go router.Start()
	orm.DB.LogMode(true)
}

func tearDown() {
	var tables = []string{}

	for _, name := range tables {
		cleanTable(name)
	}
}

func cleanTable(tbName string) {
	orm.DB.Exec(fmt.Sprintf("DELETE FROM %s WHERE 1 = 1", tbName))
}
