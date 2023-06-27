package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/cherish-chat/chatgpt-firefox"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type Task struct {
	ID             string    `json:"id"`
	InstanceId     string    `json:"instance_id"`
	ConversationId string    `json:"conversation_id"`
	Model          string    `json:"model"`
	Prompt         string    `json:"prompt"`
	Response       string    `json:"response"`
	Status         string    `json:"status"`
	ErrorMessage   string    `json:"error_message"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type Update struct {
	ID             string `json:"id"`
	ConversationId string `json:"conversation_id"`
	Status         string `json:"status"`
	Response       string `json:"response"`
	ErrorMessage   string `json:"error_message"`
}

type Config struct {
	API struct {
		Host          string `yaml:"host"`
		Authorization string `yaml:"authorization"`
	} `yaml:"api"`
}

var config Config
var helper *chatgpt.Helper

func main() {
	// 读取配置文件
	readConfig()

	// 初始化 helper
	helper = chatgpt.NewHelper("cookies.json")
	err := helper.LaunchBrowser()
	if err != nil {
		logrus.Errorf("launch browser error: %v", err)
	}

	// 启动 daemon 服务
	for {
		task, err := getPendingTask()
		if err != nil {
			log.Printf("get task error: %v", err)
			time.Sleep(1 * time.Second)
			continue
		}
		logrus.Infoln("task :", task)

		update, err := processTask(task)
		if err != nil {
			log.Printf("process task error: %v", err)
			continue
		}

		err = updateTask(update)
		if err != nil {
			log.Printf("update task error: %v", err)
			continue
		}
	}
}

func readConfig() {
	// 读取配置文件
	configFile, err := os.Open("config.yaml")
	if err != nil {
		log.Fatal(err)
	}
	defer configFile.Close()

	decoder := yaml.NewDecoder(configFile)
	err = decoder.Decode(&config)
	if err != nil {
		log.Fatal(err)
	}
}

func getPendingTask() (Task, error) {
	// 获取任务
	url := config.API.Host + "/tasks/get-pending"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return Task{}, err
	}
	req.Header.Set("Authorization", config.API.Authorization)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return Task{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == 404 {
		// 如果返回状态码是 404，返回一个特定的错误，表示没有任务
		return Task{}, errors.New("No task available")
	}

	// 解析返回的任务信息
	task := Task{}
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &task)

	return task, nil
}

func processTask(task Task) (Update, error) {
	var conversationId, reply string
	// page, err := helper.GetPageByModel(task.InstanceId, task.ConversationId, task.Model)
	page, err := helper.GetPage(task.InstanceId, task.ConversationId, true)
	if err != nil {
		return Update{
			ID:           task.ID,
			Status:       "error",
			ErrorMessage: err.Error(),
		}, err
	} else {
		reply, conversationId, err = helper.SendMsgWithAutoContinue(page, task.Prompt)
		if err != nil {
			return Update{
				ID:           task.ID,
				Status:       "error",
				ErrorMessage: err.Error(),
			}, nil
		} else {
			if strings.TrimSpace(reply) == "" {
				log.Printf("openai response error: %v", err)
				// 刷新页面
				helper.ClosePage(task.ConversationId, page)
				return Update{
					ID:           task.ID,
					Status:       "error",
					ErrorMessage: "openai response error",
				}, nil
			}
		}
	}

	// 更新任务信息
	update := Update{
		ID:             task.ID,
		ConversationId: conversationId,
		Response:       reply,
		Status:         "completed",
		ErrorMessage:   "",
	}

	return update, nil
}

func updateTask(update Update) error {
	jsonUpdate, _ := json.Marshal(update)
	url := config.API.Host + "/tasks/update"
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonUpdate))
	req.Header.Set("Authorization", config.API.Authorization)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 可以检查服务器返回的响应
	body, _ := ioutil.ReadAll(resp.Body)
	log.Println(string(body))

	return nil
}
