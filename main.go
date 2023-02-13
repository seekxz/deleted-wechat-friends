package main

import (
	"deleted-wechat-friends/utils"
	"deleted-wechat-friends/wechat"
	"log"
)

func main() {
	log.Println("server start")

	// 获取 uuid
	uuid, err := utils.GetUUID()
	if err != nil {
		log.Println("uuid 获取失败", err)
		return
	}

	// 获取二维码
	err = utils.NewQrCode(uuid)
	if err != nil {
		log.Println("创建二维码失败：", err)
		return
	}
	log.Println("请扫描二维码登录")

	// 扫描二维码，获取登录地址
	redirectUrl, err := utils.GetLoginRedirectUrl(uuid)
	if err != nil {
		log.Println("获取 redirectUrl 失败：", err)
		return
	}

	// 登录
	reqs, err := Login(redirectUrl)
	if err != nil {
		log.Println("登录失败：", err)
		return
	}

	// 获取 BaseRedirectUri
	baseRedirectUri := utils.GetBaseRedirectUri(redirectUrl)

	// init webwxinit
	if err = wechat.WebWxInit(reqs, baseRedirectUri); err != nil {
		log.Println("webwxinit 失败：", err.Error())
		return
	}

	// 获取联系人
	list, count, err := GetContact(reqs, baseRedirectUri)
	if err != nil {
		log.Println("获取联系人失败：", err.Error())
		return
	}

	log.Println("联系人数量：", count)
	log.Println("联系人列表：", list)

	log.Println("server end")
}
