package sl

import (
	"fmt"
	"log"
	"strings"
)

// containMention returns if text is mentioned bot itself.
func (client *Client) containMention(text string) bool {
	log.Printf("judge if contain. mention: %#v text: %#v", client.config.mentionedText, text)
	return strings.Contains(text, client.config.mentionedText)
}

//
func (client *Client) parseCommand(text string) (c command, rawURL string, err error) {
	params := strings.Fields(text)
	tooFew := fmt.Errorf("too few args: %#v", params)
	tooMany := fmt.Errorf("too many args: %#v", params)
	if len(params) < 2 {
		return "", "", tooFew
	}

	mayCommand := strings.ToLower(params[1])
	switch mayCommand {
	// assumes ["@gs", "subscribe", "https://github.com/tbistr/gs-linker/1"]
	case string(subscribe):
		if len(params) < 3 {
			return "", "", tooFew
		} else if 3 < len(params) {
			return "", "", tooMany
		}
		// URL link is posted in slack like <https://github.com/tbistr/gs-linker/1>.
		rawURL := strings.Trim(params[2], "<>")
		return subscribe, rawURL, nil

	// assumes ["@gs", "unsubscribe"]
	case string(unsubscribe):
		if len(params) < 2 {
			return "", "", tooFew
		} else if 2 < len(params) {
			return "", "", tooMany
		}
		return unsubscribe, "", nil

	// assumes ["@gs", "summary"]
	case string(summary):
		if len(params) < 2 {
			return "", "", tooFew
		} else if 2 < len(params) {
			return "", "", tooMany
		}
		return summary, "", nil

	default:
		return "", "", fmt.Errorf("unknown command: %#v", params[1])
	}
}
