package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const (
	AddedToWishlistType EventType = "added_to_wishlist"
	PurchasedType       EventType = "purchased"
)

type EventType string

type event struct {
	UserID    uint      `json:"userId"`
	ProductID uint      `json:"productId"`
	Type      EventType `json:"type"`
}

type UserWalletService interface {
	AddBonusPoints(userId uint, points int) error
}

type EventsRepository interface {
	GetEventsStream() <-chan event
}

type WalletService struct {
	mu     sync.Mutex
	points map[uint]int
}

func NewWalletService() *WalletService {
	return &WalletService{points: make(map[uint]int)}
}

func (ws *WalletService) AddBonusPoints(userId uint, points int) error {
	ws.mu.Lock()
	defer ws.mu.Unlock()
	ws.points[userId] += points
	fmt.Printf("User %d: added %d points, total: %d\n", userId, points, ws.points[userId])
	return nil
}

type MockEventsRepo struct{}

func (r *MockEventsRepo) GetEventsStream() <-chan event {
	ch := make(chan event)

	go func() {
		var wg sync.WaitGroup
		wg.Add(3)

		send := func(e event) {
			defer wg.Done()
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
			ch <- e
		}

		go send(event{UserID: 1, ProductID: 101, Type: AddedToWishlistType})
		go send(event{UserID: 1, ProductID: 102, Type: PurchasedType})
		go send(event{UserID: 2, ProductID: 103, Type: PurchasedType})

		wg.Wait()
		close(ch)
	}()

	return ch
}

func processEvents(repo EventsRepository, wallet UserWalletService) {
	for e := range repo.GetEventsStream() {
		switch e.Type {
		case AddedToWishlistType:
			_ = wallet.AddBonusPoints(e.UserID, 1)
		case PurchasedType:
			_ = wallet.AddBonusPoints(e.UserID, 10)
		}
	}
}

func main() {
	repo := &MockEventsRepo{}
	wallet := NewWalletService()

	processEvents(repo, wallet)
}
