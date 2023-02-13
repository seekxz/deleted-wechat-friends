package main

import (
	"deleted-wechat-friends/utils"
	"encoding/json"
	"fmt"
	"strings"
)

type ContaceRequest struct {
	Response
	MemberCount int
	MemberList  []*Member
}

func (r *Response) IsSuccess() bool {
	return r.BaseResponse.Ret == 0
}

type Member struct {
	UserName   string
	NickName   string
	RemarkName string
	VerifyFlag int
}

func (m *Member) IsNormal() bool {
	return m.VerifyFlag&8 == 0 && // 公众号/服务号
		!strings.HasPrefix(m.UserName, "@@") && // 群聊
		// m.UserName != "" && // 自己
		!m.IsSpecail() // 特殊账号
}

func (m *Member) IsSpecail() bool {
	for i, count := 0, len(SpecialUsers); i < count; i++ {
		if m.UserName == SpecialUsers[i] {
			return true
		}
	}
	return false
}

func GetContact(req interface{}, baseRedirectUri string) (list []*Member, count int, err error) {
	reader, err := utils.CallRequestName(baseRedirectUri, "webwxgetcontact", req, nil)
	if err != nil {
		return
	}

	r := &ContaceRequest{}
	if err = json.NewDecoder(reader).Decode(r); err != nil {
		return
	}

	if !r.IsSuccess() {
		err = fmt.Errorf("message:[%s]", r.BaseResponse.ErrMsg)
		return
	}

	list, count = make([]*Member, 0, r.MemberCount/5*2), r.MemberCount
	for i := 0; i < count; i++ {
		if r.MemberList[i].IsNormal() {
			list = append(list, r.MemberList[i])
		}
	}

	return
}
