package factorial

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"

	"golang.org/x/net/publicsuffix"
)

type factorialClient struct {
	http.Client
	employee_id int64
	period_id   int64
}

func NewFactorialClient(email, password string) *factorialClient {
	c := new(factorialClient)
	options := cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	}
	jar, _ := cookiejar.New(&options)
	c.Client = http.Client{Jar: jar}
	c.login(email, password)
	return c
}

func (c *factorialClient) login(email, password string) {
	getCSRFToken := func(resp *http.Response) string {
		data, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		start := strings.Index(string(data), "<meta name=\"csrf-token\" content=\"") + 33
		end := strings.Index(string(data)[start:], "\" />")
		return string(data)[start : start+end]
	}

	getLoginError := func(resp *http.Response) string {
		data, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		start := strings.Index(string(data), "<div class=\"flash flash--wrong\">") + 32
		if start < 0 {
			return ""
		}
		end := strings.Index(string(data)[start:], "</div>")
		if start < 0 || end-start > 100 {
			return ""
		}
		fmt.Println(start, end, len(string(data)))
		return string(data)[start : start+end]
	}

	resp, _ := c.Get("https://api.factorialhr.com/users/sign_in")
	csrf_token := getCSRFToken(resp)
	fmt.Println(csrf_token)
	body := url.Values{
		"authenticity_token": {csrf_token},
		"return_host":        {"factorialhr.es"},
		"user[email]":        {email},
		"user[password]":     {password},
		"user[remember_me]":  {"0"},
		"commit":             {"Sign in"},
	}
	resp, _ = c.PostForm("https://api.factorialhr.com/users/sign_in", body)
	if err := getLoginError(resp); err != "" {
		log.Fatal(err)
	}
}
