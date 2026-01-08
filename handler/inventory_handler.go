package handler

import (
	"context"
	"inventory-api/models"
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

type ProductResponse struct {
	Body *models.Product
}

type ProductListResponse struct {
	Body struct {
		Products []models.Product `json:"products"`
		Limit    int              `json:"limit"`
		Offset   int              `json:"offset"`
	}
}

type TransactionResponse struct {
	Body *models.Transaction
}

type TransactionListResponse struct {
	Body struct {
		Transactions []models.Transaction `json:"transactions"`
		Limit        int                  `json:"limit"`
		Offset       int                  `json:"offset"`
	}
}

type CreateProductRequest struct {
	Body models.CreateProductInput
}

type UpdateProductRequest struct {
	ID   uint `path:"id"`
	Body models.UpdateProductInput
}

type CreateTransactionRequest struct {
	Body models.CreateTransactionInput
}

type IDParam struct {
	ID uint `path:"id"`
}

type PaginationQuery struct {
	Limit  int `query:"limit" default:"10" minimum:"1" maximum:"100"`
	Offset int `query:"offset" default:"0" minimum:"0"`
}

func (h *InventoryHandler) RegisterRoutes(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "create-product",
		Method:      http.MethodPost,
		Path:        "/products",
		Summary:     "Create a new product",
		Tags:        []string{"Products"},
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
	}, h.UpdateProduct)

	huma.Register(api, huma.Operation{
		OperationID: "delete-product",
		Method:      http.MethodDelete,
		Path:        "/products/{id}",
		Summary:     "Delete product",
		Tags:        []string{"Products"},
	}, h.DeleteProduct)

	huma.Register(api, huma.Operation{
		OperationID: "create-transaction",
		Method:      http.MethodPost,
		Path:        "/transactions",
		Summary:     "Create a new transaction",
		Tags:        []string{"Transactions"},
	}, h.CreateTransaction)

	huma.Register(api, huma.Operation{
		OperationID: "get-transaction",
		Method:      http.MethodGet,
		Path:        "/transactions/{id}",
		Summary:     "Get transaction by ID",
		Tags:        []string{"Transactions"},
	}, h.GetTransaction)

	huma.Register(api, huma.Operation{
		OperationID: "list-transactions",
		Method:      http.MethodGet,
		Path:        "/transactions",
		Summary:     "List all transactions",
		Tags:        []string{"Transactions"},
	}, h.ListTransactions)

	huma.Register(api, huma.Operation{
		OperationID: "list-product-transactions",
		Method:      http.MethodGet,
		Path:        "/products/{id}/transactions",
		Summary:     "List transactions for a product",
		Tags:        []string{"Transactions"},
	}, h.ListProductTransactions)
}

func (h *InventoryHandler) CreateProduct(ctx context.Context, input *CreateProductRequest) (*ProductResponse, error) {
	product, err := h.service.CreateProduct(&input.Body)
	if err != nil {
		return nil, huma.Error400BadRequest(err.Error())
	}
	return &ProductResponse{Body: product}, nil
}

func (h *InventoryHandler) GetProduct(ctx context.Context, input *IDParam) (*ProductResponse, error) {
	product, err := h.service.GetProductByID(input.ID)
	if err != nil {
		return nil, huma.Error404NotFound(err.Error())
	}
	return &ProductResponse{Body: product}, nil
}

func (h *InventoryHandler) ListProducts(ctx context.Context, input *PaginationQuery) (*ProductListResponse, error) {
	products, err := h.service.GetAllProducts(input.Limit, input.Offset)
	if err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}
	resp := &ProductListResponse{}
	resp.Body.Products = products
	resp.Body.Limit = input.Limit
	resp.Body.Offset = input.Offset
	return resp, nil
}

func (h *InventoryHandler) UpdateProduct(ctx context.Context, input *UpdateProductRequest) (*ProductResponse, error) {
	product, err := h.service.UpdateProduct(input.ID, &input.Body)
	if err != nil {
		return nil, huma.Error400BadRequest(err.Error())
	}
	return &ProductResponse{Body: product}, nil
}

func (h *InventoryHandler) DeleteProduct(ctx context.Context, input *IDParam) (*struct{}, error) {
	err := h.service.DeleteProduct(input.ID)
	if err != nil {
		return nil, huma.Error404NotFound(err.Error())
	}
	return &struct{}{}, nil
}

func (h *InventoryHandler) CreateTransaction(ctx context.Context, input *CreateTransactionRequest) (*TransactionResponse, error) {
	transaction, err := h.service.CreateTransaction(&input.Body)
	if err != nil {
		return nil, huma.Error400BadRequest(err.Error())
	}
	return &TransactionResponse{Body: transaction}, nil
}

func (h *InventoryHandler) GetTransaction(ctx context.Context, input *IDParam) (*TransactionResponse, error) {
	transaction, err := h.service.GetTransactionByID(input.ID)
	if err != nil {
		return nil, huma.Error404NotFound(err.Error())
	}
	return &TransactionResponse{Body: transaction}, nil
}

func (h *InventoryHandler) ListTransactions(ctx context.Context, input *PaginationQuery) (*TransactionListResponse, error) {
	transactions, err := h.service.GetAllTransactions(input.Limit, input.Offset)
	if err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}
	resp := &TransactionListResponse{}
	resp.Body.Transactions = transactions
	resp.Body.Limit = input.Limit
	resp.Body.Offset = input.Offset
	return resp, nil
}

func (h *InventoryHandler) ListProductTransactions(ctx context.Context, input *struct {
	IDParam
	PaginationQuery
}) (*TransactionListResponse, error) {
	transactions, err := h.service.GetTransactionsByProductID(input.ID, input.Limit, input.Offset)
	if err != nil {
		return nil, huma.Error404NotFound(err.Error())
	}
	resp := &TransactionListResponse{}
	resp.Body.Transactions = transactions
	resp.Body.Limit = input.Limit
	resp.Body.Offset = input.Offset
	return resp, nil
}
