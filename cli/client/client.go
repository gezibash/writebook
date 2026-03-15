package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
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
	CoverStyle     string `json:"cover_style"`
	CoverSeed      string `json:"cover_seed"`
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
	CoverStyle     string `json:"cover_style,omitempty"`
	CoverSeed      string `json:"cover_seed,omitempty"`
	EveryoneAccess *bool  `json:"everyone_access,omitempty"`
}

type UpdateBookParams struct {
	Title          *string `json:"title,omitempty"`
	Subtitle       *string `json:"subtitle,omitempty"`
	Author         *string `json:"author,omitempty"`
	Theme          *string `json:"theme,omitempty"`
	CoverStyle     *string `json:"cover_style,omitempty"`
	CoverSeed      *string `json:"cover_seed,omitempty"`
	EveryoneAccess *bool   `json:"everyone_access,omitempty"`
	Published      *bool   `json:"published,omitempty"`
}

type CreatePageParams struct {
	Title    string
	Body     string
	Position *int
}

type UpdatePageParams struct {
	Title *string
	Body  *string
}

type CreateSectionParams struct {
	Title    string
	Body     string
	Theme    string
	Position *int
}

type UpdateSectionParams struct {
	Title *string
	Body  *string
	Theme *string
}

type CreatePictureParams struct {
	Title    string
	Caption  string
	Image    string // file path
	Position *int
}

type UpdatePictureParams struct {
	Title   *string
	Caption *string
	Image   *string // file path
}

type SearchResult struct {
	ID             int    `json:"id"`
	Title          string `json:"title"`
	TitleSnippet   string `json:"title_snippet"`
	ContentSnippet string `json:"content_snippet"`
	Type           string `json:"type"`
	BookID         int    `json:"book_id"`
	BookTitle      string `json:"book_title"`
}

