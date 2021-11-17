package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/numary/ledger/core"
	"github.com/numary/ledger/ledger"
	"github.com/numary/ledger/ledger/query"
)

// TransactionController -
type TransactionController struct {
	Controllers
}

// NewTransactionController -
func NewTransactionController() *TransactionController {
	return &TransactionController{}
}

// GetTransactions -
func (ctl *TransactionController) GetTransactions(c *gin.Context) {
	l, _ := c.Get("ledger")
	cursor, err := l.(*ledger.Ledger).FindTransactions(
		query.After(c.Query("after")),
		query.Reference(c.Query("reference")),
	)
	if err != nil {
		ctl.responseError(
			c,
			http.StatusInternalServerError,
			err,
		)
	}
	ctl.responseCollection(
		c,
		http.StatusOK,
		cursor,
	)
}

// PostTransaction -
func (ctl *TransactionController) PostTransaction(c *gin.Context) {
	l, _ := c.Get("ledger")

	var t core.Transaction
	c.ShouldBind(&t)

	err := l.(*ledger.Ledger).Commit([]core.Transaction{t})

	res := gin.H{
		"ok": err == nil,
	}

	if err != nil {
		res["err"] = err.Error()
	}

	c.JSON(200, res)
}

// RevertTransaction -
func (ctl *TransactionController) RevertTransaction(c *gin.Context) {
	l, _ := c.Get("ledger")

	err := l.(*ledger.Ledger).RevertTransaction(c.Param("transactionId"))

	res := gin.H{
		"ok": err == nil,
	}

	if err != nil {
		res["err"] = err.Error()
	}

	c.JSON(200, res)
}

// GetTransactionMetadata -
func (ctl *TransactionController) GetTransactionMetadata(c *gin.Context) {
	l, _ := c.Get("ledger")

	var m core.Metadata
	c.ShouldBind(&m)

	err := l.(*ledger.Ledger).SaveMeta(
		"transaction",
		c.Param("transactionId"),
		m,
	)

	res := gin.H{
		"ok": err == nil,
	}

	if err != nil {
		res["err"] = err.Error()
	}

	c.JSON(200, res)
}
