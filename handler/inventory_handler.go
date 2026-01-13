package handler

import (
	"context"
	"inventory-api/dtos"
	"inventory-api/middleware"
	"inventory-api/services"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

type InventoryHandler struct {
	service *services.InventoryService
}

func NewInventoryHandler(service *services.InventoryService) *InventoryHandler {
	return &InventoryHandler{service: service}
}

func (h *InventoryHandler) RegisterRoutes(api huma.API) {
	// Product routes - require authentication
	huma.Register(api, huma.Operation{
		OperationID: "create-product",
		Method:      http.MethodPost,
		Path:        "/products",
		Summary:     "Create a new product",
		Tags:        []string{"Products"},
		Security: []map[string][]string{
			{"bearerAuth": {}},
		},
	}, h.CreateProduct)

	huma.Register(api, huma.Operation{
		OperationID: "get-product",
		Method:      http.MethodGet,
		Path:        "/products/{id}",
		Summary:     "Get product by ID",
		Tags:        []string{"Products"},
	}, h.GetProduct)

	huma.Register(api, huma.Operation{
		OperationID: "list-products",
		Method:      http.MethodGet,
		Path:        "/products",
		Summary:     "List all products",
		Tags:        []string{"Products"},
	}, h.ListProducts)

	huma.Register(api, huma.Operation{
		OperationID: "update-product",
		Method:      http.MethodPut,
		Path:        "/products/{id}",
		Summary:     "Update product",
		Tags:        []string{"Products"},
		Security: []map[string][]string{
			{"bearerAuth": {}},
		},
	}, h.UpdateProduct)

	huma.Register(api, huma.Operation{
		OperationID: "delete-product",
		Method:      http.MethodDelete,
		Path:        "/products/{id}",
		Summary:     "Delete product",
		Tags:        []string{"Products"},
		Security: []map[string][]string{
			{"bearerAuth": {}},
		},
	}, h.DeleteProduct)

	// Transaction routes - require authentication
	huma.Register(api, huma.Operation{
		OperationID: "create-transaction",
		Method:      http.MethodPost,
		Path:        "/transactions",
		Summary:     "Create a new transaction",
		Tags:        []string{"Transactions"},
		Security: []map[string][]string{
			{"bearerAuth": {}},
		},
	}, h.CreateTransaction)

	huma.Register(api, huma.Operation{
		OperationID: "get-transaction",
		Method:      http.MethodGet,
		Path:        "/transactions/{id}",
		Summary:     "Get transaction by ID",
		Tags:        []string{"Transactions"},
		Security: []map[string][]string{
			{"bearerAuth": {}},
		},
	}, h.GetTransaction)

	huma.Register(api, huma.Operation{
		OperationID: "list-transactions",
		Method:      http.MethodGet,
		Path:        "/transactions",
		Summary:     "List all transactions",
		Tags:        []string{"Transactions"},
		Security: []map[string][]string{
			{"bearerAuth": {}},
		},
	}, h.ListTransactions)

	huma.Register(api, huma.Operation{
		OperationID: "list-product-transactions",
		Method:      http.MethodGet,
		Path:        "/products/{id}/transactions",
		Summary:     "List transactions for a product",
		Tags:        []string{"Transactions"},
		Security: []map[string][]string{
			{"bearerAuth": {}},
		},
	}, h.ListProductTransactions)
}

func (h *InventoryHandler) CreateProduct(ctx context.Context, input *dtos.CreateProductRequest) (*dtos.SingleProductResponse, error) {
	// Verify authentication
	auth := middleware.GetAuthContext(ctx)
	if auth == nil {
		return nil, huma.Error401Unauthorized("Authentication required")
	}

	product, err := h.service.CreateProduct(&input.Body)
	if err != nil {
		return nil, huma.Error400BadRequest(err.Error())
	}
	return &dtos.SingleProductResponse{Body: product}, nil
}

func (h *InventoryHandler) GetProduct(ctx context.Context, input *dtos.IDParam) (*dtos.SingleProductResponse, error) {
	product, err := h.service.GetProductByID(input.ID)
	if err != nil {
		return nil, huma.Error404NotFound(err.Error())
	}
	return &dtos.SingleProductResponse{Body: product}, nil
}

func (h *InventoryHandler) ListProducts(ctx context.Context, input *dtos.ProductListQuery) (*dtos.ProductListResponse, error) {
	// Convert query to filter
	filter := input.ToProductFilter()

	products, err := h.service.GetProductsWithFilter(filter, input.Limit, input.Offset)
	if err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}

	resp := &dtos.ProductListResponse{}
	resp.Body.Products = products
	resp.Body.Limit = input.Limit
	resp.Body.Offset = input.Offset
	return resp, nil
}

