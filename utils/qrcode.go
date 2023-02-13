package utils

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// https://login.weixin.qq.com/qrcode/AeGdQA6D4A==

const QR_CODE_URL = "https://login.weixin.qq.com/qrcode/"
const QR_CODE_REDIRECT_URI = "https://login.wx.qq.com/cgi-bin/mmwebwx-bin/login"

var qrImagePath = ""
var currentPath = ""

func init() {
	// create the QR image path
	var err error
	if currentPath, err = os.Getwd(); err != nil {
		log.Println("获取当前路径失败: ", err)
		return
	}

	qrImagePath = filepath.Join(currentPath, "qrcode.png")
}

func NewQrCode(uuid string) (err error) {
	qrcode, err := getQRCode(uuid)
	if err != nil {
		log.Println("获取二维码失败", err)
		return err
	}

	// save the qrcode to the file
	err = ShowQRCode(qrcode)
	if err != nil {
		log.Println("显示二维码失败", err)
		return err
	}

	// use system default application to open the image
	cmd, err := OpenQrCode(qrImagePath)
	if err != nil {
		log.Println("打开二维码失败", err)
		return err
	}

	// wait for the process to exit
	log.Println(cmd, "等待进程退出")
	go func() {
		// cmd.Wait()

		// if cmd.Process != nil {
		// 	cmd.Process.Kill()
		// }

		// os.Remove(qrImagePath)
	}()

	return
}

func getQRCode(uuid string) (qrcode []byte, err error) {
	qrurl := QR_CODE_URL + uuid
	params := url.Values{}
	params.Set("t", "webwx")
	params.Set("_", strconv.FormatInt(time.Now().Unix(), 10))

	resp, err := Client.PostForm(qrurl, params)
	if err != nil {
		log.Println("获取二维码失败", err)
		return nil, err
	}

	// parse current the response
	qrcode, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("qrcode 返回错误:", err)
		return nil, err
	}

	return qrcode, nil
}

func ShowQRCode(qrcode []byte) error {
	file, err := os.OpenFile(qrImagePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		log.Println("打开文件失败", err)
		return nil
	}
	defer file.Close()

	_, err = file.Write(qrcode)
	return err
}

func OpenQrCode(path string) (cmd *exec.Cmd, err error) {
	command := "open"
	switch os := runtime.GOOS; os {
	case "linux":
		command = "xdg-open"
	case "windows":
		command = "start"
	case "darwin":
	default:
		err = fmt.Errorf("unsupported operating system: %s", os)
		return
	}

	cmd = exec.Command(command, path)
	err = cmd.Start()

	return
}

func GetLoginRedirectUrl(uuid string) (redirectUrl string, err error) {
	redirectUrl, code, tip := "", "", 0
	for code != "200" {
		redirectUrl, code, tip, err = loginRedirectUrl(uuid, tip)
		if err != nil {
			log.Println("重新获取登录地址失败", err)
			return
		}

		time.Sleep(1 * time.Second)
	}

	return
}

func loginRedirectUrl(uuid string, tip int) (redirectUrl string, code string, defaultTip int, err error) {
	// https://login.wx.qq.com/cgi-bin/mmwebwx-bin/login?loginicon=true&uuid=QYJLbPjNKA==&tip=0&r=-1140064285&_=1676177258876
	// https://login.wx.qq.com/cgi-bin/mmwebwx-bin/login?loginicon=true&uuid=4YSbQSS1kA==&tip=0&r=-1676179648&_=1676179648":
	// net/http: timeout awaiting response headers (Client.Timeout exceeded while awaiting headers)

	loginUrl, defaultTip := fmt.Sprintf("%s?uuid=%s&tip=%d&_=%d", QR_CODE_REDIRECT_URI, uuid, tip, time.Now().Unix()), defaultTip

	// resp, err := Client.Get(loginUrl)
	resp, err := http.Get(loginUrl)
	if err != nil {
		log.Println("获取登录地址失败", err)
		return
	}
	defer resp.Body.Close()

	// parse current the response
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("登录地址返回错误:", err)
		return
	}
	ds := string(data)
	code = ParseQrcodeResult(ds, "window.code")

	// window.code=200;
	// window.redirect_uri="https://wx.qq.com/cgi-bin/mmwebwx-bin/webwxnewloginpage?ticket=AXrlMtIjUCCGfz8kOlHleUju@qrticket_0&uuid=YbxQ1twuGQ==&lang=zh_CN&scan=1676180561";

	// handle different code logic
	switch code {
	case "201":
		log.Println("扫码成功，请在手机上点击确认以登录")
		defaultTip = 0
	case "200":
		resultUrl := ParseQrcodeResult(ds, "window.redirect_uri")
		redirectUrl = resultUrl[1:len(resultUrl)-1] + "&fun=new"
	case "408":
		err = fmt.Errorf("login timeout")
	default:
		err = fmt.Errorf("unknown error")
	}

	return
}

func GetBaseRedirectUri(redirectUri string) string {
	index := strings.LastIndex(redirectUri, "/")
	if index == -1 {
		index = len(redirectUri)
	}
	baseUri := redirectUri[:index]

	return baseUri
}
