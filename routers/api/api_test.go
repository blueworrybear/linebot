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

	"github.com/blueworrybear/lineBot/config"
	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

var (
	cfg *config.Config
)

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

func testWebhook(req *http.Request, f func(*httptest.ResponseRecorder)) {
	r := gin.Default()
	r.Use(UseLineBot(cfg))
	r.POST("/webhooks", HandleWebhooks())
	testRequest(r, req, f)
}

func buildBody(t *testing.T, events... *linebot.Event) *bytes.Buffer{
	body := &struct {
		Events []*linebot.Event `json:"events"`
	} {
		events,
	}
	rawBody, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}
	return bytes.NewBuffer(rawBody)
}

func TestFollow(t *testing.T) {
	events := []*linebot.Event {
		{
			Type: linebot.EventTypeFollow,
		},
	}
	req, _ := http.NewRequest("POST", "/webhooks", buildBody(t, events...))
	testWebhook(req, func (w *httptest.ResponseRecorder) {
		res := w.Result()
		if res.StatusCode != 200 {
			t.Fail()
		}
	})
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
