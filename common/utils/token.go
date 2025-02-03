package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"time"
)

var key = genKey()

const (
	sha256Len = 16 // sha256部分保留的字节数
	aesLen    = 16 // 12 --> 30
)

// 将userId和sha256 揉到一起
// 类似于sha256(userId+time)(16字节)+aes(userId+time)(16字节)，最终64个字符
func genAccessToken(uid int64) (string, error) {
	byteInfo := make([]byte, 12)
	binary.BigEndian.PutUint64(byteInfo, uint64(uid))
	binary.BigEndian.PutUint32(byteInfo[8:], uint32(time.Now().UnixNano()))
	encryptedByte, err := AesEncrypt([]byte(key), byteInfo)
	if err != nil {
		return "", err
	}
	sha256Byte := sha256.Sum256(byteInfo)
	token := append(sha256Byte[:sha256Len], encryptedByte...)
	return hex.EncodeToString(token), nil
}

func genRefreshToken(uid int64) (string, error) {
	return genAccessToken(uid)
}

func GenUserAuthToken(uid int64) (accessToken, refreshToken string, err error) {
	accessToken, err = genAccessToken(uid)
	if err != nil {
		return
	}
	refreshToken, err = genRefreshToken(uid)
	return
}

func GenSessionId(uid int64) string {
	return fmt.Sprintf("%d-%d-%s", uid, time.Now().UnixNano(), RandNumStr(6))
}

func ParseUserIdFromToken(accessToken string) (int64, error) {
	if len(accessToken) != 2*(sha256Len+aesLen) {
		return 0, errors.New("invalid token")
	}
	encodeStr := accessToken[sha256Len*2:]
	encryptedData, err := hex.DecodeString(encodeStr)
	if err != nil {
		return 0, err
	}
	originData, err := AesDecrypt([]byte(key), encryptedData)
	if err != nil {
		return 0, err
	}
	uid := binary.BigEndian.Uint64(originData[:8])
	if uid == 0 {
		return 0, errors.New("invalid token")
	}
	return int64(uid), nil
}

func genKey() string {
	key := make([]byte, 16)
	_, err := rand.Read(key)
	if err != nil {
		panic(err)
	}
	return string(key)
}
