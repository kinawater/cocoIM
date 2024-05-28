package clinet

import (
	"cocoIM/clinet/sdk"
	"fmt"
	"github.com/gookit/color"
	"github.com/rocket049/gocui"
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

func RunMain() {
	// 测试下
	fmt.Println("这是客户端")
}
