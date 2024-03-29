package main

import (
	"os"
	"runtime"
	"strings"

	"github.com/anmitsu/goful/app"
	"github.com/anmitsu/goful/cmdline"
	"github.com/anmitsu/goful/filer"
	"github.com/anmitsu/goful/look"
	"github.com/anmitsu/goful/menu"
	"github.com/anmitsu/goful/message"
	"github.com/anmitsu/goful/widget"
	conf "github.com/linuxing3/goful/config"
	"github.com/mattn/go-runewidth"
)

func main() {
	// ui.TestListBox()
	start()
}

func start() {
	// 加载自定义设置
	conf.InitConfig()
	widget.Init()
	defer widget.Fini()

	// 修改终端标题
	if strings.Contains(os.Getenv("TERM"), "screen") {
		os.Stdout.WriteString("\033kgoful\033") // for tmux
	} else {
		os.Stdout.WriteString("\033]0;goful\007") // for otherwise
	}

	const state = "~/.goful/state.json"
	const history = "~/.goful/history/shell"

	goful := app.NewGoful(state)
	config(goful)
	_ = cmdline.LoadHistory(history)

	goful.Run()

	_ = goful.SaveState(state)
	_ = cmdline.SaveHistory(history)
}

func config(g *app.Goful) {

	look.Set("default") // default, midnight, black, white

	if runewidth.EastAsianWidth {
		// Because layout collapsing for ambiguous runes if LANG=ja_JP.
		widget.SetBorder('|', '-', '+', '+', '+', '+')
	} else {
		// Look good if environment variable RUNEWIDTH_EASTASIAN=0 and
		// ambiguous char setting is half-width for gnome-terminal.
		widget.SetBorder('│', '─', '┌', '┐', '└', '┘') // 0x2502, 0x2500, 0x250c, 0x2510, 0x2514, 0x2518
	}
	g.SetBorderStyle(widget.AllBorder) // AllBorder, ULBorder, NoBorder

	message.SetInfoLog("~/.goful/log/info.log")   // "" is not logging
	message.SetErrorLog("~/.goful/log/error.log") // "" is not logging
	message.Sec(5)                                // display second for a message

	// Setup widget keymaps.
	g.ConfigFiler(filerKeymap)
	filer.ConfigFinder(finderKeymap)
	cmdline.Config(cmdlineKeymap)
	cmdline.ConfigCompletion(completionKeymap)
	menu.Config(menuKeymap)

	filer.SetStatView(true, false, true)  // size, permission and time
	filer.SetTimeFormat("06-01-02 15:04") // ex: "Jan _2 15:04"

	console := "cmd"
	switch runtime.GOOS {
	case "windows":
		console = "cmd"
	case "linux":
		console = "zsh"
	}
	g.AddKeymap("t", func() { g.Shell(console) })

	// C-m或回车打开文件
	// 宏 %f 代表当前文件
	opener := "xdg-open %f %&"
	switch runtime.GOOS {
	case "windows":
		opener = "explorer %~f %&"
	case "darwin":
		opener = "open %f %&"
	}
	g.MergeKeymap(widget.Keymap{
		"C-m": func() { g.Spawn(opener) },
		"o":   func() { g.Spawn(opener) },
	})

	// 设置页 pager by $PAGER
	pager := os.Getenv("PAGER")
	if pager == "" {
		if runtime.GOOS == "windows" {
			pager = "bat"
		} else {
			pager = "bat"
		}
	}
	if runtime.GOOS == "windows" {
		pager += " %~f"
	} else {
		pager += " %f"
	}
	g.AddKeymap("i", func() { g.Spawn(pager) })

	// shell 或 terminal 用于执行外部命令
	// The shell is called when execute on background by the macro %&.
	// The terminal is called when the other.
	if runtime.GOOS == "windows" {
		g.ConfigShell(func(cmd string) []string {
			return []string{"cmd", "/c", cmd}
		})
		g.ConfigTerminal(func(cmd string) []string {
			return []string{"cmd", "/c", "start", "cmd", "/c", cmd + "& pause"}
		})
	} else {
		g.ConfigShell(func(cmd string) []string {
			return []string{"bash", "-c", cmd}
		})
		g.ConfigTerminal(func(cmd string) []string {
			// for not close the terminal when the shell finishes running
			const tail = `;read -p "HIT ENTER KEY"`

			if strings.Contains(os.Getenv("TERM"), "screen") { // such as screen and tmux
				return []string{"tmux", "new-window", "-n", cmd, cmd + tail}
			}
			// To execute bash in gnome-terminal of a new window or tab.
			title := "echo -n '\033]0;" + cmd + "\007';" // for change title
			return []string{"gnome-terminal", "--", "bash", "-c", title + cmd + tail}
		})
	}

	// 设置菜单并添加触发按键
	menu.Add("sort",
		"n", "sort name          ", func() { g.Dir().SortName() },
		"N", "sort name decending", func() { g.Dir().SortNameDec() },
		"s", "sort size          ", func() { g.Dir().SortSize() },
		"S", "sort size decending", func() { g.Dir().SortSizeDec() },
		"t", "sort time          ", func() { g.Dir().SortMtime() },
		"T", "sort time decending", func() { g.Dir().SortMtimeDec() },
		"e", "sort ext           ", func() { g.Dir().SortExt() },
		"E", "sort ext decending ", func() { g.Dir().SortExtDec() },
		".", "toggle priority    ", func() { filer.TogglePriority(); g.Workspace().ReloadAll() },
	)
	g.AddKeymap("s", func() { g.Menu("sort") })

	menu.Add("view",
		"s", "stat menu    ", func() { g.Menu("stat") },
		"l", "layout menu  ", func() { g.Menu("layout") },
		"L", "look menu    ", func() { g.Menu("look") },
		".", "toggle show hidden files", func() { filer.ToggleShowHiddens(); g.Workspace().ReloadAll() },
	)
	g.AddKeymap("v", func() { g.Menu("view") })

	menu.Add("layout",
		"t", "tile       ", func() { g.Workspace().LayoutTile() },
		"T", "tile-top   ", func() { g.Workspace().LayoutTileTop() },
		"b", "tile-bottom", func() { g.Workspace().LayoutTileBottom() },
		"r", "one-row    ", func() { g.Workspace().LayoutOnerow() },
		"c", "one-column ", func() { g.Workspace().LayoutOnecolumn() },
		"f", "fullscreen ", func() { g.Workspace().LayoutFullscreen() },
	)

	menu.Add("stat",
		"s", "toggle size  ", func() { filer.ToggleSizeView() },
		"p", "toggle perm  ", func() { filer.TogglePermView() },
		"t", "toggle time  ", func() { filer.ToggleTimeView() },
		"1", "all stat     ", func() { filer.SetStatView(true, true, true) },
		"0", "no stat      ", func() { filer.SetStatView(false, false, false) },
	)

	menu.Add("look",
		"d", "default      ", func() { look.Set("default") },
		"n", "midnight     ", func() { look.Set("midnight") },
		"b", "black        ", func() { look.Set("black") },
		"w", "white        ", func() { look.Set("white") },
		"a", "all border   ", func() { g.SetBorderStyle(widget.AllBorder) },
		"u", "ul border    ", func() { g.SetBorderStyle(widget.ULBorder) },
		"0", "no border    ", func() { g.SetBorderStyle(widget.NoBorder) },
	)

	menu.Add("command",
		"c", "copy         ", func() { g.Copy() },
		"m", "move         ", func() { g.Move() },
		"D", "delete       ", func() { g.Remove() },
		"k", "mkdir        ", func() { g.Mkdir() },
		"n", "newfile      ", func() { g.Touch() },
		"M", "chmod        ", func() { g.Chmod() },
		"r", "rename       ", func() { g.Rename() },
		"R", "bulk rename  ", func() { g.BulkRename() },
		"d", "chdir        ", func() { g.Chdir() },
		"g", "glob         ", func() { g.Glob() },
		"G", "globdir      ", func() { g.Globdir() },
	)
	g.AddKeymap("x", func() { g.Menu("command") })

	// 添加自定义外部命令
	conf.CustomizeConfig(g, "external-command")
	if runtime.GOOS == "windows" {
		menu.Add("external-command",
			"Z", "Config Template", func() { conf.MakeDefaultConfig(conf.ConfigTemplate) },
		)
	} else {
		menu.Add("external-command",
			"Z", "Config Template", func() { conf.MakeDefaultConfig(conf.ConfigTemplate) },
		)
	}
	g.AddKeymap("X", func() { g.Menu("external-command") })

	if runtime.GOOS == "windows" {
		menu.Add("build",
			"1", "xmake build ", func() { g.Shell("xmake") },
			"2", "xmake run ", func() { g.Shell("xmake r") },
			"3", "cmake build ", func() { g.Shell("cmake --build .") },
			"4", "compile commands ", func() { g.Shell("xmake project -k compile_commands") },
			"5", "cargo build ", func() { g.Shell("cargo build") },
			"6", "cargo run ", func() { g.Shell("cargo run") },
			"7", "go build ", func() { g.Shell("go build") },
			"8", "go run ", func() { g.Shell("go run") },
		)
	} else {
		menu.Add("build",
			"1", "xmake build ", func() { g.Shell("xmake") },
			"2", "xmake run ", func() { g.Shell("xmake r") },
			"3", "cmake build ", func() { g.Shell("cmake --build .") },
			"4", "compile commands ", func() { g.Shell("xmake project -k compile_commands") },
			"5", "cargo build ", func() { g.Shell("cargo build") },
			"6", "cargo run ", func() { g.Shell("cargo run") },
			"7", "go build ", func() { g.Shell("go build") },
			"8", "go run ", func() { g.Shell("go run") },
		)
	}
	conf.CustomizeConfig(g, "build")
	g.AddKeymap("B", func() { g.Menu("build") })

	menu.Add("archive",
		"z", "zip     ", func() { g.Shell(`zip -roD %x.zip %m`, -7) },
		"t", "tar     ", func() { g.Shell(`tar cvf %x.tar %m`, -7) },
		"g", "tar.gz  ", func() { g.Shell(`tar cvfz %x.tgz %m`, -7) },
		"b", "tar.bz2 ", func() { g.Shell(`tar cvfj %x.bz2 %m`, -7) },
		"x", "tar.xz  ", func() { g.Shell(`tar cvfJ %x.txz %m`, -7) },
		"r", "rar     ", func() { g.Shell(`rar u %x.rar %m`, -7) },

		"Z", "extract zip for %m", func() { g.Shell(`for i in %m; do unzip "$i" -d ./; done`, -6) },
		"T", "extract tar for %m", func() { g.Shell(`for i in %m; do tar xvf "$i" -C ./; done`, -6) },
		"G", "extract tgz for %m", func() { g.Shell(`for i in %m; do tar xvfz "$i" -C ./; done`, -6) },
		"B", "extract bz2 for %m", func() { g.Shell(`for i in %m; do tar xvfj "$i" -C ./; done`, -6) },
		"X", "extract txz for %m", func() { g.Shell(`for i in %m; do tar xvfJ "$i" -C ./; done`, -6) },
		"R", "extract rar for %m", func() { g.Shell(`for i in %m; do unrar x "$i" -C ./; done`, -6) },

		"1", "find . *.zip extract", func() { g.Shell(`find . -name "*.zip" -type f -prune -print0 | xargs -n1 -0 unzip -d ./`) },
		"2", "find . *.tar extract", func() { g.Shell(`find . -name "*.tar" -type f -prune -print0 | xargs -n1 -0 tar xvf -C ./`) },
		"3", "find . *.tgz extract", func() { g.Shell(`find . -name "*.tgz" -type f -prune -print0 | xargs -n1 -0 tar xvfz -C ./`) },
		"4", "find . *.bz2 extract", func() { g.Shell(`find . -name "*.bz2" -type f -prune -print0 | xargs -n1 -0 tar xvfj -C ./`) },
		"5", "find . *.txz extract", func() { g.Shell(`find . -name "*.txz" -type f -prune -print0 | xargs -n1 -0 tar xvfJ -C ./`) },
		"6", "find . *.rar extract", func() { g.Shell(`find . -name "*.rar" -type f -prune -print0 | xargs -n1 -0 unrar x -C ./`) },
	)

	// 添加自定义书签
	conf.CustomizeConfig(g, "bookmark")
	menu.Add("bookmark",
		"1", "~/Desktop  ", func() { g.Dir().Chdir("~/Desktop") },
		"2", "~/Documents", func() { g.Dir().Chdir("~/Documents") },
		"3", "~/Downloads", func() { g.Dir().Chdir("~/Downloads") },
		"4", "~/Music    ", func() { g.Dir().Chdir("~/Music") },
		"5", "~/Pictures ", func() { g.Dir().Chdir("~/Pictures") },
		"6", "~/Videos   ", func() { g.Dir().Chdir("~/Videos") },
		"7", "/etc          ", func() { g.Dir().Chdir("/etc") },
		"8", "/usr/local/etc", func() { g.Dir().Chdir("/usr/local/etc") },
		"9", "/usr          ", func() { g.Dir().Chdir("/usr") },
		"a", "/media        ", func() { g.Dir().Chdir("/media") },
		"b", "/mnt          ", func() { g.Dir().Chdir("/mnt") },
	)
	g.AddKeymap("b", func() { g.Menu("bookmark") })

	// 添加windows下自定义编辑器
	conf.CustomizeConfig(g, "editor")
	menu.Add("editor",
		"v", "nvim          ", func() { g.Spawn("nvim %f") },
		"V", "vim           ", func() { g.Spawn("vim %f") },
	)
	g.AddKeymap("e", func() { g.Menu("editor") })

	// 添加git菜单
	conf.CustomizeConfig(g, "git")
	g.AddKeymap("a", func() { g.Menu("git") })

	conf.CustomizeConfig(g, "web")
	g.AddKeymap("W", func() { g.Menu("web") })

	menu.Add("image",
		"x", "default    ", func() { g.Spawn(opener) },
		"e", "eog        ", func() { g.Spawn("eog %f %&") },
		"g", "gimp       ", func() { g.Spawn("gimp %m %&") },
	)

	menu.Add("media",
		"x", "default ", func() { g.Spawn(opener) },
		"m", "mpv     ", func() { g.Spawn("mpv %f") },
		"v", "vlc     ", func() { g.Spawn("vlc %f %&") },
	)

	var associate widget.Keymap
	if runtime.GOOS == "windows" {
		associate = widget.Keymap{
			".dir":  func() { g.Dir().EnterDir() },
			".html": func() { g.Shell("msedge %~f") },

			".md":   func() { g.Menu("editor") },
			".json": func() { g.Menu("editor") },
			".yml":  func() { g.Menu("editor") },
			".yaml": func() { g.Menu("editor") },
			".cmd":  func() { g.Menu("editor") },
			".bat":  func() { g.Menu("editor") },
			".vbs":  func() { g.Menu("editor") },
			".log":  func() { g.Menu("editor") },
			".org":  func() { g.Shell("runemacs -q %~f") },
			".el":   func() { g.Shell("runemacs -q %~f") },
			".vim":  func() { g.Shell("nvim-qt %~f") },

			".go": func() { g.Shell("go run %~f") },
			".py": func() { g.Shell("python %~f") },
			".rb": func() { g.Shell("ruby %~f") },
			".js": func() { g.Shell("node %~f") },
			".ts": func() { g.Shell("deno %~f") },
			".rs": func() { g.Shell("cargo run %~f") },
		}
	} else {
		associate = widget.Keymap{
			".dir":  func() { g.Dir().EnterDir() },
			".exec": func() { g.Shell(" ./" + g.File().Name()) },

			".zip": func() { g.Shell("unzip %f -d %D") },
			".tar": func() { g.Shell("tar xvf %f -C %D") },
			".gz":  func() { g.Shell("tar xvfz %f -C %D") },
			".tgz": func() { g.Shell("tar xvfz %f -C %D") },
			".bz2": func() { g.Shell("tar xvfj %f -C %D") },
			".xz":  func() { g.Shell("tar xvfJ %f -C %D") },
			".txz": func() { g.Shell("tar xvfJ %f -C %D") },
			".rar": func() { g.Shell("unrar x %f -C %D") },

			".html": func() { g.Shell("elinks %f") },
			".md":   func() { g.Shell("micro %f") },
			".json": func() { g.Shell("micro %f") },
			".yml":  func() { g.Shell("micro %f") },
			".yaml": func() { g.Shell("micro %f") },
			".log":  func() { g.Shell("micro %f") },
			".org":  func() { g.Shell("emacs -q %f") },
			".el":   func() { g.Shell("emacs -q %f") },
			".vim":  func() { g.Shell("nvim %f") },

			".go": func() { g.Shell("go run %f") },
			".py": func() { g.Shell("python %f") },
			".rb": func() { g.Shell("ruby %f") },
			".js": func() { g.Shell("node %f") },
			".ts": func() { g.Shell("deno %f") },
			".rs": func() { g.Shell("cargo run %f") },

			".jpg":  func() { g.Menu("image") },
			".jpeg": func() { g.Menu("image") },
			".gif":  func() { g.Menu("image") },
			".png":  func() { g.Menu("image") },
			".bmp":  func() { g.Menu("image") },

			".avi":  func() { g.Menu("media") },
			".mp4":  func() { g.Menu("media") },
			".mkv":  func() { g.Menu("media") },
			".wmv":  func() { g.Menu("media") },
			".flv":  func() { g.Menu("media") },
			".mp3":  func() { g.Menu("media") },
			".flac": func() { g.Menu("media") },
			".tta":  func() { g.Menu("media") },
		}
	}

	g.MergeExtmap(widget.Extmap{
		"C-m": associate,
		"o":   associate,
	})
}

