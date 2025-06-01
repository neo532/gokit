package client

/*
 * @abstract The operation method of the client of the transmission protocol http
 * @mail neo532@126.com
 * @date 2022-05-30
 */

import (
	"strings"
)

const (
	ContentTypeHeaderKey          = "Content-Type"
	ContentTypeHeaderDefaultValue = "application/json;"
)

// multipart/form-data => ""
// application/x-www-form-urlencoded;charset=utf-8 => x-www-form-urlencoded
// Content-Type: application/json;charset=utf-8 => json
func ContentSubtype(contentType string) (subType string) {
	subType = "json"
	contentType = strings.ToLower(contentType)
	cts := strings.SplitN(contentType, "application/", 2)
	if len(cts) <= 1 {
		return
	}
	sts := strings.SplitN(cts[1], ";", 2)
	if len(sts) <= 1 {
		return
	}
	subType = sts[0]
	return
}
