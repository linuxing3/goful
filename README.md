G# Goful

[![Go Report Card](https://goreportcard.com/badge/github.com/linuxing3/goful)](https://goreportcard.com/report/github.com/linuxing3/goful)
[![Go Reference](https://pkg.go.dev/badge/github.com/linuxing3/goful.svg)](https://pkg.go.dev/github.com/linuxing3/goful)
[![MIT License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/linuxing3/goful/blob/master/LICENSE)

> 本项目基于`anmitsu/goful`，进行了自己的配置

Goful 跨平台的简单快捷终端文件管理器

* 跨平台
* 多窗口，多工作区
* 命令行执行bash和tmux
* 模糊查找, 异步拷贝, 文件块, 批量重命名等

## 安装

### Go version >= 1.16

    $ go install github.com/linuxing3/goful@latest
    ...
    $ goful

### Go version < 1.16

    $ go get github.com/linuxing3/goful
    ...
    $ goful

## 用法

### [Tutorial Demos](.github/demo.md)

key                  | function
---------------------|-----------
`C-n` `down` `j`     | Move cursor down
`C-p` `up` `k`       | Move cursor up
`C-a` `home` `u`     | Move cursor top
`C-e` `end` `G`      | Move cursor bottom
`C-f` `C-i` `right` `l`| Move cursor right
`C-b` `left` `h`     | Move cursor left
`C-d`                | More move cursor down
`C-u`                | More move cursor up
`C-v` `pgdn`         | Page down
`M-v` `pgup`         | Page up
`M-n`                | Scroll down
`M-p`                | Scroll up
`C-h` `backspace` `u`| Change to upper directory
`~`                  | Change to home directory
`\`                  | Change to root directory
`w`                  | Change to neighbor directory
`C-o`                | Create directory window
`C-w`                | Close directory window
`M-f`                | Move next workspace
`M-b`                | Move previous workspace
`M-C-o`              | Create workspace
`M-C-w`              | Close workspace
`space`              | Toggle mark
`C-space`            | Invert mark
`C-l`                | Reload
`C-m` `o`            | Open
`i`                  | Open by pager
`s`                  | Sort
`v`                  | View
`b`                  | Bookmark
`e`                  | Editor
`x`                  | Command
`X`                  | External command
`f` `/`              | Find
`:`                  | Shell
`;`                  | Shell suspend
`n`                  | Make file
`K`                  | Make directory
`c`                  | Copy
`m`                  | Move
`r`                  | Rename
`R`                  | Bulk rename by regexp
`D`                  | Remove
`d`                  | Change directory
`g`                  | Glob
`$`                  | Glob recursive
`C-g` `C-[`          | Cancel
`q` `Q`              | Quit

更多信息请查看 [main.go](main.go)

## 自定义

Goful内有配置文件，可以直接修改`main.go`.

* 修改添加快捷键 Change and add keybindings
* 更改命令行程序 Change terminal and shell
* 修改文件打开程序 Change file opener (editor, pager and more)
* 添加书签 Adding bookmarks
* 设置颜色和外观 Setting colors and looks

克隆原仓库到本地，删除所有子文件夹，仅保留`main.go`文件等

    $ cd $GOPATH/src/github.com/linuxing3/goful
    $ go mod init github.com/<yourname>/goful 
    $ go install

## 贡献

[Contributing Guide](.github/CONTRIBUTING.md)