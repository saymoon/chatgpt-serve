package chatgpt

import (
	"encoding/json"
	"github.com/atotto/clipboard"
	"github.com/playwright-community/playwright-go"
	"github.com/sirupsen/logrus"
	"strings"
	"time"
)

func (h *Helper) SendMsg(page playwright.Page, inputText string) (string, string, error) {
	var conversationId string
	// 等待 //textarea
	textarea, err := page.WaitForSelector("//textarea", playwright.PageWaitForSelectorOptions{
		Timeout: playwright.Float(1000 * 60),
	})
	if err != nil {
		logrus.Errorf("Error while waiting for selector: %v", err)
		return "", "", err
	}

	// 输入
	//inputText := "你好，我是OpenAI生产的复读机，可以复读你说的话。"
	// 复制 inputText 粘贴到 textarea
	h.sendMsgLock.Lock()
	defer h.sendMsgLock.Unlock()
	err = clipboard.WriteAll(inputText)
	if err != nil {
		logrus.Errorf("Error while writing to clipboard: %v", err)
		return "", "", err
	}
	err = textarea.Focus()
	if err != nil {
		logrus.Errorf("Error while focusing textarea: %v", err)
		return "", "", err
	}
	err = textarea.Press(controlA())
	if err != nil {
		logrus.Errorf("Error while pressing control+a: %v", err)
		return "", "", err
	}
	err = textarea.Press(controlV())
	if err != nil {
		logrus.Errorf("Error while pressing control+v: %v", err)
		return "", "", err
	}
	go func() {
		for {
			time.Sleep(time.Second)
			// 回车
			// //textarea/../button/@disabled 是否有值 如果有说明此时不能回车
			selector, err := page.QuerySelector("//textarea/../button/@disabled")
			if err != nil {
				logrus.Errorf("Error while querying selector: %v", err)
				continue
			}
			if selector != nil {
				continue
			}
			time.Sleep(time.Second * 2)
			err = textarea.Press("Enter")
			if err != nil {
				logrus.Errorf("Error while pressing enter: %v", err)
			}
			break
		}
	}()
	// 等待响应
	response := page.WaitForResponse("https://chat.openai.com/backend-api/conversation", playwright.PageWaitForResponseOptions{
		Timeout: playwright.Float(1000 * 60),
	})
	// 解析 text/event-stream
	{
		text, _ := response.Text()
		// 换行符分割，去掉 data:
		lines := strings.Split(text, "\n")
		var finalLine *ConversationStreamResponseItem
		for _, line := range lines {
			if strings.HasPrefix(line, "data:") {
				line = strings.TrimPrefix(line, "data:")
				// 解析json
				var data = &ConversationStreamResponseItem{}
				err := json.Unmarshal([]byte(line), data)
				if err != nil {
					continue
				}
				finalLine = data
			}
		}
		if finalLine != nil {
			inputText = finalLine.Text()
			conversationId = finalLine.ConversationId
			logrus.Infof("AI: %s %s", finalLine.Text(), conversationId)
		} else {
			inputText = ""
		}
	}
	return inputText, conversationId, nil
}

func (h *Helper) SendContinue(page playwright.Page) (string, string, error) {
	var inputText, conversationId string
	h.sendMsgLock.Lock()
	defer h.sendMsgLock.Unlock()

	logrus.Infoln("Waiting for button...")
	continueBtn, err := page.WaitForSelector("//button[contains(., 'Continue generating')]", playwright.PageWaitForSelectorOptions{
		Timeout: playwright.Float(1000 * 4),
	})

	if err != nil {
		logrus.Infoln("Button not found...")
		return "", "", err
	}

	logrus.Infoln("Button found...")
	retries := 3
	for i := 0; i < retries; i++ {
		time.Sleep(time.Second)
		err = continueBtn.Click()

		if err != nil {
			logrus.Errorf("Error while tapping continue button: %v", err)
            if i < retries-1 {  // if it's not the last retry, print "Retrying click..."
				logrus.Infoln("Retrying click...")
			} else {
                return inputText, conversationId, err
            }
		} else {
			logrus.Infoln("Button clicked successfully.")
            break
		}
	}
	// 等待响应
	response := page.WaitForResponse("https://chat.openai.com/backend-api/conversation", playwright.PageWaitForResponseOptions{
		Timeout: playwright.Float(1000 * 20),
	})
	// 解析 text/event-stream
	{
		text, _ := response.Text()
		// 换行符分割，去掉 data:
		lines := strings.Split(text, "\n")
		var finalLine *ConversationStreamResponseItem
		for _, line := range lines {
			if strings.HasPrefix(line, "data:") {
				line = strings.TrimPrefix(line, "data:")
				// 解析json
				var data = &ConversationStreamResponseItem{}
				err := json.Unmarshal([]byte(line), data)
				if err != nil {
					continue
				}
				finalLine = data
			}
		}
		if finalLine != nil {
			inputText = finalLine.Text()
			conversationId = finalLine.ConversationId
			logrus.Infof("AI: %s %s", finalLine.Text(), conversationId)
		} else {
			inputText = ""
		}
	}
	return inputText, conversationId, nil
}

func (h *Helper) SendMsgWithAutoContinue(page playwright.Page, inputText string) (string, string, error) {
	var replyCombined, conversationId string
	reply, conversationId, err := h.SendMsg(page, inputText)
	if err != nil {
		return replyCombined, conversationId, err
	}
	replyCombined = reply
    logrus.Infoln("aaaaaaaaaaaaa")
	for i := 0; i < 3; i++ {
		reply, _, err = h.SendContinue(page)
		if err == nil {
			replyCombined += reply
		} else {
            break
        }
	}
	return replyCombined, conversationId, nil
}

