package handler

import (
	"net/http"

	"github.com/prostiate/bank-statement-viewer/internal/helper"
	"github.com/prostiate/bank-statement-viewer/internal/model"
	"github.com/prostiate/bank-statement-viewer/internal/service"
)

type TransactionHandler struct {
	svc *service.TransactionService
}

func NewTransactionHandler(svc *service.TransactionService) *TransactionHandler {
	return &TransactionHandler{
		svc: svc,
	}
}

func (h *TransactionHandler) Upload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		helper.ResponseError(w, http.StatusMethodNotAllowed, "Method Not Allowed")
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		helper.ResponseError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer file.Close()

	if header.Filename == "" {
		helper.ResponseError(w, http.StatusBadRequest, "no file uploaded")
		return
	}

	if header.Header.Get("Content-Type") != "text/csv" {
		helper.ResponseError(w, http.StatusBadRequest, "invalid file type, expected csv")
		return
	}

	transactions, err := h.svc.ProcessCSV(file)
	if err != nil {
		helper.ResponseError(w, http.StatusBadRequest, err.Error())
		return
	}

	helper.ResponseSuccess(w, http.StatusOK, transactions)
}

func (h *TransactionHandler) GetBalance(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		helper.ResponseError(w, http.StatusMethodNotAllowed, "Method Not Allowed")
		return
	}

	balance, err := h.svc.GetBalance()
	if err != nil {
		helper.ResponseError(w, http.StatusInternalServerError, err.Error())
		return
	}

	helper.ResponseSuccess(w, http.StatusOK, map[string]interface{}{"balance": balance})
}

func (h *TransactionHandler) GetIssues(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		helper.ResponseError(w, http.StatusMethodNotAllowed, "Method Not Allowed")
		return
	}

	issues, err := h.svc.GetIssues()
	if err != nil {
		helper.ResponseError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if issues == nil {
		issues = []model.Transaction{}
	}

	helper.ResponseSuccess(w, http.StatusOK, map[string]interface{}{"issues": issues})
}
