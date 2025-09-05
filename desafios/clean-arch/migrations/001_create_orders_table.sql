USE orders;

CREATE TABLE IF NOT EXISTS orders (
    id varchar(255) NOT NULL,
    price float NOT NULL,
    tax float NOT NULL,
    final_price float NOT NULL,
    PRIMARY KEY (id)
);

INSERT INTO orders (id, price, tax, final_price)
VALUES ('migration-test-001', 100.0, 20.0, 120.0);