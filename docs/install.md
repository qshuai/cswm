Prerequisites
===
- Go environment: https://go.dev/
- Mysql (docker recommended for development)
- Redis (docker recommended for development)

Install
===
1. Setup mysql
```shell
$ mkdir -p conf/conf.d
$ echo -e "[mysqld]\nmysql-native-password = ON" > conf/conf.d/my.cnf

$ docker run --name mysql --privileged=true -d -p 3306:3306 --restart unless-stopped -v mysql_data:/var/lib/mysql -v $(pwd)/conf:/etc/mysql -e MYSQL_ROOT_PASSWORD=your-mysql-password mysql:8.4
```

2. Setup redis
```shell
$ docker run -d --name redis -p 6379:6379 redis --requirepass "your-redis-password"
```


3. Check the status of mysql and redis:
```shell
$ docker ps
CONTAINER ID   IMAGE       COMMAND                   CREATED        STATUS       PORTS                                                  NAMES
8947cd86bde0   mysql:8.4   "docker-entrypoint.s…"   2 hours ago    Up 2 hours   0.0.0.0:3306->3306/tcp, :::3306->3306/tcp, 33060/tcp   mysql
ac056a59cf6e   redis       "docker-entrypoint.s…"   18 hours ago   Up 4 hours   0.0.0.0:6379->6379/tcp, :::6379->6379/tcp              redis
```

4. Config mysql:
```shell
$ docker exec -it mysql /bin/bash
$ mysql -p
Enter password: your-mysql-password
mysql> ALTER USER 'root'@'%' IDENTIFIED WITH mysql_native_password BY 'your-mysql-password';
mysql> CREATE DATABASE `erp`;
mysql> exit
$ exit
```

Now, migrate mysql table.
```shell
cd your-cwsm-repository-directory
go run .
```
Wait a moment before cwsm completing initiation and stop it by press `Ctrl + C`

Prepare genesis data:
```shell
$ docker exec -it mysql /bin/bash
$ mysql -p
Enter password: your-mysql-password
mysql> insert into user (username, password, name, tel, position, created, updated) VALUES ('admin', '21232f297a57a5a743894a0e4a801fc3', '管理员', '13051666281', '超级管理员', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);
mysql> insert into permission (`user_id`, `add_member`, `edit_member`, `active_member`, `add_consumer`, `edit_consumer`, `view_consumer`, `add_brand`, `add_dealer`, `view_dealer`, `add_supplier`, `view_supplier`, `add_product`, `input_in_price`, `view_product_store`, `view_stock`, `view_in_price`, `edit_product`, `delete_product`, `output_product`, `view_sale`, `view_sale_consumer`, `view_sale_in_price`, `edit_sale`, `operate_category`, `request_move`, `response_move`, `view_move`, `add_store`, `view_store`, `operate_other_store`) VALUES (1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1);
mysql> exit
$ exit
```

5. Run & browse
```shell
cd your-cwsm-repository-directory
cp conf/app-sample.conf conf/app.conf
# 修改配置内容：app.conf（包括mysql、redis配置等）
go run .
```
Open http://localhost:9090 in browser and enjoy it.
