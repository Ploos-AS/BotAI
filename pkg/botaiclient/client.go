package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const APIVersion = "1.0.0"

type Message struct {
	Role string `json:"role"`
	Content string `json:"content"`
}

type Expert struct {
	ID string `json:"id"`
	Description string `json:"description"`
}

type ChatRequest struct {
	Expert string `json:"expert,omitempty"`
	History []Message `json:"history,omitempty"`
	Message string `json:"message"`
}

type ChatResponse struct {
	Expert string `json:"expert"`
	Text string `json:"text"`
	Provider string `json:"provider"`
}

type Client struct {
	baseURL string
	http *http.Client
}

func New(baseURL string) (*Client, error) {
	baseURL = strings.TrimRight(strings.TrimSpace(baseURL), "/")
	if baseURL == "" { return nil, fmt.Errorf("base URL is required") }
	return &Client{baseURL:baseURL, http:&http.Client{Timeout:60*time.Second}}, nil
}

func (c *Client) do(ctx context.Context, method, path string, in, out any) error {
	var body io.Reader
	if in != nil {
		b, err := json.Marshal(in)
		if err != nil { return err }
		body = bytes.NewReader(b)
	}
	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, body)
	if err != nil { return err }
	if in != nil { req.Header.Set("Content-Type","application/json") }
	resp, err := c.http.Do(req)
	if err != nil { return err }
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		io.Copy(io.Discard, io.LimitReader(resp.Body, 64<<10))
		return fmt.Errorf("BotAI returned HTTP %d", resp.StatusCode)
	}
	if out == nil { return nil }
	return json.NewDecoder(io.LimitReader(resp.Body, 1<<20)).Decode(out)
}

func (c *Client) Version(ctx context.Context) (string, error) {
	var out struct{ APIVersion string `json:"api_version"` }
	if err := c.do(ctx,http.MethodGet,"/v1/version",nil,&out); err != nil { return "",err }
	return out.APIVersion,nil
}

func (c *Client) Compatible(ctx context.Context) error {
	v, err := c.Version(ctx)
	if err != nil { return err }
	if v != APIVersion { return fmt.Errorf("unsupported BotAI API version %q",v) }
	return nil
}

func (c *Client) Experts(ctx context.Context) ([]Expert,error) {
	var out []Expert
	if err:=c.do(ctx,http.MethodGet,"/v1/experts",nil,&out);err!=nil{return nil,err}
	return out,nil
}

func (c *Client) Chat(ctx context.Context, in ChatRequest) (ChatResponse,error) {
	var out ChatResponse
	if strings.TrimSpace(in.Message)=="" { return out,fmt.Errorf("message is required") }
	if len(in.History)>20 { return out,fmt.Errorf("history exceeds 20 messages") }
	err:=c.do(ctx,http.MethodPost,"/v1/chat",in,&out)
	return out,err
}
