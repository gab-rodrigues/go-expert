To run this project locally, follow these steps:

1. **Clone the Repository**:

2. Run the docker compose command to build and start the containers:

   ```bash
   docker compose up -d
   ```
   
  2.1 - Run mysql container
  ```bash
  Docker exec -it {container_mysql}
   ```

3. Create the database tables by running the following command:

   ```bash
   sqlite3 db.sqlite
   create table categories (id string, name string, description string);
   ```
4. Run go main.go wire.go to run the application:

   ```bash
   go run main.go wire.go
   ```
5. Use the file create_order.http to create orders using the web server (port 8000)
6. Execute evans -r repl to interact with the gRPC server and create orders (port 50051)

    ```bash
    evans -r repl
    ```

7. Use the graphql playground (http://localhost:8080/) to interact with the GraphQL server and create orders (port 8080)
```graphql
mutation createOrder{
  createOrder(input:{
    id: "order_4_gql",
    Price: 12.0,
    Tax: 2.0
  }) {
    id
    Price
    Tax
    FinalPrice
  }
}
```
8. Use the file list_orders.http to list orders using the web server
9. Execute evans -r repl to interact with the gRPC server and list orders and select the call listOrders
10. Use the graphql playground (http://localhost:8080/) to interact with the GraphQL server and list orders
```graphql
query ListAllOrders {
  orders {
    id
    Price
    Tax
    FinalPrice
  }
}
```