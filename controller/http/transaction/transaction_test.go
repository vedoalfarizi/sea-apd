package transaction

import (
	"encoding/json"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo"
	message "github.com/williamchang80/sea-apd/common/constants/response"
	"github.com/williamchang80/sea-apd/common/constants/transaction_status"
	domain "github.com/williamchang80/sea-apd/domain/transaction"
	"github.com/williamchang80/sea-apd/dto/request/transaction"
	request "github.com/williamchang80/sea-apd/dto/request/transaction"
	"github.com/williamchang80/sea-apd/dto/response/base"
	transaction_repository "github.com/williamchang80/sea-apd/mocks/repository/transaction"
	transaction_mock_usecase "github.com/williamchang80/sea-apd/mocks/usecase/transaction"
	transaction_usecase "github.com/williamchang80/sea-apd/usecase/transaction"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"
)

var (
	mockUpdateTransactionStatusRequest = transaction.UpdateTransactionRequest{
		TransactionId: "1",
		Status:        transaction_status.ACCEPTED,
	}
	mockCreateTransactionRequest = transaction.TransactionRequest{
		BankNumber: "123456789",
		BankName:   "Mock Bank",
		Amount:     10000,
		MerchantId: "1",
		CustomerId: "1",
	}
	mockId = "1"
)

