package wechat

import (
	"bytes"
	"deleted-wechat-friends/utils"
	"encoding/json"
	"encoding/xml"
	"fmt"
)

type Request struct {
	BaseRequest *BaseRequest
}

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

type Response struct {
	BaseRequest *BaseRequest
}

type User struct {
	UserName string
}

type InitReqs struct {
	Response
	User *User
}

func (r *Response) IsSuccess() bool {
	return r.BaseRequest.Ret == 0
}

// TODO: reqs type
func WebWxInit(req interface{}, baseRedirectUri string) error {
	br := &Request{
		BaseRequest: req.(*BaseRequest),
	}

	data, err := json.Marshal(br) // 序列化
	if err != nil {
		return err
	}

	reader, err := utils.CallRequestName(baseRedirectUri, "webwxinit", req, bytes.NewReader(data))
	if err != nil {
		return err
	}

	r := &InitReqs{}
	if err = json.NewDecoder(reader).Decode(r); err != nil {
		return err
	}

	if !r.IsSuccess() {
		return fmt.Errorf("message:[%s]", r.Response.BaseRequest.Message)
	}

	// r.User.UserName

	return nil
}
