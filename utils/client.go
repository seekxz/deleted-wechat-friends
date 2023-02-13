package utils

import (
	"crypto/tls"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"time"
)

type BaseRequest struct {
	// XMLName struct{} `xml:"xml" json:"-"`
	XMLName xml.Name `xml:"xml" json:"-"`

	Ret        int    `xml:"Ret" json:"ret"`
	Message    string `xml:"Message" json:"message"`
	Skey       string `xml:"Skey" json:"skey"`
	Wxsid      string `xml:"Wxsid" json:"wxsid"`
	Wxuin      string `xml:"Wxuin" json:"wxuin"`
	PassTicket string `xml:"PassTicket" json:"passTicket"`

	DeviceID string `xml:"DeviceID" json:"deviceID"`
}

func NewClient() (client *http.Client) {
	// transport
	transport := &http.Transport{}
	transport.ResponseHeaderTimeout = 1 * time.Second
	transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	// cookie jar
	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Println(err.Error())
	}

	client = &http.Client{
		Transport: transport,
		Jar:       jar,
		Timeout:   1 * time.Second,
	}

	return client
}

// CallRequest
// uri: https://login.wx.qq.com/cgi-bin/mmwebwx-bin/login
// name: login
// reqs: nil
// body: nil
func CallRequestName(uri, name string, reqs interface{}, body io.Reader) (reader io.Reader, err error) {
	apiUrl := fmt.Sprintf("%s/%s?pass_ticket=%s&skey=%s&r=%v", uri, name, reqs.(*BaseRequest).PassTicket, reqs.(*BaseRequest).Skey, time.Now().Unix())

	method := "GET"
	if body != nil {
		method = "POST"
	}

	req, err := http.NewRequest(method, apiUrl, body)
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	resp, err := Client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	reader = resp.Body.(io.Reader)

	return
}
