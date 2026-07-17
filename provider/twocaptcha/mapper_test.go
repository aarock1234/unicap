package twocaptcha

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/aarock1234/unicap"
	"github.com/aarock1234/unicap/tasks"
)

func marshalTask(t *testing.T, task unicap.Task) map[string]any {
	t.Helper()

	mapped, err := mapTask(task)
	if err != nil {
		t.Fatalf("mapTask: %v", err)
	}

	data, err := json.Marshal(mapped)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var out map[string]any
	if err := json.Unmarshal(data, &out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	return out
}

func TestMapTaskType(t *testing.T) {
	proxy := &unicap.Proxy{Type: unicap.ProxyTypeHTTP, Address: "1.2.3.4", Port: 8080}

	tests := []struct {
		name     string
		task     unicap.Task
		wantType string
		wantKeys []string
	}{
		{
			name:     "recaptcha v2 proxyless",
			task:     &tasks.ReCaptchaV2Task{WebsiteURL: "u", WebsiteKey: "k"},
			wantType: "RecaptchaV2TaskProxyless",
		},
		{
			name:     "recaptcha v2 proxied",
			task:     &tasks.ReCaptchaV2Task{WebsiteURL: "u", WebsiteKey: "k", Proxy: proxy},
			wantType: "RecaptchaV2Task",
			wantKeys: []string{"proxyAddress", "proxyPort"},
		},
		{
			name:     "aws waf",
			task:     &tasks.AWSWAFTask{WebsiteURL: "u", Key: "k"},
			wantType: "AmazonTaskProxyless",
		},
		{
			name:     "mtcaptcha",
			task:     &tasks.MTCaptchaTask{WebsiteURL: "u", WebsiteKey: "k"},
			wantType: "MtCaptchaTaskProxyless",
		},
		{
			name:     "friendly captcha",
			task:     &tasks.FriendlyCaptchaTask{WebsiteURL: "u", WebsiteKey: "k"},
			wantType: "FriendlyCaptchaTaskProxyless",
		},
		{
			name:     "lemin",
			task:     &tasks.LeminTask{WebsiteURL: "u", CaptchaID: "c", DivID: "d"},
			wantType: "LeminTaskProxyless",
			wantKeys: []string{"captchaId", "divId"},
		},
		{
			name:     "cutcaptcha",
			task:     &tasks.CutCaptchaTask{WebsiteURL: "u", MiseryKey: "m", APIKey: "a"},
			wantType: "CutCaptchaTaskProxyless",
			wantKeys: []string{"miseryKey", "apiKey"},
		},
		{
			name:     "text",
			task:     &tasks.TextCaptchaTask{Question: "q"},
			wantType: "TextCaptchaTask",
			wantKeys: []string{"comment"},
		},
		{
			name:     "prosopo proxyless",
			task:     &tasks.ProsopoTask{WebsiteURL: "u", WebsiteKey: "k"},
			wantType: "ProsopoTaskProxyless",
			wantKeys: []string{"websiteKey"},
		},
		{
			name:     "prosopo proxied",
			task:     &tasks.ProsopoTask{WebsiteURL: "u", WebsiteKey: "k", Proxy: proxy},
			wantType: "ProsopoTask",
			wantKeys: []string{"proxyAddress", "proxyPort"},
		},
		{
			name:     "altcha with challenge url",
			task:     &tasks.AltchaTask{WebsiteURL: "u", ChallengeURL: "c"},
			wantType: "AltchaTaskProxyless",
			wantKeys: []string{"challengeURL"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := marshalTask(t, tt.task)

			if out["type"] != tt.wantType {
				t.Errorf("type = %v, want %v", out["type"], tt.wantType)
			}

			for _, key := range tt.wantKeys {
				if _, ok := out[key]; !ok {
					t.Errorf("missing key %q in %v", key, out)
				}
			}
		})
	}
}

func TestMapTaskUnsupported(t *testing.T) {
	_, err := mapTask(&tasks.CloudflareChallengeTask{WebsiteURL: "u"})
	if !errors.Is(err, unicap.ErrUnsupportedTask) {
		t.Fatalf("errors.Is(%v, ErrUnsupportedTask) = false, want true", err)
	}
}
