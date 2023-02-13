package main

import "encoding/xml"

type Request struct {
	BaseRequest *BaseRequest
}

type BaseRequest struct {
	// XMLName struct{} `xml:"xml" json:"-"`
	XMLName xml.Name `xml:"error" json:"-"`

	Ret        int    `xml:"Ret" json:"ret"`
	Message    string `xml:"Message" json:"message"`
	Skey       string `xml:"Skey" json:"skey"`
	Wxsid      string `xml:"Wxsid" json:"wxsid"`
	Wxuin      string `xml:"Wxuin" json:"wxuin"`
	PassTicket string `xml:"PassTicket" json:"passTicket"`

	DeviceID string `xml:"DeviceID" json:"deviceID"`
}

type Response struct {
	BaseResponse *BaseResponse
}

type BaseResponse struct {
	Ret    int
	ErrMsg string
}

var SpecialUsers = []string{
	"newsapp", "fmessage", "filehelper", "weibo", "qqmail",
	"tmessage", "qmessage", "qqsync", "floatbottle", "lbsapp",
	"shakeapp", "medianote", "qqfriend", "readerapp", "blogapp",
	"facebookapp", "masssendapp", "meishiapp", "feedsapp", "voip",
	"blogappweixin", "weixin", "brandsessionholder", "weixinreminder", "wxid_novlwrv3lqwv11",
	"gh_22b87fa7cb3c", "officialaccounts", "notification_messages", "wxitil", "userexperience_alarm",
}
