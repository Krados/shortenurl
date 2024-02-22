# shortenurl

shortenurl is a service that can help you shorten URL from long to short.

## 說明

大家應該對短網址服務不陌生，例如: TinyURL、reurl 等等，使用者輸入長網址，後端服務產出一個短網址，使用者將產出的短網址輸入至瀏覽器，後端服務接收到請求並取出短網址對應的長網址，並返回 HTTP Status Code 301 or 302，如果返回的是 301 (Moved Permanently) 瀏覽器會快取這個短網址並把你導向至長網址的網站，如果是 302 (Found) 瀏覽器也會把你導向到長網址的網站， 301 跟 302 的差別在於快取的部分，如果你是返回 301 則下一次使用者在瀏覽器上輸入短網址時會直接讀取瀏覽器的快取並導向至長網址的網站不會再一次請求短網址的服務，但如果是 302 則一樣會請求短網址的服務。

## 流程

![shortener_put.drawio.png](https://github.com/Krados/shortenurl/blob/master/shortener_put.drawio.png)

使用短網址 code 反查長網址就相對簡單，code 查不到即返回錯誤，我就不另外畫圖了。

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
