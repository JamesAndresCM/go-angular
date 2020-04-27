package main

import (
	"database/sql"
	"encoding/json"
	"go-angular/database"
	"net/http"

	"github.com/go-chi/chi"
	_ "github.com/go-sql-driver/mysql"
)

var databaseConnection *sql.DB

type Product struct {
	ID           int    `json:"id"`
	Product_Code string `json:"product_code"`
	Description  string `json:"description"`
}

func catch(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	databaseConnection = database.InitDB()
	defer databaseConnection.Close()
	r := chi.NewRouter()
	r.Get("/products", AllProducts)
	r.Post("/products", CreateProduct)
	r.Put("/products/{id}", UpdateProduct)
	r.Delete("/products/{id}", DeleteProduct)
	r.Get("/products/{id}", ShowProduct)
	http.ListenAndServe(":3000", r)
}

func DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	query, err := databaseConnection.Prepare("DELETE FROM PRODUCTS WHERE id=?")
	catch(err)
	_, er := query.Exec(id)
	catch(er)
	query.Close()
	respondwithJSON(w, http.StatusOK, map[string]string{"message": "product was deleted"})
}

func UpdateProduct(w http.ResponseWriter, r *http.Request) {
	var product Product
	id := chi.URLParam(r, "id")
	json.NewDecoder(r.Body).Decode(&product)
	query, err := databaseConnection.Prepare("UPDATE products SET product_code=?, description=? WHERE id=?")
	catch(err)
	_, er := query.Exec(product.Product_Code, product.Description, id)
	catch(er)
	defer query.Close()
	respondwithJSON(w, http.StatusOK, map[string]string{"message": "product was successflully updated"})
}

func CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product Product
	json.NewDecoder(r.Body).Decode(&product)
	query, err := databaseConnection.Prepare("INSERT products SET product_code=?, description=?")
	catch(err)
	_, er := query.Exec(product.Product_Code, product.Description)
	catch(er)
	defer query.Close()

	respondwithJSON(w, http.StatusCreated, map[string]string{"message": "product successfully created"})
}

func AllProducts(w http.ResponseWriter, r *http.Request) {
	const sql = `SELECT id, product_code, COALESCE(description, '') FROM products`

	results, err := databaseConnection.Query(sql)
	catch(err)
	var products []*Product

	for results.Next() {
		product := &Product{}
		err = results.Scan(&product.ID, &product.Product_Code, &product.Description)
		catch(err)
		products = append(products, product)
	}
	respondwithJSON(w, http.StatusOK, products)
}

func ShowProduct(w http.ResponseWriter, r *http.Request) {
	var product Product
	id := chi.URLParam(r, "id")
	query := databaseConnection.QueryRow("SELECT id, product_code, COALESCE(description, '') FROM products WHERE id=?", id)
	error := query.Scan(&product.ID, &product.Product_Code, &product.Description)
	catch(error)
	respondwithJSON(w, http.StatusOK, product)
}

func respondwithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
