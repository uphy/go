package e2e

import (
	"testing"
	"time"

	pubnub "github.com/pubnub/go"
	"github.com/pubnub/go/tests/stubs"
	"github.com/stretchr/testify/assert"
)

func TestListAllChannelGroupNotStubbed(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(configCopy())
	_, _, err := pn.ListAllChannelsChannelGroup().
		ChannelGroup("cg").
		Execute()

	assert.Nil(err)
}

func TestListAllChannelGroupMissingGroup(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(configCopy())

	_, _, err := pn.DeleteChannelGroup().
		Execute()

	assert.Contains(err.Error(), "Missing Channel Group")
}

func TestListAllChannelGroupSuperCall(t *testing.T) {
	assert := assert.New(t)

	config := pamConfigCopy()

	// Not allowed characters:
	// .,:*

	validCharacters := "-_~?#[]@!$&'()+;=`|"

	config.Uuid = SPECIAL_CHARACTERS
	config.AuthKey = SPECIAL_CHARACTERS

	pn := pubnub.NewPubNub(config)

	_, _, err := pn.ListAllChannelsChannelGroup().
		ChannelGroup(validCharacters).
		Execute()

	assert.Nil(err)
}

func TestListAllChannelGroupSuccess(t *testing.T) {
	assert := assert.New(t)
	myChannel := "my-channel"
	myGroup := randomized("my-group")

	interceptor := stubs.NewInterceptor()
	interceptor.AddStub(&stubs.Stub{
		Method:             "GET",
		Path:               "/v1/channel-registration/sub-key/sub-c-5c4fdcc6-c040-11e5-a316-0619f8945a4f/channel-group/" + myGroup,
		Query:              "add=my-channel",
		ResponseBody:       `{"status": 200, "message": "OK", "service": "channel-registry", "error": false}`,
		IgnoreQueryKeys:    []string{"uuid", "pnsdk"},
		ResponseStatusCode: 200,
	})
	interceptor.AddStub(&stubs.Stub{
		Method:             "GET",
		Path:               "/v1/channel-registration/sub-key/sub-c-5c4fdcc6-c040-11e5-a316-0619f8945a4f/channel-group/" + myGroup,
		Query:              "",
		ResponseBody:       `{"status": 200, "payload": {"channels": ["my-channel"], "group": "` + myGroup + `"}, "service": "channel-registry", "error": false}`,
		IgnoreQueryKeys:    []string{"uuid", "pnsdk"},
		ResponseStatusCode: 200,
	})
	interceptor.AddStub(&stubs.Stub{
		Method:             "GET",
		Path:               "/v1/channel-registration/sub-key/sub-c-5c4fdcc6-c040-11e5-a316-0619f8945a4f/channel-group/" + myGroup,
		Query:              "remove=my-channel",
		ResponseBody:       `{"status": 200, "message": "OK", "service": "channel-registry", "error": false}`,
		IgnoreQueryKeys:    []string{"uuid", "pnsdk"},
		ResponseStatusCode: 200,
	})

	pn := pubnub.NewPubNub(configCopy())
	pn.SetClient(interceptor.GetClient())

	_, _, err := pn.AddChannelChannelGroup().
		Channels([]string{myChannel}).
		Group(myGroup).
		Execute()

	assert.Nil(err)

	// await for adding channel
	time.Sleep(2 * time.Second)

	res, _, err := pn.ListAllChannelsChannelGroup().
		ChannelGroup(myGroup).
		Execute()

	assert.Nil(err)

	assert.Equal(myChannel, res.Channels[0])
	assert.Equal(myGroup, res.Group)

	_, _, err = pn.RemoveChannelChannelGroup().
		Channels([]string{myChannel}).
		Group(myGroup).
		Execute()

	assert.Nil(err)
}