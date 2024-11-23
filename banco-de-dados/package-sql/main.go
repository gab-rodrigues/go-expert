package bd

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

type Product struct {
	ID    string
	Name  string
	Price float64
}

func NewProduct(name string, price float64) *Product {
	return &Product{
		ID:    uuid.New().String(),
		Name:  name,
		Price: price,
	}
}

func main() {
	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/goexpert")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	product := NewProduct("Notebook", 1899.90)
	err = insertProduct(db, product)
	if err != nil {
		panic(err)
	}
	product.Price = 210.80
	err = updateProduct(db, product)
	if err != nil {
		panic(err)
	}

	products, err := selectProducts(db)
	if err != nil {
		panic(err)
	}

	for _, product := range products {
		fmt.Printf("Product: %s possui preco %.2f\n", product.Name, product.Price)
	}

	err = deleteProductByID(db, product.ID)
	if err != nil {
		panic(err)
	}

	//ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	//defer cancel()
	//
	//p, err := getProductByIdWithContext(ctx, db, product.ID)
	//select {
	//case <-ctx.Done():
	//	fmt.Println("timeout exceeded")
	//	return
	//case <-time.After(1 * time.Millisecond):
	//	fmt.Printf("Product %s, possui o preco de %.2f", p.Name, p.Price)
	//}
}

func insertProduct(db *sql.DB, product *Product) error {
	stmt, err := db.Prepare("INSERT INTO products(id, name, price) VALUES(?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(product.ID, product.Name, product.Price)
	if err != nil {
		return err
	}
	return nil
}

func updateProduct(db *sql.DB, product *Product) error {
	stmt, err := db.Prepare("update products set name=?, price=? where id=?")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(product.Name, product.Price, product.ID)
	if err != nil {
		return err
	}
	return nil
}

func getProductByID(db *sql.DB, ID string) (*Product, error) {
	stmt, err := db.Prepare("SELECT id, name, price FROM products WHERE id=?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var p Product
	err = stmt.QueryRow(ID).Scan(&p.ID, &p.Name, &p.Price)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func getProductByIdWithContext(ctx context.Context, db *sql.DB, ID string) (*Product, error) {
	stmt, err := db.Prepare("SELECT id, name, price FROM products WHERE id=?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var p Product
	err = stmt.QueryRowContext(ctx, ID).Scan(&p.ID, &p.Name, &p.Price)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func selectProducts(db *sql.DB) ([]Product, error) {
	rows, err := db.Query("SELECT id, name, price FROM products")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var p Product
		err = rows.Scan(&p.ID, &p.Name, &p.Price)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}

func deleteProductByID(db *sql.DB, ID string) error {
	stmt, err := db.Prepare("DELETE FROM products WHERE id=?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(ID)
	if err != nil {
		return err
	}

	return nil
}
