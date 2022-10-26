package logic

import (
	"NechatService/models"
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/dlclark/regexp2"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/argon2"
	"log"
	"net/http"
	"strings"
)

func Register(userInfo *models.Userinfo) (int, gin.H) {
	// 账户名为4-16位, 可以包含字符串a-z A-Z 0-9 - _
	usernameOK, _ := regexp2.MustCompile(`^[a-zA-Z0-9-_]{4,16}$`, 0).MatchString(userInfo.Username)
	// 密码是否符合规范,(数字 字母 小写或小写字母 且长度大于6)
	passwordOK, _ := regexp2.MustCompile(`^.*(?=.{6,})(?=.*\d)(?=.*[a-zA-Z])(?=.*[!@#$%^&*?.]).*$`, 0).MatchString(userInfo.Password)
	if usernameOK && passwordOK {
		userInfo.CryptoPassword = argonCrypto(userInfo.Password)
		err := models.SaveUserInfo(userInfo)
		if err != nil {
			return http.StatusInternalServerError, gin.H{"code": 2, "msg": "Save userinfo failed."}
		}
		err = models.GetUserID(userInfo)
		if err != nil {
			return http.StatusInternalServerError, gin.H{"code": 2, "msg": "Get UserID failed."}
		}
		token, err := generateToken(userInfo.UserID, userInfo.Username)
		if err != nil {
			return http.StatusInternalServerError, gin.H{"code": 2, "msg": "Generate TOKEN failed."}
		}
		return http.StatusCreated, gin.H{"code": 0, "msg": "Register success.", "uid": userInfo.UserID, "token": token}
	} else {
		return http.StatusNotAcceptable, gin.H{"code": 1101, "msg": "Username or password is invalid."}
	}
}

func argonCrypto(plaintext string) (encodeHash string) {
	var p = &models.CryptoParam{
		Memory:      64 * 1024,
		Iterations:  3,
		Parallelism: 2,
		SaltLength:  16,
		KeyLength:   32,
	}
	encodeHash, err := generateFromPassword(plaintext, p)
	if err != nil {
		log.Fatal(err)
	}
	return encodeHash
}

func generateFromPassword(password string, p *models.CryptoParam) (encodedHash string, err error) {
	salt, err := generateRandomBytes(p.SaltLength)
	if err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, p.Iterations, p.Memory, p.Parallelism, p.KeyLength)

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	encodedHash = fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, p.Memory, p.Iterations, p.Parallelism, b64Salt, b64Hash)

	return encodedHash, nil
}

func generateRandomBytes(n uint32) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func comparePasswordAndHash(password, encodedHash string) (match bool, err error) {
	p, salt, hash, err := decodeHash(encodedHash)
	if err != nil {
		return false, err
	}

	otherHash := argon2.IDKey([]byte(password), salt, p.Iterations, p.Memory, p.Parallelism, p.KeyLength)

	if subtle.ConstantTimeCompare(hash, otherHash) == 1 {
		return true, nil
	}
	return false, nil
}

func decodeHash(encodedHash string) (p *models.CryptoParam, salt, hash []byte, err error) {
	values := strings.Split(encodedHash, "$")
	if len(values) != 6 {
		return nil, nil, nil, errors.New("the encoded hash is not in the correct format")
	}

	var version int
	_, err = fmt.Sscanf(values[2], "v=%d", &version)
	if err != nil {
		return nil, nil, nil, err
	}
	if version != argon2.Version {
		return nil, nil, nil, errors.New("incompatible version of argon2")
	}

	p = &models.CryptoParam{}
	_, err = fmt.Sscanf(values[3], "m=%d,t=%d,p=%d", &p.Memory, &p.Iterations, &p.Parallelism)
	if err != nil {
		return nil, nil, nil, err
	}

	salt, err = base64.RawStdEncoding.Strict().DecodeString(values[4])
	if err != nil {
		return nil, nil, nil, err
	}
	p.SaltLength = uint32(len(salt))

	hash, err = base64.RawStdEncoding.Strict().DecodeString(values[5])
	if err != nil {
		return nil, nil, nil, err
	}
	p.KeyLength = uint32(len(hash))

	return p, salt, hash, nil
}

func Login(userInfo *models.Userinfo) (int, gin.H) {
	trueUserInfo := models.Userinfo{}
	trueUserInfo.UserID = userInfo.UserID
	err := models.GetUserInfoByID(&trueUserInfo)
	if err != nil {
		return http.StatusInternalServerError, gin.H{"code": 2, "msg": "Get userInfo failed."}
	}
	ok, err := comparePasswordAndHash(userInfo.Password, trueUserInfo.CryptoPassword)
	if err != nil {
		return http.StatusInternalServerError, gin.H{"code": 2, "msg": err}
	}
	if ok {
		have, err := models.HaveApplyFriend(trueUserInfo.UserID)
		if err != nil {
			return http.StatusInternalServerError, gin.H{"code": 2, "msg": err}
		}

		token, err := generateToken(trueUserInfo.UserID, trueUserInfo.Username)
		if err != nil {
			return http.StatusInternalServerError, gin.H{"code": 2, "msg": "Generate TOKEN failed."}
		}
		return http.StatusOK, gin.H{"code": 0, "msg": "Login Success.", "token": token, "Username": trueUserInfo.Username, "HaveMessage": have}
	} else {
		return http.StatusInternalServerError, gin.H{"code": 2, "msg": "UserID or Password is wrong."}
	}
}
