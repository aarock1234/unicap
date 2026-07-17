package solverapi

import (
	"encoding/json"
	"testing"
)

func TestTaskIDUnmarshalJSON(t *testing.T) {
	tests := []struct {
		name string
		json string
		want string
	}{
		{
			name: "string id",
			json: `{"taskId":"61138bb6-19fb-11ec"}`,
			want: "61138bb6-19fb-11ec",
		},
		{
			name: "numeric id",
			json: `{"taskId":123456}`,
			want: "123456",
		},
		{
			name: "missing id",
			json: `{}`,
			want: "",
		},
		{
			name: "null id",
			json: `{"taskId":null}`,
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var resp struct {
				TaskID TaskID `json:"taskId"`
			}
			if err := json.Unmarshal([]byte(tt.json), &resp); err != nil {
				t.Fatalf("unmarshal: %v", err)
			}

			if got := resp.TaskID.String(); got != tt.want {
				t.Errorf("TaskID = %q, want %q", got, tt.want)
			}
		})
	}
}
