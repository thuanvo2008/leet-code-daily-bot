package client

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/thuanvo2008/leet-code-daily-bot/model"
)

// LeetCodeClient handles communication with LeetCode API
type LeetCodeClient struct {
	requestSender RequestSender
	baseURL       string
}

// NewLeetCodeClient creates a new client to interact with LeetCode
func NewLeetCodeClient(client *http.Client) *LeetCodeClient {
	return &LeetCodeClient{
		requestSender: NewHTTPRequestClient(client),
		baseURL:       "https://leetcode.com",
	}
}

// FetchDailyProblem retrieves the daily coding challenge from LeetCode
func (c *LeetCodeClient) FetchDailyProblem() (*model.LeetCodeResponse, error) {
	resp, err := c.makeGraphQLRequest()
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}

	lcResp, err := c.parseResponse(resp)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	c.enhanceResponse(lcResp)
	return lcResp, nil
}

// makeGraphQLRequest sends a GraphQL request to LeetCode API
func (c *LeetCodeClient) makeGraphQLRequest() ([]byte, error) {
	query := map[string]string{
		"query": `query questionOfToday {
			activeDailyCodingChallengeQuestion {
				date
				link
				question {
					title
					titleSlug
					difficulty
					content
					exampleTestcases
				}
			}
		}`,
	}

	jsonData, err := json.Marshal(query)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal query: %w", err)
	}

	req, err := CreateRequest("POST", c.baseURL+"/graphql", jsonData)
	if err != nil {
		return nil, err
	}

	c.setRequestHeaders(req)

	return c.requestSender.SendRequest(req)
}

// setRequestHeaders sets common headers for LeetCode API requests
func (c *LeetCodeClient) setRequestHeaders(req *http.Request) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Referer", c.baseURL)
	//req.Header.Set("User-Agent", "GoLeetBot/1.0")
}

// parseResponse parses the JSON response into a LeetCodeResponse object
func (c *LeetCodeClient) parseResponse(data []byte) (*model.LeetCodeResponse, error) {
	var lcResp model.LeetCodeResponse
	if err := json.Unmarshal(data, &lcResp); err != nil {
		return nil, err
	}
	return &lcResp, nil
}

// enhanceResponse adds additional information to the response
func (c *LeetCodeClient) enhanceResponse(resp *model.LeetCodeResponse) {
	titleSlug := resp.Data.ActiveDailyCodingChallengeQuestion.Question.TitleSlug
	resp.Data.ActiveDailyCodingChallengeQuestion.Question.Url =
		fmt.Sprintf("%s/problems/%s/", c.baseURL, titleSlug)
}
