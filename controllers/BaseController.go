package controllers

//定义响应json数据格式
type ResponseInfo struct {
	Code    string
	Message string
	Data    interface{}
}
