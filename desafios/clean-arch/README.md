To run this project locally, follow these steps:

1. **Clone the Repository**:

2. Run the docker compose command to build and start the containers:

   ```bash
   docker compose up
   ```

3. Use the file create_order.http to create orders using the web server (port 8000)
4. Execute evans -r repl to interact with the gRPC server and create orders (port 50051)

    ```bash
    evans -r repl
    ```

5. Use the graphql playground (http://localhost:8080/) to interact with the GraphQL server and create orders (port 8080)
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
6. Use the file list_orders.http to list orders using the web server
7. Execute evans -r repl to interact with the gRPC server and list orders and select the call listOrders
8. Use the graphql playground (http://localhost:8080/) to interact with the GraphQL server and list orders
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