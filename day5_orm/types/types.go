package types

import "fmt"

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

type TestPerson struct {
	Id   int
	Name string
}

type TestUser struct {
	Name string
	age  int
}

func NewTestUser(name string, age int) TestUser {
	return TestUser{
		name,
		age,
	}
}

func NewTestUserPtr(name string, age int) *TestUser {
	return &TestUser{
		name,
		age,
	}
}

func (u TestUser) GetAge() int {
	return u.age
}

func (u *TestUser) ChangeName(name string) {
	u.Name = name
}

func (u TestUser) private() {
	fmt.Println("private func")
}
