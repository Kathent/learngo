package timer_client

import (
	"regexp"
	"fmt"
)

const(
	CALL_BACK_PROTO_PATTERN     = "(?:(direct|http|rpc|https?)://)?"
	CALL_BACK_HOST_NAME_PATTERN = "(?:[a-z0-9](?:[-a-z0-9]*[a-z0-9])?\\.)+(?:com|net|edu|biz|gov|org|in(?:t|fo)|(?-i:[a-z][a-z]))"
	CALL_BACK_IP_PATTERN        = "(?:[01]?\\d\\d?|2[0-4]\\d|25[0-5])\\.(?:[01]?\\d\\d?|2[0-4]\\d|25[0-5])\\.(?:[01]?\\d\\d?|2[0-4]\\d|25[0-5])\\.(?:[01]?\\d\\d?|2[0-4]\\d|25[0-5])"
	CALL_BACK_PORT_PATTERN      = "(?::(\\d{1,5}))?"
	CALL_BACK_PATH_PATTERN      = "(/.*)?"
	CALL_BACK_PATTERN           = "%s(%s|%s%s)%s"
)

var pattern *regexp.Regexp

//parseCallBack 解析callBack返回call类型和接口
func parseCallBack(callBack string) (string ,string){
	arr := pattern.FindStringSubmatch(callBack)
	if len(arr) > 1 {
		return arr[0], arr[1]
	}else if len(arr) > 0 {
		return arr[0], ""
	}
	return "", ""
}

func init() {
	ps := fmt.Sprintf(CALL_BACK_PATTERN, CALL_BACK_PROTO_PATTERN, CALL_BACK_HOST_NAME_PATTERN,
		CALL_BACK_IP_PATTERN, CALL_BACK_PORT_PATTERN, CALL_BACK_PATH_PATTERN)
	p, err := regexp.Compile(ps)
	if err != nil {
		panic(err)
	}
	pattern = p
}

