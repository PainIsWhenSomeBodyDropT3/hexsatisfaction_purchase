package model

import "time"

type (

	// CreatePurchaseRequest represents a request to create purchase.
	CreatePurchaseRequest struct {
		// required: true
		UserID string `json:"userID"`
		// required: true
		Date time.Time `json:"date"`
		// required: true
		FileID string `json:"fileID"`
	}

	// IDPurchaseRequest represents a request to find the purchase by id.
	IDPurchaseRequest struct {
		// required: true
		ID string `json:"-"`
	}

	// DeletePurchaseRequest represents a request to delete purchase.
	DeletePurchaseRequest struct {
		// required: true
		ID string `json:"-"`
	}

	// UserIDPurchaseRequest represents a request to find last added purchase by user id.
	UserIDPurchaseRequest struct {
		// required: true
		ID string `json:"-"`
	}

	// UserIDPeriodPurchaseRequest represents a request to find all purchases by user id and date period.
	UserIDPeriodPurchaseRequest struct {
		// required: true
		ID string `json:"-"`
		// required: true
		Start time.Time `json:"start"`
		// required: true
		End time.Time `json:"end"`
	}

	// UserIDAfterDatePurchaseRequest represents a request to find all purchases by user id after date.
	UserIDAfterDatePurchaseRequest struct {
		// required: true
		ID string `json:"-"`
		// required: true
		Start time.Time `json:"start"`
	}

	// UserIDBeforeDatePurchaseRequest represents a request to find all purchases by user id before date.
	UserIDBeforeDatePurchaseRequest struct {
		// required: true
		ID string `json:"-"`
		// required: true
		End time.Time `json:"end"`
	}

	// UserIDFileIDPurchaseRequest represents a request to find all purchases by user id and file name.
	UserIDFileIDPurchaseRequest struct {
		// required: true
		UserID string `json:"-"`
		// required: true
		FileID string `json:"fileID"`
	}

	// PeriodPurchaseRequest represents a request to find all purchases by date period.
	PeriodPurchaseRequest struct {
		// required: true
		Start time.Time `json:"start"`
		// required: true
		End time.Time `json:"end"`
	}

	// AfterDatePurchaseRequest represents a request to find all purchases after date.
	AfterDatePurchaseRequest struct {
		// required: true
		Start time.Time `json:"start"`
	}

	// BeforeDatePurchaseRequest represents a request to find all purchases before date.
	BeforeDatePurchaseRequest struct {
		// required: true
		End time.Time `json:"end"`
	}

	// FileIDPurchaseRequest represents a request to find all purchases by file name.
	FileIDPurchaseRequest struct {
		// required: true
		FileID string `json:"-"`
	}
)

type (

	// CreateCommentRequest represents a request to create comment.
	CreateCommentRequest struct {
		// required: true
		UserID string `json:"userID"`
		// required: true
		PurchaseID string `json:"purchaseID"`
		// required: true
		Date time.Time `json:"Date"`
		// required: true
		Text string `json:"Text"`
	}

	// UpdateCommentRequest represents a request to update comment.
	UpdateCommentRequest struct {
		// required: true
		ID string `json:"-"`
		// required: true
		UserID string `json:"userID"`
		// required: true
		PurchaseID string `json:"purchaseID"`
		// required: true
		Date time.Time `json:"date"`
		// required: true
		Text string `json:"text"`
	}

	// DeleteCommentRequest represents a request to delete comment.
	DeleteCommentRequest struct {
		// required: true
		ID string `json:"-"`
	}

	// IDCommentRequest represents a request to find comment by id.
	IDCommentRequest struct {
		// required: true
		ID string `json:"-"`
	}

	// UserIDCommentRequest represents a request to find comments by user id.
	UserIDCommentRequest struct {
		// required: true
		ID string `json:"-"`
	}

	// PurchaseIDCommentRequest represents a request to find comments by purchase id.
	PurchaseIDCommentRequest struct {
		// required: true
		ID string `json:"-"`
	}

	// UserPurchaseIDCommentRequest represents a request to find comments by purchase and user ids.
	UserPurchaseIDCommentRequest struct {
		// required: true
		UserID     string `json:"-"`
		PurchaseID string `json:"-"`
	}

	// TextCommentRequest represents a request to find comments by text.
	TextCommentRequest struct {
		// required: true
		Text string `json:"text"`
	}

	// PeriodCommentRequest represents a request to find comments by date period.
	PeriodCommentRequest struct {
		// required: true
		Start time.Time `json:"start"`
		// required: true
		End time.Time `json:"end"`
	}
)

type (
	// CreateFileRequest represents a request to create file.
	CreateFileRequest struct {
		// required: true
		Name string `json:"name"`
		// required: true
		Description string `json:"description"`
		// required: true
		Size int `json:"size"`
		// required: true
		Path string `json:"path"`
		// required: true
		AddDate time.Time `json:"addDate"`
		// required: true
		UpdateDate time.Time `json:"updateDate"`
		// required: true
		Actual bool `json:"actual"`
		// required: true
		AuthorID string `json:"authorID"`
	}

	// UpdateFileRequest represents a request to update file.
	UpdateFileRequest struct {
		// required: true
		ID string `json:"-"`
		// required: true
		Name string `json:"name"`
		// required: true
		Description string `json:"description"`
		// required: true
		Size int `json:"size"`
		// required: true
		Path string `json:"path"`
		// required: true
		AddDate time.Time `json:"addDate"`
		// required: true
		UpdateDate time.Time `json:"updateDate"`
		// required: true
		Actual bool `json:"actual"`
		// required: true
		AuthorID string `json:"authorID"`
	}

	// DeleteFileRequest represents a request to delete file.
	DeleteFileRequest struct {
		// required: true
		ID string `json:"-"`
	}

	// IDFileRequest represents a request to find file by id.
	IDFileRequest struct {
		// required: true
		ID string `json:"-"`
	}

	// NameFileRequest represents a request to find files by name.
	NameFileRequest struct {
		// required: true
		Name string `json:"-"`
	}

	// AuthorIDFileRequest represents a request to find files by author id.
	AuthorIDFileRequest struct {
		// required: true
		ID string `json:"-"`
	}

	// AddedPeriodFileRequest represents a request to find added files by date period.
	AddedPeriodFileRequest struct {
		// required: true
		Start time.Time `json:"start"`
		// required: true
		End time.Time `json:"end"`
	}

	// UpdatedPeriodFileRequest represents a request to find updated files by date period.
	UpdatedPeriodFileRequest struct {
		// required: true
		Start time.Time `json:"start"`
		// required: true
		End time.Time `json:"end"`
	}
)
