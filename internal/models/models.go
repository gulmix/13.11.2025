package models

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
