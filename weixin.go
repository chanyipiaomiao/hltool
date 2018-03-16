package hltool

import (
	"fmt"

	"github.com/levigross/grequests"
	"github.com/tidwall/gjson"
)

var (

	// WeixinErr 微信错误信息码
	WeixinErr = func(errcode int64, errmsg string) error {
		return fmt.Errorf("weixin return error, errcode: %d, errmsg: %s", errcode, errmsg)
	}
)

const (
	// TEXT 微信文本消息
	TEXT = "text"

	// VIDEO 微信视频消息
	VIDEO = "video"

	// IMAGE 微信图片消息
	IMAGE = "image"

	// VOICE 微信语音消息
	VOICE = "voice"

	// FILE 微信文件消息
	FILE = "file"

	// TEXTCARD 微信文本卡片消息
	TEXTCARD = "textcard"

	// NEWS 微信图文消息
	NEWS = "news"

	// MPNEWS 微信图文消息
	MPNEWS = "mpnews"
)

// WeixinText 文本消息
type WeixinText struct {
	Content map[string]interface{}
}

// NewWeixinText new 文本消息,
// content  文本内容
func NewWeixinText(content string) *WeixinText {
	return &WeixinText{
		Content: map[string]interface{}{
			"content": content,
		},
	}
}

// WeixinImageVoiceFile 图片语音文件消息 统一用这个
type WeixinImageVoiceFile struct {
	MediaID map[string]interface{}
}

// NewWeixinImageVoiceFile new 图片音频文件消息,
// mediaID 图片媒体文件id，可以调用上传临时素材接口获取
func NewWeixinImageVoiceFile(mediaID string) *WeixinImageVoiceFile {
	return &WeixinImageVoiceFile{
		MediaID: map[string]interface{}{
			"media_id": mediaID,
		},
	}
}

// WeixinVideo 视频消息
type WeixinVideo struct {
	Video map[string]interface{}
}

// NewWeixinVideo new 视频消息,
// mediaID 图片媒体文件id，可以调用上传临时素材接口获取,
// title 视频消息的标题,
// desc 视频消息的描述
func NewWeixinVideo(mediaID, title, desc string) *WeixinVideo {
	return &WeixinVideo{
		Video: map[string]interface{}{
			"media_id":    mediaID,
			"title":       title,
			"description": desc,
		},
	}
}

// WeixinTextCard 文本卡片消息
type WeixinTextCard struct {
	Textcard map[string]interface{}
}

// NewWeixinTextCard new文本卡片消息,
// title 标题,
// desc  描述,
// url   点击后跳转的链接,
// btntxt 按钮文字。 默认为“详情”， 不超过4个文字，超过自动截断
func NewWeixinTextCard(title, desc, url, btntxt string) *WeixinTextCard {
	return &WeixinTextCard{
		Textcard: map[string]interface{}{
			"title":       title,
			"description": desc,
			"url":         url,
			"btntxt":      btntxt,
		},
	}
}

// WeixinNews 图文消息
type WeixinNews struct {
	Articles []map[string]interface{}
}

// NewWeixinNews new 图文消息
// title 标题,
// desc  描述,
// url   点击后跳转的链接,
// btntxt 按钮文字。 默认为“详情”， 不超过4个文字，超过自动截断,
// picurl 图文消息的图片链接，支持JPG、PNG格式，较好的效果为大图640320，小图8080。
func NewWeixinNews(title, desc, url, picurl, btntxt string) *WeixinNews {
	return &WeixinNews{
		Articles: []map[string]interface{}{
			map[string]interface{}{
				"title":       title,
				"description": desc,
				"url":         url,
				"picurl":      picurl,
				"btntxt":      btntxt,
			},
		},
	}
}

// WeixinMPNews 图文消息 跟普通的图文消息一致，唯一的差异是图文内容存储在企业微信
// 多次发送mpnews，会被认为是不同的图文，阅读、点赞的统计会被分开计算
type WeixinMPNews struct {
	Articles []map[string]interface{}
}

