package twocaptcha

import (
	"github.com/aarock1234/unicap/v3/internal/solverapi"
	"github.com/aarock1234/unicap/v3/tasks"
)

type vkCaptchaTask struct {
	Type        string `json:"type"`
	RedirectURI string `json:"redirectUri"`
	UserAgent   string `json:"userAgent"`
	solverapi.ProxyFields
}

type yidunTask struct {
	Type                    string `json:"type"`
	WebsiteURL              string `json:"websiteURL"`
	WebsiteKey              string `json:"websiteKey"`
	UserAgent               string `json:"userAgent,omitempty"`
	YidunGetLib             string `json:"yidunGetLib,omitempty"`
	YidunAPIServerSubdomain string `json:"yidunApiServerSubdomain,omitempty"`
	Challenge               string `json:"challenge,omitempty"`
	HCG                     string `json:"hcg,omitempty"`
	HCT                     int64  `json:"hct,omitempty"`
	solverapi.ProxyFields
}

type binanceTask struct {
	Type       string `json:"type"`
	WebsiteURL string `json:"websiteURL"`
	WebsiteKey string `json:"websiteKey"`
	ValidateID string `json:"validateId"`
	UserAgent  string `json:"userAgent,omitempty"`
	solverapi.ProxyFields
}

type huntTask struct {
	Type       string `json:"type"`
	WebsiteURL string `json:"websiteURL"`
	APIGetLib  string `json:"apiGetLib"`
	UserAgent  string `json:"userAgent,omitempty"`
	Data       string `json:"data,omitempty"`
	solverapi.ProxyFields
}

type tspdTask struct {
	Type           string `json:"type"`
	WebsiteURL     string `json:"websiteURL"`
	TSPDCookie     string `json:"tspdCookie"`
	HTMLPageBase64 string `json:"htmlPageBase64"`
	UserAgent      string `json:"userAgent,omitempty"`
	solverapi.ProxyFields
}

type incapsulaTask struct {
	Type               string `json:"type"`
	WebsiteURL         string `json:"websiteURL"`
	IncapsulaScriptURL string `json:"incapsulaScriptUrl"`
	IncapsulaCookies   string `json:"incapsulaCookies"`
	UserAgent          string `json:"userAgent,omitempty"`
	Reese84URLEndpoint string `json:"reese84UrlEndpoint,omitempty"`
	solverapi.ProxyFields
}

type alibabaTask struct {
	Type          string `json:"type"`
	WebsiteURL    string `json:"websiteURL"`
	SceneID       string `json:"sceneId"`
	Prefix        string `json:"prefix"`
	UserID        string `json:"userId,omitempty"`
	UserUserID    string `json:"userUserId,omitempty"`
	VerifyType    string `json:"verifyType,omitempty"`
	Region        string `json:"region,omitempty"`
	UserCertifyID string `json:"userCertifyId,omitempty"`
	APIGetLib     string `json:"apiGetLib,omitempty"`
	UserAgent     string `json:"userAgent,omitempty"`
	solverapi.ProxyFields
}

// binanceVariant spells the proxyless type exactly as 2Captcha documents it.
var binanceVariant = solverapi.Variant{
	Proxyless: "BinanceTaskproxyless",
	Proxy:     "BinanceTask",
}

func mapVKCaptcha(task *tasks.VKCaptchaTask) (vkCaptchaTask, error) {
	taskType, err := proxyOnly("VKCaptchaTask").For(task.Proxy)
	if err != nil {
		return vkCaptchaTask{}, err
	}

	return vkCaptchaTask{
		Type:        taskType,
		RedirectURI: task.RedirectURI,
		UserAgent:   task.UserAgent,
		ProxyFields: solverapi.NewProxyFields(task.Proxy),
	}, nil
}

