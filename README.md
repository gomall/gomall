# GoMall  Application

The base project is cloned from  [qor-example](https://github.com/qor/qor-example).


## Quick Started

```shell
# Get example app
$ go get -u github.com/gomall/gomall

# Setup database
$ mysql -uroot -p
mysql> CREATE DATABASE qor_example;

# Run Application
$ cd $GOPATH/src/github.com/gomall/gomall
$ go run main.go
```

#### Generate sample data

```go
$ go run db/seeds/main.go
```

## Admin Management Interface

[Qor Example admin configuration](https://github.com/gomall/gomall/blob/master/config/admin/admin.go)


## RESTful API

[Qor Example API configuration](https://github.com/gomall/gomall/blob/master/config/api/api.go)



## License

Released under the MIT License.
