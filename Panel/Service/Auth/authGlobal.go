/*
 * 创建人： deadmau5v
 * 创建时间： 2024-0-0
 * 文件作用：权限坚定 全局
 */

package Auth

import (
	"LoongPanel/Panel/Service/Database"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"github.com/casbin/casbin/v2"
	"strings"
)

const (
	minUserNameLength = 3
	maxUserNameLength = 20

	minPasswordLength = 8
	maxPasswordLength = 64
)

var (
	ErrInvalidPassword = errors.New("invalid password")
	ErrInvalidUsername = errors.New("invalid username")
	Authenticator      *casbin.Enforcer
)

// ValidateCredentials 验证用户名和密码
func ValidateCredentials(username, password string) error {
	if err := validateUsername(username); err != nil {
		return err
	}
	return validatePassword(password)
}

// validateUsername 验证用户名
func validateUsername(username string) error {
	if len(username) < minUserNameLength || len(username) > maxUserNameLength {
		return ErrInvalidUsername
	}
	return nil
}

// validatePassword 验证密码
func validatePassword(password string) error {
	if len(password) < minPasswordLength || len(password) > maxPasswordLength {
		return ErrInvalidPassword
	}
	return nil
}

// HashPassword 对密码进行哈希处理
func HashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}

// SanitizeUsername 清理用户名
func SanitizeUsername(username string) string {
	return strings.TrimSpace(strings.ToLower(username))
}

// CreateUser 创建新用户
func CreateUser(user Database.User) error {
	if err := ValidateCredentials(user.Name, user.Password); err != nil {
		return err
	}

	hashedPassword := HashPassword(user.Password)
	sanitizedUsername := SanitizeUsername(user.Name)
	user.Name = sanitizedUsername
	user.Password = hashedPassword

	user.Save()

	return nil
}

// AuthenticateUser 验证用户身份
func AuthenticateUser(username, password string) (bool, error) {
	sanitizedUsername := SanitizeUsername(username)
	hashedPassword := HashPassword(password)
	var storedPassword string
	err := Database.DB.Model(&Database.User{}).
		Where("name = ?", sanitizedUsername).
		Select("password").
		Row().Scan(&storedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	return storedPassword == hashedPassword, nil
}

// AuthenticateUserByEmail 通过邮箱验证用户身份
func AuthenticateUserByEmail(email, password string) (bool, error) {
	hashedPassword := HashPassword(password)

	var storedPassword string
	err := Database.DB.Model(&Database.User{}).Where("mail = ?", email).Select("password").Row().Scan(&storedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	return storedPassword == hashedPassword, nil
}
