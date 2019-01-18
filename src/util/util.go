package util

import (
	"bytes"
	"fmt"
	"strings"
)

// Trim 去除空格
func Trim(s string) string {
	return strings.TrimSpace(s)
}

// GetLevelCode 获取分级码
func GetLevelCode(orderLevelCodes []string) string {
	l := len(orderLevelCodes)

	if l == 0 {
		return "01"
	} else if l == 1 {
		return orderLevelCodes[0] + "01"
	}

	root := orderLevelCodes[0]
	toValue := func(i int) string {
		if i < 10 {
			return fmt.Sprintf("%s0%d", root, i)
		}
		return fmt.Sprintf("%s%d", root, i)
	}

	for i := 1; i < 100; i++ {
		code := toValue(i)
		if i < l &&
			orderLevelCodes[i] == code {
			continue
		}
		return code
	}

	return ""
}

// ParseLevelCodes 解析分级码（去重）
func ParseLevelCodes(levelCodes ...string) []string {
	var allCodes []string

	for _, levelCode := range levelCodes {
		codes := parseLevelCode(levelCode)

		for _, code := range codes {
			var exists bool
			for _, c := range allCodes {
				if code == c {
					exists = true
					break
				}
			}

			if !exists {
				allCodes = append(allCodes, code)
			}
		}
	}

	return allCodes
}

func parseLevelCode(levelCode string) []string {
	if len(levelCode) < 2 {
		return nil
	}
	var (
		codes []string
		root  bytes.Buffer
	)

	for i := range levelCode {
		idx := i + 1
		if idx%2 == 0 {
			root.WriteString(levelCode[idx-2 : idx])
			codes = append(codes, root.String())
		}
	}

	root.Reset()
	return codes
}
