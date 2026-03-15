package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Client struct {
	BaseURL    string
	Token      string
	HTTPClient *http.Client
}

func NewClient(baseURL, token string) *Client {
	return &Client{
		BaseURL:    strings.TrimRight(baseURL, "/"),
		Token:      token,
		HTTPClient: &http.Client{},
	}
}

// API types

type TokenResponse struct {
	Token string       `json:"token"`
	User  TokenUser    `json:"user"`
}

type TokenUser struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Role string `json:"role"`
}

type Book struct {
	ID             int    `json:"id"`
	Title          string `json:"title"`
	Subtitle       string `json:"subtitle"`
	Author         string `json:"author"`
	Slug           string `json:"slug"`
	Published      bool   `json:"published"`
	Theme          string `json:"theme"`
	EveryoneAccess bool   `json:"everyone_access"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
}

type BookDetail struct {
	Book
	Leaves []Leaf `json:"leaves"`
}

type Leaf struct {
	ID        int     `json:"id"`
	Title     string  `json:"title"`
	Type      string  `json:"type"`
	Body      string  `json:"body,omitempty"`
	Theme     string  `json:"theme,omitempty"`
	Caption   string  `json:"caption,omitempty"`
	HasImage  bool    `json:"has_image,omitempty"`
	Status    string  `json:"status"`
	Position  float64 `json:"position"`
	Slug      string  `json:"slug"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}

type CreateBookParams struct {
	Title          string `json:"title"`
	Subtitle       string `json:"subtitle,omitempty"`
	Author         string `json:"author,omitempty"`
	Theme          string `json:"theme,omitempty"`
	EveryoneAccess *bool  `json:"everyone_access,omitempty"`
}

type UpdateBookParams struct {
	Title          *string `json:"title,omitempty"`
	Subtitle       *string `json:"subtitle,omitempty"`
	Author         *string `json:"author,omitempty"`
	Theme          *string `json:"theme,omitempty"`
	EveryoneAccess *bool   `json:"everyone_access,omitempty"`
}

type CreatePageParams struct {
	Title string
	Body  string
}

type UpdatePageParams struct {
	Title *string
	Body  *string
}

type CreateSectionParams struct {
	Title string
	Body  string
	Theme string
}

type UpdateSectionParams struct {
	Title *string
	Body  *string
	Theme *string
}

type APIError struct {
	ErrorMsg string `json:"error"`
}

func (e *APIError) Error() string {
	return e.ErrorMsg
}

// HTTP helpers

func (c *Client) doRequest(method, path string, body interface{}) ([]byte, int, error) {
	url := c.BaseURL + path

	var reqBody io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, 0, fmt.Errorf("marshaling request: %w", err)
		}
		reqBody = bytes.NewReader(data)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, 0, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, fmt.Errorf("reading response: %w", err)
	}

	return respBody, resp.StatusCode, nil
}

func (c *Client) checkError(body []byte, statusCode int) error {
	if statusCode >= 200 && statusCode < 300 {
		return nil
	}
	var apiErr APIError
	if err := json.Unmarshal(body, &apiErr); err == nil && apiErr.ErrorMsg != "" {
		return fmt.Errorf("API error (%d): %s", statusCode, apiErr.ErrorMsg)
	}
	return fmt.Errorf("API error (%d): %s", statusCode, string(body))
}

// Token

func (c *Client) Login(email, password string) (*TokenResponse, error) {
	body := map[string]string{"email": email, "password": password}
	respBody, status, err := c.doRequest("POST", "/api/v1/tokens", body)
	if err != nil {
		return nil, err
	}
	if err := c.checkError(respBody, status); err != nil {
		return nil, err
	}
	var result TokenResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}
	return &result, nil
}

// Books

func (c *Client) ListBooks() ([]Book, error) {
	respBody, status, err := c.doRequest("GET", "/api/v1/books", nil)
	if err != nil {
		return nil, err
	}
	if err := c.checkError(respBody, status); err != nil {
		return nil, err
	}
	var books []Book
	if err := json.Unmarshal(respBody, &books); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}
	return books, nil
}

func (c *Client) GetBook(id int) (*BookDetail, error) {
	respBody, status, err := c.doRequest("GET", fmt.Sprintf("/api/v1/books/%d", id), nil)
	if err != nil {
		return nil, err
	}
	if err := c.checkError(respBody, status); err != nil {
		return nil, err
	}
	var book BookDetail
	if err := json.Unmarshal(respBody, &book); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}
	return &book, nil
}

func (c *Client) CreateBook(params CreateBookParams) (*Book, error) {
	body := map[string]interface{}{"book": params}
	respBody, status, err := c.doRequest("POST", "/api/v1/books", body)
	if err != nil {
		return nil, err
	}
	if err := c.checkError(respBody, status); err != nil {
		return nil, err
	}
	var book Book
	if err := json.Unmarshal(respBody, &book); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}
	return &book, nil
}

