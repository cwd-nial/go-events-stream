package main

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
