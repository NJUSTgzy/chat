package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"net"
	"sync"
)

type WinConfig struct {
	App fyne.App
	Win fyne.Window
}

type Server struct {
	ConnType  string
	Ip        string
	Port      string
	isStart   bool
	OnlineMap map[string]*User
	Message   chan string
	file      string
	btnSend   *widget.Button
	sendTo    []string
	sendP     int
	nowPeople int
	UserLock  sync.RWMutex
	win       *WinConfig
}

func Init(con, Ip, Port string) *Server {
	server := &Server{
		ConnType:  con,
		Ip:        Ip,
		Port:      Port,
		isStart:   false,
		OnlineMap: make(map[string]*User),
		nowPeople: 0,
		sendP:     0,
		sendTo:    make([]string, 16),
		Message:   make(chan string),

		win: &WinConfig{
			App: app.New(),
		},
	}

	fmt.Println("server Init success.....")
	return server
}

func (s *Server) Start() {

	s.win.Win = s.win.App.NewWindow("server")
	makeServer(s.win.Win, s)
	fmt.Println("make serverUI success ....")
	s.win.Win.ShowAndRun()

}

func (s *Server) Handle(conn net.Conn) {
	s.online(conn)
	go s.listenUser(conn)
}

func (s *Server) listenUser(conn net.Conn) {

	for {
		buf := make([]byte, 4096)
		n, err := conn.Read(buf)
		if n == 0 {
			s.offline(conn)
		}
		if err != nil {
			s.Message <- string(buf)

			fmt.Println("s.Message: " + string(buf))
		}

	}
}

func (s *Server) offline(conn net.Conn) {

}

func (s *Server) online(conn net.Conn) {
	User := NewUser(conn)
	go User.ListenMsg()
	fmt.Println(User.name, " online...")
	s.UserLock.Lock()
	s.OnlineMap[User.name] = User
	s.nowPeople++
	s.UserLock.Unlock()
}

func (s *Server) BroadCast(u *User, msg string) {
	sendMsg := "[" + u.Addr + "]" + u.name + ":" + msg
	s.Message <- sendMsg
}

func serverStart(s *Server) func() {
	return func() {
		if s.isStart == false {
			fmt.Println("prepared to start Server")
			go s.serverStart()
			s.isStart = true
		} else {
			dialog.ShowInformation("Warning", "Already start ", s.win.Win)
		}
	}
}

func (s *Server) ListenMessager() {
	for {
		fmt.Println(s.Message)
		msg := <-s.Message
		s.Message = nil
		fmt.Println(s.Message)

		if msg == "--all" {
			msg = GetKeys(s.OnlineMap, 1)[0]
		}
		for _, u := range s.OnlineMap {
			u.c <- msg
			fmt.Println(u.c)
		}
	}
}

func (s *Server) serverStart() {
	listen, err := net.Listen(s.ConnType, fmt.Sprintf("%s:%s", s.Ip, s.Port))
	if err != nil {
		dialog.ShowError(err, s.win.Win)
		fmt.Println("server Listen error", err)
		return
	}
	fmt.Println("server start to accept")
	defer listen.Close()
	go s.ListenMessager()

	for {
		conn, err := listen.Accept()
		if err != nil {
			dialog.ShowError(err, s.win.Win)
			fmt.Println("accept error", err)
			continue
		}

		go s.Handle(conn)
	}
}
