package utilx

import (
	"reflect"
)

// CollectString walk s then return all union strings selected by fieldFn
func CollectString(s interface{}, fieldFn func(interface{}) string) []string {
	sv := reflect.ValueOf(s)
	if !(sv.Kind() == reflect.Slice || sv.Kind() == reflect.Array) {
		panic("CollectString: s must be slice or array")
	}

	slen := sv.Len()
	strs := make([]string, 0, slen)
	smap := make(map[string]bool, slen)
	for i := 0; i < slen; i++ {
		siv := sv.Index(i)
		// fieldFn should receive a pointer to an item in s to avoid memory copy
		if siv.Kind() != reflect.Pointer && siv.CanAddr() {
			siv = siv.Addr()
			siv.Elem()
		}

		str := fieldFn(siv.Interface())
		if str == "" {
			continue
		}
		if smap[str] {
			continue
		}
		smap[str] = true
		strs = append(strs, str)
	}

	return strs
}

// Returns empty string if s is nil
func StringElem(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
