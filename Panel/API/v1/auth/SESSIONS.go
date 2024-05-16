/*
 * 创建人： deadmau5v
 * 创建时间： 2024-0-0
 * 文件作用：
 */

package auth

import "github.com/google/uuid"

var SESSIONS = map[string]bool{}

func RandomSESSION() string {
	uuid_ := uuid.New().String()
	SESSIONS[uuid_] = true
	return uuid_
}