// NewWeixinMPNews new new 图文消息
// title 标题
// thumbMediaID 图文消息缩略图的media_id
// author 图文消息的作者
// contentSourceURL 图文消息点击“阅读原文”之后的页面链接
// content 图文消息的内容，支持html标签，不超过666 K个字节
// digest 图文消息的描述
func NewWeixinMPNews(title, thumbMediaID, author, contentSourceURL, content, digest string) *WeixinMPNews {
	return &WeixinMPNews{
		Articles: []map[string]interface{}{
			map[string]interface{}{
				"title":              title,
				"thumb_media_id":     thumbMediaID,
				"author":             author,
				"content_source_url": contentSourceURL,
				"content":            content,
				"digest":             digest,
			},
		},
	}
}

// WeixinMessage 微信消息
type WeixinMessage struct {
	MsgType  string                `json:"msgtype"`
	ToUser   string                `json:"touser"`
	ToParty  string                `json:"toparty"`
	ToTag    string                `json:"totag"`
	AgentID  int64                 `json:"agentid"`
	Safe     int64                 `json:"safe"`
	Text     *WeixinText           `json:"text"`
	Image    *WeixinImageVoiceFile `json:"image"`
	Voice    *WeixinImageVoiceFile `json:"voice"`
	File     *WeixinImageVoiceFile `json:"file"`
	Video    *WeixinVideo          `json:"video"`
	TextCard *WeixinTextCard       `json:"textcard"`
	News     *WeixinNews           `json:"news"`
	MPNews   *WeixinMPNews         `json:"mpnews"`
}

// NewWeixinMessage new 消息对象
func NewWeixinMessage(msgtype, toUser, toParty, toTag string, agentID, safe int64, message interface{}) *WeixinMessage {
	msg := &WeixinMessage{
		MsgType: msgtype,
		ToUser:  toUser,
		ToParty: toParty,
		ToTag:   toTag,
		AgentID: agentID,
		Safe:    safe,
	}
	switch message.(type) {
	case *WeixinText:
		msg.Text = message.(*WeixinText)
	case *WeixinImageVoiceFile:
		switch msgtype {
		case IMAGE:
			msg.Image = message.(*WeixinImageVoiceFile)
		case VOICE:
			msg.Voice = message.(*WeixinImageVoiceFile)
		case FILE:
			msg.File = message.(*WeixinImageVoiceFile)
		}
	case *WeixinVideo:
		msg.Video = message.(*WeixinVideo)
	case *WeixinTextCard:
		msg.TextCard = message.(*WeixinTextCard)
	case *WeixinNews:
		msg.News = message.(*WeixinNews)
	case *WeixinMPNews:
		msg.MPNews = message.(*WeixinMPNews)
	default:
		return nil
	}
	return msg
}

// WeixinClient 微信
type WeixinClient struct {
	API string
}

// NewWeixinClient new weixin对象
func NewWeixinClient(api string) *WeixinClient {
	return &WeixinClient{
		API: api,
	}
}

// GetAccessToken 获取AccessToken
func (w *WeixinClient) GetAccessToken(corpid, corpsecret string) (string, error) {
	api := w.API + "/gettoken"

	o := &grequests.RequestOptions{
		Params: map[string]string{
			"corpid":     corpid,
			"corpsecret": corpsecret,
		},
	}

	resp, err := grequests.Get(api, o)
	if err != nil {
		return "", err
	}

	respJSON := resp.String()
	errcode := gjson.Get(respJSON, "errcode")
	token := gjson.Get(respJSON, "access_token")
	if errcode.Int() == 0 {
		return token.String(), nil
	}
	return "", WeixinErr(errcode.Int(), gjson.Get(respJSON, "errmsg").String())
}

// SendMessage 发送消息
func (w *WeixinClient) SendMessage() (bool, error) {

	return false, nil
}