func (h *InventoryHandler) UpdateProduct(ctx context.Context, input *dtos.UpdateProductRequest) (*dtos.SingleProductResponse, error) {
	// Verify authentication
	auth := middleware.GetAuthContext(ctx)
	if auth == nil {
		return nil, huma.Error401Unauthorized("Authentication required")
	}

	product, err := h.service.UpdateProduct(input.ID, &input.Body)
	if err != nil {
		return nil, huma.Error400BadRequest(err.Error())
	}
	return &dtos.SingleProductResponse{Body: product}, nil
}

func (h *InventoryHandler) DeleteProduct(ctx context.Context, input *dtos.IDParam) (*dtos.EmptyResponse, error) {
	// Only admins can delete products
	if !middleware.IsAdmin(ctx) {
		return nil, huma.Error403Forbidden("Only admins can delete products")
	}

	err := h.service.DeleteProduct(input.ID)
	if err != nil {
		return nil, huma.Error404NotFound(err.Error())
	}
	return &dtos.EmptyResponse{}, nil
}

func (h *InventoryHandler) CreateTransaction(ctx context.Context, input *dtos.CreateTransactionRequest) (*dtos.SingleTransactionResponse, error) {
	// Verify authentication
	auth := middleware.GetAuthContext(ctx)
	if auth == nil {
		return nil, huma.Error401Unauthorized("Authentication required")
	}

	transaction, err := h.service.CreateTransaction(&input.Body)
	if err != nil {
		return nil, huma.Error400BadRequest(err.Error())
	}
	return &dtos.SingleTransactionResponse{Body: transaction}, nil
}

func (h *InventoryHandler) GetTransaction(ctx context.Context, input *dtos.IDParam) (*dtos.SingleTransactionResponse, error) {
	// Verify authentication
	auth := middleware.GetAuthContext(ctx)
	if auth == nil {
		return nil, huma.Error401Unauthorized("Authentication required")
	}

	transaction, err := h.service.GetTransactionByID(input.ID)
	if err != nil {
		return nil, huma.Error404NotFound(err.Error())
	}
	return &dtos.SingleTransactionResponse{Body: transaction}, nil
}

func (h *InventoryHandler) ListTransactions(ctx context.Context, input *dtos.PaginationQuery) (*dtos.TransactionListResponse, error) {
	// Verify authentication
	auth := middleware.GetAuthContext(ctx)
	if auth == nil {
		return nil, huma.Error401Unauthorized("Authentication required")
	}

	transactions, err := h.service.GetAllTransactions(input.Limit, input.Offset)
	if err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}
	resp := &dtos.TransactionListResponse{}
	resp.Body.Transactions = transactions
	resp.Body.Limit = input.Limit
	resp.Body.Offset = input.Offset
	return resp, nil
}

func (h *InventoryHandler) ListProductTransactions(ctx context.Context, input *dtos.ProductTransactionsQuery) (*dtos.TransactionListResponse, error) {
	// Verify authentication
	auth := middleware.GetAuthContext(ctx)
	if auth == nil {
		return nil, huma.Error401Unauthorized("Authentication required")
	}

	transactions, err := h.service.GetTransactionsByProductID(input.ID, input.Limit, input.Offset)
	if err != nil {
		return nil, huma.Error404NotFound(err.Error())
	}
	resp := &dtos.TransactionListResponse{}
	resp.Body.Transactions = transactions
	resp.Body.Limit = input.Limit
	resp.Body.Offset = input.Offset
	return resp, nil
}
