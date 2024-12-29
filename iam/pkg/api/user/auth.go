package user

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"io"
	"time"
)

var sKey = []byte("bfmFCucA1fP2EWh7idTd")

// GenerateHashPwd 对密码进行哈希加密
func GenerateHashPwd(pwd string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("Error hashing password:", err)
		return "", err
	}
	return string(hashedPassword), nil
}

// Compare 验证 hashedPassword 解密是否为 password, 注: hashedPassword 为 GenerateHashPwd 获取
func Compare(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// GenerateSymmetryPwd 将pwd进行
func GenerateSymmetryPwd(pwd string) (string, error) {
	// 使用 AES 创建一个新的加密块
	block, err := aes.NewCipher(sKey)
	if err != nil {
		fmt.Println("Error creating cipher block:", err)
		return "", err
	}

	// 创建一个使用 AES CBC 模式的加密器
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		fmt.Println("Error generating random iv:", err)
		return "", err
	}
	stream := cipher.NewCTR(block, iv)

	// 对数据进行加密
	ciphertext := make([]byte, aes.BlockSize+len(pwd))
	stream.XORKeyStream(ciphertext[aes.BlockSize:], []byte(pwd))

	return string(ciphertext), nil

	// 输出加密后的数据
	//fmt.Printf("Ciphertext: %x\n", ciphertext)
	//
	//// 创建一个使用 AES CBC 模式的解密器
	//stream = cipher.NewCTR(block, iv)
	//
	//// 对加密后的数据进行解密
	//plaintextCopy := make([]byte, len(pwd))
	//stream.XORKeyStream(plaintextCopy, ciphertext[aes.BlockSize:])
	//
	//// 输出解密后的数据
	//fmt.Println("Plaintext:", string(plaintextCopy))
}

// Sign 通过secretKey进行生成token，使用secretID和secretKey进行对应
func Sign(secretID string, secretKey string, iss, aud string) string {
	claims := jwt.MapClaims{
		"exp": time.Now().Add(time.Minute).Unix(),
		"iat": time.Now().Unix(),
		"nbf": time.Now().Add(0).Unix(),
		"aud": aud,
		"iss": iss,
	}

	// create a new token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token.Header["kid"] = secretID

	// Sign the token with the specified secret.
	tokenString, _ := token.SignedString([]byte(secretKey))

	return tokenString
}
