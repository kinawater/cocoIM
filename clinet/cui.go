package clinet

import (
	"cocoIM/common/sdk"
	"fmt"
	"github.com/gookit/color"
	"github.com/rocket049/gocui"
	"io"
	"log"
	"net"
	"os"
)

var (
	buf     string
	chat    *sdk.Chat
	step    int
	verbose bool
)

// 设置对话框标题
func setHeadText(gui *gocui.Gui, msg string) {
	v, err := gui.View("head")
	if err == nil {
		v.Clear()
		fmt.Fprintf(v, color.FgGreen.Text(msg))
	}

}

// 定义单条消息要展示的信息
type VOT struct {
	Name, Msg, Sep string
}

// 单条对话展示
func (self VOT) Show(g *gocui.Gui) error {
	view, err := g.View("out")
	if err != nil {
		return nil
	}
	fmt.Fprintf(view, "%v:%v%v\n", color.FgGreen.Text(self.Name), self.Sep, color.FgYellow.Text(self.Msg))
	return nil
}
func viewPrint(g *gocui.Gui, name, msg string, newLine bool) {
	var out VOT
	out.Name = name
	out.Msg = msg
	if newLine {
		out.Sep = "\n"
	} else {
		out.Sep = " "
	}
	g.Update(out.Show)
}

func doRecv(g *gocui.Gui) {
	recvChannel := chat.Recv()
	for msg := range recvChannel {
		switch msg.Type {
		case sdk.MsgTypeText:
			// 直接展示
			viewPrint(g, msg.Name, msg.Content, true)
		}
	}
	g.Close()
}

func quit(g *gocui.Gui, view *gocui.View) error {
	chat.Close()
	ov, _ := g.View("out")
	buf = ov.Buffer()
	g.Close()
	return gocui.ErrQuit
}

func doSay(g *gocui.Gui, cv *gocui.View) {
	v, err := g.View("out")
	if cv != nil && err == nil {
		// 从输入框取出输入的文字
		p := cv.ReadEditor()
		if p != nil {
			var msg = &sdk.Message{
				Type:       sdk.MsgTypeText,
				Name:       "花椒鱼",
				FromUserID: "123456",
				ToUserID:   "654321",
				Content:    string(p),
			}
			// 自己输入的话直接显示到消息中
			viewPrint(g, "我", msg.Content, false)
			// 再发送到服务器
			chat.Send(msg)
		}
		v.Autoscroll = true
	}
}
func viewUpdate(g *gocui.Gui, cv *gocui.View) error {
	doSay(g, cv)
	l := len(cv.Buffer())
	cv.MoveCursor(0-l, 0, true)
	cv.Clear()
	return nil
}
func viewUpScroll(g *gocui.Gui, cv *gocui.View) error {
	v, err := g.View("out")
	_, y := v.Size()
	ox, oy := v.Origin()
	lnum := len(v.BufferLines())
	if err == nil {
		if oy > lnum-y-1 {
			v.Autoscroll = true
		} else {
			v.SetOrigin(ox, oy+1)
		}
	}
	return nil
}
func viewDownScroll(g *gocui.Gui, cv *gocui.View) error {
	v, err := g.View("out")
	_, y := v.Size()
	ox, oy := v.Origin()
	lnum := len(v.BufferLines())
	if err == nil {
		if oy > lnum-y-1 {
			v.Autoscroll = true
		} else {
			v.SetOrigin(ox, oy+1)
		}
	}
	return nil
}
func viewOutput(g *gocui.Gui, x0, y0, x1, y1 int) error {
	v, err := g.SetView("out", x0, y0, x1, y1)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Wrap = true
		v.Overwrite = false
		v.Autoscroll = true
		v.SelBgColor = gocui.ColorRed
		v.Title = "Messages"
	}
	return nil
}
func viewInput(g *gocui.Gui, x0, y0, x1, y1 int) error {
	if v, err := g.SetView("main", x0, y0, x1, y1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		//当 err == gocui.ErrUnknownView 时运行
		v.Editable = true
		v.Wrap = true
		v.Overwrite = false
		if _, err := g.SetCurrentView("main"); err != nil {
			return err
		}
	}
	return nil
}
func viewHead(g *gocui.Gui, x0, y0, x1, y1 int) error {
	if v, err := g.SetView("head", x0, y0, x1, y1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Wrap = false
		v.Overwrite = true
		msg := "开始聊天了!"
		setHeadText(g, msg)
	}
	return nil
}
func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	//maxX = maxX - 100
	//maxY = maxY - 100
	if err := viewHead(g, 1, 1, maxX-1, 3); err != nil {
		return err
	}
	if err := viewOutput(g, 1, 4, maxX-1, maxY-4); err != nil {
		return err
	}
	if err := viewInput(g, 1, maxY-3, maxX-1, maxY-1); err != nil {
		return err
	}
	return nil
}

var pos int

func pasteUP(g *gocui.Gui, cv *gocui.View) error {
	v, err := g.View("out")
	if err != nil {
		fmt.Fprintf(cv, "error:%s", err)
		return nil
	}
	bls := v.BufferLines()
	lnum := len(bls)
	if pos < lnum-1 {
		pos++
	}
	cv.Clear()
	fmt.Fprintf(cv, "%s", bls[lnum-pos-1])
	return nil
}

func pasteDown(g *gocui.Gui, cv *gocui.View) error {
	v, err := g.View("out")
	if err != nil {
		fmt.Fprintf(cv, "error:%s", err)
		return nil
	}
	if pos > 0 {
		pos--
	}
	bls := v.BufferLines()
	lnum := len(bls)
	cv.Clear()
	fmt.Fprintf(cv, "%s", bls[lnum-pos-1])
	return nil
}

func RunMain() {
	// 测试下
	fmt.Println("这是客户端")
	// 创建chat
	chat = sdk.MakeNewChat(net.ParseIP("0.0.0.0"), 8900, "logic", "123456", "2131")
	// step2 创建GUI
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		// TODO 记录日志
		log.Panicln(err)
	}
	g.Cursor = true
	g.Mouse = false
	g.ASCII = false
	// 设置编排函数
	g.SetManagerFunc(layout)
	// 注册回调事件
	if err := g.SetKeybinding("main", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("main", gocui.KeyEnter, gocui.ModNone, viewUpdate); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("main", gocui.KeyPgup, gocui.ModNone, viewUpScroll); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("main", gocui.KeyPgdn, gocui.ModNone, viewDownScroll); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("main", gocui.KeyArrowDown, gocui.ModNone, pasteDown); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("main", gocui.KeyArrowUp, gocui.ModNone, pasteUP); err != nil {
		log.Panicln(err)
	}

	// 启动消费函数
	// 会一直查redv通道里有没有数据
	go doRecv(g)

	if err = g.MainLoop(); err != nil {
		log.Println(err)
	}

	file, err := os.Create("chat.log")
	if err != nil {
		log.Panicln(err)
	}
	defer file.Close()
	_, err = io.WriteString(file, string([]byte(buf)))
	if err != nil {
		log.Panicln(err)
	}
}
