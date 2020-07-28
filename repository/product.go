package repository

import (
	"code.sentiments/config"
	"code.sentiments/models"
)

type ProductModel struct {
	Db *config.Connection
}

func (r ProductModel) GetProduct(name string) (models.Product, error) {
	product := models.Product{}
	rows, err := r.Db.Conn.Query(`SELECT ID, post_title, average_rating 
		FROM wp_posts p INNER JOIN wp_wc_product_meta_lookup pr ON p.ID = pr.product_id 
		WHERE p.post_title LIKE ?`, "%"+name+"%")
	if err != nil {
		return product, err
	}

	for rows.Next() {
		if err := rows.Scan(&product.ID, &product.Name, &product.Average_Rating); err != nil {
			return product, err
		}
	}

	return product, nil
}

func (r ProductModel) UpdateProduct(product models.Product) error {
	stmt, err := r.Db.Conn.Prepare("UPDATE  wp_wc_product_meta_lookup SET average_rating = ? WHERE product_id = ?")
	if err != nil {
		return err
	}
	if _, err := stmt.Exec(product.Average_Rating, product.ID); err != nil {
		return err
	}

	return nil
}

func (r ProductModel) ListProducts() ([]models.Product, error) {
	rows, err := r.Db.Conn.Query(`SELECT ID, post_title, average_rating 
		FROM wp_posts p INNER JOIN wp_wc_product_meta_lookup pr ON p.ID = pr.product_id`)
	if err != nil {
		return nil, err
	}
	products := []models.Product{}
	product := models.Product{}
	for rows.Next() {
		if err := rows.Scan(&product.ID, &product.Name, &product.Average_Rating); err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	return products, nil
}
