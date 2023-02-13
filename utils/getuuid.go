package utils

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"strconv"
	"time"
)

const APP_ID = "wx782c26e4c19acffb"
const REDIRECT_URI = "https://wx.qq.com/cgi-bin/mmwebwx-bin/webwxnewloginpage"
const JS_LOGIN_URL = "https://login.wx.qq.com/jslogin"

// https://login.wx.qq.com/jslogin?appid=wx782c26e4c19acffb&redirect_uri=https%3A%2F%2Fwx.qq.com%2Fcgi-bin%2Fmmwebwx-bin%2Fwebwxnewloginpage&fun=new&lang=zh_CN&_=1676127661127

var Client = NewClient()

func GetUUID() (uuid string, err error) {
	params := url.Values{}
	params.Add("appid", APP_ID)
	params.Add("redirect_uri", REDIRECT_URI)
	params.Add("fun", "new")
	params.Add("lang", "zh_CN")
	params.Add("_", strconv.FormatInt(time.Now().Unix(), 10))

	resp, err := Client.PostForm(JS_LOGIN_URL, params)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// parse current the response
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	ds := string(data)

	// window.QRLogin.code = 200; window.QRLogin.uuid = "geRvdUP_NA==";
	code := ParseQrcodeResult(ds, "window.QRLogin.code")

	codeInt, _ := strconv.Atoi(code)
	if codeInt != 200 {
		return "", fmt.Errorf("get uuid error, code: %s", code)
	}

	uuid = ParseQrcodeResult(ds, "window.QRLogin.uuid")
	uuid = uuid[1 : len(uuid)-1]

	return uuid, nil
}
