package product

type Product struct {
	ID   int
	Name string
}

type ProductCreate struct {
	Name string
}

type ProductCreateRequest struct {
	Name string
}
