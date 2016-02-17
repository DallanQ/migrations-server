#!/usr/bin/env bash
git subtree pull --prefix _vendor/src/github.com/jtolds/gls https://github.com/jtolds/gls master --squash
git subtree pull --prefix _vendor/src/github.com/smartystreets/assertions https://github.com/smartystreets/assertions master --squash
git subtree pull --prefix _vendor/src/github.com/smartystreets/goconvey https://github.com/smartystreets/goconvey master --squash

git subtree pull --prefix _vendor/src/github.com/go-sql-driver/mysql https://github.com/go-sql-driver/mysql master --squash

git subtree pull --prefix _vendor/src/github.com/jmoiron/sqlx https://github.com/jmoiron/sqlx master --squash

git subtree pull --prefix _vendor/src/github.com/kelseyhightower/envconfig https://github.com/kelseyhightower/envconfig master --squash

git subtree pull --prefix _vendor/src/github.com/bradfitz/http2 https://github.com/bradfitz/http2 master --squash
git subtree pull --prefix _vendor/src/github.com/labstack/gommon https://github.com/labstack/gommon master --squash
git subtree pull --prefix _vendor/src/golang.org/x/net https://github.com/golang/net master --squash
git subtree pull --prefix _vendor/src/github.com/labstack/echo https://github.com/labstack/echo master --squash
git subtree pull --prefix _vendor/src/github.com/rs/cors https://github.com/rs/cors master --squash
