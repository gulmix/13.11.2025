package service

import (
	"bytes"
	"context"
	"fmt"
	"url-checker/internal/models"
	"url-checker/internal/repository"
	"url-checker/internal/utils"

	"github.com/jung-kurt/gofpdf"
)

type Service struct {
	repo repository.Repository
}

func NewService(repo repository.Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) CheckLinks(ctx context.Context, req models.CheckRequest) (models.CheckResponse, error) {
	linksID, err := s.repo.SaveProcessing(ctx, req.Links)
	if err != nil {
		return models.CheckResponse{}, err
	}

	res := make(map[string]string)
	for _, link := range req.Links {
		res[link] = utils.CheckLink(link)
	}

	s.repo.SaveStorage(linksID, res)

	return models.CheckResponse{Links: res, LinksID: linksID}, nil
}

func (s *Service) GetReport(ctx context.Context, req models.ReportRequest) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "URL Status Report")
	pdf.Ln(12)

	pdf.SetFont("Arial", "", 12)
	for _, id := range req.LinksNum {
		links, res, err := s.repo.Get(ctx, id)
		fmt.Println(links, res, err)
		if err != nil {
			continue
		}
		pdf.SetFont("Arial", "B", 12)
		pdf.Cell(40, 10, fmt.Sprintf("ID: %d", id))
		pdf.Ln(10)
		pdf.SetFont("Arial", "", 12)
		for _, link := range links {
			status := res[link]

			pdf.Cell(0, 10, fmt.Sprintf("%s: %s", link, status))
			pdf.Ln(8)
		}
		pdf.Ln(4)
	}

	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
