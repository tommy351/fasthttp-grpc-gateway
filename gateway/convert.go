package gateway

import "strconv"

func ConvertStringToFloat64(s string) float64 {
	v, _ := strconv.ParseFloat(s, 64)
	return v
}

func ConvertStringToFloat32(s string) float32 {
	v, _ := strconv.ParseFloat(s, 32)
	return float32(v)
}

func ConvertStringToInt64(s string) int64 {
	v, _ := strconv.ParseInt(s, 10, 64)
	return v
}

func ConvertStringToInt32(s string) int32 {
	v, _ := strconv.ParseInt(s, 10, 32)
	return int32(v)
}

func ConvertStringToInt16(s string) int16 {
	v, _ := strconv.ParseInt(s, 10, 16)
	return int16(v)
}

func ConvertStringToInt8(s string) int8 {
	v, _ := strconv.ParseInt(s, 10, 8)
	return int8(v)
}

func ConvertStringToInt(s string) int {
	v, _ := strconv.ParseInt(s, 10, 0)
	return int(v)
}

func ConvertStringToUint64(s string) uint64 {
	v, _ := strconv.ParseUint(s, 10, 64)
	return v
}

func ConvertStringToUint32(s string) uint32 {
	v, _ := strconv.ParseUint(s, 10, 32)
	return uint32(v)
}

func ConvertStringToUint16(s string) uint16 {
	v, _ := strconv.ParseUint(s, 10, 16)
	return uint16(v)
}

func ConvertStringToUint8(s string) uint8 {
	v, _ := strconv.ParseUint(s, 10, 8)
	return uint8(v)
}

func ConvertStringToUint(s string) uint {
	v, _ := strconv.ParseUint(s, 10, 0)
	return uint(v)
}

func ConvertStringToBool(s string) bool {
	if s == "on" {
		return true
	}

	v, _ := strconv.ParseBool(s)
	return v
}
