package crawler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

const (
	ghBaseURL = "https://api.github.com"
)

type githubClient struct {
	token string
}

func (c githubClient) get(url string, modifiers ...func(req *http.Request)) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	for i := range modifiers {
		modifiers[i](req)
	}

	req.Header.Add("Authorization", fmt.Sprintf("token %s", c.token))
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if res.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("error request at %s", url))
	}

	buf, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return buf, nil
}

func (c githubClient) getRepository(path string) (object, error) {
	var o object

	buf, err := c.get(fmt.Sprintf("%s/repos/%s", ghBaseURL, path))
	if err != nil {
		return o, err
	}

	if err := json.Unmarshal(buf, &o); err != nil {
		return o, errors.WithStack(err)
	}

	return o, nil
}

func (c githubClient) getRepositoryStargazer(path string, page int) ([]object, error) {
	buf, err := c.get(
		fmt.Sprintf("%s/repos/%s/stargazers?page=%d&per_page=100", ghBaseURL, path, page),
		func(req *http.Request) { req.Header.Add("Accept", "application/vnd.github.v3.star+json") },
	)
	if err != nil {
		return nil, err
	}

	var os []object
	if err := json.Unmarshal(buf, &os); err != nil {
		return nil, errors.WithStack(err)
	}

	return os, nil
}

func (c githubClient) getUser(login string) (object, error) {
	var o object

	buf, err := c.get(fmt.Sprintf("%s/users/%s", ghBaseURL, login))
	if err != nil {
		return o, err
	}

	if err := json.Unmarshal(buf, &o); err != nil {
		return o, errors.WithStack(err)
	}

	return o, nil
}
