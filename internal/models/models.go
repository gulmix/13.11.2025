package models

import "time"

type CheckRequest struct {
	Links []string `json:"links"`
}

type CheckResponse struct {
	Links   map[string]string `json:"links"`
	LinksID int               `json:"links_num"`
}

type ReportRequest struct {
	LinksNum []int `json:"links_num"`
}

type Task struct {
	ID        int               `json:"id"`
	Links     []string          `json:"links"`
	Results   map[string]string `json:"results"`
	CreatedAt time.Time         `json:"created_at"`
	Status    string            `json:"status"`
}
