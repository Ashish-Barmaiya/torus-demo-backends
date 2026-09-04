package data

import "fmt"

func Users() []User {
	users := make([]User, 20)

	for i := range users {
		id := i + 1

		users[i] = User{
			ID:     fmt.Sprintf("usr_%06d", id),
			Name:   fmt.Sprintf("Demo User %d", id),
			Email:  fmt.Sprintf("user%d@example.com", id),
			Plan:   []string{"free", "pro", "enterprise"}[i%3],
			Status: "active",
		}
	}

	return users
}

func UserByID(id int) (User, bool) {
	users := Users()

	if id < 1 || id > len(users) {
		return User{}, false
	}

	return users[id-1], true
}

func Orders() []Order {
	orders := make([]Order, 20)

	statuses := []string{
		"pending",
		"processing",
		"shipped",
		"completed",
	}

	for i := range orders {
		id := i + 1

		orders[i] = Order{
			ID:         fmt.Sprintf("ord_%06d", id),
			CustomerID: fmt.Sprintf("usr_%06d", ((id-1)%20)+1),
			Status:     statuses[i%len(statuses)],
			Currency:   "USD",
			Total:      int64(4999 + (i * 1250)),
		}
	}

	return orders
}

func OrderByID(id int) (Order, bool) {
	orders := Orders()

	if id < 1 || id > len(orders) {
		return Order{}, false
	}

	return orders[id-1], true
}
