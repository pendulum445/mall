package order

import (
	"database/sql"
)

type OrderService struct {
	UnimplementedOrderServiceServer
	Db *sql.DB
}

//func (s *OrderService) PlaceOrder(ctx *context.Context, req *PlaceOrderReq) (*PlaceOrderResp, error) {
//
//}
