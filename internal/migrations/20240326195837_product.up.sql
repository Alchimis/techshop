CREATE TABLE IF NOT EXISTS product (
	id SERIAL PRIMARY KEY NOT NULL,
	title VARCHAR(64) NOT NULL,
	main_rack_id INT NOT NULL REFERENCES rack(id)
);

INSERT INTO product(id, title, main_rack_id)
VALUES 
	(1, 'Ноутбук', 1),
	(2, 'Телевизор', 1),
	(3, 'Телефон', 2),
	(4, 'Системный блок', 3),
	(5, 'Часы', 3),
	(6, 'Микрофон', 3);