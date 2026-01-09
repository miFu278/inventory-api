package dtos

type CreateProductRequest struct {
	Body CreateProductInput
}

type UpdateProductRequest struct {
	ID   uint `path:"id"`
	Body UpdateProductInput
}

type CreateTransactionRequest struct {
	Body CreateTransactionInput
}

type IDParam struct {
	ID uint `path:"id"`
}

type PaginationQuery struct {
	Limit  int `query:"limit" default:"10" minimum:"1" maximum:"100"`
	Offset int `query:"offset" default:"0" minimum:"0"`
}

type ProductListQuery struct {
	SKU      string  `query:"sku" doc:"Filter by SKU (exact match)"`
	Name     string  `query:"name" doc:"Filter by name (partial match)"`
	MinPrice float64 `query:"min_price" doc:"Filter by minimum price"`
	MaxPrice float64 `query:"max_price" doc:"Filter by maximum price"`
	Limit    int     `query:"limit" default:"10" minimum:"1" maximum:"100"`
	Offset   int     `query:"offset" default:"0" minimum:"0"`
}

func (q *ProductListQuery) ToProductFilter() *ProductFilter {
	filter := &ProductFilter{}

	// Only set pointer if value is not empty/zero
	if q.SKU != "" {
		filter.SKU = &q.SKU
	}
	if q.Name != "" {
		filter.Name = &q.Name
	}
	if q.MinPrice > 0 {
		filter.MinPrice = &q.MinPrice
	}
	if q.MaxPrice > 0 {
		filter.MaxPrice = &q.MaxPrice
	}

	return filter
}

type ProductTransactionsQuery struct {
	IDParam
	PaginationQuery
}
