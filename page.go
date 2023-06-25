package chatgpt

import (
	"fmt"
	"log"

	"github.com/cherish-chat/chatgpt-firefox/config"
	"github.com/playwright-community/playwright-go"
	"github.com/sirupsen/logrus"
)

// Deprecated 请使用GetPage
func (h *Helper) MustGetPage(id string) playwright.Page {
	if page, ok := h.pageMap.Load(id); ok {
		return page.(playwright.Page)
	}
	var page playwright.Page
	var err error
	page, err = h.NewPage(id)
	if err != nil {
		logrus.Errorf("Error while creating new page: %v", err)
		log.Fatalf("Error while creating new page: %v", err)
	}
	return page
}

//	func (h *Helper) GetPage(id string) (playwright.Page, error) {
//		if page, ok := h.pageMap.Load(id); ok {
//			return page.(playwright.Page), nil
//		}
//		return h.NewPage(id)
//	}
func (h *Helper) GetPageByModel(id, conversationId, model string) (pg playwright.Page, err error) {
	if page, ok := h.pageMap.Load(id); ok {
		pg = page.(playwright.Page)
		if _, err = pg.QuerySelector("//a[contains(., 'New chat')]"); err != nil {
			logrus.Info("page closed, reopen it", err.Error())
			pg.Close()
			pg, err = h.NewPageByModel(id, model)
		}
	} else {
		pg, err = h.NewPageByModel(id, model)
	}
	if err != nil {
		return nil, err
	}

	err = loadConversationByModel(pg, conversationId, model)

	return pg, err
}

func (h *Helper) NewPageByModel(id, model string) (playwright.Page, error) {
	page, err := h.browser.NewPage()
	if err != nil {
		logrus.Errorf("Error while creating new page: %v", err)
		return nil, err
	}
	logrus.Info("New page created successfully")
	{
		// 设置cookie
		cookies, err := config.LoadCookies("cookies.json")
		if err != nil {
			return nil, err
		}
		err = h.browser.AddCookies(cookies...)
		if err != nil {
			logrus.Errorf("Error while adding cookies: %v", err)
			return nil, err
		}
		logrus.Info("Cookies added successfully")
	}
	targetUrl := fmt.Sprintf("https://chat.openai.com/chat?model=%s", model)
	if err = gotoTargetUrl(page, targetUrl); err != nil {
		return nil, err
	}

	logrus.Info("Navigated to openai successfully, target URL: " + targetUrl)
	h.pageMap.Store(id, page)
	closeTips(page)
	return page, nil
}

func loadConversationByModel(page playwright.Page, conversationId, model string) error {
	var initUrl, targetUrl string

	initUrl = fmt.Sprintf("https://chat.openai.com/chat?model=%s", model)
	// if isPlus {
	// 	initUrl = "https://chat.openai.com/chat?model=gpt-4"
	// } else {
	// 	initUrl = "https://chat.openai.com/?model=text-davinci-002-render-sha"
	// }

	if conversationId != "" {
		targetUrl = "https://chat.openai.com/c/" + conversationId
	} else {
		targetUrl = initUrl
	}
	// 如果当前页面已经是需要打开对话直接返回
	if page.URL() == targetUrl {
		return nil
	}
	// 否则尝试跳转
	if err := gotoTargetUrl(page, targetUrl); err != nil {
		return err
	}
	// 如果跳转对话失败，说明对应的 ConversationId 错误或对话以被删除，自动开启新对话
	if page.URL() == "https://chat.openai.com/" {
		targetUrl = initUrl
		if err := gotoTargetUrl(page, targetUrl); err != nil {
			return err
		}
	}
	return nil
}

func (h *Helper) GetPage(id, conversationId string, isPlus bool) (pg playwright.Page, err error) {
	if page, ok := h.pageMap.Load(id); ok {
		pg = page.(playwright.Page)
		if _, err = pg.QuerySelector("//a[contains(., 'New chat')]"); err != nil {
			logrus.Info("page closed, reopen it", err.Error())
			pg.Close()
			pg, err = h.NewPlusPage(id)
		}
	} else if isPlus {
		pg, err = h.NewPlusPage(id)
	} else {
		pg, err = h.NewPage(id)
	}
	if err != nil {
		return nil, err
	}

	err = loadConversation(pg, conversationId, isPlus)

	return pg, err
}

