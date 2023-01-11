package main

import (
	"net"
	"strings"
	"time"
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

//func (u *User) ListenMsg() {
//	for {
//		fmt.Println(u.name + " is waiting")
//		msg := <-u.c
//		u.conn.Write([]byte(msg + "\n"))
//		fmt.Println("server send to" + u.name)
//	}
//}

func GetKeys(m map[string]*User, style int) []string {
	keys := make([]string, len(m))
	if style == 0 {
		i := 0
		for k, _ := range m {
			keys[i] = k
			i++
		}
	} else if style == 2 {
		for _, k := range m {
			keys[0] += k.name + ";  "
		}
	}

	return keys
}

func (u *User) SendMsg(msg string) {
	u.conn.Write([]byte(msg))
}

func (u *User) DoMsg(msg string, s *Server) {
	le := len(msg)
	if msg[:le/2] == msg[le/2:] {
		msg = msg[:le/2]
	}
	if msg == "who" {
		u.conn.Write([]byte(GetKeys(s.OnlineMap, 2)[0]))
	} else if msg[:3] == "to|" {
		remoteName := strings.Split(msg, "|")[1]
		for _, v := range s.OnlineMap {
			if v.name == remoteName {
				v.SendMsg(u.name + "said to you : " + strings.Split(msg, "|")[2])
				return
			}
		}
		u.SendMsg("server replay: no such person")
	} else if msg[:7] == "modify|" {
		u.name = strings.Split(msg, "|")[2]
		u.SendMsg("you name change to: " + u.name)
	} else if msg[:4] == "off|" {
		s.offline(u.conn)
	} else if msg == "who am i" {
		u.SendMsg(u.name)
	} else {
		msg = u.name + "  " + time.Now().Format("2006-01-02 15:04:05") + "  " + msg
		s.Message <- msg
	}
}
