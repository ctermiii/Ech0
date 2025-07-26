package ssh

import (
	"context"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/lin-snow/ech0/internal/config"
	"github.com/lin-snow/ech0/internal/tui"
)

var SSHServer *ssh.Server

func SSHStart() {
	host := config.Config.SSH.Host
	port := config.Config.SSH.Port
	key := config.Config.SSH.Key
	
	var err error

	SSHServer, err = wish.NewServer(
		wish.WithAddress(net.JoinHostPort(host, port)),
		wish.WithHostKeyPath(key),
		wish.WithMiddleware(
			BubbleteaMiddleware(teaHandler),
			ActivetermMiddleware(), // Bubble Tea apps usually require a PTY.
		),
	)
	if err != nil {
		// log.Error("Could not start server", "error", err)
	}

	// done := make(chan os.Signal, 1)
	// signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	// log.Info("Starting SSH server", "host", host, "port", port)
	go func() {
		fmt.Println("🚀 Ech0 SSH已启动，监听端口", port)
		if err = SSHServer.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
			// log.Error("Could not start server", "error", err)
			// done <- nil
		}
	}()

	// <-done
	// // log.Info("Stopping SSH server")
	// ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	// defer func() { cancel() }()
	// if err := s.Shutdown(ctx); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
	// 	// log.Error("Could not stop server", "error", err)
	// }
}

func SSHStop() error {
	if SSHServer == nil {
		return nil
	}

	// When it arrives, we create a context with a timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer func() { cancel() }()

	// When we start the shutdown, the server will no longer accept new
	// connections, but will wait as much as the given context allows for the
	// active connections to finish.
	// After the timeout, it shuts down anyway.
	if err := SSHServer.Shutdown(ctx); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
		// 强制关闭服务器
		SSHServer.Close()

		return err
	}

	SSHServer = nil // Clear the server instance
	return nil
}

// Middleware will exit 1 connections trying with no active terminals.
func ActivetermMiddleware() wish.Middleware {
	return func(next ssh.Handler) ssh.Handler {
		return func(sess ssh.Session) {
			_, _, active := sess.Pty()
			if active {
				next(sess)
				return
			}
			wish.Println(sess, "Requires an active PTY")
			_ = sess.Exit(1)
		}
	}
}

// You can wire any Bubble Tea model up to the middleware with a function that
// handles the incoming ssh.Session. Here we just grab the terminal info and
// pass it to the new model. You can also return tea.ProgramOptions (such as
// tea.WithAltScreen) on a session by session basis.
func teaHandler(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	// This should never fail, as we are using the activeterm middleware.
	pty, _, _ := s.Pty()

	// When running a Bubble Tea app over SSH, you shouldn't use the default
	// lipgloss.NewStyle function.
	// That function will use the color profile from the os.Stdin, which is the
	// server, not the client.
	// We provide a MakeRenderer function in the bubbletea middleware package,
	// so you can easily get the correct renderer for the current session, and
	// use it to create the styles.
	// The recommended way to use these styles is to then pass them down to
	// your Bubble Tea model.
	renderer := MakeRenderer(s)
	txtStyle := renderer.NewStyle().Foreground(lipgloss.Color("10"))
	quitStyle := renderer.NewStyle().Foreground(lipgloss.Color("8"))

	bg := "light"
	if renderer.HasDarkBackground() {
		bg = "dark"
	}

	m := model{
		term:      pty.Term,
		profile:   renderer.ColorProfile().Name(),
		width:     pty.Window.Width,
		height:    pty.Window.Height,
		bg:        bg,
		txtStyle:  txtStyle,
		quitStyle: quitStyle,
		logo: tui.GetLogoBanner(),
		textarea: textarea.New(),
	}

	m.textarea.Placeholder = "请输入..."
	m.textarea.Focus()

	return m, []tea.ProgramOption{tea.WithAltScreen()}
}

// Just a generic tea.Model to demo terminal information of ssh.
type model struct {
	term      string
	profile   string
	width     int
	height    int
	bg        string
	txtStyle  lipgloss.Style
	quitStyle lipgloss.Style
	textarea textarea.Model
	logo      string
}

func (m model) Init() tea.Cmd {
	return textarea.Blink
}

// func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
// 	switch msg := msg.(type) {
// 	case tea.WindowSizeMsg:
// 		m.height = msg.Height
// 		m.width = msg.Width
// 	case tea.KeyMsg:
// 		switch msg.String() {
// 		case "q", "ctrl+c":
// 			return m, tea.Quit
// 		}
// 	}
// 	return m, nil
// }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    var cmds []tea.Cmd // 声明 cmds 切片
    var cmd tea.Cmd    // 声明 cmd 变量

    switch msg := msg.(type) {
    case tea.WindowSizeMsg:
        // 处理窗口大小变化消息。当 SSH 客户端窗口大小改变时，wish 中间件会发送此消息。
        m.height = msg.Height
        m.width = msg.Width
        // 将消息传递给 textarea，让它也能调整自身大小（如果需要）
        m.textarea, cmd = m.textarea.Update(msg)
        cmds = append(cmds, cmd) // 将 textarea 返回的命令添加到列表中
    case tea.KeyMsg:
        switch msg.String() {
        case "ctrl+c":
            // 处理 'Ctrl+C' 退出命令
            return m, tea.Quit
        case "esc": // 添加对 Esc 键的处理
            if m.textarea.Focused() {
                m.textarea.Blur()
            }
        default:
            // 如果 textarea 没有焦点，按下任意键使其获得焦点
            if !m.textarea.Focused() {
                cmd = m.textarea.Focus()
                cmds = append(cmds, cmd)
            }
        }
		// 将键盘消息传递给 textarea，让它处理文本输入
		m.textarea, cmd = m.textarea.Update(msg)
		cmds = append(cmds, cmd) // 将 textarea 返回的命令添加到列表中
	}

	// 返回更新后的模型和所有累积的命令
	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	return m.txtStyle.Render(m.logo) +
	"\n\n" +
	// fmt.Sprintf(
	// 	"\n\n"+"🧙 请把你的故事写在这里.\n\n%s",
	// 	m.textarea.View(),
	// ) +
	tui.GetEch0Info() +
	"\n\n🧙 欢迎使用 Ech0 SSH (更多功能开发中...)" +
	"\n\n" + m.quitStyle.Render("按 'ctrl+c' 退出\n")
}


