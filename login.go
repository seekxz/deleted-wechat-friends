package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"net/http"
)

const REQUEST_SUCCESSED = 0

var deviceID = flag.String("deviceID", "e000000000000000", "设备ID")

func Login(redirectUrl string) (reqs *BaseRequest, err error) {
	// https://wx.qq.com/cgi-bin/mmwebwx-bin/webwxnewloginpage?ticket=AbxCQubR0pE0oMLOoN_I-vEv@qrticket_0&uuid=IYJmoHlUIA==&lang=zh_CN&scan=1676249998&fun=new

	resp, err := http.Get(redirectUrl)
	if err != nil {
		return reqs, err
	}
	defer resp.Body.Close()

	// parse current the response body to BaseRequest
	// body, err := ioutil.ReadAll(resp.Body)
	// fmt.Println(string(body), "resp")
	reader := resp.Body.(io.Reader)

	reqs = &BaseRequest{}
	// 由于微信网页版无法登录，故这边会报错
	if err = xml.NewDecoder(reader).Decode(reqs); err != nil {
		return
	}

	if reqs.Ret != REQUEST_SUCCESSED {
		err = fmt.Errorf("message:[%s]", reqs.Message)
		return
	}

	reqs.DeviceID = *deviceID
	return
}