// Widget keymap functions.

func filerKeymap(g *app.Goful) widget.Keymap {
	return widget.Keymap{
		"M-C-o": func() { g.CreateWorkspace() },
		"M-C-w": func() { g.CloseWorkspace() },
		"C-n":   func() { g.CreateWorkspace() },
		"C-q":   func() { g.CloseWorkspace() },
		"C-f":   func() { g.MoveWorkspace(1) },
		"C-b":   func() { g.MoveWorkspace(-1) },
		"M-f":   func() { g.MoveWorkspace(1) },
		"M-b":   func() { g.MoveWorkspace(-1) },
		"C-o":   func() { g.Workspace().CreateDir() },
		"C-w":   func() { g.Workspace().CloseDir() },
		"C-l":   func() { g.Workspace().ReloadAll() },
		// "C-f":       func() { g.Workspace().MoveFocus(1) },
		// "C-b":       func() { g.Workspace().MoveFocus(-1) },
		"right": func() { g.Workspace().MoveFocus(1) },
		"left":  func() { g.Workspace().MoveFocus(-1) },
		// "C-i":       func() { g.Workspace().MoveFocus(1) },
		"l": func() { g.Workspace().MoveFocus(1) },
		"h": func() { g.Workspace().MoveFocus(-1) },
		"F": func() { g.Workspace().SwapNextDir() },
		// "B":         func() { g.Workspace().SwapPrevDir() },
		"w":         func() { g.Workspace().ChdirNeighbor() },
		"C-h":       func() { g.Dir().Chdir("..") },
		"backspace": func() { g.Dir().Chdir("..") },
		"^":         func() { g.Dir().Chdir("..") },
		"~":         func() { g.Dir().Chdir("~") },
		"\\":        func() { g.Dir().Chdir("/") },
		// "C-n":       func() { g.Dir().MoveCursor(1) },
		// "C-p":       func() { g.Dir().MoveCursor(-1) },
		"down":    func() { g.Dir().MoveCursor(1) },
		"up":      func() { g.Dir().MoveCursor(-1) },
		"j":       func() { g.Dir().MoveCursor(1) },
		"k":       func() { g.Dir().MoveCursor(-1) },
		"C-d":     func() { g.Dir().MoveCursor(5) },
		"C-u":     func() { g.Dir().MoveCursor(-5) },
		"C-a":     func() { g.Dir().MoveTop() },
		"C-e":     func() { g.Dir().MoveBottom() },
		"home":    func() { g.Dir().MoveTop() },
		"end":     func() { g.Dir().MoveBottom() },
		"u":       func() { g.Dir().MoveTop() },
		"G":       func() { g.Dir().MoveBottom() },
		"M-n":     func() { g.Dir().Scroll(1) },
		"M-p":     func() { g.Dir().Scroll(-1) },
		"C-v":     func() { g.Dir().PageDown() },
		"M-v":     func() { g.Dir().PageUp() },
		"pgdn":    func() { g.Dir().PageDown() },
		"pgup":    func() { g.Dir().PageUp() },
		" ":       func() { g.Dir().ToggleMark() },
		"C-space": func() { g.Dir().InvertMark() },
		"C-g":     func() { g.Dir().Reset() },
		"C-[":     func() { g.Dir().Reset() }, // C-[ means ESC
		"f":       func() { g.Dir().Finder() },
		"/":       func() { g.Dir().Finder() },
		"q":       func() { g.Quit() },
		"Q":       func() { g.Quit() },
		":":       func() { g.Shell("") },
		"M-x":     func() { g.Shell("") },
		";":       func() { g.ShellSuspend("") },
		"M-W":     func() { g.ChangeWorkspaceTitle() },
		"n":       func() { g.Touch() },
		"K":       func() { g.Mkdir() },
		"c":       func() { g.Copy() },
		"m":       func() { g.Move() },
		"r":       func() { g.Rename() },
		"R":       func() { g.BulkRename() },
		"D":       func() { g.Remove() },
		"d":       func() { g.Chdir() },
		"g":       func() { g.Glob() },
		"$":       func() { g.Globdir() },
	}
}

