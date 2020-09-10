CREATE TABLE IF NOT EXISTS productimages (
	product_id int NOT NULL,
	url varchar(20) NOT NULL,
	FOREIGN KEY(product_Id) 
REFERENCES products(id) ON DELETE CASCADE ON UPDATE CASCADE
);
