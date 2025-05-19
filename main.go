package main

import (
	"fmt"
	"time"
)

// Interfaces and event types

type UserWalletService interface {
	AddBonusPoints(userId uint, points int) error
}

type Event interface {
	Name() string
}

type EventsRepository interface {
	GetEventsStream() <-chan Event
}

type ProductAddedToWishlist struct {
	UserID    uint `json:"userId"`
	ProductID uint `json:"productId"`
}

func (p ProductAddedToWishlist) Name() string {
	return "product.addedToWishlist"
}

type ProductPurchased struct {
	UserID    uint `json:"userId"`
	ProductID uint `json:"productId"`
}

func (p ProductPurchased) Name() string {
	return "product.purchased"
}

// App implementation

type App struct {
	walletService UserWalletService
	eventRepo     EventsRepository
}

func NewApp(walletService UserWalletService, eventRepo EventsRepository) *App {
	return &App{
		walletService: walletService,
		eventRepo:     eventRepo,
	}
}

func (a *App) Run() {
	for event := range a.eventRepo.GetEventsStream() {
		switch e := event.(type) {
		case ProductAddedToWishlist:
			_ = a.walletService.AddBonusPoints(e.UserID, 1)
		case ProductPurchased:
			_ = a.walletService.AddBonusPoints(e.UserID, 10)
		}
	}
}

// Mock implementations for testing

type MockWalletService struct{}

func (m *MockWalletService) AddBonusPoints(userId uint, points int) error {
	fmt.Printf("User %d awarded %d points\n", userId, points)
	return nil
}

type MockEventsRepository struct {
	stream chan Event
}

func NewMockEventsRepository() *MockEventsRepository {
	ch := make(chan Event)
	go func() {
		defer close(ch)
		ch <- ProductAddedToWishlist{UserID: 1, ProductID: 101}
		time.Sleep(100 * time.Millisecond)
		ch <- ProductPurchased{UserID: 1, ProductID: 101}
		time.Sleep(100 * time.Millisecond)
		ch <- ProductAddedToWishlist{UserID: 2, ProductID: 202}
		time.Sleep(100 * time.Millisecond)
	}()
	return &MockEventsRepository{stream: ch}
}

func (m *MockEventsRepository) GetEventsStream() <-chan Event {
	return m.stream
}

// Main

func main() {
	wallet := &MockWalletService{}
	events := NewMockEventsRepository()
	app := NewApp(wallet, events)
	app.Run()
}
