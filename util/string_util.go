package util

import (
	"fmt"
	"strconv"
)

const MAX_KEY_LENGTH = 12

func BuildLineWithLength(leftItem string, rightItem string, leftItemSize int) string {
	return fmt.Sprintf("%-"+strconv.Itoa(leftItemSize)+"s   %s", leftItem, rightItem)
}

func BuildLineOfFixedLength(key string, value string) string {
	return BuildLineWithLength(key, value, MAX_KEY_LENGTH)
}
