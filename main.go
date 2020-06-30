package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

//Product Model Schema
type Product struct {
	ProductID      int    `json:"productId"`
	Manufacturer   string `json:"manufacturer"`
	Sku            string `json:"sku"`
	Upc            string `json:"upc"`
	PricePerUnit   string `json:"pricePerUnit"`
	QuantityOnHand int    `json:"quantityOnHand"`
	ProductName    string `json:"productName"`
}

var ProductList []Product

//Initialize the product array
func init() {
	var ProductJSON = `[
  {
    "productId": 1,
    "manufacturer": "Johns-Jenkins",
    "sku": "p5z343vdS",
    "upc": "939581000000",
    "pricePerUnit": "497.45",
    "quantityOnHand": 9703,
    "productName": "sticky note"
  },
  {
    "productId": 2,
    "manufacturer": "Hessel, Schimmel and Feeney",
    "sku": "i7v300kmx",
    "upc": "740979000000",
    "pricePerUnit": "282.29",
    "quantityOnHand": 9217,
    "productName": "leg warmers"
  },
  {
    "productId": 3,
    "manufacturer": "Swaniawski, Bartoletti and Bruen",
    "sku": "q0L657ys7",
    "upc": "111730000000",
    "pricePerUnit": "436.26",
    "quantityOnHand": 5905,
    "productName": "lamp shade"
  }
]`

	err := json.Unmarshal([]byte(ProductJSON), &ProductList)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(ProductList)
}

//Add a product or get a list of products
func HandleProducts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		ProductJSON, err := json.Marshal(ProductList)

		if err != nil {
			log.Fatal(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write(ProductJSON)

	case http.MethodPost:
		var newProduct Product
		bodyBytes, err := ioutil.ReadAll(r.Body)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = json.Unmarshal(bodyBytes, &newProduct)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		newProduct.ProductID = getNextId()
		ProductList = append(ProductList, newProduct)

		w.WriteHeader(http.StatusOK)
	}
}

//Calculate the new ProductId while adding a product
func getNextId() int {
	nextId := 0
	for _, v := range ProductList {
		if nextId <= v.ProductID {
			nextId = v.ProductID + 1
		}
	}
	return nextId
}

//Get a product by Id or update a product by Id
func GetProduct(w http.ResponseWriter, r *http.Request) {
	urlSegments := strings.Split(r.URL.Path, "products/")

	productId, err := strconv.Atoi(urlSegments[len(urlSegments)-1])

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	product, listIndex := FindProductById(productId)

	if product == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	switch r.Method {
	case http.MethodGet:
		bytes, err := json.Marshal(product)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Write(bytes)
		return
	case http.MethodPut:
		requestBody, err := ioutil.ReadAll(r.Body)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var productBody Product

		err = json.Unmarshal(requestBody, &productBody)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if productBody.ProductID != productId {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		ProductList[listIndex] = productBody
		w.WriteHeader(http.StatusOK)
		return
	}

}

//Find a product from products array
func FindProductById(productId int) (*Product, int) {
	for i, v := range ProductList {
		if v.ProductID == productId {
			return &v, i
		}
	}
	return nil, 0
}

func main() {
	http.HandleFunc("/products", HandleProducts)
	http.HandleFunc("/products/", GetProduct)
	http.ListenAndServe(":5000", nil)
}
