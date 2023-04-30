package inventory

import (
	"golang.org/x/net/context"
	"sync"
)

type ServerInventory struct {
	sum            int64
	inventoryCount int64
	mutex          sync.Mutex
}

func (s ServerInventory) UpdateProductCount(ctx context.Context, info *ProductInfo) (*NumberResult, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.sum < s.inventoryCount {
		s.sum += 1
		return &NumberResult{IsInventorySuccess: true}, nil
	}
	return &NumberResult{IsInventorySuccess: false}, nil
}

func (s ServerInventory) mustEmbedUnimplementedInventoryServiceServer() {
	panic("implement me")
}

func NewServerInventory(productCount int64) *ServerInventory {
	return &ServerInventory{
		sum:            0,
		inventoryCount: productCount,
	}
}
