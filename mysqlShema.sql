CREATE TABLE books (
        id int  NOT NULL AUTO_INCREMENT ,
        title TEXT NOT NULL,
        autor TEXT NOT NULL,
        price DECIMAL(10,2) NOT NULL DEFAULT 0.00,
        CONSTRAINT books_bkey PRIMARY KEY (id)
    ) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;