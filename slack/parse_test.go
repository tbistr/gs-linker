package sl

import "testing"

func Test_parseCommand(t *testing.T) {
	tests := map[string]struct {
		text       string
		wantC      command
		wantRawURL string
		wantErr    bool
	}{
		// normal scenario
		"subscribe":     {"<@bot_id> subscribe <https://github.com/tbistr/gs-linker/1>", subscribe, "https://github.com/tbistr/gs-linker/1", false},
		"unsubscribe":   {"<@bot_id> unsubscribe", unsubscribe, "", false},
		"summary":       {"<@bot_id> summary", summary, "", false},
		"shuffle case":  {"<@bot_id> SuBscRibE <url>", subscribe, "url", false},
		"insert spases": {"  <@bot_id> \t\t\n  subscribe  \n\n <url>   ", subscribe, "url", false},

		// exeption scenario
		"many(sub)":     {"<@bot_id> subscribe arg arg", "", "", true},
		"many(unsub)":   {"<@bot_id> unsubscribe arg", "", "", true},
		"many(summary)": {"<@bot_id> summary arg", "", "", true},
		"few":           {"<@bot_id>", "", "", true},
		"few(sub)":      {"<@bot_id> subscribe", "", "", true},
		"unknown":       {"<@bot_id> unknown command", "", "", true},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			gotC, gotRawURL, err := parseCommand(tc.text)
			if (err != nil) != tc.wantErr {
				t.Errorf("parseCommand() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if gotC != tc.wantC {
				t.Errorf("parseCommand() gotC = %v, want %v", gotC, tc.wantC)
			}
			if gotRawURL != tc.wantRawURL {
				t.Errorf("parseCommand() gotRawURL = %v, want %v", gotRawURL, tc.wantRawURL)
			}
		})
	}
}
