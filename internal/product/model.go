package product

type Product struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
}

type CreateProductRequest struct {
	Name        string  `json: "name"`
	Description string  `json: "description"`
	Price       float64 `json: "price"`
	Stock       int     `json: "stock"`
}

type UpdateProductRequest struct {
	Name        string  `json: "name"`
	Description string  `json: "description"`
	Price       float64 `json: "price"`
	Stock       int     `json: "stock"`
}
