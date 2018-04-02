package utils

import "math"

func GetByte4(val int) []byte {
	res := make([]byte, 0)
	res = append(res, byte(val & 0xFF))
	res = append(res, byte(val >> 8 ))
	res = append(res, byte(val >> 16))
	res = append(res, byte(val >> 24))
	return res
}

func GetByte8(val int64) []byte {
	res := make([]byte, 0)
	res = append(res, GetByte4(int(val))...)
	res = append(res, GetByte4(int(val >> 32))...)
	return res
}

func GetByte3(val int) []byte {
	res := make([]byte, 0)
	val = val & 0xFFFFFF
	res = append(res, byte(val))
	res = append(res, byte(val >> 8 ))
	res = append(res, byte(val >> 16))
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

func GetLongToInt(val int64) []byte {
	if val < 251 {
		return []byte{byte(val)}
	}

	bts := make([]byte, 0)
	if val >= 251 && float64(val) < math.Pow(2, 16) {
		bts = append(bts, 0xfc)
		bts = append(bts, GetByte2(int(val))...)
		return bts
	}

	if float64(val) >= math.Pow(2, 16) && float64(val) < math.Pow(2, 24) {
		bts = append(bts, 0xfd)
		bts = append(bts, GetByte4(int(val))...)
		return bts
	}

	bts = append(bts, 0xfe)
	bts = append(bts, GetByte8(val)...)
	return bts
}

func GetStringLenBts(val string) []byte {
	res := make([]byte, 0)
	if len(val) <= 0 {
		res = append(res, 0)
		return res
	}

	res = append(res, GetLongToInt(int64(len(val)))...)
	res = append(res, []byte(val)...)
	return res
}
