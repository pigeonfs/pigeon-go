package pigeon

import (
	"context"
	"net/http"
	"net/url"
)

type EmailsSvc struct {
	client *Client
}

type SendEmailRequest struct {
	From      string            `json:"from"`
	To        []string          `json:"to"`
	Subject   string            `json:"subject"`
	Bcc       []string          `json:"bcc,omitempty"`
	Cc        []string          `json:"cc,omitempty"`
	ReplyTo   []string          `json:"reply_to,omitempty"`
	Html      string            `json:"html,omitempty"`
	Text      string            `json:"text,omitempty"`
	Template  string            `json:"template,omitempty"`
	Variables map[string]any    `json:"variables,omitempty"`
}

type Email struct {
	Id        string   `json:"id"`
	From      string   `json:"from"`
	To        []string `json:"to"`
	Cc        []string `json:"cc"`
	Bcc       []string `json:"bcc"`
	ReplyTo   []string `json:"reply_to"`
	Subject   string   `json:"subject"`
	Html      string   `json:"html"`
	Text      string   `json:"text"`
	Status    string   `json:"status"`
	CreatedAt string   `json:"created_at"`
}

type ListEmailsResponse struct {
	Data []Email `json:"data"`
}

func (s *EmailsSvc) Send(ctx context.Context, params *SendEmailRequest) (*Email, error) {
	path := "api/emails"
	req, err := s.client.NewRequest(ctx, http.MethodPost, path, params)
	if err != nil {
		return nil, err
	}
	email := new(Email)
	_, err = s.client.Perform(req, email)
	return email, err
}

func (s *EmailsSvc) List(ctx context.Context, status, q string) (*ListEmailsResponse, error) {
	query := url.Values{}
	if status != "" {
		query.Set("status", status)
	}
	if q != "" {
		query.Set("q", q)
	}
	path := "api/emails"
	if encoded := query.Encode(); encoded != "" {
		path += "?" + encoded
	}
	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	out := new(ListEmailsResponse)
	_, err = s.client.Perform(req, out)
	return out, err
}

func (s *EmailsSvc) Get(ctx context.Context, id string) (*Email, error) {
	req, err := s.client.NewRequest(ctx, http.MethodGet, "api/emails/"+id, nil)
	if err != nil {
		return nil, err
	}
	email := new(Email)
	_, err = s.client.Perform(req, email)
	return email, err
}

func (s *EmailsSvc) Cancel(ctx context.Context, id string) (*Email, error) {
	req, err := s.client.NewRequest(ctx, http.MethodPost, "api/emails/"+id+"/cancel", nil)
	if err != nil {
		return nil, err
	}
	email := new(Email)
	_, err = s.client.Perform(req, email)
	return email, err
}
