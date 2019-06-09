package publish

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Git interface {
	AddFile(path string, content string)
	Commit(branch string, message string)
}

type GitLabRequest struct {
	Branch        string         `json:"branch"`
	CommitMessage string         `json:"commit_message"`
	Actions       []GitLabAction `json:"actions"`
}

type GitLabAction struct {
	Action  string `json:"action"`
	Path    string `json:"file_path"`
	Content string `json:"content"`
}

type GitLabResponse struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Committer string    `json:"committer_name"`
	CreatedAt time.Time `json:"created_at"`
}

func NewGitLab(projectID string) GitLab {
	return GitLab{
		ProjectID: projectID,
		Request: GitLabRequest{
			Actions: []GitLabAction{},
		},
	}
}

type GitLab struct {
	ProjectID     string
	PersonalToken string
	Request       GitLabRequest
}

func (g *GitLab) Auth(personalToken string) {
	g.PersonalToken = personalToken
}

func (g *GitLab) AddFile(path string, content string) {
	g.Request.Actions = append(g.Request.Actions, GitLabAction{
		Action:  "create",
		Path:    path,
		Content: content,
	})
}

func (g *GitLab) Commit(branch string, message string) *GitLabResponse {
	url := fmt.Sprintf("https://gitlab.com/api/v4/projects/%v/repository/commits", g.ProjectID)

	g.Request.Branch = branch
	g.Request.CommitMessage = message

	body, _ := json.Marshal(g.Request)

	timeout := time.Duration(60 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	request, _ := http.NewRequest("POST", url, bytes.NewBuffer(body))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Private-Token", g.PersonalToken)

	response, err := client.Do(request)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	defer response.Body.Close()

	if response.StatusCode > 201 {
		responseData, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Println(err)
			return nil
		}

		fmt.Println(string(responseData))
	}

	var content GitLabResponse
	json.NewDecoder(response.Body).Decode(&content)

	return &content
}

/*
type GitHub struct {
}

func (g *GitHub) Auth(username string, personalToken string) {

}

func (g *GitHub) AddFile(path string, content string) {
	// base64 encoding!
}

func (g *GitHub) Commit(branch string, message string) {
	// PUT
	url := fmt.Sprintf("/repos/%v/%v/contents/%v", g.Owner, g.Repository, path)

	body, _ := json.Marshal(g.Request)

	timeout := time.Duration(60 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	request, _ := http.NewRequest("POST", url, bytes.NewBuffer(body))
	request.header.Set("Content-Type", "application/json")
	request.header.Set("Authorization", g.Token)

	response, err := client.Do(request)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer response.Body.Close()
}
*/
