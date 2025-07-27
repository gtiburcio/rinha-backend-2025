package repository

import (
	"log"
	"rinha-backend-2025-gtiburcio/src/model"
	"sync"
	"time"
)

type ConcurrentSlice struct {
	sync.Mutex
	Data       []model.PaymentSummaryDTO
	repository Repository
}

func NewConcurrentSlice(repository Repository) *ConcurrentSlice {
	return &ConcurrentSlice{
		repository: repository,
	}
}

func (c *ConcurrentSlice) Start() {
	go c.refresh()
}

func (c *ConcurrentSlice) refresh() {
	for {
		list, err := c.repository.FindAll()
		if err != nil {
			log.Fatalf("find all failed %v", err)
		}
		c.setRefresh(list)
		time.Sleep(50 * time.Millisecond)
	}
}

func (c *ConcurrentSlice) setRefresh(list []model.PaymentSummaryDTO) {
	c.Lock()
	defer c.Unlock()
	c.Data = list
}

func (c *ConcurrentSlice) GetData() []model.PaymentSummaryDTO {
	c.Lock()
	defer c.Unlock()

	results := make([]model.PaymentSummaryDTO, len(c.Data))

	copy(results, c.Data)

	return results
}
