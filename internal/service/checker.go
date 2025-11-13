package service

import (
	"context"
	"net/http"
	"strings"
	"sync"
	"time"
	"url-checker/internal/models"
	"url-checker/internal/repository"
)

type CheckerService struct {
	repo       repository.Repository
	httpClient *http.Client
	taskQueue  chan *models.Task
	wg         sync.WaitGroup
	ctx        context.Context
	cancel     context.CancelFunc
}

func NewCheckerService(repo repository.Repository) *CheckerService {
	ctx, cancel := context.WithCancel(context.Background())

	service := &CheckerService{
		repo: repo,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
		taskQueue: make(chan *models.Task, 100),
		ctx:       ctx,
		cancel:    cancel,
	}

	for i := 0; i < 5; i++ {
		service.wg.Add(1)
		go service.worker()
	}

	service.recoverTasks()

	return service
}

func (s *CheckerService) worker() {
	defer s.wg.Done()

	for {
		select {
		case <-s.ctx.Done():
			return
		case task := <-s.taskQueue:
			s.processTask(task)
		}
	}
}

func (s *CheckerService) recoverTasks() {
	tasks := s.repo.GetAllProcessingTasks()
	for _, task := range tasks {
		select {
		case s.taskQueue <- task:
		case <-s.ctx.Done():
			return
		}
	}
}

func (s *CheckerService) CheckLinks(links []string) (*models.Task, error) {
	task := &models.Task{
		Links:     links,
		Results:   make(map[string]string),
		CreatedAt: time.Now(),
		Status:    "processing",
	}

	if err := s.repo.SaveTask(task); err != nil {
		return nil, err
	}

	select {
	case s.taskQueue <- task:
	case <-s.ctx.Done():
		return task, nil
	default:
	}

	return task, nil
}

func (s *CheckerService) processTask(task *models.Task) {
	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, link := range task.Links {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()

			status := s.checkURL(url)

			mu.Lock()
			task.Results[url] = status
			mu.Unlock()
		}(link)
	}

	wg.Wait()

	task.Status = "completed"
	s.repo.UpdateTask(task)
}

func (s *CheckerService) checkURL(url string) string {
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = "http://" + url
	}

	req, err := http.NewRequestWithContext(s.ctx, "GET", url, nil)
	if err != nil {
		return "not available"
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "not available"
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 400 {
		return "available"
	}

	return "not available"
}

func (s *CheckerService) Shutdown(timeout time.Duration) error {
	s.cancel()

	done := make(chan struct{})
	go func() {
		s.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-time.After(timeout):
		return nil
	}
}

func (s *CheckerService) GetTask(id int) (*models.Task, bool) {
	return s.repo.GetTask(id)
}

func (s *CheckerService) GetTasks(ids []int) []*models.Task {
	return s.repo.GetTasks(ids)
}
