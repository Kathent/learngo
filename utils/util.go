package utils

func GetByte4(val int) []byte {
	res := make([]byte, 0)
	res = append(res, byte(val & 0xFF))
	res = append(res, byte(val >> 8 ))
	res = append(res, byte(val >> 16))
	res = append(res, byte(val >> 24))
	return res
}

func GetStringNul(val string) []byte {
	res := make([]byte, 0)
	res = append(res, []byte(val)...)
	res = append(res, 0)
	return res
}

func GetByte2(val int) []byte {
	res := make([]byte, 0)
	val = val & 0xFFFF
	res = append(res, byte(val & 0xFF))
	res = append(res, byte(val >> 8 ))
	return res
}
