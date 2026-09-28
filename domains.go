package pigeon

import (
	"context"
	"net/http"
)

type DomainsSvc struct {
	client *Client
}

type CreateDomainRequest struct {
	Name   string `json:"name"`
	Region string `json:"region,omitempty"`
}

type Domain struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Status  string `json:"status"`
	Region  string `json:"region"`
	Records any    `json:"records"`
}

type ListDomainsResponse struct {
	Data []Domain `json:"data"`
}

func (s *DomainsSvc) List(ctx context.Context) (*ListDomainsResponse, error) {
	req, err := s.client.NewRequest(ctx, http.MethodGet, "api/domains", nil)
	if err != nil {
		return nil, err
	}
	out := new(ListDomainsResponse)
	_, err = s.client.Perform(req, out)
	return out, err
}

func (s *DomainsSvc) Create(ctx context.Context, params *CreateDomainRequest) (*Domain, error) {
	req, err := s.client.NewRequest(ctx, http.MethodPost, "api/domains", params)
	if err != nil {
		return nil, err
	}
	out := new(Domain)
	_, err = s.client.Perform(req, out)
	return out, err
}

func (s *DomainsSvc) Verify(ctx context.Context, id string) (*Domain, error) {
	req, err := s.client.NewRequest(ctx, http.MethodPost, "api/domains/"+id+"/verify", nil)
	if err != nil {
		return nil, err
	}
	out := new(Domain)
	_, err = s.client.Perform(req, out)
	return out, err
}

func (s *DomainsSvc) Remove(ctx context.Context, id string) error {
	req, err := s.client.NewRequest(ctx, http.MethodDelete, "api/domains/"+id, nil)
	if err != nil {
		return err
	}
	_, err = s.client.Perform(req, nil)
	return err
}
