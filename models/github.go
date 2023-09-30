package models

type GitHubResponse struct {
	ID       int `json:"id"`
	Email    string `json:"email"`
  Name    string `json:"name"`
	AvatarURL string `json:"avatar_url"`
}

