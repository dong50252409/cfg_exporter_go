package main

import (
	"bufio"
	"bytes"
	"cfg_exporter/config"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"image/color"
	"os"
	"path/filepath"
	"strings"
)

type fileInfo struct {
	fullFilename string
	filename     string
}

var (
	searchEntryWidget    *widget.Entry
	fileListWidget       *widget.List
	allButtonWidget      *widget.Button
	clientButtonWidget   *widget.Button
	serverButtonWidget   *widget.Button
	checkButtonWidget    *widget.Check
	multiLineEntryWidget *MultiLineEntryEx
)

var (
	logBuffers     bytes.Buffer
	filenameList   []fileInfo
	selectFilename string
)

type MultiLineEntryEx struct {
	*widget.Entry
}

func NewMultiLineEntry() *MultiLineEntryEx {
	e := &MultiLineEntryEx{&widget.Entry{MultiLine: true, Wrapping: fyne.TextWrapBreak}}
	e.ExtendBaseWidget(e)
	return e
}

func (e *MultiLineEntryEx) Append(text string) {
	if len(multiLineEntryWidget.Text) > 1024*50 { // 50KB 就需要清空一下
		logBuffers.Reset()
		logBuffers.WriteString("清空日志缓存...\n")
	} else {
		logBuffers.WriteString(text + "\n")
	}
	onChanged := e.OnChanged
	e.OnChanged = nil
	multiLineEntryWidget.SetText(logBuffers.String())
	e.OnChanged = onChanged
	multiLineEntryWidget.CursorRow = len(multiLineEntryWidget.Text) - 1
}

func (e *MultiLineEntryEx) CleanText() {
	logBuffers.Reset()
	onChanged := e.OnChanged
	e.OnChanged = nil
	multiLineEntryWidget.SetText(logBuffers.String())
	e.OnChanged = onChanged
	multiLineEntryWidget.CursorRow = len(multiLineEntryWidget.Text) - 1
}

//func (e *MultiLineEntryEx) TappedSecondary(_ *fyne.PointEvent) {
//	// 防止右键崩溃
//}

func startUI() {
	fyneApp := app.New()
	window := fyneApp.NewWindow("配置表导出工具")
	window.Resize(fyne.NewSize(800, 700))
	window.SetFixedSize(true)
	// 输入栏
	initSearchEntry()

	// 文件搜索结果列表
	initList()

	// 初始化按钮
	initBottom()

	// 初始化文本框
	initText()

	left := createLeft(window)
	right := createRight(window)

	content := container.NewHBox(left, right)
	window.SetContent(content)

	window.ShowAndRun()
}

func createLeft(window fyne.Window) *fyne.Container {
	// 创建带有颜色的矩形作为边框
	topBorder := canvas.NewRectangle(color.RGBA{R: 0, G: 0, B: 255, A: 255}) // 蓝色边框
	topBorder.Resize(fyne.NewSize(window.Canvas().Size().Width, 10))         // 设置边框高度

	leftBorder := canvas.NewRectangle(color.RGBA{R: 0, G: 0, B: 255, A: 255})
	leftBorder.Resize(fyne.NewSize(10, window.Canvas().Size().Height-20)) // 减去上下边框的高度

	rightBorder := canvas.NewRectangle(color.RGBA{R: 0, G: 0, B: 255, A: 255})
	rightBorder.Resize(fyne.NewSize(10, window.Canvas().Size().Height-20))

	bottomBorder := canvas.NewRectangle(color.RGBA{R: 0, G: 0, B: 255, A: 255})
	bottomBorder.Resize(fyne.NewSize(window.Canvas().Size().Width, 10))

	searchEntryContainer := container.NewGridWrap(fyne.NewSize(300, 40), searchEntryWidget)
	fileListContainer := container.NewGridWrap(fyne.NewSize(300, 635), fileListWidget)
	left := container.NewBorder(topBorder, bottomBorder, leftBorder, rightBorder, container.NewVBox(searchEntryContainer, fileListContainer))
	return left
}

func initSearchEntry() {
	searchEntryWidget = widget.NewEntry()
	searchEntryWidget.SetPlaceHolder("请输入配置文件名，支持模糊搜索")
	// 绑定搜索功能到输入栏
	searchEntryWidget.OnChanged = func(s string) {
		searchFiles(searchEntryWidget.Text)
	}
}

