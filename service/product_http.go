package service

import (
	"RohitNalePUCSD/ecommeres/db"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	logger "github.com/sirupsen/logrus"
)

// @Title ListProducts
// @Description list all Products
// @Router /products [get]
// @Accept json
// @Success 200 {object}
// @Failure 400 {object}

func listProductsHandler(deps Dependencies) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		products, err := deps.Store.ListProducts(req.Context())
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error fetching data")
			rw.WriteHeader(http.StatusInternalServerError)
			repsonse(rw, http.StatusBadRequest, errorResponse{
				Error: messageObject{
					Message: "No Record Founds Products",
				},
			})
			return
		}

		respBytes, err := json.Marshal(products)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error marshaling products data")
			rw.WriteHeader(http.StatusInternalServerError)
			repsonse(rw, http.StatusBadRequest, errorResponse{
				Error: messageObject{
					Message: "Internal server error",
				},
			})
			return
		}
		rw.Header().Add("Content-Type", "application/json")
		rw.Write(respBytes)

	})
}

// @ Title getProductById
// @ Description get single product by its id
// @ Router /product/product_id [get]
// @ Accept json
// @ Success 200 {object}
// @ Failure 400 {object}

func getProductByIdHandler(deps Dependencies) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		vars := mux.Vars(req)
		id, err := strconv.Atoi(vars["product_id"])
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error id key is missing")
			rw.WriteHeader(http.StatusBadRequest)
			repsonse(rw, http.StatusBadRequest, errorResponse{
				Error: messageObject{
					Message: "Error product_id is invalid",
				},
			})
			return
		}

		product, err := deps.Store.GetProductById(req.Context(), id)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			logger.WithField("err", err.Error()).Error("Error fetching data")
			repsonse(rw, http.StatusBadRequest, errorResponse{
				Error: messageObject{
					Message: "Error feching data No Row Found",
				},
			})
			return
		}

		respBytes, err := json.Marshal(product)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			logger.WithField("err", err.Error()).Error("Error marshaling products data")
			repsonse(rw, http.StatusBadRequest, errorResponse{
				Error: messageObject{
					Message: "Internal server error",
				},
			})
			return
		}
		rw.Header().Add("Content-Type", "application/json")
		rw.Write(respBytes)
		return

	})
}

// @ Title deleteProductById
// @ Description delete product by its id
// @ Router /product/product_id [delete]
// @ Accept json
// @ Success 200 {object}
// @ Failure 400 {object}

func deleteProductByIdHandler(deps Dependencies) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		vars := mux.Vars(req)
		id, err := strconv.Atoi(vars["product_id"])
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error id key is missing")
			rw.WriteHeader(http.StatusBadRequest)
			repsonse(rw, http.StatusBadRequest, errorResponse{
				Error: messageObject{
					Message: "Error id is missing/invalid",
				},
			})
			return
		}

		err = deps.Store.DeleteProductById(req.Context(), id)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			logger.WithField("err", err.Error()).Error("Error fetching data no row found")
			repsonse(rw, http.StatusBadRequest, errorResponse{
				Error: messageObject{
					Message: "Internal server error  (Error feching data)",
				},
			})
			return
		}

		rw.WriteHeader(http.StatusOK)
		rw.Header().Add("Content-Type", "application/json")

	})
}

// @ Title updateProductById
// @ Description update product by its id
// @ Router /product/product_id [put]
// @ Accept json
// @ Success 200 {object}
// @ Failure 400 {object}

func updateProductByIdHandler(deps Dependencies) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		vars := mux.Vars(req)
		id, err := strconv.Atoi(vars["product_id"])
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error id key is missing")
			rw.WriteHeader(http.StatusBadRequest)
			repsonse(rw, http.StatusBadRequest, errorResponse{
				Error: messageObject{
					Message: "Error id is missing/invalid",
				},
			})
			return
		}

		var product db.Product
		err = json.NewDecoder(req.Body).Decode(&product)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			logger.WithField("err", err.Error()).Error("Error while decoding user")
			repsonse(rw, http.StatusBadRequest, errorResponse{
				Error: messageObject{
					Message: "Internal server error",
				},
			})
			return
		}

		errRes, valid := product.Validate()
		if !valid {
			respBytes, err := json.Marshal(errRes)
			if err != nil {
				logger.WithField("err", err.Error()).Error("Error marshaling product data")
				repsonse(rw, http.StatusBadRequest, errorResponse{
					Error: messageObject{
						Message: "Invalid json body",
					},
				})
				rw.WriteHeader(http.StatusInternalServerError)
				return
			}
			rw.Header().Add("Content-Type", "application/json")
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write(respBytes)
			return
		}

		var updatedProduct db.Product
		updatedProduct, err = deps.Store.UpdateProductById(req.Context(), product, id)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			repsonse(rw, http.StatusInternalServerError, errorResponse{
				Error: messageObject{
					Message: "Internal server error",
				},
			})
			logger.WithField("err", err.Error()).Error("Error while updating product attribute")
			return
		}

		repsonse(rw, http.StatusOK, successResponse{Data: updatedProduct})

		return
	})
}
