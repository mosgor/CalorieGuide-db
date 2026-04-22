CREATE TABLE client(
    id SERIAL NOT NULL,
    user_name varchar(125),
    surname varchar(124),
    email varchar(255),
    password varchar(123),
    picture integer[],
    PRIMARY KEY(id)
);

INSERT INTO client (user_name, surname, email, password)
VALUES ('Максим', 'Михайлов', 'example@list.ru', '123');

CREATE TABLE food(
    id SERIAL NOT NULL,
    food_name varchar(255) NOT NULL,
    description varchar(255),
    calories integer DEFAULT 0,
    proteins integer DEFAULT 0,
    carbohydrates integer DEFAULT 0,
    fats integer DEFAULT 0,
    likes integer DEFAULT 0,
    author_id integer NOT NULL,
    picture integer[],
    PRIMARY KEY(id),
    CONSTRAINT food_author_id_fkey FOREIGN key(author_id) REFERENCES client(id)
);

INSERT INTO food (food_name, description, calories, proteins, carbohydrates, fats, likes, author_id)
VALUES ('Огурец', 'Базовый овощь', 100, 200, 400, 200, 1, 1);

CREATE TABLE meal(
    id SERIAL NOT NULL,
    meal_name varchar(124) NOT NULL,
    total_calories integer DEFAULT 0,
    total_proteins integer DEFAULT 0,
    total_fats integer DEFAULT 0,
    total_carbohydrates integer DEFAULT 0,
    author_id integer NOT NULL,
    description varchar(255),
    likes integer DEFAULT 0,
    picture integer[],
    PRIMARY KEY(id),
    CONSTRAINT meal_author_id_fkey FOREIGN key(author_id) REFERENCES client(id)
);

INSERT INTO meal (meal_name, description, total_calories, total_proteins, total_carbs, total_fats, likes, author_id)
VALUES ('Битый огурец', 'Странное блюдо', 100, 200, 400, 200, 1, 1);

CREATE TABLE goal(
    id SERIAL NOT NULL,
    calories_goal integer,
    fats_goal integer,
    proteins_goal integer,
    carbohydrates_goal integer,
    PRIMARY KEY(id)
);

INSERT INTO goal (calories_goal, fats_goal, proteins_goal, carbohydrates_goal)
VALUES (1000, 1000, 1000, 1000);

CREATE TABLE diet(
    id SERIAL NOT NULL,
    breakfast_id integer,
    lunch_id integer,
    dinner_id integer,
    PRIMARY KEY(id),
    CONSTRAINT diet_breakfast_id_fkey FOREIGN key(breakfast_id) REFERENCES meal(id),
    CONSTRAINT diet_lunch_id_fkey FOREIGN key(lunch_id) REFERENCES meal(id),
    CONSTRAINT diet_dinner_id_fkey FOREIGN key(dinner_id) REFERENCES meal(id)
);

INSERT INTO diet (breakfast_id, lunch_id, dinner_id)
VALUES (1, 1, 1);

CREATE TABLE meal_food(
    id SERIAL NOT NULL,
    meal_id integer,
    food_id integer,
    quantity integer,
    PRIMARY KEY(id),
    CONSTRAINT meal_food_meal_id_fkey FOREIGN key(meal_id) REFERENCES meal(id),
    CONSTRAINT meal_food_food_id_fkey FOREIGN key(food_id) REFERENCES food(id)
);

INSERT INTO meal_food (meal_id, food_id, quantity)
VALUES (1, 1, 1);

CREATE TABLE meal_client(
    id SERIAL NOT NULL,
    meal_id integer,
    user_id integer,
    PRIMARY KEY(id),
    CONSTRAINT meal_client_meal_id_fkey FOREIGN key(meal_id) REFERENCES meal(id),
    CONSTRAINT meal_client_user_id_fkey FOREIGN key(user_id) REFERENCES client(id)
);

INSERT INTO meal_client (meal_id, user_id)
VALUES (1, 1);

CREATE TABLE food_client(
    id SERIAL NOT NULL,
    food_id integer,
    user_id integer,
    PRIMARY KEY(id),
    CONSTRAINT food_client_food_id_fkey FOREIGN key(food_id) REFERENCES food(id),
    CONSTRAINT food_client_user_id_fkey FOREIGN key(user_id) REFERENCES client(id)
);

INSERT INTO food_client (food_id, user_id)
VALUES (1, 1);