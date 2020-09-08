package db

import (
	"context"

	logger "github.com/sirupsen/logrus"
)

const (
	getProductByIdQuery  = `SELECT * FROM products WHERE id=$1 LIMIT 1`
	listProduct          = `SELECT * FROM products ORDER BY name ASC`
	deleteProductIdQuery = `DELETE FROM products WHERE id = $1`
	getProductQuery      = `SELECT id, name, des,price,discount,available_quantity, category_id
		FROM products WHERE id=$1 `
	updateProductQuery = `UPDATE products SET (
			name, des, price, discount, available_quantity, category_id
			) =  ($1, $2, $3, $4, $5, $6) where id = $7`
)

type Product struct {
	Id                 int     `db:"id" json:"Id"`
	Name               string  `db:"name" json:"Name"`
	Des                string  `db:"des" json:"Des"`
	Price              float32 `db:"price" json:"Price"`
	Discount           float32 `db:"discount" json:"Discount"`
	Available_quantity int     `db:"available_quantity" json:"Available_quantity"`
	Category_Id        int     `db:"category_id" json:"Category"`
}

func (s *pgStore) ListProducts(ctx context.Context) (products []Product, err error) {
	err = s.db.Select(&products, listProduct)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error listing Produts")
		return
	}
	return
}

func (s *pgStore) GetProductById(ctx context.Context, id int) (products Product, err error) {

	err = s.db.Get(&products, getProductByIdQuery, id)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error selecting product from database by id ", id)
		return
	}
	return
}

func (s *pgStore) DeleteProductById(ctx context.Context, id int) (err error) {

	_, err = s.db.Exec(deleteProductIdQuery, id)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error deleting product")
		return
	}
	return
}

/* func (s *pgStore) UpdateProductById(ctx context.Context, product Product, id int) (updatedProduct Product, err error) {

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		logger.WithField("err:", err.Error()).Error("Error while initiating transaction")
		return
	}

	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		tx.Commit()
	}()

	var dbProduct Product
	err = s.db.Get(&dbProduct, getProductQuery, id)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error while getting product ")
		return
	}

	_, err = tx.ExecContext(ctx,
		updateProductQuery,
		product.Name,
		product.Des,
		product.Price,
		product.Discount,
		product.Available_quantity,
		product.Category_Id,
		id,
	)

	if err != nil {
		logger.WithField("err", err.Error()).Error("Error updating product attribute")
		return
	}

	updatedProduct, err = s.GetProductById(ctx, id)

	if err != nil {
		logger.WithField("err", err.Error()).Error("Error selecting user from database with productID: ", id)
		return
	}
	return

} */

func (s *pgStore) UpdateProductById(ctx context.Context, product Product, id int) (updatedProduct Product, err error) {

	var dbProduct Product
	err = s.db.Get(&dbProduct, getProductQuery, id)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error while getting product ")
		return
	}

	_, err = s.db.Exec(updateProductQuery,
		product.Name,
		product.Des,
		product.Price,
		product.Discount,
		product.Available_quantity,
		product.Category_Id,
		id,
	)

	if err != nil {
		logger.WithField("err", err.Error()).Error("Error updating product attribute")
		return
	}

	err = s.db.Get(&updatedProduct, getProductQuery, id)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error while getting product ")
		return
	}
	return

}

func (product *Product) Validate() (errorResponse map[string]ErrorResponse, valid bool) {

	fieldErrors := make(map[string]string)

	if product.Name == "" {
		fieldErrors["name"] = "Can't be blank Name"
	}

	if product.Des == "" {
		fieldErrors["discription"] = "Can't be blank Description"
	}

	if product.Price < 0 {
		fieldErrors["price"] = "Can't be Price less than zero"
	}

	if product.Discount < 0 {
		fieldErrors["discount"] = "Can't be Discount less than zero"
	}

	if product.Available_quantity < 0 {
		fieldErrors["available_quantity"] = "Can't be Available_quantity less than zero"
	}

	if len(fieldErrors) == 0 {
		valid = true
		return
	}

	errorResponse = map[string]ErrorResponse{"error": ErrorResponse{
		Code:    "invalid_data",
		Message: "Please provide valid Products's data",
		Fields:  fieldErrors,
	},
	}
	return
}