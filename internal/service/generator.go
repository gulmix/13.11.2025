package service

import (
	"bytes"
	"fmt"
	"url-checker/internal/models"

	"github.com/jung-kurt/gofpdf"
)

type Generator struct {
}

func NewGenerator() *Generator {
	return &Generator{}
}

func (g *Generator) GetReport(tasks []*models.Task) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "URL Status Report")
	pdf.Ln(12)

	pdf.SetFont("Arial", "", 12)
	for _, task := range tasks {
		pdf.Cell(40, 10, fmt.Sprintf("Task ID: %d", task.ID))
		pdf.Ln(15)

		for link, status := range task.Results {
			pdf.Cell(0, 10, fmt.Sprintf("  %s - %s", link, status))
			pdf.Ln(12)
		}
		pdf.Ln(10)
	}

	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
