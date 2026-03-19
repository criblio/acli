package bitbucket

import (
	"encoding/json"
	"fmt"
	"net/url"
)

type Workspace struct {
	UUID      string `json:"uuid"`
	Name      string `json:"name"`
	Slug      string `json:"slug"`
	IsPrivate bool   `json:"is_private"`
	CreatedOn string `json:"created_on"`
	UpdatedOn string `json:"updated_on"`
	Links     struct {
		HTML struct {
			Href string `json:"href"`
		} `json:"html"`
		Avatar struct {
			Href string `json:"href"`
		} `json:"avatar"`
	} `json:"links"`
}

type WorkspaceMember struct {
	User struct {
		DisplayName string `json:"display_name"`
		UUID        string `json:"uuid"`
		Nickname    string `json:"nickname"`
		AccountID   string `json:"account_id"`
	} `json:"user"`
	Workspace struct {
		Slug string `json:"slug"`
		UUID string `json:"uuid"`
	} `json:"workspace"`
}

type WorkspacePermission struct {
	Permission string `json:"permission"`
	User       struct {
		DisplayName string `json:"display_name"`
		UUID        string `json:"uuid"`
		Nickname    string `json:"nickname"`
	} `json:"user"`
}

func (c *Client) ListWorkspaces(opts *PaginationOptions) ([]Workspace, error) {
	params := url.Values{}
	if opts != nil {
		opts.applyParams(params)
	}
	ensurePageLen(params)

	path := "/workspaces"
	if len(params) > 0 {
		path += "?" + params.Encode()
	}

	if opts != nil && opts.All {
		pages, err := c.getAll(path)
		if err != nil && len(pages) == 0 {
			return nil, err
		}
		var workspaces []Workspace
		for _, pg := range pages {
			var pageWorkspaces []Workspace
			if err := json.Unmarshal(pg.Values, &pageWorkspaces); err != nil {
				return workspaces, fmt.Errorf("parsing workspaces: %w", err)
			}
			workspaces = append(workspaces, pageWorkspaces...)
		}
		return workspaces, nil
	}

	data, err := c.get(path)
	if err != nil {
		return nil, err
	}
	var page PaginatedResponse
	if err := json.Unmarshal(data, &page); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}
	var workspaces []Workspace
	if err := json.Unmarshal(page.Values, &workspaces); err != nil {
		return nil, fmt.Errorf("parsing workspaces: %w", err)
	}
	return workspaces, nil
}

func (c *Client) GetWorkspace(workspace string) (*Workspace, error) {
	path := fmt.Sprintf("/workspaces/%s", url.PathEscape(workspace))
	data, err := c.get(path)
	if err != nil {
		return nil, err
	}
	var ws Workspace
	if err := json.Unmarshal(data, &ws); err != nil {
		return nil, fmt.Errorf("parsing workspace: %w", err)
	}
	return &ws, nil
}

func (c *Client) ListWorkspaceMembers(workspace string, opts *PaginationOptions) ([]WorkspaceMember, error) {
	params := url.Values{}
	if opts != nil {
		opts.applyParams(params)
	}
	ensurePageLen(params)

	path := fmt.Sprintf("/workspaces/%s/members", url.PathEscape(workspace))
	if len(params) > 0 {
		path += "?" + params.Encode()
	}

	if opts != nil && opts.All {
		pages, err := c.getAll(path)
		if err != nil && len(pages) == 0 {
			return nil, err
		}
		var members []WorkspaceMember
		for _, pg := range pages {
			var pageMembers []WorkspaceMember
			if err := json.Unmarshal(pg.Values, &pageMembers); err != nil {
				return members, fmt.Errorf("parsing members: %w", err)
			}
			members = append(members, pageMembers...)
		}
		return members, nil
	}

	data, err := c.get(path)
	if err != nil {
		return nil, err
	}
	var page PaginatedResponse
	if err := json.Unmarshal(data, &page); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}
	var members []WorkspaceMember
	if err := json.Unmarshal(page.Values, &members); err != nil {
		return nil, fmt.Errorf("parsing members: %w", err)
	}
	return members, nil
}

func (c *Client) ListWorkspacePermissions(workspace string, opts *PaginationOptions) ([]WorkspacePermission, error) {
	params := url.Values{}
	if opts != nil {
		opts.applyParams(params)
	}
	ensurePageLen(params)

	path := fmt.Sprintf("/workspaces/%s/permissions", url.PathEscape(workspace))
	if len(params) > 0 {
		path += "?" + params.Encode()
	}

	if opts != nil && opts.All {
		pages, err := c.getAll(path)
		if err != nil && len(pages) == 0 {
			return nil, err
		}
		var perms []WorkspacePermission
		for _, pg := range pages {
			var pagePerms []WorkspacePermission
			if err := json.Unmarshal(pg.Values, &pagePerms); err != nil {
				return perms, fmt.Errorf("parsing permissions: %w", err)
			}
			perms = append(perms, pagePerms...)
		}
		return perms, nil
	}

	data, err := c.get(path)
	if err != nil {
		return nil, err
	}
	var page PaginatedResponse
	if err := json.Unmarshal(data, &page); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}
	var perms []WorkspacePermission
	if err := json.Unmarshal(page.Values, &perms); err != nil {
		return nil, fmt.Errorf("parsing permissions: %w", err)
	}
	return perms, nil
}
