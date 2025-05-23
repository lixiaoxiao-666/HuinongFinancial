package models

// Config yaml 文件结构体
var Config *config

// config 结构体
type config struct {
	Database Database
	Server Server
}

// Database 结构体
type Database struct {
	Host string
	Port string
	User string
	Password string
	Database string
}

// Server 结构体
type Server struct {
	Port string
	ReadTimeout string
	WriteTimeout string
}
