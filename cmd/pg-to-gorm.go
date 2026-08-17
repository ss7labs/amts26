package main

import (
	"gorm.io/driver/postgres"
	"gorm.io/gen"
	"gorm.io/gorm"
)

func main() {
	// Initialize GORM Gen configurations
	g := gen.NewGenerator(gen.Config{
		// Directory path where generated structs will be saved
		OutPath: "./internal/models",
		// Generator options: generate type-safe DAO APIs or context queries
		Mode: gen.WithoutContext | gen.WithDefaultQuery,
	})

	// Establish a temporary connection to your PostgreSQL database
	dsn := "host=localhost user=postgres password=pwd dbname=amts port=5434 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database: " + err.Error())
	}

	// Tell the generator to use this active database connection
	g.UseDB(db)

	// Automatically map and generate structs for ALL existing tables
	g.ApplyBasic(g.GenerateAllTable()...)

	// Execute code generation
	g.Execute()
}