type BookSearchGroup struct {
	BookID    int            `json:"book_id"`
	BookTitle string         `json:"book_title"`
	Results   []SearchResult `json:"results"`
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

func (c *Client) Join(joinCode, name, email, password string) (*TokenResponse, error) {
	body := map[string]string{
		"join_code": joinCode,
		"name":      name,
		"email":     email,
		"password":  password,
	}
	respBody, status, err := c.doRequest("POST", "/api/v1/join", body)
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

func (c *Client) ListPages(bookID int) ([]Leaf, error) {
	respBody, status, err := c.doRequest("GET", fmt.Sprintf("/api/v1/books/%d/pages", bookID), nil)
	if err != nil {
		return nil, err
	}
	if err := c.checkError(respBody, status); err != nil {
		return nil, err
	}
	var pages []Leaf
	if err := json.Unmarshal(respBody, &pages); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}
	return pages, nil
}

func (c *Client) CreatePage(bookID int, params CreatePageParams) (*Leaf, error) {
	body := map[string]interface{}{
		"page": map[string]string{"body": params.Body},
		"leaf": map[string]string{"title": params.Title},
	}
	if params.Position != nil {
		body["position"] = *params.Position
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

func (c *Client) ListSections(bookID int) ([]Leaf, error) {
	respBody, status, err := c.doRequest("GET", fmt.Sprintf("/api/v1/books/%d/sections", bookID), nil)
	if err != nil {
		return nil, err
	}
	if err := c.checkError(respBody, status); err != nil {
		return nil, err
	}
	var sections []Leaf
	if err := json.Unmarshal(respBody, &sections); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}
	return sections, nil
}

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
	if params.Position != nil {
		body["position"] = *params.Position
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

// Pictures

func (c *Client) doMultipartRequest(method, path string, fields map[string]string, filePath string, fileField string) ([]byte, int, error) {
	url := c.BaseURL + path

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	for k, v := range fields {
		if err := writer.WriteField(k, v); err != nil {
			return nil, 0, fmt.Errorf("writing field %s: %w", k, err)
		}
	}

	if filePath != "" {
		file, err := os.Open(filePath)
		if err != nil {
			return nil, 0, fmt.Errorf("opening file: %w", err)
		}
		defer file.Close()

		part, err := writer.CreateFormFile(fileField, filepath.Base(filePath))
		if err != nil {
			return nil, 0, fmt.Errorf("creating form file: %w", err)
		}
		if _, err := io.Copy(part, file); err != nil {
			return nil, 0, fmt.Errorf("copying file: %w", err)
		}
	}

	if err := writer.Close(); err != nil {
		return nil, 0, fmt.Errorf("closing multipart writer: %w", err)
	}

	req, err := http.NewRequest(method, url, &body)
	if err != nil {
		return nil, 0, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
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

func (c *Client) ListPictures(bookID int) ([]Leaf, error) {
	respBody, status, err := c.doRequest("GET", fmt.Sprintf("/api/v1/books/%d/pictures", bookID), nil)
	if err != nil {
		return nil, err
	}
	if err := c.checkError(respBody, status); err != nil {
		return nil, err
	}
	var pictures []Leaf
	if err := json.Unmarshal(respBody, &pictures); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}
	return pictures, nil
}

func (c *Client) CreatePicture(bookID int, params CreatePictureParams) (*Leaf, error) {
	fields := map[string]string{
		"leaf[title]":      params.Title,
		"picture[caption]": params.Caption,
	}
	if params.Position != nil {
		fields["position"] = fmt.Sprintf("%d", *params.Position)
	}

	respBody, status, err := c.doMultipartRequest(
		"POST",
		fmt.Sprintf("/api/v1/books/%d/pictures", bookID),
		fields,
		params.Image,
		"picture[image]",
	)
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

func (c *Client) GetPicture(bookID, pictureID int) (*Leaf, error) {
	respBody, status, err := c.doRequest("GET", fmt.Sprintf("/api/v1/books/%d/pictures/%d", bookID, pictureID), nil)
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

func (c *Client) UpdatePicture(bookID, pictureID int, params UpdatePictureParams) (*Leaf, error) {
	fields := map[string]string{}
	if params.Title != nil {
		fields["leaf[title]"] = *params.Title
	}
	if params.Caption != nil {
		fields["picture[caption]"] = *params.Caption
	}

	imagePath := ""
	if params.Image != nil {
		imagePath = *params.Image
	}

	respBody, status, err := c.doMultipartRequest(
		"PATCH",
		fmt.Sprintf("/api/v1/books/%d/pictures/%d", bookID, pictureID),
		fields,
		imagePath,
		"picture[image]",
	)
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

func (c *Client) DeletePicture(bookID, pictureID int) error {
	respBody, status, err := c.doRequest("DELETE", fmt.Sprintf("/api/v1/books/%d/pictures/%d", bookID, pictureID), nil)
	if err != nil {
		return err
	}
	if status == 204 {
		return nil
	}
	return c.checkError(respBody, status)
}

// Search

func (c *Client) SearchGlobal(query string) ([]BookSearchGroup, error) {
	respBody, status, err := c.doRequest("GET", "/api/v1/search?q="+url.QueryEscape(query), nil)
	if err != nil {
		return nil, err
	}
	if err := c.checkError(respBody, status); err != nil {
		return nil, err
	}
	var groups []BookSearchGroup
	if err := json.Unmarshal(respBody, &groups); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}
	return groups, nil
}

func (c *Client) SearchBook(bookID int, query string) ([]SearchResult, error) {
	respBody, status, err := c.doRequest("GET", fmt.Sprintf("/api/v1/books/%d/search?q=%s", bookID, url.QueryEscape(query)), nil)
	if err != nil {
		return nil, err
	}
	if err := c.checkError(respBody, status); err != nil {
		return nil, err
	}
	var results []SearchResult
	if err := json.Unmarshal(respBody, &results); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}
	return results, nil
}
