

## Run locally
Golang should be installed in your local machine
    https://golang.org/doc/install

MySql should be installed in your machine. If you have not installed MySQL please install it.
 
 Create database books
 
 MySQL table that our CRUD operations operate on:
 
    CREATE TABLE books (
        id int  NOT NULL AUTO_INCREMENT ,
        title TEXT NOT NULL,
        author TEXT NOT NULL,
        price DECIMAL(10,2) NOT NULL DEFAULT 0.00,
        CONSTRAINT books_bkey PRIMARY KEY (id)
    ) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;

In IDE the database console, create a table
    or you can use my data for review.

Make sure this table is up and running.

Сlone this repository to the machine

Open project in IDE

In file .env update your database details

    DB_HOST=127.0.0.1 #Default mysql host
    DB_DRIVER = mysql 
    # API_SECRET=98hbun98h #Used for creating a JWT. Can be anything 
    DB_USER={name}
    DB_PASSWORD={password}
    DB_NAME=books
    DB_PORT=3306 #Default mysql port
or you can use my data for review

Run the project by the following the command:

    go run main.go

Now the server is running on port 8002.


Test the routes in postman tool.

https://www.getpostman.com

Make sure you run the server:
    go run main.go

Getting all books(/books)

 GET  http://localhost:8002/books 

    example response :
 
     [
        {
            "id": 2,
            "title": "Anna Karenina",
            "author": "Lev Tolstoy",
            "price": 47
        },
        {
            "id": 3,
            "title": "Alice in Wonderland",
            "author": "Lewis Karoll",
            "price": 28
        },
     ]

Getting one bookat id (/book/{id})

 GET  http://localhost:8002/book/9 

    example response :{"id":9,"title":"Язык GO","author":"Керниган, Ритчи","price":1585}


Create new book in DB (/book)

 POST  http://localhost:8002/books 

    example reqest:
    Body/row json /{"title":"Nose","author":"M.Gogol", "price": 18.47}
    example response:19

Update book at id (/book/{id})

 PUT  http://localhost:8002/book/7 

    example reqest:
    Body/row json /{"title":"Gone with the Wind","author":"M.Mitchell", "price": 18.47}
    example response:{"id":7,"title":"Gone with the Wind","author":"M.Mitchell","price":18.47}

Delete book at id (/book/{id})

 DELETE  http://localhost:8002/book/7 

    example reqest:{"result":"success"}


## TESTS

## models_test

Running all tests in the models_test package:

To run the test suite in the modeltests package, make sure in your terminal, you are in the path:

    $ cd models

Then run:

    go test -v

 All tests in the models_test package are run and all passed.
 
    We are connected to the mysql database
    === RUN   TestCreateBook
    --- PASS: TestCreateBook (0.01s)
    === RUN   TestGetAllBooks
    books [{1 testBook testAuthor 125.2}]
    books [{1 testBook testAuthor 125.2} {2 testBook testAuthor 125}]
    --- PASS: TestGetAllBooks (0.01s)
    === RUN   TestGetBook
    --- PASS: TestGetBook (0.02s)
    === RUN   TestUpdateBook
    --- PASS: TestUpdateBook (0.02s)
    === RUN   TestDeleteBook
    --- PASS: TestDeleteBook (0.02s)
    PASS
    ok      github.com/NM211077/testTask_techTechnorely/models      0
    
  
  ## controllers_test
  
  Running all tests in the controllers_test package:
  
  To run the test suite in the controllers_test package, make sure in your terminal, you are in the path:
  
      $ cd controllers
  
  Then run:
  
      go test -v
  
   All tests in the models_test package are run and all passed.
      
      We are connected to the mysql database
      === RUN   TestCreateProduct
      --- PASS: TestCreateProduct (0.02s)
      === RUN   TestGetBooks
      books [{1 testTitle testAuthor 12.25}]
      --- PASS: TestGetBooks (0.00s)
      === RUN   TestGetBookByID
      --- PASS: TestGetBookByID (0.00s)
      === RUN   TestUpdateBookByID
      --- PASS: TestUpdateBookByID (0.02s)
      === RUN   TestDeleteBookByID
      --- PASS: TestDeleteBookByID (0.11s)
      PASS
      ok      github.com/NM211077/testTask_techTechnorely/controllers 0.276s