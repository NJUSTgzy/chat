package main

import (
	"fmt"
	"net"
)

type User struct {
	name      string
	c         chan string
	Addr      string
	fileList  map[string]int
	FriendMap map[string]*User
	conn      net.Conn
}

func NewUser(conn net.Conn) *User {
	addr := conn.RemoteAddr().String()
	cl := &User{
		name:      addr,
		Addr:      addr,
		conn:      conn,
		FriendMap: make(map[string]*User),
		fileList:  make(map[string]int),
		c:         make(chan string),
	}

	return cl
}

func (u *User) ListenMsg() {
	for {
		fmt.Println(u.name + " is waiting")
		msg := <-u.c
		u.conn.Write([]byte(msg + "\n"))
		fmt.Println("server send to" + u.name)
	}
}

func GetKeys(m map[string]*User, style int) []string {
	keys := make([]string, len(m))
	if style == 0 {
		i := 0
		for k, _ := range m {
			keys[i] = k
			i++
		}
	} else if style == 2 {
		for k, _ := range m {
			keys[0] += k + ";"
		}
	}

	return keys
}

func (u *User) SendMsg(msg string) {
	//buf := byte[]()
}

func (u *User) DoMsg(msg string) {

}
