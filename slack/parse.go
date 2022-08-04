package sl

import "strings"

// containMention returns if text is mentioned bot itself.
func (client *Client) containMention(text string) bool {
	return strings.Contains(text, client.config.mentionedText)
}
