package pubnub

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/pubnub/go/utils"
)

const SUBSCRIBE_PATH = "/v2/subscribe/%s/%s/0"

func newSubscribeRequest(ctx Context) *SubscribeResponse {
	return &SubscribeResponse{}
}

type SubscribeResponse struct {
}

type SubscribeOpts struct {
	pubnub *PubNub

	Channels []string
	Groups   []string

	Heartbeat        int
	Region           string
	Timetoken        int64
	FilterExpression string
	WithPresence     bool

	Transport http.RoundTripper

	ctx Context
}

func (o *SubscribeOpts) config() Config {
	return *o.pubnub.Config
}

func (o *SubscribeOpts) client() *http.Client {
	return o.pubnub.GetSubscribeClient()
}

func (o *SubscribeOpts) context() Context {
	return o.ctx
}

func (o *SubscribeOpts) validate() error {
	if o.config().PublishKey == "" {
		return newValidationError(o, StrMissingPubKey)
	}

	if o.config().SubscribeKey == "" {
		return newValidationError(o, StrMissingSubKey)
	}

	if len(o.Channels) == 0 && len(o.Groups) == 0 {
		return newValidationError(o, StrMissingChannel)
	}

	return nil
}

func (o *SubscribeOpts) buildPath() (string, error) {
	channels := utils.JoinChannels(o.Channels)

	return fmt.Sprintf(SUBSCRIBE_PATH,
		o.pubnub.Config.SubscribeKey,
		channels,
	), nil
}

func (o *SubscribeOpts) buildQuery() (*url.Values, error) {
	q := defaultQuery(o.pubnub.Config.Uuid)

	if len(o.Groups) > 0 {
		channelGroup := utils.JoinChannels(o.Groups)
		q.Set("channel-group", string(channelGroup))
	}

	if o.Timetoken != 0 {
		q.Set("tt", strconv.FormatInt(o.Timetoken, 10))
	}

	if o.Region != "" {
		q.Set("tr", o.Region)
	}

	if o.FilterExpression != "" {
		q.Set("filter-expr", utils.UrlEncode(o.FilterExpression))
	}

	// hb timeout should be at least 4 seconds
	if o.Heartbeat >= 4 {
		q.Set("heartbeat", fmt.Sprintf("%d", o.Heartbeat))
	}

	return q, nil
}

func (o *SubscribeOpts) buildBody() ([]byte, error) {
	return []byte{}, nil
}

func (o *SubscribeOpts) httpMethod() string {
	return "GET"
}

func (o *SubscribeOpts) isAuthRequired() bool {
	return true
}

func (o *SubscribeOpts) requestTimeout() int {
	return o.pubnub.Config.SubscribeRequestTimeout
}

func (o *SubscribeOpts) connectTimeout() int {
	return o.pubnub.Config.ConnectTimeout
}

func (o *SubscribeOpts) operationType() OperationType {
	return PNSubscribeOperation
}