package pigeon

import (
	"context"
	"net/http"
)

type ContactsSvc struct {
	client *Client
}

type CreateContactRequest struct {
	Email string `json:"email"`
	Name  string `json:"name,omitempty"`
}

type Contact struct {
	Id    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

func (s *ContactsSvc) Create(ctx context.Context, params *CreateContactRequest) (*Contact, error) {
	req, err := s.client.NewRequest(ctx, http.MethodPost, "api/contacts", params)
	if err != nil {
		return nil, err
	}
	out := new(Contact)
	_, err = s.client.Perform(req, out)
	return out, err
}

type AutomationsSvc struct {
	client *Client
}

type TriggerAutomationRequest struct {
	To        string         `json:"to,omitempty"`
	Variables map[string]any `json:"variables,omitempty"`
}

func (s *AutomationsSvc) Trigger(ctx context.Context, id string, params *TriggerAutomationRequest) (*Email, error) {
	req, err := s.client.NewRequest(ctx, http.MethodPost, "api/automations/"+id+"/trigger", params)
	if err != nil {
		return nil, err
	}
	out := new(Email)
	_, err = s.client.Perform(req, out)
	return out, err
}
