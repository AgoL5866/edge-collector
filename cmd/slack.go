package cmd

import (
	"fmt"
	"github.com/slack-go/slack"
	"strings"
)

const (
	k_chain  = "chain"
	k_node   = "node"
	k_number = "number"
	k_hash   = "hash"
	k_reason = "reason"
)

var _gListKeys = []string{k_chain, k_node, k_number, k_hash, k_reason}

func newBlockField(k string, v any) *slack.TextBlockObject {
	o := slack.NewTextBlockObject(slack.MarkdownType, fmt.Sprintf("*%s* : %s", strings.ToUpper(k), v), false, false)
	return o
}

func pushSlackMessage(hookUrl string, header string, kv map[string]any, skipNotice bool) bool {
	icon := ":information_source:"
	if strings.HasPrefix(kv[k_node].(string), "rockx") {
		icon = ":warning:"
	} else { //is notice
		if skipNotice == true {
			return false
		}
	}
	headerText := slack.NewTextBlockObject(slack.MarkdownType, fmt.Sprintf("%s  *%s*", icon, header), false, false)
	headerSection := slack.NewSectionBlock(headerText, nil, nil)

	fieldSlice := make([]*slack.TextBlockObject, 0, 4)
	for _, k := range _gListKeys {
		v, ok := kv[k]
		if !ok || v == nil || v == "" {
			continue
		}
		fieldSlice = append(fieldSlice, newBlockField(k, v))
	}

	fieldsSection := slack.NewSectionBlock(nil, fieldSlice, nil)
	msg := slack.WebhookMessage{
		Blocks: &slack.Blocks{BlockSet: []slack.Block{headerSection, fieldsSection}},
	}
	if err := slack.PostWebhook(hookUrl, &msg); err != nil {
		fmt.Println(err)
		return false
	}
	return true
}
