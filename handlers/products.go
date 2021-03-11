package handlers

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/alazo8807/jackson_tut/data"
	"github.com/gorilla/mux"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

// This is only needed if implementing our server handlers.
// With Gorilla Mux this is not needed anymore
// func (p *Products) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
// 	if r.Method == http.MethodGet {
// 		p.getProducts(rw, r)
// 		return
// 	}

// 	p.l.Println("method: ", r.Method)

// 	if r.Method == http.MethodPost {
// 		p.addProduct(rw, r)
// 		return
// 	}

// 	if r.Method == http.MethodPut {
// 		// expect the id in the url
// 		reg := regexp.MustCompile(`/([0-9])+`)
// 		g := reg.FindAllStringSubmatch(r.URL.Path, -1)

// 		if len(g) != 1 {
// 			http.Error(rw, "Invalid URI", http.StatusBadRequest)
// 			return
// 		}

// 		if len(g[0]) != 2 {
// 			http.Error(rw, "Invalid URI", http.StatusBadRequest)
// 			return
// 		}

// 		idString := g[0][1]
// 		id, err := strconv.Atoi(idString)

// 		if err != nil {
// 			http.Error(rw, "Invalid URI", http.StatusBadRequest)
// 		}
// 		p.l.Println("got id", id)

// 		p.updateProduct(id, rw, r)
// 		return
// 	}

// 	// Return not allowed method for everything else
// 	rw.WriteHeader(http.StatusMethodNotAllowed)
// }

// // Update Product (Without gorilla mux)
// func (p *Products) UpdateProduct(id int, rw http.ResponseWriter, r *http.Request) {
// 	p.l.Println("Handle PUT Product")

// 	prod := &data.Product{}

// 	err := prod.FromJSON(r.Body)
// 	if err != nil {
// 		http.Error(rw, "Unable to marshal json", http.StatusBadRequest)
// 		return
// 	}

// 	err = data.UpdateProduct(id, prod)

// 	if err == data.ErrProductNotFound {
// 		http.Error(rw, "Product not found", http.StatusNotFound)
// 		return
// 	}

// 	if err != nil {
// 		http.Error(rw, "Product not found", http.StatusInternalServerError)
// 		return
// 	}

// }

// Add Product
func (p *Products) AddProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle post request")

	// Deserialize body into a Product object (Logic moved to middleware)
	// prod := &data.Product{}
	// err := prod.FromJSON(r.Body)
	// if err != nil {
	// 	http.Error(rw, "Unable to unmarshar json", http.StatusBadRequest)
	// }

	// Get deserialized product object from context. (Comes from middleware)
	prod := r.Context().Value(KeyProduct{}).(data.Product)

	data.AddProduct(&prod)
}

// Update Product
func (p *Products) UpdateProduct(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(rw, "Unable to convert id to int", http.StatusBadRequest)
	}

	p.l.Println("Handle PUT Product")

	// Get deserialized product object from context. (Comes from middleware)
	prod := r.Context().Value(KeyProduct{}).(data.Product)

	err = data.UpdateProduct(id, &prod)

	if err == data.ErrProductNotFound {
		http.Error(rw, "Product not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(rw, "Product not found", http.StatusInternalServerError)
		return
	}
}

// Get Products
func (p *Products) GetProducts(rw http.ResponseWriter, r *http.Request) {
	lp := data.GetProducts()
	// We defined a ToJSON method for the list of products that it will use an
	// Encoder instead of Marshal. This method is more efficient than using marhsal
	err := lp.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to read json", http.StatusInternalServerError)
	}

	// Using Marshal example:
	// Marshal returns a byte buffer. Which allocates the result in memory.
	// d, err := json.Marshal(lp)

	// if err != nil {
	// 	http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	// }

	// rw.Write(d)
	//
}

type KeyProduct struct{}

func (p *Products) MiddlewareProductValidateProduct(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		prod := data.Product{}

		err := prod.FromJSON(r.Body)
		if err != nil {
			p.l.Println("[ERROR] desirializing the product", err)
			http.Error(rw, "Unable to unmarshar json", http.StatusBadRequest)
			return
		}

		// add the product to the context
		ctx := context.WithValue(r.Context(), KeyProduct{}, prod)
		r = r.WithContext(ctx)

		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(rw, r)
	})
}
