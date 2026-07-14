package bot

import (
	"encoding/json"
	"testing"
)

func TestExtractText(t *testing.T) {
	cases := []struct {
		json string
		want string
	}{
		{`[{"type":"text","data":{"text":"/ping"}}]`, "/ping"},
		{`[{"type":"text","data":{"text":"hello"}},{"type":"text","data":{"text":" world"}}]`, "hello world"},
		{`[{"type":"image","data":{"url":"x.png"}}]`, ""},
	}
	for _, c := range cases {
		got := ExtractText(json.RawMessage(c.json))
		if got != c.want {
			t.Errorf("ExtractText(%q) = %q, want %q", c.json, got, c.want)
		}
	}
}
