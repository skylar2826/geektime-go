package types

const (
	UserName = "root"
	Password = "15271908767Aa!"
	Ip       = "127.0.0.1"
	Port     = "3306"
	DbName   = "test1"
)

type User struct {
	Id    int
	Name  string
	Age   int
	Sex   int
	Phone string
}
