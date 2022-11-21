package convert

import (
	"strconv"
)

// String2Int int
func String2Int(s string) (int, error) {
	return strconv.Atoi(s)
}

// String2Int32 int32
func String2Int32(s string) (int32, error) {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	return int32(i), nil
}

// String2Int64 int64
func String2Int64(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}
