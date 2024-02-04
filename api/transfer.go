package api

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/msarifin29/simple_bank/db/sqlc"
)

type TransferRequest struct {
	FromAccountId int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountId   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,gt=1"`
	Currency      string `json:"currency" binding:"required"`
}

func (server *Server) createTransfer(ctx *gin.Context) {
	var tr TransferRequest
	if err := ctx.ShouldBindJSON(&tr); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if !server.validAccount(ctx, tr.FromAccountId, tr.Currency) {
		return
	}
	if !server.validAccount(ctx, tr.ToAccountId, tr.Currency) {
		return
	}

	arg := db.TransferTxParams{
		FromAccountID: tr.FromAccountId,
		ToAccountID:   tr.ToAccountId,
		Amount:        tr.Amount,
	}

	result, err := server.store.TransferTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (server *Server) validAccount(ctx *gin.Context, accountID int64, currnecy string) bool {
	account, err := server.store.GetAccount(ctx, accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return false
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return false
	}
	if account.Currency != currnecy {
		err := fmt.Errorf("account [%d] currency mishatch %s vs %s", account.ID, account.Currency, currnecy)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return false
		}
	}
	return true
}
