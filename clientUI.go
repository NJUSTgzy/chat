package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"net"
)

func makeClientUI(w fyne.Window, c *Client) {
	c.win.Win.Resize(fyne.NewSize(600, 800))
	btnConnect := widget.NewButton("CONNECT", connectToServer(c))
	c.win.config = container.NewVBox(btnConnect)

	data := GetKeys(c.usr.FriendMap, 0)

	list := widget.NewList(
		func() int {
			return len(data)
		},
		func() fyne.CanvasObject {
			return widget.NewButton("test", nil)
		},
		func(id widget.ListItemID, object fyne.CanvasObject) {
			object.(*widget.Button).OnTapped = chatTo(c, data[id])
			object.(*widget.Button).SetText(data[id])
		},
	)
	msg := widget.NewLabel("you can chose a friend to chat")
	c.win.friends = container.NewVBox(msg, list)

	inPutEntry := widget.NewMultiLineEntry()
	btnSend := widget.NewButton("send", snedMsg(c, inPutEntry))
	c.win.chatRoom = container.NewVBox(inPutEntry, btnSend)

	btnAll := widget.NewButton("List All", listAll(c))

	c.win.allUser = container.NewVBox(btnAll)

	tabs := container.NewAppTabs(
		container.NewTabItem("config", c.win.config),
		container.NewTabItem("friends", c.win.friends),
		container.NewTabItem("chatRoom", c.win.chatRoom),
		container.NewTabItem("allUser", c.win.allUser),
	)

	tabs.SetTabLocation(container.TabLocationLeading)
	c.win.Win.SetContent(tabs)
	//connectToServer(w)

}

func snedMsg(c *Client, entry *widget.Entry) func() {
	return func() {
		msg := entry.Text
		c.usr.SendMsg(msg)
		entry.Text = ""
		entry.Refresh()
		sendMssg(c, msg)

	}
}

func listAll(c *Client) func() {
	return func() {
		allUsr := c.win.App.NewWindow("ALLUSERS")
		allUsr.Resize(fyne.NewSize(800, 600))

		data := GetAllUsr(c)
		list := widget.NewList(
			func() int {
				return len(data)
			},
			func() fyne.CanvasObject {
				return widget.NewButton("test", nil)
			},
			func(id widget.ListItemID, object fyne.CanvasObject) {
				object.(*widget.Button).OnTapped = chatTo(c, data[id])
				object.(*widget.Button).SetText(data[id])
			},
		)

		allUsr.SetContent(list)

		allUsr.Show()
	}
}

func (c *Client) listFriend() {
	data := GetKeys(c.usr.FriendMap, 0)

	list := widget.NewList(
		func() int {
			return len(data)
		},
		func() fyne.CanvasObject {
			return widget.NewButton("test", nil)
		},
		func(id widget.ListItemID, object fyne.CanvasObject) {
			object.(*widget.Button).OnTapped = chatTo(c, data[id])
			object.(*widget.Button).SetText(data[id])
		},
	)
	c.win.Win.SetContent(list)
}

func chatTo(c *Client, s string) func() {
	return func() {
		win := c.win.App.NewWindow("Chatting")
		win.Resize(fyne.NewSize(800, 600))
		inPutEntry := widget.NewMultiLineEntry()
		btnSend := widget.NewButton("send", nil)
		con := container.NewVBox(inPutEntry, btnSend)
		go func() {
			buf := make([]byte, 4096)
			for win != nil {
				c.usr.conn.Read(buf)
				con.Add(widget.NewRichTextWithText(string(buf)))
				con.Refresh()
			}

		}()

		win.SetContent(con)
		win.Show()
	}
}

func sendMssg(c *Client, text string) {
	c.usr.conn.Write([]byte(text))
}

func connectToServer(c *Client) func() {
	return func() {
		win := c.win.App.NewWindow("connectToServer")

		IP := widget.NewEntry()
		IP.SetText("127.0.0.2")
		IP.SetPlaceHolder("Please input ip")
		PORT := widget.NewEntry()
		PORT.SetText("9090")
		PORT.SetPlaceHolder("Please input port")

		form := &widget.Form{
			Items: []*widget.FormItem{
				{Text: "IP", Widget: IP},
				{Text: "PORT", Widget: PORT},
			},

			OnSubmit: func() {
				conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", IP.Text, PORT.Text))
				if err != nil {
					dialog.ShowError(err, c.win.Win)
				} else {
					fmt.Println("client no err in connection")
					c.usr = NewUser(conn)
					dialog.ShowInformation("success", "you successfully connect to server", c.win.Win)
					win.Close()
				}
			},
		}
		win.SetContent(form)
		win.Resize(fyne.NewSize(800, 600))
		go win.Show()
	}
}
