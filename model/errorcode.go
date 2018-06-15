package model

//errorCode 错误代码
type errorCode struct {
	SUCCESS      int
	ERROR        int
	NotFound     int
	LoginError   int
	LoginTimeOut int
	InActive     int
}

//ErrorCode 错误赋值
var ErrorCode = errorCode{
	SUCCESS:      0,
	ERROR:        1,
	NotFound:     404,
	LoginError:   1000, //用户名或密码错误
	LoginTimeOut: 1001, //登录超时
	InActive:     1002, //未激活账号
}
