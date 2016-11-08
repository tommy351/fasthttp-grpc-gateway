package gateway

import "strconv"

func ConvertStringToFloat64(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}

func ConvertStringToFloat32(s string) (float32, error) {
	v, err := strconv.ParseFloat(s, 32)
	return float32(v), err
}

func ConvertStringToInt64(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}

func ConvertStringToInt32(s string) (int32, error) {
	v, err := strconv.ParseInt(s, 10, 32)
	return int32(v), err
}

func ConvertStringToInt16(s string) (int16, error) {
	v, err := strconv.ParseInt(s, 10, 16)
	return int16(v), err
}

func ConvertStringToInt8(s string) (int8, error) {
	v, err := strconv.ParseInt(s, 10, 8)
	return int8(v), err
}

func ConvertStringToInt(s string) (int, error) {
	v, err := strconv.ParseInt(s, 10, 0)
	return int(v), err
}

func ConvertStringToUint64(s string) (uint64, error) {
	return strconv.ParseUint(s, 10, 64)
}

func ConvertStringToUint32(s string) (uint32, error) {
	v, err := strconv.ParseUint(s, 10, 32)
	return uint32(v), err
}

func ConvertStringToUint16(s string) (uint16, error) {
	v, err := strconv.ParseUint(s, 10, 16)
	return uint16(v), err
}

func ConvertStringToUint8(s string) (uint8, error) {
	v, err := strconv.ParseUint(s, 10, 8)
	return uint8(v), err
}

func ConvertStringToUint(s string) (uint, error) {
	v, err := strconv.ParseUint(s, 10, 0)
	return uint(v), err
}

func ConvertStringToBool(s string) (bool, error) {
	if s == "on" {
		return true, nil
	}

	v, err := strconv.ParseBool(s)
	return v, err
}
