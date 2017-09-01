# README



# my.cnf

> https://github.com/go-sql-driver/mysql/issues/583

```
net_read_timeout         = 14400
net_write_timeout        = 14400

wait_timeout             = 28800
interactive_timeout      = 28800
```

net_read_timeout ：数据读取超时时间。在终止读之前，从一个连接获得数据而等待的时间秒数；当服务正在从客户端读取数据时，net_read_timeout控制何时超时。即客户端执行数据读取，等待多少秒仍未执行成功时自动断开连接。 
net_write_timeout ：数据库写超时时间。和net_read_timeout意义类似，在终止写之前，等待多少秒把block写到连接；当服务正在写数据到客户端时，net_write_timeout控制何时超时。
wait_timeout：连接空闲超时时间。与服务器端无交互状态的连接，直到被服务器端强制关闭而等待的时间。可以认为是服务器端连接空闲的时间，空闲超过这个时间将自动关闭。
interactive_timeout ：连接空闲超时时间。与服务器端无交互状态的连接，直到被服务器端强制关闭而等待的时间。
interactive_timeout和wait_timeoutu意义虽然相同，但是有使用对象有本质的区别。interactive_timeout针对交互式连接（比如通过mysql客户端连接数据库），wait_timeout针对非交互式连接（比如一般在PHP中使用PDO连接数据库，当然你可以设置CLIENT_INTERACTIVE选项来改变）。所谓的交互式连接，即在mysql_real_connect()函数中使用了CLIENT_INTERACTIVE选项。


```
db.SetMaxIdleConns(10)                    //连接池中最大空闲连接数
db.SetMaxOpenConns(150)                   //打开的最大连接数, 注意要小于mysql配置文件my.cnf中的配置
db.SetConnMaxLifetime(7200 * time.Second) //连接的最大空闲时间(可选)
```