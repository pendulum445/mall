package product

import (
	"context"
	"database/sql"
	"log"
	"strings"
)

type ProductService struct {
	UnimplementedProductCatalogServiceServer
	Db *sql.DB
}

func (s *ProductService) ListProducts(ctx context.Context, req *ListProductsReq) (*ListProductsResp, error) {
	query := `
		SELECT p.id, p.name, p.description, p.picture, p.price, GROUP_CONCAT(c.name) AS categories
		FROM products p
		JOIN product_categories pc ON p.id = pc.product_id
		JOIN categories c ON pc.category_id = c.id
		WHERE c.name = ? GROUP BY p.id
	`
	rows, err := s.Db.Query(query, req.CategoryName)
	if err != nil {
		return nil, err
	}
	resp := &ListProductsResp{}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Fatalf("failed to close rows: %v", err)
		}
	}(rows)
	for rows.Next() {
		var (
			id          uint32
			name        string
			description string
			picture     string
			price       float32
			categories  string
		)
		if err := rows.Scan(&id, &name, &description, &picture, &price, &categories); err != nil {
			return nil, err
		}
		resp.Products = append(resp.Products, &Product{
			Id:          id,
			Name:        name,
			Description: description,
			Picture:     picture,
			Price:       price,
			Categories:  strings.Split(categories, ","),
		})
	}
	return resp, nil
}

func (s *ProductService) GetProduct(ctx context.Context, req *GetProductReq) (*GetProductResp, error) {
	query := `
		SELECT p.id, p.name, p.description, p.picture, p.price, GROUP_CONCAT(c.name) AS categories
		FROM products p
		JOIN product_categories pc ON p.id = pc.product_id
		JOIN categories c ON pc.category_id = c.id
		WHERE p.id = ? GROUP BY p.id
	`
	var (
		id          uint32
		name        string
		description string
		picture     string
		price       float32
		categories  string
	)
	err := s.Db.QueryRow(query, req.Id).Scan(&id, &name, &description, &picture, &price, &categories)
	if err != nil {
		return nil, err
	}
	return &GetProductResp{
		Product: &Product{
			Id:          id,
			Name:        name,
			Description: description,
			Picture:     picture,
			Price:       price,
			Categories:  strings.Split(categories, ","),
		},
	}, nil
}

func (s *ProductService) SearchProducts(ctx context.Context, req *SearchProductsReq) (*SearchProductsResp, error) {
	query := `
		SELECT p.id, p.name, p.description, p.picture, p.price, GROUP_CONCAT(c.name) AS categories
		FROM products p
		JOIN product_categories pc ON p.id = pc.product_id
		JOIN categories c ON pc.category_id = c.id
		WHERE p.name LIKE ? or p.description LIKE ? GROUP BY p.id
	`
	rows, err := s.Db.Query(query, "%"+req.Query+"%", "%"+req.Query+"%")
	if err != nil {
		return nil, err
	}
	resp := &SearchProductsResp{}
	for rows.Next() {
		var (
			id          uint32
			name        string
			description string
			picture     string
			price       float32
			categories  string
		)
		if err := rows.Scan(&id, &name, &description, &picture, &price, &categories); err != nil {
			return nil, err
		}
		resp.Results = append(resp.Results, &Product{
			Id:          id,
			Name:        name,
			Description: description,
			Picture:     picture,
			Price:       price,
			Categories:  strings.Split(categories, ","),
		})
	}
	return resp, nil
}
