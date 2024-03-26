CREATE TABLE IF NOT EXISTS rack_has_product (
	rack_id INT NOT NULL REFERENCES rack(id),
	product_id INT NOT NULL REFERENCES product(id),
	PRIMARY KEY(rack_id, product_id)
);

INSERT INTO rack_has_product(rack_id, product_id)
VALUES 
	(5, 3),
	(4, 3),
	(1, 5);