func finderKeymap(w *filer.Finder) widget.Keymap {
	return widget.Keymap{
		"C-h":       func() { w.DeleteBackwardChar() },
		"backspace": func() { w.DeleteBackwardChar() },
		"M-p":       func() { w.MoveHistory(1) },
		"M-n":       func() { w.MoveHistory(-1) },
		"C-g":       func() { w.Exit() },
		"C-[":       func() { w.Exit() },
	}
}

func cmdlineKeymap(w *cmdline.Cmdline) widget.Keymap {
	return widget.Keymap{
		"C-a":       func() { w.MoveTop() },
		"C-e":       func() { w.MoveBottom() },
		"C-f":       func() { w.ForwardChar() },
		"C-b":       func() { w.BackwardChar() },
		"right":     func() { w.ForwardChar() },
		"left":      func() { w.BackwardChar() },
		"M-f":       func() { w.ForwardWord() },
		"M-b":       func() { w.BackwardWord() },
		"C-d":       func() { w.DeleteChar() },
		"delete":    func() { w.DeleteChar() },
		"C-h":       func() { w.DeleteBackwardChar() },
		"backspace": func() { w.DeleteBackwardChar() },
		"M-d":       func() { w.DeleteForwardWord() },
		"M-h":       func() { w.DeleteBackwardWord() },
		"C-k":       func() { w.KillLine() },
		"C-i":       func() { w.StartCompletion() },
		"C-m":       func() { w.Run() },
		"C-g":       func() { w.Exit() },
		"C-[":       func() { w.Exit() },
		"C-n":       func() { w.History.CursorDown() },
		"C-p":       func() { w.History.CursorUp() },
		"down":      func() { w.History.CursorDown() },
		"up":        func() { w.History.CursorUp() },
		"C-v":       func() { w.History.PageDown() },
		"M-v":       func() { w.History.PageUp() },
		"pgdn":      func() { w.History.PageDown() },
		"pgup":      func() { w.History.PageUp() },
		"M-<":       func() { w.History.MoveTop() },
		"M->":       func() { w.History.MoveBottom() },
		"home":      func() { w.History.MoveTop() },
		"end":       func() { w.History.MoveBottom() },
		"M-n":       func() { w.History.Scroll(1) },
		"M-p":       func() { w.History.Scroll(-1) },
		"C-x":       func() { w.History.Delete() },
	}
}