func searchFiles(text string) {

	// 调用filepath.Walk函数遍历目录
	err := filepath.Walk(config.Config.Source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 检查是否是文件
		if !info.IsDir() {
			if strings.Contains(info.Name(), text) {
				filenameList = append(filenameList, fileInfo{fullFilename: path, filename: info.Name()})
			}
		}
		return nil
	})
	if err != nil {
		return
	}

	fileListWidget.Length = func() int { return len(filenameList) }
	fileListWidget.UpdateItem = func(id widget.ListItemID, item fyne.CanvasObject) {
		item.(*widget.Label).SetText(filenameList[id].filename)
	}
	fileListWidget.Refresh()
}

func initList() {
	fileListWidget = widget.NewList(
		func() int { return 0 },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(i widget.ListItemID, o fyne.CanvasObject) {},
	)
	searchFiles("")
	fileListWidget.OnSelected = func(id widget.ListItemID) {
		selectFilename = filenameList[id].fullFilename
	}
}

func createRight(window fyne.Window) *fyne.Container {
	// 创建带有颜色的矩形作为边框
	topBorder := canvas.NewRectangle(color.RGBA{R: 255, G: 0, B: 0, A: 255}) // 红色边框
	topBorder.Resize(fyne.NewSize(window.Canvas().Size().Width, 10))         // 设置边框高度

	leftBorder := canvas.NewRectangle(color.RGBA{R: 255, G: 0, B: 0, A: 255})
	leftBorder.Resize(fyne.NewSize(10, window.Canvas().Size().Height-20)) // 减去上下边框的高度

	rightBorder := canvas.NewRectangle(color.RGBA{R: 255, G: 0, B: 0, A: 255})
	rightBorder.Resize(fyne.NewSize(10, window.Canvas().Size().Height-20))

	bottomBorder := canvas.NewRectangle(color.RGBA{R: 255, G: 0, B: 0, A: 255})
	bottomBorder.Resize(fyne.NewSize(window.Canvas().Size().Width, 10))

	buttonContainer := container.NewGridWrap(fyne.NewSize(100, 40), allButtonWidget, clientButtonWidget, serverButtonWidget, checkButtonWidget)
	textContainer := container.NewGridWrap(fyne.NewSize(470, 635), multiLineEntryWidget)

	right := container.NewBorder(topBorder, bottomBorder, leftBorder, rightBorder, container.NewVBox(buttonContainer, textContainer))
	return right
}

func initBottom() {
	allButtonWidget = widget.NewButton("全部导出", func() {
		multiLineEntryWidget.CleanText()
		//// TODO 复用数据
		clientTappedFunc()
		serverTappedFunc()
	})

	clientButtonWidget = widget.NewButton("客户端导出", func() {
		multiLineEntryWidget.CleanText()
		clientTappedFunc()
	})

	serverButtonWidget = widget.NewButton("服务端导出", func() {
		multiLineEntryWidget.CleanText()
		serverTappedFunc()
	})

	checkButtonWidget = widget.NewCheck("是否检查", func(b bool) { config.Config.Verify = b })
	checkButtonWidget.SetChecked(config.Config.Verify)
}

func clientTappedFunc() {
	if selectFilename == "" {
		return
	}
	config.Config.SchemaName = "flatbuffers"
	err := run(selectFilename)
	if err != nil {
		return
	}
}

func serverTappedFunc() {
	if selectFilename == "" {
		return
	}
	config.Config.SchemaName = "erlang"
	err := run(selectFilename)
	if err != nil {
		return
	}
}

func initText() {
	multiLineEntryWidget = NewMultiLineEntry()
	multiLineEntryWidget.OnChanged = func(s string) {
		before, found := strings.CutSuffix(multiLineEntryWidget.Text, s)
		if found {
			multiLineEntryWidget.SetText(before)
		}
	}

	r, w, _ := os.Pipe()

	//defer func() { _ = r.Close() }()
	//defer func() { _ = w.Close() }()

	os.Stdout = w
	os.Stderr = w
	go func() {
		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			multiLineEntryWidget.Append(scanner.Text())
		}
	}()
}
