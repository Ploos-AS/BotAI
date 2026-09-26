package client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestReferenceClient(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
		w.Header().Set("Content-Type","application/json")
		switch r.URL.Path {
		case "/v1/version": w.Write([]byte(`{"api_version":"1.0.0"}`))
		case "/v1/experts": w.Write([]byte(`[{"id":"irc","description":"IRC"}]`))
		case "/v1/chat": w.Write([]byte(`{"expert":"irc","text":"hello","provider":"test"}`))
		default: http.NotFound(w,r)
		}
	}))
	defer s.Close()
	c,err:=New(s.URL);if err!=nil{t.Fatal(err)}
	ctx:=context.Background()
	if err:=c.Compatible(ctx);err!=nil{t.Fatal(err)}
	e,err:=c.Experts(ctx);if err!=nil||len(e)!=1||e[0].ID!="irc"{t.Fatalf("experts=%#v err=%v",e,err)}
	out,err:=c.Chat(ctx,ChatRequest{Expert:"irc",Message:"hi"});if err!=nil{t.Fatal(err)}
	if out.Text!="hello"||out.Expert!="irc"{t.Fatalf("out=%#v",out)}
}

func TestReferenceClientRejectsLongHistory(t *testing.T) {
	c,_:=New("http://127.0.0.1")
	h:=make([]Message,21)
	if _,err:=c.Chat(context.Background(),ChatRequest{Message:"hi",History:h});err==nil{t.Fatal("expected history error")}
}

func TestCompatibilityRejectsDifferentVersion(t *testing.T) {
	s:=httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){w.Write([]byte(`{"api_version":"2.0.0"}`))}))
	defer s.Close()
	c,_:=New(s.URL)
	if err:=c.Compatible(context.Background());err==nil{t.Fatal("expected compatibility error")}
}

func TestUnavailableReturnsError(t *testing.T) {
	c,_:=New("http://127.0.0.1:1")
	if _,err:=c.Chat(context.Background(),ChatRequest{Message:"hi"});err==nil{t.Fatal("expected unavailable error")}
}

func TestHTTPErrorDoesNotExposeBody(t *testing.T) {
	const secret="provider-secret-body"
	s:=httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
		http.Error(w,secret,http.StatusBadGateway)
	}))
	defer s.Close()
	c,_:=New(s.URL)
	_,err:=c.Chat(context.Background(),ChatRequest{Message:"hi"})
	if err==nil{t.Fatal("expected HTTP error")}
	if strings.Contains(err.Error(),secret){t.Fatalf("upstream body leaked: %v",err)}
	if !strings.Contains(err.Error(),"502"){t.Fatalf("status missing: %v",err)}
}

func TestOversizedResponseFails(t *testing.T) {
	s:=httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
		w.Header().Set("Content-Type","application/json")
		w.Write([]byte(`{"expert":"irc","text":"`))
		w.Write([]byte(strings.Repeat("x",(1<<20)+1024)))
		w.Write([]byte(`","provider":"test"}`))
	}))
	defer s.Close()
	c,_:=New(s.URL)
	if _,err:=c.Chat(context.Background(),ChatRequest{Message:"hi"});err==nil{t.Fatal("expected bounded decode error")}
}
