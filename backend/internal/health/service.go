package health

import (
	"context"
	"sync"
)

type Status string

const (
	StatusUp   Status = "UP"
	StatusDown Status = "DOWN"
)

type ComponentStatus struct {
	Status Status `json:"status"`
	Error  string `json:"error,omitempty"`
}

type Result struct {
	Status     Status                     `json:"status"`
	Components map[string]ComponentStatus `json:"components"`
}

// Service định nghĩa nghiệp vụ kiểm tra sức khỏe hệ thống (DIP).
type Service interface {
	Check(ctx context.Context) Result
}

type healthService struct {
	checkers []Checker
}

// NewService nhận danh sách các Checker đa hình (OCP, DIP).
func NewService(checkers []Checker) Service {
	return &healthService{
		checkers: checkers,
	}
}

func (s *healthService) Check(ctx context.Context) Result {
	res := Result{
		Status:     StatusUp,
		Components: make(map[string]ComponentStatus, len(s.checkers)),
	}

	var (
		wg sync.WaitGroup
		mu sync.Mutex
	)

	for _, checker := range s.checkers {
		wg.Add(1)
		go func(c Checker) {
			defer wg.Done()

			err := c.Check(ctx)
			compStatus := ComponentStatus{
				Status: StatusUp,
			}

			if err != nil {
				compStatus.Status = StatusDown
				compStatus.Error = err.Error()
			}

			mu.Lock()
			res.Components[c.Name()] = compStatus
			if compStatus.Status == StatusDown {
				res.Status = StatusDown
			}
			mu.Unlock()
		}(checker)
	}

	wg.Wait()
	return res
}
