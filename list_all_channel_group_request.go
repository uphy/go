package pubnub

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/pubnub/go/pnerr"
	"github.com/pubnub/go/utils"
)

const ALL_CHANNEL_GROUP = "/v1/channel-registration/sub-key/%s/channel-group/%s"

var emptyAllChannelGroupResponse *AllChannelGroupResponse

type allChannelGroupBuilder struct {
	opts *allChannelGroupOpts
}

func newAllChannelGroupBuilder(pubnub *PubNub) *allChannelGroupBuilder {
	builder := allChannelGroupBuilder{
		opts: &allChannelGroupOpts{
			pubnub: pubnub,
		},
	}

	return &builder
}

func newAllChannelGroupBuilderWithContext(pubnub *PubNub,
	context Context) *allChannelGroupBuilder {
	builder := allChannelGroupBuilder{
		opts: &allChannelGroupOpts{
			pubnub: pubnub,
			ctx:    context,
		},
	}

	return &builder
}

func (b *allChannelGroupBuilder) ChannelGroup(
	cg string) *allChannelGroupBuilder {
	b.opts.ChannelGroup = cg
	return b
}

func (b *allChannelGroupBuilder) Execute() (
	*AllChannelGroupResponse, StatusResponse, error) {
	rawJson, status, err := executeRequest(b.opts)
	if err != nil {
		return emptyAllChannelGroupResponse, status, err
	}

	return newAllChannelGroupResponse(rawJson, status)
}

type allChannelGroupOpts struct {
	pubnub *PubNub

	ChannelGroup string

	Transport http.RoundTripper

	ctx Context
}

func (o *allChannelGroupOpts) config() Config {
	return *o.pubnub.Config
}

func (o *allChannelGroupOpts) client() *http.Client {
	return o.pubnub.GetClient()
}

func (o *allChannelGroupOpts) context() Context {
	return o.ctx
}

func (o *allChannelGroupOpts) validate() error {
	if o.config().SubscribeKey == "" {
		return ErrMissingSubKey
	}

	if o.ChannelGroup == "" {
		return ErrMissingChannelGroup
	}

	return nil
}

func (o *allChannelGroupOpts) buildPath() (string, error) {
	return fmt.Sprintf(ALL_CHANNEL_GROUP,
		o.pubnub.Config.SubscribeKey,
		utils.UrlEncode(o.ChannelGroup)), nil
}

func (o *allChannelGroupOpts) buildQuery() (*url.Values, error) {
	q := defaultQuery(o.pubnub.Config.Uuid)

	return q, nil
}

func (o *allChannelGroupOpts) buildBody() ([]byte, error) {
	return []byte{}, nil
}

func (o *allChannelGroupOpts) httpMethod() string {
	return "GET"
}

func (o *allChannelGroupOpts) isAuthRequired() bool {
	return true
}

func (o *allChannelGroupOpts) requestTimeout() int {
	return o.pubnub.Config.NonSubscribeRequestTimeout
}

func (o *allChannelGroupOpts) connectTimeout() int {
	return o.pubnub.Config.ConnectTimeout
}

func (o *allChannelGroupOpts) operationType() PNOperationType {
	return PNChannelsForGroupOperation
}

type AllChannelGroupResponse struct {
	Channels []string
	Group    string
}

func newAllChannelGroupResponse(jsonBytes []byte, status StatusResponse) (
	*AllChannelGroupResponse, StatusResponse, error) {
	resp := &AllChannelGroupResponse{}

	var value interface{}

	err := json.Unmarshal(jsonBytes, &value)
	if err != nil {
		e := pnerr.NewResponseParsingError("Error unmarshalling response",
			ioutil.NopCloser(bytes.NewBufferString(string(jsonBytes))), err)

		return emptyAllChannelGroupResponse, status, e
	}

	if parsedValue, ok := value.(map[string]interface{}); ok {
		if payload, ok := parsedValue["payload"].(map[string]interface{}); ok {
			if group, ok := payload["group"].(string); ok {
				resp.Group = group
			}

			if channels, ok := payload["channels"].([]interface{}); ok {
				parsedChannels := []string{}

				for _, channel := range channels {
					if ch, ok := channel.(string); ok {
						parsedChannels = append(parsedChannels, ch)
					}
				}

				resp.Channels = parsedChannels
			}
		}
	}

	return resp, status, nil
}