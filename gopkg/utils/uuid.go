package utils

import (
	"encoding/hex"

	"github.com/bwmarrin/snowflake"
	"github.com/google/uuid"
)

// node 实例
var node, _ = snowflake.NewNode(1)

// GenUUID 生成UUID
func GenUUID() string {
	return uuid.New().String()
}

// GenUUIDWithoutUnderline 不带下划线的UUID
func GenUUIDWithoutUnderline() string {
	var u = uuid.New()
	var dst [32]byte
	hex.Encode(dst[:], u[:])
	return string(dst[:])
}

// SnowflakeGenUUID 雪花算法UUID
func SnowflakeGenUUID() string {
	return node.Generate().String()
}

func SnowflakeGenIntUUID() int64 {
	return node.Generate().Int64()
}
