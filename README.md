This repo contains a simple blog web app. The app has two modules: category and post. 
You can create, update and delete a category. When creating a post you will need to give a title, author, select a category and provide the content.
You can read, craete, edit and delete a post.
You can also search for a post based on the title.  
gRPC was used to communicate with the services.
Please change the environment variables of the server and client accordingly in the config file.  
Find the config file of the server navigate to **todo/env**.  
Find the config file of the client navigate to **cms/env**.  
To run the DB migration:  
```
cd todo
go run migrations/migrate.go up
```
To run the server: *go run todo/main.go*  
To run the client: *go run cms/main.go*  
To view the web app in the browser go to: localhost:8080 (or the port number you specified in the config file).  
