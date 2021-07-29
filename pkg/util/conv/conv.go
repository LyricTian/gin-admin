package conv

import "strconv"

func ParseStringSliceToUint64(s []string) []uint64 {
	iv := make([]uint64, len(s))
	for i, v := range s {
		iv[i], _ = strconv.ParseUint(v, 10, 64)
	}
	return iv
}