func mapYidun(task *tasks.YidunTask) (yidunTask, error) {
	taskType, err := variant("YidunTask").For(task.Proxy)
	if err != nil {
		return yidunTask{}, err
	}

	return yidunTask{
		Type:                    taskType,
		WebsiteURL:              task.WebsiteURL,
		WebsiteKey:              task.WebsiteKey,
		UserAgent:               task.UserAgent,
		YidunGetLib:             task.GetLibURL,
		YidunAPIServerSubdomain: task.APIServerSubdomain,
		Challenge:               task.Challenge,
		HCG:                     task.HCG,
		HCT:                     task.HCT,
		ProxyFields:             solverapi.NewProxyFields(task.Proxy),
	}, nil
}

func mapBinance(task *tasks.BinanceTask) (binanceTask, error) {
	taskType, err := binanceVariant.For(task.Proxy)
	if err != nil {
		return binanceTask{}, err
	}

	return binanceTask{
		Type:        taskType,
		WebsiteURL:  task.WebsiteURL,
		WebsiteKey:  task.WebsiteKey,
		ValidateID:  task.ValidateID,
		UserAgent:   task.UserAgent,
		ProxyFields: solverapi.NewProxyFields(task.Proxy),
	}, nil
}

func mapHunt(task *tasks.HuntTask) (huntTask, error) {
	taskType, err := proxyOnly("HuntTask").For(task.Proxy)
	if err != nil {
		return huntTask{}, err
	}

	return huntTask{
		Type:        taskType,
		WebsiteURL:  task.WebsiteURL,
		APIGetLib:   task.APIGetLib,
		UserAgent:   task.UserAgent,
		Data:        task.Data,
		ProxyFields: solverapi.NewProxyFields(task.Proxy),
	}, nil
}

func mapTSPD(task *tasks.TSPDTask) (tspdTask, error) {
	taskType, err := proxyOnly("TspdTask").For(task.Proxy)
	if err != nil {
		return tspdTask{}, err
	}

	return tspdTask{
		Type:           taskType,
		WebsiteURL:     task.WebsiteURL,
		TSPDCookie:     task.TSPDCookie,
		HTMLPageBase64: task.HTMLPageBase64,
		UserAgent:      task.UserAgent,
		ProxyFields:    solverapi.NewProxyFields(task.Proxy),
	}, nil
}

func mapBasilisk(task *tasks.BasiliskTask) (siteKeyTask, error) {
	taskType, err := variant("BasiliskTask").For(task.Proxy)
	if err != nil {
		return siteKeyTask{}, err
	}

	return siteKeyTask{
		Type:        taskType,
		WebsiteURL:  task.WebsiteURL,
		WebsiteKey:  task.WebsiteKey,
		UserAgent:   task.UserAgent,
		ProxyFields: solverapi.NewProxyFields(task.Proxy),
	}, nil
}

func mapImperva(task *tasks.ImpervaTask) (incapsulaTask, error) {
	taskType, err := proxyOnly("IncapsulaTask").For(task.Proxy)
	if err != nil {
		return incapsulaTask{}, err
	}

	return incapsulaTask{
		Type:               taskType,
		WebsiteURL:         task.WebsiteURL,
		IncapsulaScriptURL: task.ScriptURL,
		IncapsulaCookies:   task.Cookies,
		UserAgent:          task.UserAgent,
		Reese84URLEndpoint: task.Reese84URLEndpoint,
		ProxyFields:        solverapi.NewProxyFields(task.Proxy),
	}, nil
}

func mapAlibaba(task *tasks.AlibabaTask) (alibabaTask, error) {
	taskType, err := variant("AlibabaTask").For(task.Proxy)
	if err != nil {
		return alibabaTask{}, err
	}

	return alibabaTask{
		Type:          taskType,
		WebsiteURL:    task.WebsiteURL,
		SceneID:       task.SceneID,
		Prefix:        task.Prefix,
		UserID:        task.UserID,
		UserUserID:    task.UserUserID,
		VerifyType:    task.VerifyType,
		Region:        task.Region,
		UserCertifyID: task.UserCertifyID,
		APIGetLib:     task.APIGetLib,
		UserAgent:     task.UserAgent,
		ProxyFields:   solverapi.NewProxyFields(task.Proxy),
	}, nil
}
