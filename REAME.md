# Mongodb and mongo-express

All config and code in the docker-compose.yml
You need to create database **master-music** before run code

mongo express [Mongo Express](http://localhost:8081/) with user/password are admin/pass

# API

Configuration in .env file
if you change some thing in config, you need to run code behind to build new image

``` bash
make build-image
```

When all done, you can run code behind to run all image

``` bash
docker compose up
```

When you run success
you can access the swagger [Swaggerui](http://localhost:8191/swaggerui/#/)

If you want to check api
You can use request.http file