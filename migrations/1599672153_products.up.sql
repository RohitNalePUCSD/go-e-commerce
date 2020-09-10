CREATE TABLE IF NOT EXISTS products (
	id INT NOT NULL PRIMARY KEY,
	name varchar(50)   NOT NULL,
description  varchar(50)   NOT NULL,
	price	float NOT NULL,
	discount   float NOT NULL,
quantity	int NOT NULL ,
category_Id	int NOT NULL ,
FOREIGN KEY(category_id) 
REFERENCES Category(id) ON DELETE CASCADE ON UPDATE CASCADE
);
