package handler

import (
	"fmt"
	"sort"
	"strings"
	"vid2txt/config"
	"vid2txt/internal/service"
	"vid2txt/pkg/cli"
)

func ModifyLanguage() {
	currentLanguage := config.AppConfig.Settings.Language
	supportedLanguages := map[string]string{
		"cn": "中文",
		"en": "英语",
		"ja": "日语",
		"ko": "韩语",
	}

		// 提取所有语言代码并排序
	languageCodes := make([]string, 0, len(supportedLanguages))
	for code := range supportedLanguages {
		languageCodes = append(languageCodes, code)
	}
	sort.Strings(languageCodes)

	fmt.Printf("当前的语言: %s\n", currentLanguage)

	// 按排序后的顺序显示支持的语言
	fmt.Println("支持设置:")
	for _, code := range languageCodes {
		name := supportedLanguages[code]
		fmt.Printf("  %s（%s）", name, code)
		if code == currentLanguage {
			fmt.Print(" ✓")
		}
		fmt.Println()
	}

	// 获取用户输入
	newLanguage := strings.ToLower(cli.GetUserInput("请输入新的语言代码: "))

	// 检查输入是否有效并更新
	if _, ok := supportedLanguages[newLanguage]; ok {
		config.ModifyConfig("settings.language", newLanguage)
		// 同时更新内存中的配置
		config.AppConfig.Settings.Language = newLanguage
		fmt.Printf("✅ 语言已成功修改为: %s\n", supportedLanguages[newLanguage])
	} else {
		fmt.Println("❌ 输入的语言不受支持，请重新输入!")
	}
}

func RunCLI() {
	for {
		fmt.Println(config.PromptMenu)

		choice := cli.GetUserInput("💡 请输入你的选择: ")

		switch choice {
		case "1":
			// 将视频转换成文本
			videoPath := cli.GetUserInput(config.ColorCyan+ "🎬 请输入视频文件路径: " + config.ColorReset)
			err := service.StartVideoToTextTask(videoPath)
			if err != nil {
				fmt.Println("❌", err)
				continue
			}
		case "2":
			// 将音频转换成文本
			audioPath := cli.GetUserInput(config.ColorCyan+ "🎵 请输入音频文件路径: " + config.ColorReset)
			err := service.StartAudioToTextTask(audioPath)
			if err != nil {
				fmt.Println("❌", err)
				continue
			}
		case "3":
			// 获取 taskinfo
			taskInfo, err := service.QueryTaskInfo(config.TaskResponseFilePath())
			if err != nil {
				fmt.Println("❌", err)
				continue
			}

			// 判断任务状态  todo
			taskStatus, err := service.QueryTaskStatus(taskInfo)
			if err != nil {
				fmt.Println("❌", err)
				continue
			}

			if taskStatus == config.TaskStatusOngoing {
				fmt.Println("⏳ 任务正在进行中，请稍后再试。")
				continue
			}
		
			if taskStatus == config.TaskStatusFailed {
				fmt.Println("❌ 任务处理失败，请检查错误日志。")
				continue
			}
		
			if taskStatus != config.TaskStatusCompleted {
				fmt.Println("⚠️ 未知任务状态:", taskStatus)
				continue
			}
			
			// 提示用户是否保存结果 
			input := cli.GetUserInput("✅ 转录已完成，是否保存转录结果? (yes/no): ")
	
			if !cli.ShouldProceed(input) {
				fmt.Println("🙅‍♂️ 用户取消操作。")
				continue
			}

			err = service.SaveTranscription([]byte(taskInfo), config.ResultOutputFilePath())
			if err != nil {
				fmt.Println("❌", err)
				continue
			}
			fmt.Println("🙆‍♂️ 转录结果已保存！")
		case "4":
			// 修改语言配置
			ModifyLanguage()
		default:
			fmt.Println("❌ 无效输入，请输入 1、2、3、4、q")
		}
	}
}
