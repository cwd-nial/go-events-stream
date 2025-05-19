TASK
====

Using the following types as a starting point, create a simple golang application that listens to the events stream and
updates the user wallet based on the following rules:

Add 1 point for adding a product to the wishlist;
Add 10 points for purchasing a product;
Add mock implementations for testing. Don’t explain what you are doing, simply output the code.
---

Interfaces, structs, helper functions:

```
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
    UserID uint `json:"userId"`
    ProductID uint `json:"productId"`
}

func (p ProductAddedToWishlist) Name() string {
    return "product.addedToWishlist"
}

type ProductPurchased struct {
    UserID uint `json:"userId"`
    ProductID uint `json:"productId"`
}

func (p ProductPurchased) Name() string {
    return "product.purchased"
}
```