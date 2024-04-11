
<h1 align="center">go-authorization boilerplate</h1>

<div align="center">
 Echo + Gorm + Casbin  based scaffolding

<br/>
<br/>

<div align=center>
<img src="https://img.shields.io/badge/golang-1.22-blue"/>
<img src="https://img.shields.io/badge/echo-4.11.4-lightBlue"/>
<img src="https://img.shields.io/badge/gorm-1.25.9-red"/>
<img src="https://img.shields.io/badge/casbin-2.87.1-brightgreen"/>
</div>

<br/>
</div>

## Feature

* Follow RESTful API design specifications
* Provides rich middleware support based on Echo Api framework (jwt, authority, request-level transaction, access log, cors...)
* Casbin based RBAC access control model
* GORM based database storage that can expand many types of databases
* implements dependency injection pattern
* Support Swagger documentation (based on swaggo)
* Configuration, modularization

## Synopsis

welcome PR and Issue.

```
# readonly account
username: test
password: 123123
```

## Getting started

```
golang >= 1.22
```

**Use git to clone this project**

```
git clone https://github.com/sub-rat/go-authorization
```

**API docs generation**

```
make swagger
```

**Initialize the database and start the service**

```
make migrate # create tables
make setup # setup menu data
make # start
```