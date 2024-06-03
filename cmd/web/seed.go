package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"

	"github.com/dimitargrozev5/expenses-go-1/internal/config"
	"github.com/dimitargrozev5/expenses-go-1/internal/driver"
	"github.com/dimitargrozev5/expenses-go-1/internal/models"
	"github.com/dimitargrozev5/expenses-go-1/internal/repository/dbrepo"
	"golang.org/x/crypto/bcrypt"
)

func Seed(DBPath string) {

	// Get DB name
	userEmil := "asd@asd.asd"
	dbName := DBPath + userEmil

	// Migrate db
	err := Migrate(dbName)
	if err != nil {
		log.Fatal(err)
	}

	// Open db
	db, err := sql.Open("sqlite3", fmt.Sprintf("%s.db?_fk=%s", dbName, url.QueryEscape("true")))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Seed user
	stmt := `INSERT INTO user (email, password, db_version) VALUES ($1, $2, $3)`

	// Get password
	password, _ := bcrypt.GenerateFromPassword([]byte("asd"), bcrypt.DefaultCost)

	// Execute query
	_, err = db.Exec(stmt, "asd@asd.asd", password, 0)
	if err != nil {
		log.Printf("%q: %s\n", err, stmt)
		return
	}

	// Create user connection
	dbconn, err := driver.ConnectSQL(dbrepo.GetUserDBPath(DBPath, userEmil, false))
	if err != nil {
		log.Fatal(err)
	}

	// Get db repo
	repo := dbrepo.NewSqliteRepo(&config.AppConfig{}, userEmil, dbconn.SQL)

	// Add accounts
	repo.AddAccount("ДСК Спестовна") // 1
	repo.AddAccount("ЦКБ")           // 2
	repo.AddAccount("ДСК")           // 3
	repo.AddAccount("Кеш")           // 4

	// Add free funds
	repo.ModifyFreeFunds(1469, 4, "init")
	repo.ModifyFreeFunds(20226, 3, "init")
	repo.ModifyFreeFunds(432, 2, "init")
	repo.ModifyFreeFunds(4573, 1, "init")

	// Add categories
	repo.AddCategory("ЦКБ От Татко", 0, 0, 999, 1)            //  1
	repo.AddCategory("ДСК Спестявания", 0, 0, 999, 1)         //  2
	repo.AddCategory("От Къща", 0, 0, 999, 1)                 //  3
	repo.AddCategory("Спестявания", 300, 100, 1, 2)           //  4
	repo.AddCategory("Инвестиции", 100, 100, 1, 2)            //  5
	repo.AddCategory("Други", 120, 120, 1, 2)                 //  6
	repo.AddCategory("Техника", 70, 0, 1, 2)                  //  7
	repo.AddCategory("Дрехи", 80, 80, 1, 2)                   //  8
	repo.AddCategory("Здраве", 50, 50, 1, 2)                  //  9
	repo.AddCategory("Забавления", 50, 50, 1, 2)              // 10
	repo.AddCategory("Почивки", 200, 200, 1, 2)               // 11
	repo.AddCategory("Застраховка Автомобил", 100, 100, 1, 2) // 12
	repo.AddCategory("Автомобил", 200, 200, 1, 2)             // 13
	repo.AddCategory("Сметки", 300, 300, 1, 2)                // 14
	repo.AddCategory("Наем", 450, 350, 1, 2)                  // 15
	repo.AddCategory("Гориво", 200, 200, 1, 2)                // 16
	repo.AddCategory("Потреби", 80, 80, 1, 2)                 // 17
	repo.AddCategory("Ресторанти", 150, 150, 1, 2)            // 18
	repo.AddCategory("Храна", 550, 550, 1, 2)                 // 19

	// Setup reset values
	resetValues := []models.ResetCategoryData{
		{Amount: 779, CategoryId: 1, BudgetInput: 0, SpendingLimit: 0, InputInterval: 999, InputPeriod: 1},
		{Amount: 2714, CategoryId: 2, BudgetInput: 0, SpendingLimit: 0, InputInterval: 999, InputPeriod: 1},
		{Amount: 10000, CategoryId: 3, BudgetInput: 0, SpendingLimit: 0, InputInterval: 999, InputPeriod: 1},
		{Amount: 3260, CategoryId: 4, BudgetInput: 300, SpendingLimit: 100, InputInterval: 1, InputPeriod: 2},
		{Amount: 1150, CategoryId: 5, BudgetInput: 100, SpendingLimit: 100, InputInterval: 1, InputPeriod: 2},
		{Amount: 317, CategoryId: 6, BudgetInput: 120, SpendingLimit: 120, InputInterval: 1, InputPeriod: 2},
		{Amount: 210, CategoryId: 7, BudgetInput: 70, SpendingLimit: 0, InputInterval: 1, InputPeriod: 2},
		{Amount: 385, CategoryId: 8, BudgetInput: 80, SpendingLimit: 80, InputInterval: 1, InputPeriod: 2},
		{Amount: 734, CategoryId: 9, BudgetInput: 50, SpendingLimit: 50, InputInterval: 1, InputPeriod: 2},
		{Amount: 14, CategoryId: 10, BudgetInput: 50, SpendingLimit: 50, InputInterval: 1, InputPeriod: 2},
		{Amount: 953, CategoryId: 11, BudgetInput: 200, SpendingLimit: 200, InputInterval: 1, InputPeriod: 2},
		{Amount: 1403, CategoryId: 12, BudgetInput: 100, SpendingLimit: 100, InputInterval: 1, InputPeriod: 2},
		{Amount: 1163, CategoryId: 13, BudgetInput: 200, SpendingLimit: 200, InputInterval: 1, InputPeriod: 2},
		{Amount: 1647, CategoryId: 14, BudgetInput: 300, SpendingLimit: 300, InputInterval: 1, InputPeriod: 2},
		{Amount: 839, CategoryId: 15, BudgetInput: 450, SpendingLimit: 350, InputInterval: 1, InputPeriod: 2},
		{Amount: 152, CategoryId: 16, BudgetInput: 200, SpendingLimit: 200, InputInterval: 1, InputPeriod: 2},
		{Amount: 137, CategoryId: 17, BudgetInput: 80, SpendingLimit: 80, InputInterval: 1, InputPeriod: 2},
		{Amount: 743, CategoryId: 18, BudgetInput: 150, SpendingLimit: 150, InputInterval: 1, InputPeriod: 2},
		{Amount: 100, CategoryId: 19, BudgetInput: 550, SpendingLimit: 550, InputInterval: 1, InputPeriod: 2},
	}

	// Reset categories
	err = repo.ResetCategories(resetValues)
	if err != nil {
		log.Fatal(err)
	}
}
