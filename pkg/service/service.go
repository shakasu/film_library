//package service
//
//import (
//	"database/sql"
//	"net/http"
//)
//
//type Service struct {
//	Film
//}
//
//func NewService(repo *repository.Repository) *Service {
//	return &Service{Film: NewUserService(repo)}
//}
//
//type Film interface {
//	Balance(userId int) (*avitoTech.User, error)
//	TopUp(userId int, amount float64) (*avitoTech.User, error)
//	Debit(userId int, amount float64) (*avitoTech.User, error)
//	Transfer(userId int, toId int, amount float64) (*avitoTech.User, error)
//	ConvertBalance(user *avitoTech.User, currency string) (*avitoTech.User, error)
//	Transaction(userId int, sort string) (*[]avitoTech.Transaction, error)
//}

package main

import (
	"database/sql"
	"film_library/utils"
	"net/http"
)

type Service struct {
	db *sql.DB
}

func NewService(db *sql.DB) *Service {
	return &Service{db: db}
}

func (s *Service) CreateProduct(writer http.ResponseWriter, request *http.Request) {

	//Read request body
	productRequest := Product{}
	err := utils.ReadFromRequestBody(request, &productRequest)
	if err != nil {
		utils.WriteErrToResponseBody(writer, err)
		return
	}

	//Query to insert data
	SQL := `INSERT INTO "products" (name, price, stock) VALUES ($1, $2, $3) RETURNING id`
	err = s.db.QueryRow(SQL, productRequest.Name, productRequest.Price, productRequest.Stock).Scan(&productRequest.ID)
	if err != nil {
		utils.WriteErrToResponseBody(writer, err)
		return
	}

	//Write response
	utils.WriteToResponseBody(writer, productRequest)
}

func (s *Service) GetProducts(writer http.ResponseWriter, request *http.Request) {

	// Query to get all products
	SQL := `SELECT id, name, price, stock FROM "products"`
	rows, err := s.db.Query(SQL)
	if err != nil {
		utils.WriteErrToResponseBody(writer, err)
		return
	}
	defer rows.Close()

	// Iterate over the rows, appending each product to a slice
	products := []Product{}
	for rows.Next() {
		product := Product{}
		err = rows.Scan(&product.ID, &product.Name, &product.Price, &product.Stock)
		if err != nil {
			utils.WriteErrToResponseBody(writer, err)
			return
		}
		products = append(products, product)
	}

	// Check for errors from iterating over rows
	if err = rows.Err(); err != nil {
		utils.WriteErrToResponseBody(writer, err)
		return
	}

	// Write the products to the response body
	utils.WriteToResponseBody(writer, products)
}
