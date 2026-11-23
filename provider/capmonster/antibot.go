package capmonster

import (
	"github.com/aarock1234/unicap/v3/internal/solverapi"
	"github.com/aarock1234/unicap/v3/tasks"
)

// tencentTask is CapMonster Cloud's Tencent task, a CustomTask of class
// "TenDI". The website key is the captcha's app ID.
type tencentTask struct {
	Type       string           `json:"type"`
	Class      string           `json:"class"`
	WebsiteURL string           `json:"websiteURL"`
	WebsiteKey string           `json:"websiteKey"`
	Metadata   *tencentMetadata `json:"metadata,omitempty"`
	solverapi.ProxyFields
}

// tencentMetadata is the metadata object of a Tencent task.
type tencentMetadata struct {
	CaptchaURL string `json:"captchaUrl"`
}

// huntTask is CapMonster Cloud's Hunt task, a CustomTask of class "HUNT",
// which is solved only through the caller's proxy.
type huntTask struct {
	Type       string       `json:"type"`
	Class      string       `json:"class"`
	WebsiteURL string       `json:"websiteURL"`
	Metadata   huntMetadata `json:"metadata"`
	UserAgent  string       `json:"userAgent,omitempty"`
	solverapi.ProxyFields
}

// huntMetadata is the metadata object of a Hunt task. Data is the page's
// meta.token value.
type huntMetadata struct {
	APIGetLib string `json:"apiGetLib"`
	Data      string `json:"data,omitempty"`
}

// tspdTask is CapMonster Cloud's F5 TSPD task, a CustomTask of class "tspd",
// which is solved only through the caller's proxy.
type tspdTask struct {
	Type       string       `json:"type"`
	Class      string       `json:"class"`
	WebsiteURL string       `json:"websiteURL"`
	Metadata   tspdMetadata `json:"metadata"`
	UserAgent  string       `json:"userAgent,omitempty"`
	solverapi.ProxyFields
}

// tspdMetadata is the metadata object of a TSPD task.
type tspdMetadata struct {
	TSPDCookie     string `json:"tspdCookie"`
	HTMLPageBase64 string `json:"htmlPageBase64"`
}

// basiliskTask is CapMonster Cloud's Basilisk task, a CustomTask of class
// "Basilisk". It has no metadata object.
type basiliskTask struct {
	Type       string `json:"type"`
	Class      string `json:"class"`
	WebsiteURL string `json:"websiteURL"`
	WebsiteKey string `json:"websiteKey"`
	UserAgent  string `json:"userAgent,omitempty"`
	solverapi.ProxyFields
}

// impervaTask is CapMonster Cloud's Imperva (Incapsula) task, a CustomTask
// of class "Imperva", which is solved only through the caller's proxy.
type impervaTask struct {
	Type       string          `json:"type"`
	Class      string          `json:"class"`
	WebsiteURL string          `json:"websiteURL"`
	Metadata   impervaMetadata `json:"metadata"`
	UserAgent  string          `json:"userAgent,omitempty"`
	solverapi.ProxyFields
}

// impervaMetadata is the metadata object of an Imperva task.
type impervaMetadata struct {
	IncapsulaScriptURL string `json:"incapsulaScriptUrl"`
	IncapsulaCookies   string `json:"incapsulaCookies"`
	Reese84URLEndpoint string `json:"reese84UrlEndpoint,omitempty"`
}

// alibabaTask is CapMonster Cloud's Alibaba Cloud captcha task, a CustomTask
// of class "alibaba".
type alibabaTask struct {
	Type       string          `json:"type"`
	Class      string          `json:"class"`
	WebsiteURL string          `json:"websiteURL"`
	Metadata   alibabaMetadata `json:"metadata"`
	UserAgent  string          `json:"userAgent,omitempty"`
	solverapi.ProxyFields
}

// alibabaMetadata is the metadata object of an Alibaba task. CapMonster
// Cloud capitalises UserCertifyId, unlike the other keys.
type alibabaMetadata struct {
	SceneID       string `json:"sceneId"`
	Prefix        string `json:"prefix"`
	UserID        string `json:"userId,omitempty"`
	UserUserID    string `json:"userUserId,omitempty"`
	VerifyType    string `json:"verifyType,omitempty"`
	Region        string `json:"region,omitempty"`
	UserCertifyID string `json:"UserCertifyId,omitempty"`
	APIGetLib     string `json:"apiGetLib,omitempty"`
}

