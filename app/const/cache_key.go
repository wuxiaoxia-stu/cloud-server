package _const

import "fmt"

//token缓存
func TOKEN_CACHE_KEY(token string) string {
	return fmt.Sprintf("token:%s", token)
}

//登录验证码缓存
func LOGIN_CODE_CACHE_KEY(username, code string) string {
	return fmt.Sprintf("login:code:%s:%s", username, code)
}
