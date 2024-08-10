package util

import (
	"crypto/rand"
	"fmt"
)

func UUID() string {
	arr := make([]byte, 16)
	_, err := rand.Read(arr)
	if err != nil {
		panic(err)
	}

	// 设置 UUID 版本和变体
	arr[8] = arr[8]&^0xc0 | 0x80
	arr[6] = arr[6]&^0xf0 | 0x40

	return fmt.Sprintf("%x%x%x%x%x", arr[0:4], arr[4:6], arr[6:8], arr[8:10], arr[10:])
}
