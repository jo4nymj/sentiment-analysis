package repository

import (
	"database/sql"

	"code.sentiments/models"
)

type ProductModel struct {
	Db *sql.DB
}

func (r ProductModel) GetProduct(name string) (models.Product, error) {
	product := models.Product{}
	rows, err := r.Db.Query("SELECT ID, post_title, average_rating FROM wp_posts p INNER JOIN wp_wc_product_meta_lookup pr ON p.ID = pr.product_id WHERE p.post_title LIKE '%'||?||'%'", name)
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
