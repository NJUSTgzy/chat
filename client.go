package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
	"net"
	"strings"
)

type cWinConfig struct {
	App      fyne.App
	Win      fyne.Window
	config   *fyne.Container
	friends  *fyne.Container
	chatRoom *fyne.Container
	allUser  *fyne.Container
}

type Client struct {
	win *cWinConfig
	usr *User
}

func cInit() *Client {

	client := &Client{
		win: &cWinConfig{
			App: app.New(),
		},
	}

	fmt.Println("client Init success.....")
	return client
}

func (c *Client) cStart(IP string, PORT string) {

	c.win.Win = c.win.App.NewWindow("client")
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", IP, PORT))
	if err != nil {
		fmt.Println("connect err")
	} else {
		fmt.Println("client connected")
		c.usr = NewUser(conn)
	}

	go c.listenMsg()
	makeClientUI(c.win.Win, c)
	fmt.Println("make clientUI success ....")
	c.win.Win.ShowAndRun()

}

func (c *Client) listenMsg() {

	for {
		buf := make([]byte, 4096)
		read, err := c.usr.conn.Read(buf)
		if err != nil {
			return
		}
		if read == 0 {
			continue
		} else {
			a := string(buf)
			msg := widget.NewLabel(a)
			c.win.chatRoom.Add(msg)
			c.win.chatRoom.Refresh()
		}
	}
}

func GetAllUsr(c *Client) []string {
	c.usr.conn.Write([]byte("--all"))
	recvBuf := make([]byte, 4096)
	c.usr.conn.Read(recvBuf)
	allUsers := string(recvBuf)
	return strings.Split(allUsers, ";")
}
