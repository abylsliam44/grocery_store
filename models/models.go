package models

type Product struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Stock    int     `json:"stock"`
	ImageURL string  `json:"image_url"`
	Category string  `json:"category"`
}

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CartItem struct {
	ProductID   int
	ProductName string
	Price       float64
	Quantity    int
	Total       float64
}

type Category struct {
	ID       int
	Name     string
	Products []Product
}

type Review struct {
	Rating int
	Review string
	User   User
}
