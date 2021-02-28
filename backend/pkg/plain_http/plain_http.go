package plain_http

import (
	"bytes"
	"fmt"
	"strconv"
	"text/template"
)

// MakeHTTPTemplates check if data can be converted to HTTP request and return templates mask
func MakeHTTPTemplates(data map[string]string) (map[string]*template.Template, error) {
	templates := make(map[string]*template.Template, 4)
	methodTemplate, err := template.New("method").Parse(data["methodTemplate"])
	if err != nil {
		return templates, err
	}
	templates["methodTemplate"] = methodTemplate
	urlTemplate, err := template.New("url").Parse(data["urlTemplate"])
	if err != nil {
		return templates, err
	}
	templates["urlTemplate"] = urlTemplate
	headersTemplate, err := template.New("headers").Parse(data["headersTemplate"])
	if err != nil {
		return templates, err
	}
	templates["headersTemplate"] = headersTemplate
	bodyTemplate, err := template.New("body").Parse(data["bodyTemplate"])
	if err != nil {
		return templates, err
	}
	templates["bodyTemplate"] = bodyTemplate
	return templates, nil
}

// BuildHTTP use templates mask and data to make plain http request
func BuildHTTP(templates map[string]*template.Template, data map[string]interface{}) (plainRequest string, err error) {
	var methodBuff bytes.Buffer
	err = templates["methodTemplate"].Execute(&methodBuff, data)
	if err != nil {
		return
	}
	method := methodBuff.String()
	var urlBuff bytes.Buffer
	err = templates["urlTemplate"].Execute(&urlBuff, data)
	if err != nil {
		return
	}
	url := urlBuff.String()
	var headersBuff bytes.Buffer
	err = templates["headersTemplate"].Execute(&headersBuff, data)
	if err != nil {
		return
	}
	headers := headersBuff.String()
	var bodyBuff bytes.Buffer
	err = templates["bodyTemplate"].Execute(&bodyBuff, data)
	if err != nil {
		return
	}
	body := bodyBuff.String()
	cl := strconv.Itoa(len(body))
	plainRequest = fmt.Sprintf("%v %v HTTP/1.1\nContent-Length: %v\n%v\n\n%v", method, url, cl, headers, body)
	return
}
