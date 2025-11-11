CREATE TABLE client
(
	id SERIAL PRIMARY KEY,
	user_name VARCHAR(125), 
	surname VARCHAR(124), 
	email VARCHAR(255), 
	password VARCHAR(123), 
	picture VARCHAR(155)
);

INSERT INTO client (user_name, surname, email, password)
VALUES ('Максим', 'Михайлов', 'example@list.ru', '123');

CREATE TABLE food
(
    id SERIAL PRIMARY KEY,
	food_name VARCHAR(255), 
	description VARCHAR(255), 
	calories INT, 
	proteins INT, 
	carbohydrates INT, 
	fats INT, 
	likes INT, 
	author_id INT, 
	picture VARCHAR(155)
);

INSERT INTO food (food_name, description, calories, proteins, carbohydrates, fats, likes, author_id)
VALUES ('Огурец', 'Базовый овощь', 100, 200, 400, 200, 100, 1);