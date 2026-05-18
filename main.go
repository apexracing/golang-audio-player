package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/apexracing/golang-audio-player/audio"
)

func main() {
	var engine audio.AudioPlayerEngine
	if err := engine.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "初始化音频引擎失败: %v\n", err)
		os.Exit(1)
	}
	defer engine.Destroy()

	if len(os.Args) < 2 {
		listDevices(&engine)
		return
	}

	switch os.Args[1] {
	case "list":
		listDevices(&engine)
	case "play":
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, "用法: go run main.go play <wav文件> [设备ID]")
			os.Exit(1)
		}
		deviceID := ""
		if len(os.Args) > 3 {
			deviceID = os.Args[3]
		}
		playFile(&engine, os.Args[2], deviceID)
	default:
		fmt.Fprintf(os.Stderr, "未知命令: %s\n", os.Args[1])
		fmt.Fprintln(os.Stderr, "可用命令: list, play")
		os.Exit(1)
	}
}

func listDevices(engine *audio.AudioPlayerEngine) {
	devices, err := engine.ListDevices()
	if err != nil {
		fmt.Fprintf(os.Stderr, "枚举设备失败: %v\n", err)
		os.Exit(1)
	}
	if len(devices) == 0 {
		fmt.Println("未找到任何播放设备")
		return
	}
	fmt.Println("音频播放设备:")
	for _, d := range devices {
		defaultMark := ""
		if d.IsDefault {
			defaultMark = " [默认]"
		}
		fmt.Printf("  %s  %s%s\n", d.ID, d.Name, defaultMark)
	}
}

func playFile(engine *audio.AudioPlayerEngine, wavPath, deviceID string) {
	name := wavPath
	if err := engine.Preload(name, wavPath); err != nil {
		fmt.Fprintf(os.Stderr, "预加载失败: %v\n", err)
		os.Exit(1)
	}

	player, err := engine.PlaySound(name, deviceID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "播放失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("正在播放,按 Ctrl+C 停止...")

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	select {
	case <-player.Done():
		fmt.Println("播放完毕")
	case <-sigCh:
		fmt.Println("播放已停止")
	}
}
