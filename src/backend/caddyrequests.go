package main

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func postCaddy(path string, data string) {
	resp, err := http.Post("http://127.0.0.1:2019"+path, "application/json", strings.NewReader(data))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		return
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	bodyString := string(bodyBytes)
	err = errors.New("server respond not 200.\npath=" + path + "\nmsg=" + bodyString)
	panic(err)
}

func getCaddy(path string) (string, error) {
	resp, err := http.Get("http://127.0.0.1:2019" + path)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}
		bodyString := string(bodyBytes)
		return bodyString, nil
	}
	err = errors.New("server respond not 200")
	return "", err
}

func initCaddy() {
	postCaddy("/load", `{"apps": {"http": {"servers": {
		"server": {
			"listen": [":80"],
			"routes": []
		}
	}}}}`)
	postCaddy("/config/apps/http/servers/server/routes/...", `[{
		"group": "",
		"terminal": false,
		"match": [
			{"host": ["direct-to-me.com"], "path": ["/api/*"]}
		],
		"handle": [
			{
				"handler": "reverse_proxy",
				"upstreams": [{"dial": "127.0.0.1:1323"}]
			}
		]
	},
	{
		"group": "",
		"terminal": false,
		"match": [
			{"host": ["direct-to-me.com"]}
		],
		"handle": [
			{
				"handler": "file_server",
				"root": "/var/caddy/",
				"index_names": ["index.html"]
			}
		]
	}]`)
}

func checkUrl(urlString string) *url.URL {
	urlParsed, err := url.Parse(urlString)
	if err != nil {
		panic(err)
	}
	if urlParsed.Scheme != "http" && urlParsed.Scheme != "https" {
		panic(errors.New("scheme is not correct"))
	}
	if urlParsed.Hostname() == "" {
		panic(errors.New("scheme is empty"))
	}
	return urlParsed
}

func baseAddRedirect(from string, to string, handlerSettingText string) {
	urlFrom := checkUrl(from)
	checkUrl(to)
	postCaddy("/config/apps/http/servers/server/routes/...", `[{
		"group": "",
		"terminal": false,
		"match": [
			{"host": ["`+urlFrom.Hostname()+`"], "path": ["`+urlFrom.Path+`"]}
		],
		"handle": [
			{
				"handler": "static_response",
				`+handlerSettingText+`,
				"close": true
			}	
		]
	}]`)
}

func addSafeRedirect(from string, to string) {
	baseAddRedirect(from, to, `
		"status_code": "307",
		"headers": {
			"Location": ["`+to+`"],
			"Cache-Control": ["no-cache"]
		},
		"body": "<head><meta http-equiv=\"refresh\" content=\"0; URL=`+to+`\"/></head>"
	`)
}

func addStableRedirect(from string, to string) {
	baseAddRedirect(from, to, `
		"status_code": "308",
		"headers": {
			"Location": ["`+to+`"],
			"Cache-Control": ["public"]
		},
		"body": "<head><meta http-equiv=\"refresh\" content=\"0; URL=`+to+`\"/></head>"
	`)
}

func addNoSniffRedirect(from string, to string, canonical string) {
	baseAddRedirect(from, to, `
		"status_code": "200",
		"headers": {
			"Cache-Control": ["no-cache"]
		},
		"body": "<head><link rel=\"canonical\" href=\"`+canonical+`\"/><meta http-equiv=\"refresh\" content=\"0; URL=`+to+`\"/></head>"
	`)
}

func checkCaddy() string {
	resp, err := getCaddy("/config/")
	if err != nil {
		panic(err)
	}
	return resp
}
