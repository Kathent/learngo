package utils

import "testing"

func TestGetByte4(t *testing.T) {
	byte4 := GetByte4(100)
	t.Log(byte4)
	if byte4[0] != 100 {
		t.FailNow()
	}
}
