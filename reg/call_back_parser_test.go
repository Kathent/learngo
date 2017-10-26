package timer_client

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"regexp"
)

func TestParseCallBack(t *testing.T) {
	str := "direct://http://192.168.0.0.1:1234/url"
	compile, err := regexp.Compile(CALL_BACK_PROTO_PATTERN)
	assert.Nil(t, err)

	submatch := compile.FindStringSubmatch(str)
	t.Logf("submatch is:%+v, lenght:%d", submatch, len(submatch))

	stringSubmatch := pattern.FindStringSubmatch(str)
	t.Logf("parse result is :%+v, length:%d", stringSubmatch, len(stringSubmatch))
}