func completionKeymap(w *cmdline.Completion) widget.Keymap {
	return widget.Keymap{
		"C-n":   func() { w.CursorDown() },
		"C-p":   func() { w.CursorUp() },
		"down":  func() { w.CursorDown() },
		"up":    func() { w.CursorUp() },
		"C-f":   func() { w.CursorToRight() },
		"C-b":   func() { w.CursorToLeft() },
		"right": func() { w.CursorToRight() },
		"left":  func() { w.CursorToLeft() },
		"C-i":   func() { w.CursorToRight() },
		"C-v":   func() { w.PageDown() },
		"M-v":   func() { w.PageUp() },
		"pgdn":  func() { w.PageDown() },
		"pgup":  func() { w.PageUp() },
		"M-<":   func() { w.MoveTop() },
		"M->":   func() { w.MoveBottom() },
		"home":  func() { w.MoveTop() },
		"end":   func() { w.MoveBottom() },
		"M-n":   func() { w.Scroll(1) },
		"M-p":   func() { w.Scroll(-1) },
		"C-m":   func() { w.InsertCompletion() },
		"C-g":   func() { w.Exit() },
		"C-[":   func() { w.Exit() },
	}
}

func menuKeymap(w *menu.Menu) widget.Keymap {
	return widget.Keymap{
		"C-n":  func() { w.MoveCursor(1) },
		"C-p":  func() { w.MoveCursor(-1) },
		"down": func() { w.MoveCursor(1) },
		"up":   func() { w.MoveCursor(-1) },
		"C-v":  func() { w.PageDown() },
		"M-v":  func() { w.PageUp() },
		"M->":  func() { w.MoveBottom() },
		"M-<":  func() { w.MoveTop() },
		"C-m":  func() { w.Exec() },
		"C-g":  func() { w.Exit() },
		"C-[":  func() { w.Exit() },
	}
}
