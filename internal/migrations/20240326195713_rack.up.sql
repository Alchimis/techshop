CREATE TABLE IF NOT EXISTS rack (
	id SERIAL PRIMARY KEY NOT NULL,
	name VARCHAR(64) NOT NULL
);

INSERT INTO rack(id, name)
VALUES 
	(default, 'А'),
	(default, 'Б'),
	(default, 'Ж'),
	(default, 'В'),
	(default, 'З');