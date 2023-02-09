# shortenurl

shortenurl is a service that can help you shorten URL from long to short.

## Prerequisites

### Install MySQL, Redis

```
$ docker run -itd --name mysql -p 3306:3306 -e MYSQL_ROOT_PASSWORD={YOUROOTPASSWD} mysql

$ docker run -itd --name redis -p 6379:6379 redis
```

### Create database and table

SQL:
```
# create database
CREATE DATABASE shortenurl;

# create table
CREATE TABLE IF NOT EXISTS `short_url` (
  `id` bigint NOT NULL,
  `code` varchar(255) NOT NULL,
  `hash_url` varchar(255) NOT NULL,
  `long_url` text NOT NULL,
  PRIMARY KEY (`id`),
  KEY `short_url` (`code`),
  KEY `hash_url` (`hash_url`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

You can SSH into container and execute the above sql

```
$ docker exec -it mysql sh
$ mysql -u root -p
$ CREATE DATABASE shortenurl;
$ USE shortenurl;
$ CREATE TABLE IF NOT EXISTS `short_url` (
  `id` bigint NOT NULL,
  `code` varchar(255) NOT NULL,
  `hash_url` varchar(255) NOT NULL,
  `long_url` text NOT NULL,
  PRIMARY KEY (`id`),
  KEY `short_url` (`code`),
  KEY `hash_url` (`hash_url`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

## Build the service

```
$ make build
```

## Run the service

```
$ .\bin\cmd.exe -conf .\configs\
or 
$ ./bin/cmd -conf ./configs
```

## Run the service inside a container

```
$ make docker-image
$ docker run -itd --name shortenurl -p 8080:8080 shortenurl
```

## Network isolation is important then do like this

```
# create a user-defined bridge network
$ docker network create asgard

$ docker run -itd --name mysql --network asgard -e MYSQL_ROOT_PASSWORD={YOUROOTPASSWD} mysql
$ docker run -itd --name redis --network asgard redis
$ docker run -itd --name shortenurl --network asgard -p 8080:8080 shortenurl
```

## Try it out

```
# create a new code
$ curl --location --request POST 'localhost:8080/api/v1/shorten' \
--header 'Content-Type: application/json' \
--data-raw '{
    "long_url":"http://www.google.com"
}'


# get url from code
$ curl --location --request GET 'localhost:8080/api/v1/shorten/4LAoxKG2tAY'
```