// Copyright 2026 QuickWrite
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/QuickWrite/Coroki/internal/db"
	"github.com/QuickWrite/Coroki/internal/server"
	"github.com/joho/godotenv"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Loading the environment file. If the file cannot be found, the error should be ignored.
	_ = godotenv.Load()

	dsnCfg := db.Config{
		Host:     env("DB_HOST", "localhost"),
		Port:     envInt("DB_PORT", 5432),
		User:     env("DB_USER", "postgres"),
		Password: env("DB_PASSWORD", "postgres"),
		Name:     env("DB_NAME", "app"),
		SSLMode:  env("DB_SSLMODE", "disable"),
	}

	sqlDB, err := db.Connect(dsnCfg)
	if err != nil {
		log.Fatalf("db connect: %v", err)
	}
	defer sqlDB.Close()

	// Migrate this bad boy to the correct version
	err = db.Migrate(sqlDB)
	if err != nil {
		log.Fatalf("database couldn't be migrated: %v", err)
	}

	dependencies := createDependencies(sqlDB)
	handler := server.Routes(&dependencies)

	srv, err := server.Start(ctx, server.Config{
		Addr:         env("ADDR", ":8080"),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}, handler)

	if err != nil {
		log.Fatalf("server start: %v", err)
	}

	_ = srv
}

func env(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}

	return fallback
}

func envInt(key string, fallback int) int {
	if v, ok := os.LookupEnv(key); ok {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}

	return fallback
}
