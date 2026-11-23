package nopecha

import (
	"github.com/aarock1234/unicap/v3/internal/solverapi"
	"github.com/aarock1234/unicap/v3/tasks"
)

// rqdataKey is the only tasks.HCaptchaTask.EnterprisePayload key NopeCHA has
// a field for.
const rqdataKey = "rqdata"

// sizeInvisible is the widget size of an invisible hCaptcha.
const sizeInvisible = "invisible"

// hCaptchaRequest is the body NopeCHA's hCaptcha token endpoint takes.
type hCaptchaRequest struct {
	Sitekey    string       `json:"sitekey"`
	URL        string       `json:"url"`
	Enterprise bool         `json:"enterprise,omitempty"`
	Data       hCaptchaData `json:"data,omitzero"`
	Proxy      *proxy       `json:"proxy,omitempty"`
	UserAgent  string       `json:"useragent,omitempty"`
}

// hCaptchaData holds the widget parameters of an hCaptcha job. An empty size
// leaves NopeCHA's default, the normal widget.
type hCaptchaData struct {
	RQData string `json:"rqdata,omitempty"`
	Size   string `json:"size,omitempty"`
}

// mapHCaptcha builds the body of NopeCHA's hCaptcha token endpoint. unicap
// has no enterprise flag, so a challenge that comes with rqdata is taken to
// be hCaptcha Enterprise. NopeCHA's cookie objects need attributes that the
// task's cookie string does not carry, so cookies are rejected.
func mapHCaptcha(task *tasks.HCaptchaTask) (hCaptchaRequest, error) {
	if err := solverapi.Check(
		solverapi.Unsupported(
			solverapi.Setting("Cookies", task.Cookies != ""),
		),
		checkProxy(task.Proxy),
	); err != nil {
		return hCaptchaRequest{}, err
	}

	rqdata, err := enterpriseValue(task.EnterprisePayload, rqdataKey)
	if err != nil {
		return hCaptchaRequest{}, err
	}

	return hCaptchaRequest{
		Sitekey:    task.WebsiteKey,
		URL:        task.WebsiteURL,
		Enterprise: rqdata != "",
		Data: hCaptchaData{
			RQData: rqdata,
			Size:   widgetSize(task.IsInvisible),
		},
		Proxy:     newProxy(task.Proxy),
		UserAgent: task.UserAgent,
	}, nil
}

// widgetSize names the size of an invisible widget, and leaves the size of
// any other widget unset.
func widgetSize(isInvisible bool) string {
	if isInvisible {
		return sizeInvisible
	}

	return ""
}
