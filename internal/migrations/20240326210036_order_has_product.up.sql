CREATE TABLE IF NOT EXISTS order_has_product (
	order_id INT NOT NULL REFERENCES client_order(id),
	product_id INT NOT NULL REFERENCES product(id),
	quantity INT NOT NULL
);


INSERT INTO order_has_product(order_id, product_id, quantity)
VALUES 
	(10, 1, 2),
	(10, 3, 1),
	(10, 6, 1),
	(11, 2, 3),
	(14, 1, 3),
	(14, 4, 4),
	(15, 5, 1);