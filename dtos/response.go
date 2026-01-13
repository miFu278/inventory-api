package dtos

type SingleProductResponse struct {
	Body *ProductResponse
}

type ProductListResponse struct {
	Body struct {
		Products []ProductResponse `json:"products"`
		Limit    int               `json:"limit"`
		Offset   int               `json:"offset"`
	}
}

type SingleTransactionResponse struct {
	Body *TransactionResponse
}

type TransactionListResponse struct {
	Body struct {
		Transactions []TransactionResponse `json:"transactions"`
		Limit        int                   `json:"limit"`
		Offset       int                   `json:"offset"`
	}
}

type EmptyResponse struct{}

// User responses
type LoginResponse struct {
	Body struct {
		Token string        `json:"token"`
		User  *UserResponse `json:"user"`
	}
}

type SingleUserResponse struct {
	Body *UserResponse
}

type UserListResponse struct {
	Body struct {
		Users  []UserResponse `json:"users"`
		Limit  int            `json:"limit"`
		Offset int            `json:"offset"`
	}
}
