// Copyright 2017 Emir Ribic. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

// GORSK - Go(lang) restful starter kit
//
// API Docs for GORSK v1
//
// 	 Terms Of Service:  N/A
//     Schemes: http
//     Version: 1.0.0
//     License: MIT http://opensource.org/licenses/MIT
//     Contact: Emir Ribic <ribice@gmail.com> https://ribice.ba
//     Host: localhost:8080
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Security:
//     - bearer: []
//
//     SecurityDefinitions:
//     bearer:
//          type: apiKey
//          name: Authorization
//          in: header
//
// swagger:meta
package main

import (
	"gorsk/internal/platform/postgres"

	"github.com/labstack/echo"

	"gorsk/cmd/api/config"
	"gorsk/cmd/api/mw"
	"gorsk/cmd/api/server"
	"gorsk/cmd/api/service"
	_ "gorsk/cmd/api/swagger"
	"gorsk/internal/account"
	"gorsk/internal/auth"
	"gorsk/internal/rbac"
	"gorsk/internal/user"

	"github.com/go-pg/pg"
)

func main() {

	cfg, err := config.Load("dev")
	checkErr(err)

	e := server.New()

	db, err := pgsql.New(cfg.DB)
	checkErr(err)

	addV1Services(cfg, e, db)

	server.Start(e, cfg.Server)
}

func addV1Services(cfg *config.Configuration, e *echo.Echo, db *pg.DB) {

	// Initialize DB interfaces

	userDB := pgsql.NewUserDB(e.Logger)
	accDB := pgsql.NewAccountDB(e.Logger)

	// Initialize services

	jwt := mw.NewJWT(cfg.JWT)
	authSvc := auth.New(db, userDB, jwt)
	service.NewAuth(authSvc, e, jwt.MWFunc())

	e.Static("/swaggerui", "cmd/api/swaggerui")

	rbacSvc := rbac.New(userDB)

	v1Router := e.Group("/v1")

	v1Router.Use(jwt.MWFunc())

	// Workaround for Echo's issue with routing.
	// v1Router should be passed to service normally, and then the group name created there
	uR := v1Router.Group("/users")
	service.NewAccount(account.New(db, accDB, userDB, rbacSvc), uR)
	service.NewUser(user.New(db, userDB, rbacSvc, authSvc), uR)
}

func checkErr(err error) {
	if err != nil {
		panic(err.Error())
	}
}