func (c *Client) UpdateBook(id int, params UpdateBookParams) (*Book, error) {
	body := map[string]interface{}{"book": params}
	respBody, status, err := c.doRequest("PATCH", fmt.Sprintf("/api/v1/books/%d", id), body)
	if err != nil {
		return nil, err
	}
	if err := c.checkError(respBody, status); err != nil {
		return nil, err
	}
	var book Book
	if err := json.Unmarshal(respBody, &book); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}
	return &book, nil
}

func (c *Client) DeleteBook(id int) error {
	respBody, status, err := c.doRequest("DELETE", fmt.Sprintf("/api/v1/books/%d", id), nil)
	if err != nil {
		return err
	}
	if status == 204 {
		return nil
	}
	return c.checkError(respBody, status)
}

// Pages

func (c *Client) CreatePage(bookID int, params CreatePageParams) (*Leaf, error) {
	body := map[string]interface{}{
		"page": map[string]string{"body": params.Body},
		"leaf": map[string]string{"title": params.Title},
	}
	respBody, status, err := c.doRequest("POST", fmt.Sprintf("/api/v1/books/%d/pages", bookID), body)
	if err != nil {
		return nil, err
	}
	if err := c.checkError(respBody, status); err != nil {
		return nil, err
	}
	var leaf Leaf
	if err := json.Unmarshal(respBody, &leaf); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}
	return &leaf, nil
}

func (c *Client) GetPage(bookID, pageID int) (*Leaf, error) {
	respBody, status, err := c.doRequest("GET", fmt.Sprintf("/api/v1/books/%d/pages/%d", bookID, pageID), nil)
	if err != nil {
		return nil, err
	}
	if err := c.checkError(respBody, status); err != nil {
		return nil, err
	}
	var leaf Leaf
	if err := json.Unmarshal(respBody, &leaf); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}
	return &leaf, nil
}

func (c *Client) UpdatePage(bookID, pageID int, params UpdatePageParams) (*Leaf, error) {
	page := map[string]string{}
	leaf := map[string]string{}
	if params.Body != nil {
		page["body"] = *params.Body
	}
	if params.Title != nil {
		leaf["title"] = *params.Title
	}
	body := map[string]interface{}{"page": page, "leaf": leaf}
	respBody, status, err := c.doRequest("PATCH", fmt.Sprintf("/api/v1/books/%d/pages/%d", bookID, pageID), body)
	if err != nil {
		return nil, err
	}
	if err := c.checkError(respBody, status); err != nil {
		return nil, err
	}
	var result Leaf
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}
	return &result, nil
}

func (c *Client) DeletePage(bookID, pageID int) error {
	respBody, status, err := c.doRequest("DELETE", fmt.Sprintf("/api/v1/books/%d/pages/%d", bookID, pageID), nil)
	if err != nil {
		return err
	}
	if status == 204 {
		return nil
	}
	return c.checkError(respBody, status)
}

// Sections

func (c *Client) CreateSection(bookID int, params CreateSectionParams) (*Leaf, error) {
	section := map[string]string{}
	if params.Body != "" {
		section["body"] = params.Body
	}
	if params.Theme != "" {
		section["theme"] = params.Theme
	}
	body := map[string]interface{}{
		"section": section,
		"leaf":    map[string]string{"title": params.Title},
	}
	respBody, status, err := c.doRequest("POST", fmt.Sprintf("/api/v1/books/%d/sections", bookID), body)
	if err != nil {
		return nil, err
	}
	if err := c.checkError(respBody, status); err != nil {
		return nil, err
	}
	var leaf Leaf
	if err := json.Unmarshal(respBody, &leaf); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}
	return &leaf, nil
}

func (c *Client) GetSection(bookID, sectionID int) (*Leaf, error) {
	respBody, status, err := c.doRequest("GET", fmt.Sprintf("/api/v1/books/%d/sections/%d", bookID, sectionID), nil)
	if err != nil {
		return nil, err
	}
	if err := c.checkError(respBody, status); err != nil {
		return nil, err
	}
	var leaf Leaf
	if err := json.Unmarshal(respBody, &leaf); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}
	return &leaf, nil
}

func (c *Client) UpdateSection(bookID, sectionID int, params UpdateSectionParams) (*Leaf, error) {
	section := map[string]string{}
	leaf := map[string]string{}
	if params.Body != nil {
		section["body"] = *params.Body
	}
	if params.Theme != nil {
		section["theme"] = *params.Theme
	}
	if params.Title != nil {
		leaf["title"] = *params.Title
	}
	body := map[string]interface{}{"section": section, "leaf": leaf}
	respBody, status, err := c.doRequest("PATCH", fmt.Sprintf("/api/v1/books/%d/sections/%d", bookID, sectionID), body)
	if err != nil {
		return nil, err
	}
	if err := c.checkError(respBody, status); err != nil {
		return nil, err
	}
	var result Leaf
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}
	return &result, nil
}

func (c *Client) DeleteSection(bookID, sectionID int) error {
	respBody, status, err := c.doRequest("DELETE", fmt.Sprintf("/api/v1/books/%d/sections/%d", bookID, sectionID), nil)
	if err != nil {
		return err
	}
	if status == 204 {
		return nil
	}
	return c.checkError(respBody, status)
}
