package product

import (
	"database/sql"

	helper "github.com/JamesAndresCM/go-angular/helpers"
)

type Repository interface {
	GetProductByID(productID int) (*Product, error)
	GetProducts(params *getProductsRequest) ([]*Product, error)
	GetTotalProducts() (int, error)
	InsertProduct(params *getAddProductRequest) (int64, error)
	UpdateProduct(params *updateProductRequest) (int64, error)
	DeleteProduct(params *deleteProductRequest) (int64, error)
	GetBestSellers() ([]*ProductTop, error)
	GetTotalSold() (float64, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(databaseConnection *sql.DB) Repository {
	return &repository{db: databaseConnection}
}

func (repo *repository) GetProductByID(productID int) (*Product, error) {
	const sql = `SELECT id, product_code, product_name, COALESCE(description, ''),
                standard_cost, list_price, category FROM products WHERE id=?`

	row := repo.db.QueryRow(sql, productID)
	product := &Product{}

	err := row.Scan(&product.Id, &product.ProductCode, &product.ProductName, &product.Description, &product.StandardCost, &product.ListPrice, &product.Category)

	helper.Catch(err)

	return product, err
}

func (repo *repository) GetProducts(params *getProductsRequest) ([]*Product, error) {
	const sql = `SELECT id, product_code, product_name, COALESCE(description, ''),
							 standard_cost,list_price, category FROM products LIMIT ? OFFSET ?`

	results, err := repo.db.Query(sql, params.Limit, params.Offset)

	helper.Catch(err)

	var products []*Product
	for results.Next() {
		product := &Product{}
		err = results.Scan(&product.Id, &product.ProductCode, &product.ProductName,
			&product.Description, &product.StandardCost, &product.ListPrice, &product.Category)
		helper.Catch(err)

		products = append(products, product)
	}
	return products, nil
}

func (repo *repository) GetTotalProducts() (int, error) {
	const sql = `SELECT COUNT(*) FROM products`
	var total int
	row := repo.db.QueryRow(sql)
	err := row.Scan(&total)
	helper.Catch(err)
	return total, nil
}

func (repo *repository) InsertProduct(params *getAddProductRequest) (int64, error) {
	const sql = `INSERT INTO products (product_code, product_name, category, description, list_price, standard_cost)
							 VALUES(?,?,?,?,?,?)`
	result, err := repo.db.Exec(sql, params.ProductCode, params.ProductName,
		params.Category, params.Description, params.ListPrice, params.StandardCost)
	helper.Catch(err)
	id, _ := result.LastInsertId()
	return id, nil
}

func (repo *repository) UpdateProduct(params *updateProductRequest) (int64, error) {
	const sql = `UPDATE products SET 
							product_code = ?,
							product_name = ?,
							category = ?,
							description = ?,
							list_price = ?,
							standard_cost = ? WHERE id = ?`
	_, err := repo.db.Exec(sql, params.ProductCode, params.ProductName,
		params.Category, params.Description, params.ListPrice, params.StandardCost, params.Id)
	helper.Catch(err)

	return params.Id, nil
}

func (repo *repository) DeleteProduct(params *deleteProductRequest) (int64, error) {
	const sql = `DELETE FROM products WHERE id = ?`
	result, err := repo.db.Exec(sql, params.ProductID)
	helper.Catch(err)
	count, err := result.RowsAffected()
	helper.Catch(err)
	return count, nil
}

func (repo *repository) GetBestSellers() ([]*ProductTop, error) {
	const sql = `select od.product_id, p.product_name,
							 SUM(od.quantity*od.unit_price) sold
							 FROM order_details od INNER JOIN products p on od.product_id = p.id
							 GROUP BY od.product_id
							 ORDER BY sold desc
							 LIMIT 10`
	results, err := repo.db.Query(sql)
	helper.Catch(err)
	var products []*ProductTop

	for results.Next() {
		product := &ProductTop{}
		err := results.Scan(&product.Id, &product.ProductName, &product.Sold)
		helper.Catch(err)
		products = append(products, product)
	}
	return products, nil
}

func (repo *repository) GetTotalSold() (float64, error) {
	const sql = `SELECT SUM(od.quantity * od.unit_price) sold FROM order_details od`
	var total float64
	row := repo.db.QueryRow(sql)
	err := row.Scan(&total)
	helper.Catch(err)
	return total, nil
}
