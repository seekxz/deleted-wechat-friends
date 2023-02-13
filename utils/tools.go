package utils

import (
	"regexp"
)

func ParseQrcodeResult(data, qrCodeKey string) string {
	reg := regexp.MustCompile(`(window[.QRLogin]*.\w+) ?= ?(.+?);`)

	m := make(map[string]string)
	for _, v := range reg.FindAllStringSubmatch(data, -1) {
		m[v[1]] = v[2]
	}

	return m[qrCodeKey]
}