// customTaskType is the wire type of the CustomTask family, whose class
// field selects the captcha.
const customTaskType = "CustomTask"

// mapTencent builds CapMonster Cloud's Tencent task. It sends metadata only
// when the task names a custom captcha script.
func mapTencent(task *tasks.TencentTask) (tencentTask, error) {
	proxy, err := proxyFields(task.Proxy)
	if err != nil {
		return tencentTask{}, err
	}

	wire := tencentTask{
		Type:        customTaskType,
		Class:       "TenDI",
		WebsiteURL:  task.WebsiteURL,
		WebsiteKey:  task.AppID,
		ProxyFields: proxy,
	}

	if task.CaptchaScript != "" {
		wire.Metadata = &tencentMetadata{
			CaptchaURL: task.CaptchaScript,
		}
	}

	return wire, nil
}

// mapHunt builds CapMonster Cloud's Hunt task.
func mapHunt(task *tasks.HuntTask) (huntTask, error) {
	proxy, err := requiredProxyFields(task.Proxy)
	if err != nil {
		return huntTask{}, err
	}

	return huntTask{
		Type:       customTaskType,
		Class:      "HUNT",
		WebsiteURL: task.WebsiteURL,
		Metadata: huntMetadata{
			APIGetLib: task.APIGetLib,
			Data:      task.Data,
		},
		UserAgent:   task.UserAgent,
		ProxyFields: proxy,
	}, nil
}

// mapTSPD builds CapMonster Cloud's F5 TSPD task.
func mapTSPD(task *tasks.TSPDTask) (tspdTask, error) {
	proxy, err := requiredProxyFields(task.Proxy)
	if err != nil {
		return tspdTask{}, err
	}

	return tspdTask{
		Type:       customTaskType,
		Class:      "tspd",
		WebsiteURL: task.WebsiteURL,
		Metadata: tspdMetadata{
			TSPDCookie:     task.TSPDCookie,
			HTMLPageBase64: task.HTMLPageBase64,
		},
		UserAgent:   task.UserAgent,
		ProxyFields: proxy,
	}, nil
}

// mapBasilisk builds CapMonster Cloud's Basilisk task.
func mapBasilisk(task *tasks.BasiliskTask) (basiliskTask, error) {
	proxy, err := proxyFields(task.Proxy)
	if err != nil {
		return basiliskTask{}, err
	}

	return basiliskTask{
		Type:        customTaskType,
		Class:       "Basilisk",
		WebsiteURL:  task.WebsiteURL,
		WebsiteKey:  task.WebsiteKey,
		UserAgent:   task.UserAgent,
		ProxyFields: proxy,
	}, nil
}

// mapImperva builds CapMonster Cloud's Imperva (Incapsula) task.
func mapImperva(task *tasks.ImpervaTask) (impervaTask, error) {
	proxy, err := requiredProxyFields(task.Proxy)
	if err != nil {
		return impervaTask{}, err
	}

	return impervaTask{
		Type:       customTaskType,
		Class:      "Imperva",
		WebsiteURL: task.WebsiteURL,
		Metadata: impervaMetadata{
			IncapsulaScriptURL: task.ScriptURL,
			IncapsulaCookies:   task.Cookies,
			Reese84URLEndpoint: task.Reese84URLEndpoint,
		},
		UserAgent:   task.UserAgent,
		ProxyFields: proxy,
	}, nil
}

// mapAlibaba builds CapMonster Cloud's Alibaba Cloud captcha task.
func mapAlibaba(task *tasks.AlibabaTask) (alibabaTask, error) {
	proxy, err := proxyFields(task.Proxy)
	if err != nil {
		return alibabaTask{}, err
	}

	return alibabaTask{
		Type:       customTaskType,
		Class:      "alibaba",
		WebsiteURL: task.WebsiteURL,
		Metadata: alibabaMetadata{
			SceneID:       task.SceneID,
			Prefix:        task.Prefix,
			UserID:        task.UserID,
			UserUserID:    task.UserUserID,
			VerifyType:    task.VerifyType,
			Region:        task.Region,
			UserCertifyID: task.UserCertifyID,
			APIGetLib:     task.APIGetLib,
		},
		UserAgent:   task.UserAgent,
		ProxyFields: proxy,
	}, nil
}
