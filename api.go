package gostrawpoll

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

// GetRequest is a structure for creating a GET request.
type GetRequest struct {
	ID int
}

// GetResponse is a structure for returning a GET response.
type GetResponse struct {
	ID       int      `json:"id"`
	Title    string   `json:"title"`
	Options  []string `json:"options"`
	Votes    []int    `json:"votes"`
	Multi    bool     `json:"multi"`
	Dupcheck string   `json:"dupcheck"`
	Captcha  bool     `json:"captcha"`
}

// PostRequest is a structure for creating a POST request.
type PostRequest struct {
	Title    string   `json:"title"`   // REQUIRED
	Options  []string `json:"options"` // REQUIRED
	Multi    bool     `json:"multi"`
	Dupcheck string   `json:"dupcheck"`
	Captcha  bool     `json:"captcha"`
}

// PostResponse is the POST response from the server.
type PostResponse struct {
	ID       int      `json:"id"`
	Title    string   `json:"title"`
	Options  []string `json:"options"`
	Multi    bool     `json:"multi"`
	Dupcheck string   `json:"dupcheck"`
	Captcha  bool     `json:"captcha"`
}

const endpoint = "https://strawpoll.me/api/v2/polls"

func Get(req *GetRequest) (*GetResponse, error) {
	if req == nil {
		return nil, errors.New("req must not be nil")
	}
	resp, err := http.Get(fmt.Sprintf("%s/%d", endpoint, req.ID))
	defer resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to send req: err")
	} else if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("req failed: %v", resp.Status)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read response body: %v", err)
	}
	ret := &GetResponse{}
	if err := json.Unmarshal(body, ret); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %v", err)
	}
	return ret, nil
}

func Post(req *PostRequest) (*PostResponse, error) {
	if req == nil {
		return nil, errors.New("req must not be nil")
	}
	if req.Title == "" {
		return nil, errors.New("Title must not be empty")
	}
	if len(req.Options) == 0 {
		return nil, errors.New("options must not be empty")
	}
	b, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("couldn't marshal json req: %v", err)
	}
	hReq, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(b))
	if err != nil {
		return nil, fmt.Errorf("failed to create req: %v", err)
	}
	resp, err := http.DefaultClient.Do(hReq)
	defer resp.Body.Close()
	if err != nil {
		fmt.Errorf("failed to send req: %v", err)
	} else if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("req failed: %v", resp.Status)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read response body: %v", err)
	}
	ret := &PostResponse{}
	if err := json.Unmarshal(body, ret); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %v", err)
	}
	return ret, nil
}