func (h *Helper) OpenBaiDu() {
    logrus.Infoln("wwwwwwww")
	page, err := h.browser.NewPage()
	if err != nil {
		logrus.Errorf("Error while creating new page: %v", err)
	}
	if _, err = page.Goto("https://www.baidu.com"); err != nil {
		log.Fatalf("could not goto: %v", err)
	}
}

func (h *Helper) NewPage(id string) (playwright.Page, error) {
	page, err := h.browser.NewPage()
	if err != nil {
		logrus.Errorf("Error while creating new page: %v", err)
		return nil, err
	}
	logrus.Info("New page created successfully")
	{
		// 设置cookie
		cookies, err := config.LoadCookies("cookies.json")
		if err != nil {
			return nil, err
		}
		err = h.browser.AddCookies(cookies...)
		if err != nil {
			logrus.Errorf("Error while adding cookies: %v", err)
			return nil, err
		}
		logrus.Info("Cookies added successfully")
	}
	targetUrl := "https://chat.openai.com/?model=text-davinci-002-render-sha"
	if err = gotoTargetUrl(page, targetUrl); err != nil {
		return nil, err
	}

	logrus.Info("Navigated to openai successfully, target URL: " + targetUrl)
	h.pageMap.Store(id, page)
	closeTips(page)
	return page, nil
}

func (h *Helper) NewPlusPage(id string) (playwright.Page, error) {
	page, err := h.browser.NewPage()
	if err != nil {
		logrus.Errorf("Error while creating new page: %v", err)
		return nil, err
	}
	logrus.Info("New page created successfully")
	{
		// 设置cookie
		cookies, err := config.LoadCookies("cookies.json")
		if err != nil {
			return nil, err
		}
		err = h.browser.AddCookies(cookies...)
		if err != nil {
			logrus.Errorf("Error while adding cookies: %v", err)
			return nil, err
		}
		logrus.Info("Cookies added successfully")
	}
	targetUrl := "https://chat.openai.com/chat?model=gpt-4"
	if err = gotoTargetUrl(page, targetUrl); err != nil {
		return nil, err
	}

	logrus.Info("Navigated to openai successfully, target URL: " + targetUrl)
	h.pageMap.Store(id, page)
	closeTips(page)
	return page, nil
}

func loadConversation(page playwright.Page, conversationId string, isPlus bool) error {
	var initUrl, targetUrl string

	if isPlus {
		initUrl = "https://chat.openai.com/chat?model=gpt-4"
	} else {
		initUrl = "https://chat.openai.com/?model=text-davinci-002-render-sha"
	}

	if conversationId != "" {
		targetUrl = "https://chat.openai.com/c/" + conversationId
	} else {
		targetUrl = initUrl
	}
	// 如果当前页面已经是需要打开对话直接返回
	if page.URL() == targetUrl {
		return nil
	}
	// 否则尝试跳转
	if err := gotoTargetUrl(page, targetUrl); err != nil {
		return err
	}
	// 如果跳转对话失败，说明对应的 ConversationId 错误或对话以被删除，自动开启新对话
	if page.URL() == "https://chat.openai.com/" {
		targetUrl = initUrl
		if err := gotoTargetUrl(page, targetUrl); err != nil {
			return err
		}
	}
	return nil
}

func gotoTargetUrl(page playwright.Page, url string) error {
	_, err := page.Goto(url)
	if err != nil {
		page.Close()
		logrus.Errorf("Error while navigating to chatgpt: %v", err)
		return err
	}
	return nil
}

func closeTips(page playwright.Page) error {
	// 检查 Next 按钮是否存在
	for {
		nextButton, err := page.WaitForSelector("//button[contains(., 'Next') or contains(., 'Done')]", playwright.PageWaitForSelectorOptions{
			Timeout: playwright.Float(1000 * 3),
		})
		if err != nil {
			return err
		}
		if nextButton == nil {
			// 如果 Next 按钮不存在，退出循环
			logrus.Infof("No 'Next' Button")
			break
		}

		// 如果 Next 按钮存在，点击它
		logrus.Infof("Found 'Next' Button")
		nextButton.Click()
	}
	return nil
}

