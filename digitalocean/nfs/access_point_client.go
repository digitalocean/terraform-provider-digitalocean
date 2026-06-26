package nfs

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/digitalocean/godo"
)

type nfsAccessPointPolicy struct {
	Anonuid                    uint64   `json:"anonuid"`
	Anongid                    uint64   `json:"anongid"`
	Protocols                  []string `json:"protocols"`
	SquashConfig               string   `json:"squash_config"`
	IdentityEnforcementEnabled bool     `json:"identity_enforcement_enabled"`
}

type nfsAccessPoint struct {
	ID           string               `json:"id"`
	Name         string               `json:"name"`
	ShareID      string               `json:"share_id"`
	Path         string               `json:"path"`
	Status       string               `json:"status"`
	AccessPolicy nfsAccessPointPolicy `json:"access_policy"`
	CreatedAt    string               `json:"created_at"`
	UpdatedAt    string               `json:"updated_at"`
	IsDefault    bool                 `json:"is_default"`
	VpcID        string               `json:"vpc_id"`
}

type nfsCreateAccessPointRequest struct {
	Name         string               `json:"name"`
	Path         string               `json:"path"`
	VpcID        string               `json:"vpc_id"`
	AccessPolicy nfsAccessPointPolicy `json:"access_policy"`
}

type nfsAccessPointRoot struct {
	AccessPoint *nfsAccessPoint `json:"access_point"`
}

type nfsAccessPointListRoot struct {
	AccessPoints []*nfsAccessPoint `json:"access_points"`
}

func createNfsAccessPoint(ctx context.Context, client *godo.Client, shareID string, reqBody *nfsCreateAccessPointRequest) (*nfsAccessPoint, *godo.Response, error) {
	path := fmt.Sprintf("v2/nfs/shares/%s/access_points", shareID)
	req, err := client.NewRequest(ctx, http.MethodPost, path, reqBody)
	if err != nil {
		return nil, nil, err
	}

	root := new(nfsAccessPointRoot)
	resp, err := client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.AccessPoint, resp, nil
}

// GetNfsAccessPoint retrieves an NFS access point by ID.
func GetNfsAccessPoint(ctx context.Context, client *godo.Client, accessPointID string) (*nfsAccessPoint, *godo.Response, error) {
	return getNfsAccessPoint(ctx, client, accessPointID)
}

func getNfsAccessPoint(ctx context.Context, client *godo.Client, accessPointID string) (*nfsAccessPoint, *godo.Response, error) {
	path := fmt.Sprintf("v2/nfs/access_points/%s", accessPointID)
	req, err := client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(nfsAccessPointRoot)
	resp, err := client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.AccessPoint, resp, nil
}

func listNfsAccessPoints(ctx context.Context, client *godo.Client, shareID string, vpcID string) ([]*nfsAccessPoint, *godo.Response, error) {
	path := fmt.Sprintf("v2/nfs/shares/%s/access_points", shareID)
	if vpcID != "" {
		path = path + "?" + url.Values{"vpc_id": {vpcID}}.Encode()
	}

	req, err := client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(nfsAccessPointListRoot)
	resp, err := client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.AccessPoints, resp, nil
}

func deleteNfsAccessPoint(ctx context.Context, client *godo.Client, accessPointID string) (*nfsAccessPoint, *godo.Response, error) {
	path := fmt.Sprintf("v2/nfs/access_points/%s", accessPointID)
	req, err := client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(nfsAccessPointRoot)
	resp, err := client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.AccessPoint, resp, nil
}