func TestNewTransactionController(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := echo.New()
	repo := transaction_repository.NewMockRepository(ctrl)
	type args struct {
		ctx *echo.Echo
	}
	tests := []struct {
		name     string
		args     args
		want     domain.TransactionController
		initMock func() domain.TransactionUsecase
	}{
		{
			name: "success",
			args: args{
				ctx: ctx,
			},
			want: &TransactionController{
				usecase: transaction_usecase.NewTransactionUsecase(repo, nil, nil),
			},
			initMock: func() domain.TransactionUsecase {
				return transaction_mock_usecase.NewMockUsecase(ctrl)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := tt.initMock()

			if got := NewTransactionController(tt.args.ctx, mock); reflect.TypeOf(got) != reflect.TypeOf(tt.want) {
				t.Errorf("NewTransactionController() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTransactionController_CreateTransaction(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		ctx     *echo.Echo
		request request.TransactionRequest
	}
	defer ctrl.Finish()
	tests := []struct {
		name     string
		args     args
		want     base.BaseResponse
		wantErr  bool
		initMock func() domain.TransactionUsecase
	}{
		{
			name: "success",
			args: args{
				ctx:     echo.New(),
				request: mockCreateTransactionRequest,
			},
			wantErr: false,
			want: base.BaseResponse{
				Code:    http.StatusOK,
				Message: message.SUCCESS,
			},
			initMock: func() domain.TransactionUsecase {
				return transaction_mock_usecase.NewMockUsecase(ctrl)
			},
		},
		{
			name: "failed with empty request",
			args: args{
				ctx:     echo.New(),
				request: request.TransactionRequest{},
			},
			wantErr: false,
			want: base.BaseResponse{
				Code:    http.StatusNotFound,
				Message: message.NOT_FOUND,
			},
			initMock: func() domain.TransactionUsecase {
				return transaction_mock_usecase.NewMockUsecase(ctrl)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := tt.initMock()
			c := echo.New()
			s, _ := json.Marshal(tt.args.request)
			req, err := http.NewRequest(echo.POST, "api/transaction", strings.NewReader(string(s)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			if err != nil {
				t.Errorf("CreateTransaction() request error= %v", tt.wantErr)
			}
			rec := httptest.NewRecorder()
			ctx := c.NewContext(req, rec)
			controller := NewTransactionController(c, mock)
			if got := controller.CreateTransaction(ctx); (got != nil) != tt.wantErr {
				t.Errorf("CreateTransaction() error= %v, want %v", got, tt.wantErr)
			}
		})
	}
}

func TestTransactionController_UpdateTransactionStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		ctx     *echo.Echo
		request request.UpdateTransactionRequest
	}
	defer ctrl.Finish()
	tests := []struct {
		name     string
		args     args
		wantErr  bool
		want     base.BaseResponse
		initMock func() domain.TransactionUsecase
	}{
		{
			name: "failed with empty request and invalid status",
			args: args{
				ctx:     echo.New(),
				request: request.UpdateTransactionRequest{},
			},
			wantErr: false,
			want: base.BaseResponse{
				Code:    http.StatusNotFound,
				Message: message.NOT_FOUND,
			},
			initMock: func() domain.TransactionUsecase {
				return transaction_mock_usecase.NewMockUsecase(ctrl)
			},
		},
		{
			name: "success",
			args: args{
				ctx:     echo.New(),
				request: mockUpdateTransactionStatusRequest,
			},
			wantErr: false,
			want: base.BaseResponse{
				Code:    http.StatusOK,
				Message: message.SUCCESS,
			},
			initMock: func() domain.TransactionUsecase {
				return transaction_mock_usecase.NewMockUsecase(ctrl)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := tt.initMock()
			c := echo.New()
			s, _ := json.Marshal(tt.args.request)
			req, err := http.NewRequest(echo.POST, "api/transaction/status", strings.NewReader(string(s)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			if err != nil && !tt.wantErr {
				t.Errorf("UpdateTransactionStatus() request error= %v", tt.wantErr)
			}
			rec := httptest.NewRecorder()
			ctx := c.NewContext(req, rec)
			controller := NewTransactionController(c, mock)
			if controller.UpdateTransactionStatus(ctx); (rec.Code != tt.want.Code) || tt.wantErr {
				t.Errorf("UpdateTransactionStatus() error= %v, want %v", rec.Code, tt.want.Code)
			}
		})
	}
}

func TestTransactionController_GetTransactionById(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		ctx       *echo.Echo
		getParams func() url.Values
	}
	defer ctrl.Finish()
	tests := []struct {
		name       string
		args       args
		wantErr    bool
		wantStatus int
		initMock   func() domain.TransactionUsecase
	}{
		{
			name: "failed with empty request and invalid status",
			args: args{
				ctx: echo.New(),
				getParams: func() url.Values {
					q := make(url.Values)
					return q
				},
			},
			wantErr:    false,
			wantStatus: http.StatusNotFound,
			initMock: func() domain.TransactionUsecase {
				return transaction_mock_usecase.NewMockUsecase(ctrl)
			},
		},
		{
			name: "success",
			args: args{
				ctx: echo.New(),
				getParams: func() url.Values {
					q := make(url.Values)
					q.Set("transactionId", mockId)
					return q
				},
			},
			wantErr:    false,
			wantStatus: http.StatusOK,
			initMock: func() domain.TransactionUsecase {
				return transaction_mock_usecase.NewMockUsecase(ctrl)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := tt.initMock()
			params := tt.args.getParams()
			c := echo.New()
			req, err := http.NewRequest(echo.GET, "api/transaction"+"?"+params.Encode(), nil)
			if err != nil && !tt.wantErr {
				t.Errorf("GetTransactionById() request error= %v", tt.wantErr)
			}
			rec := httptest.NewRecorder()
			ctx := c.NewContext(req, rec)
			controller := NewTransactionController(c, mock)
			if controller.GetTransactionById(ctx); (rec.Code != tt.wantStatus) || tt.wantErr {
				t.Errorf("GetTransactionById() error= %v, want %v", rec.Code, tt.wantStatus)
			}
		})
	}
}

func TestTransactionController_GetTransactionHistory(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		ctx       *echo.Echo
		getParams func() url.Values
	}
	defer ctrl.Finish()
	tests := []struct {
		name       string
		args       args
		wantErr    bool
		wantStatus int
		initMock   func() domain.TransactionUsecase
	}{
		{
			name: "failed with empty request and invalid status",
			args: args{
				ctx: echo.New(),
				getParams: func() url.Values {
					q := make(url.Values)
					return q
				},
			},
			wantErr:    false,
			wantStatus: http.StatusNotFound,
			initMock: func() domain.TransactionUsecase {
				return transaction_mock_usecase.NewMockUsecase(ctrl)
			},
		},
		{
			name: "success",
			args: args{
				ctx: echo.New(),
				getParams: func() url.Values {
					q := make(url.Values)
					q.Set("userId", mockId)
					return q
				},
			},
			wantErr:    false,
			wantStatus: http.StatusOK,
			initMock: func() domain.TransactionUsecase {
				return transaction_mock_usecase.NewMockUsecase(ctrl)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := tt.initMock()
			params := tt.args.getParams()
			c := echo.New()
			req, err := http.NewRequest(echo.GET, "api/transaction/history"+"?"+params.Encode(), nil)
			if err != nil && !tt.wantErr {
				t.Errorf("GetTransactionHistory() request error= %v", tt.wantErr)
			}
			rec := httptest.NewRecorder()
			ctx := c.NewContext(req, rec)
			controller := NewTransactionController(c, mock)
			if controller.GetTransactionHistory(ctx); (rec.Code != tt.wantStatus) || tt.wantErr {
				t.Errorf("GetTransactionHistory() error= %v, want %v", rec.Code, tt.wantStatus)
			}
		})
	}
}
