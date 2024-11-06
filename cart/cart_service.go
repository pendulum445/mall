package cart

import (
	"context"
	"database/sql"
	"log"
)

type CartService struct {
	UnimplementedCartServiceServer
	Db *sql.DB
}

func (s *CartService) AddItem(ctx context.Context, req *AddItemReq) (*AddItemResp, error) {
	var id int
	err := s.Db.QueryRow("SELECT id FROM users WHERE id = ?", req.UserId).Scan(&id)
	if err != nil {
		return nil, err
	}
	err = s.Db.QueryRow("SELECT id FROM products WHERE id = ?", req.Item.ProductId).Scan(&id)
	if err != nil {
		return nil, err
	}
	var quantity int32
	err = s.Db.QueryRow("SELECT quantity FROM products WHERE user_id = ? AND product_id = ?",
		req.UserId, req.Item.ProductId).Scan(&quantity)
	if err != nil {
		_, err := s.Db.Exec("INSERT INTO cart_items (user_id, product_id, quantity) VALUES (?,?,?)",
			req.UserId, req.Item.ProductId, req.Item.Quantity)
		if err != nil {
			return nil, err
		}
	} else {
		_, err := s.Db.Exec("UPDATE cart_items SET quantity = ? WHERE user_id = ? AND product_id = ?",
			quantity+req.Item.Quantity, req.UserId, req.Item.ProductId)
		if err != nil {
			return nil, err
		}
	}
	return &AddItemResp{Res: true}, nil
}

func (s *CartService) GetCart(ctx context.Context, req *GetCartReq) (*GetCartResp, error) {
	rows, err := s.Db.Query("SELECT product_id, quantity FROM cart_items WHERE user_id = ?", req.UserId)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(rows)
	resp := &GetCartResp{Cart: &Cart{UserId: req.UserId}}
	resp.Cart.Items = make([]*CartItem, 0)
	for rows.Next() {
		var (
			productId uint32
			quantity  int32
		)
		err := rows.Scan(&productId, &quantity)
		if err != nil {
			log.Fatal(err)
		}
		resp.Cart.Items = append(resp.Cart.Items, &CartItem{ProductId: productId, Quantity: quantity})
	}
	return resp, nil
}

func (s *CartService) EmptyCart(ctx context.Context, req *EmptyCartReq) (*EmptyCartResp, error) {
	_, err := s.Db.Exec("DELETE FROM cart_items WHERE user_id =?", req.UserId)
	if err != nil {
		return nil, err
	}
	return &EmptyCartResp{Res: true}, nil
}
