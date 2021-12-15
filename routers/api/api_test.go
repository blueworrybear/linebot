package api

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/blueworrybear/lineBot/config"
	"github.com/blueworrybear/lineBot/mock"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

var (
	cfg *config.Config
)

type RoundTripFunc func(req *http.Request) *http.Response

// RoundTrip .
func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

func testRequest(r *gin.Engine, req *http.Request, f func(*httptest.ResponseRecorder)) {
	rawBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Fatal(err)
	}
	h := hmac.New(sha256.New, []byte(cfg.Server.ChannelSecret))
	h.Write(rawBody)
	req.Header.Set("X-Line-Signature", base64.StdEncoding.EncodeToString(h.Sum(nil)))
	req.Body = io.NopCloser(bytes.NewBuffer(rawBody))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	f(w)
}

func testWebhook(req *http.Request, roundTrip RoundTripFunc,setup func(*gin.Engine),f func(*httptest.ResponseRecorder)) {
	r := gin.Default()
	r.Use(useMockLineBot(roundTrip))
	r.Use(ParseRequest())
	setup(r)
	testRequest(r, req, f)
}

func useMockLineBot(r RoundTripFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		client := &http.Client{
			Transport: r,
		}
		bot, _ := linebot.New(cfg.Server.ChannelSecret, cfg.Server.ChannelAccessToken, linebot.WithHTTPClient(client))
		WithLineBot(c, bot)
	}
}

func buildBody(t *testing.T, events... *linebot.Event) *bytes.Buffer{
	var err error
	rawEvents := make([]json.RawMessage, len(events))
	for i, event := range events {
		rawEvents[i], err = event.MarshalJSON()
		if err != nil {
			t.Fatal(err)
		}
	}
	body := &struct {
		Events []json.RawMessage `json:"events"`
	} {
		rawEvents,
	}
	rawBody, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}
	return bytes.NewBuffer(rawBody)
}

func TestFollow(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockStore := mock.NewMockUserStore(ctrl)
	events := []*linebot.Event {
		{
			Type: linebot.EventTypeFollow,
			Source: &linebot.EventSource{
				Type: linebot.EventSourceTypeUser,
				UserID: "123",
			},
		},
	}
	req, _ := http.NewRequest("POST", "/webhooks", buildBody(t, events...))
	testWebhook(
		req,
		RoundTripFunc(func(req *http.Request) *http.Response {
			return &http.Response{
				StatusCode: 200,
			}
		}),
		func(r *gin.Engine){
			mockStore.EXPECT().Find(gomock.Any()).Return(nil, nil)
			mockStore.EXPECT().Create(gomock.Any()).Return(nil)
			r.POST("/webhooks", HandleFollowEvent(mockStore), HandleWebhooks())
		},
		func(w *httptest.ResponseRecorder) {
			res := w.Result()
			if res.StatusCode != 200 {
				t.Fail()
			}
		},
	)
}

type Bear struct {
	t time.Time
}

func TestTime(t *testing.T) {
	now := time.Now()
	b := &Bear {
		t: now,
	}
	t.Log(b.t)
	bbb(b)
	t.Log(b.t)
}

func bbb(b *Bear) {
	b.t = time.Now()
}

func TestMain(m *testing.M) {
	var err error
	cfg, err = config.Environ()
	cfg.Server.ChannelSecret = "secret"
	cfg.Server.ChannelAccessToken = "token"
	if err != nil {
		log.Panic(err)
	}
	exit := m.Run()
	os.Exit(exit)
}
