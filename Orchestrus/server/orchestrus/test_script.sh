#!/bin/bash
go run worker_mock/worker_mock.go &
go run db_connexion_mock/db_connexion_mock.go &
