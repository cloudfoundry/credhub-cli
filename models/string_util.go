package models

import (
	"fmt"
	"strconv"
)

const MAX_KEY_LENGTH = 12

func buildLineWithLength(leftItem string, rightItem string, leftItemSize int) string {
	return fmt.Sprintf("%-"+strconv.Itoa(leftItemSize)+"s   %s", leftItem, rightItem)
}

func buildLineOfFixedLength(key string, value string) string {
	return buildLineWithLength(key, value, MAX_KEY_LENGTH)
}
