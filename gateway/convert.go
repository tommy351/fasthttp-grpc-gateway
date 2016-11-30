package gateway

import "strconv"

// ConvertStringToFloat64 converts string to float64.
func ConvertStringToFloat64(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}

// ConvertStringToFloat32 converts string to float32.
func ConvertStringToFloat32(s string) (float32, error) {
	v, err := strconv.ParseFloat(s, 32)
	return float32(v), err
}

// ConvertStringToInt64 converts string to int64.
func ConvertStringToInt64(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}

// ConvertStringToInt32 converts string to int32.
func ConvertStringToInt32(s string) (int32, error) {
	v, err := strconv.ParseInt(s, 10, 32)
	return int32(v), err
}

// ConvertStringToInt16 converts string to int16.
func ConvertStringToInt16(s string) (int16, error) {
	v, err := strconv.ParseInt(s, 10, 16)
	return int16(v), err
}

// ConvertStringToInt8 converts string to int8.
func ConvertStringToInt8(s string) (int8, error) {
	v, err := strconv.ParseInt(s, 10, 8)
	return int8(v), err
}

// ConvertStringToInt converts string to int.
func ConvertStringToInt(s string) (int, error) {
	v, err := strconv.ParseInt(s, 10, 0)
	return int(v), err
}

// ConvertStringToUint64 converts string to uint64.
func ConvertStringToUint64(s string) (uint64, error) {
	return strconv.ParseUint(s, 10, 64)
}

// ConvertStringToUint32 converts string to uint32.
func ConvertStringToUint32(s string) (uint32, error) {
	v, err := strconv.ParseUint(s, 10, 32)
	return uint32(v), err
}

// ConvertStringToUint16 converts string to uint16.
func ConvertStringToUint16(s string) (uint16, error) {
	v, err := strconv.ParseUint(s, 10, 16)
	return uint16(v), err
}

// ConvertStringToUint8 converts string to uint8.
func ConvertStringToUint8(s string) (uint8, error) {
	v, err := strconv.ParseUint(s, 10, 8)
	return uint8(v), err
}

// ConvertStringToUint converts string to uint.
func ConvertStringToUint(s string) (uint, error) {
	v, err := strconv.ParseUint(s, 10, 0)
	return uint(v), err
}

// ConvertStringToBool converts string to bool. The string "on" will be converted to "true".
func ConvertStringToBool(s string) (bool, error) {
	if s == "on" {
		return true, nil
	}

	v, err := strconv.ParseBool(s)
	return v, err
}
