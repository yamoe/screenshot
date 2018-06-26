package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// LoginData is get csrf token in form
func LoginData(urlstr string) (csrf string, cookies []*http.Cookie, err error) {
	var res *http.Response
	res, err = http.Get(urlstr)
	if err != nil {
		err = fmt.Errorf("get - %v", err)
		return
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		err = fmt.Errorf("get - status code error: %d %s", res.StatusCode, res.Status)
		return
	}

	cookies = res.Cookies()

	var doc *goquery.Document
	doc, err = goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		err = fmt.Errorf("goquery - %v", err)
		return
	}

	doc.Find("input#csrf_token").Each(func(i int, s *goquery.Selection) {
		val, exist := s.Attr("value")
		if exist {
			csrf = val
		}
	})

	if len(csrf) == 0 {
		err = fmt.Errorf("csrf is empty")
		return
	}
	return
}

// Login is get logged-in cookies using POST
func Login(urlstr string, username string, password string, csrf string, cookies []*http.Cookie) (resCookies []*http.Cookie, err error) {
	encodeData := url.Values{
		"csrf_token": []string{csrf},
		"username":   []string{username},
		"password":   []string{password},
	}.Encode()
	data := strings.NewReader(encodeData)

	var req *http.Request
	req, err = http.NewRequest("POST", urlstr, data)
	if err != nil {
		err = fmt.Errorf("new request - %v", err)
		return
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(encodeData)))
	req.Header.Add("referer", urlstr)

	for _, c := range cookies {
		cookie := http.Cookie{Name: c.Name, Value: c.Value}
		req.AddCookie(&cookie)
	}

	client := &http.Client{
		// prevent redirect
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}}

	var res *http.Response
	res, err = client.Do(req)
	if err != nil {
		err = fmt.Errorf("post - %v", err)
		return
	}
	defer res.Body.Close()

	if res.StatusCode != 302 && res.StatusCode != 200 {
		body, _ := ioutil.ReadAll(res.Body)
		err = fmt.Errorf("res - status code error: %d %s\n%s",
			res.StatusCode, res.Status, string(body))
		return
	}

	resCookies = res.Cookies()
	return
}
