package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

type Brand struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Model struct {
	ID      int    `json:"id"`
	BrandID int    `json:"brand_id"`
	Name    string `json:"name"`
}

type Car struct {
	ID        int    `json:"id"`
	BrandName string `json:"brand_name"`
	ModelName string `json:"model_name"`
	Year      int    `json:"year"`
	Price     int    `json:"price"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Ошибка загрузки .env файла")
	}

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", MainForm)
	http.HandleFunc("/search", SearchCars)
	fmt.Println("Сервер работает на http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func connectDB(ctx context.Context) (*pgx.Conn, error) {

	conn, err := pgx.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к базе данных: %v", err)
	}
	return conn, nil
}

func MainForm(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	conn, err := connectDB(ctx)
	if err != nil {
		http.Error(w, "Ошибка подключения к БД: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer conn.Close(ctx)

	brandsRows, err := getBrands(conn, ctx)
	if err != nil {
		http.Error(w, "Ошибка получения марок: "+err.Error(), http.StatusInternalServerError)
		return
	}
	modelsRows, err := getModels(conn, ctx)
	if err != nil {
		http.Error(w, "Ошибка получения моделей: "+err.Error(), http.StatusInternalServerError)
		return
	}

	Data := struct {
		Brands []Brand
		Models []Model
	}{
		Brands: brandsRows,
		Models: modelsRows,
	}

	//JSON
	var funcMap = template.FuncMap{
		"json": func(v interface{}) (string, error) {
			b, err := json.Marshal(v)
			return string(b), err
		},
	}

	tmpl := template.Must(template.New("").Funcs(funcMap).ParseFiles("templates/body.html", "templates/header.html"))

	tmpl.ExecuteTemplate(w, "body", Data)
}

func SearchCars(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	conn, err := connectDB(ctx)
	if err != nil {
		http.Error(w, "Ошибка подключения к БД: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer conn.Close(ctx)

	brandID := r.URL.Query().Get("brand")
	modelID := r.URL.Query().Get("model")
	year := r.URL.Query().Get("year")
	price := r.URL.Query().Get("price")

	cars, err := getCars(conn, ctx, brandID, modelID, year, price)
	if err != nil {
		http.Error(w, "Ошибка поиска автомобилей: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(cars)

}

func getBrands(connection *pgx.Conn, ctx context.Context) ([]Brand, error) {
	brandsRows, err := connection.Query(ctx, "SELECT id, name FROM brands")
	if err != nil {
		return nil, fmt.Errorf("ошибка загрузки марок: %w", err)
	}
	defer brandsRows.Close()

	var dataRows []Brand
	for brandsRows.Next() {
		var b Brand
		err := brandsRows.Scan(&b.ID, &b.Name)
		if err != nil {
			return nil, fmt.Errorf("ошибка сканирования марки: %v", err)
		}
		dataRows = append(dataRows, b)

	}

	//логи
	if err := brandsRows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при обработке строк: %v", err)
	}

	fmt.Printf("Отладка: загружено марок — %d\n", len(dataRows))
	return dataRows, nil

}

func getModels(connection *pgx.Conn, ctx context.Context) ([]Model, error) {
	modelsRows, err := connection.Query(ctx, "SELECT id, brand_id, name FROM models")
	if err != nil {
		return nil, fmt.Errorf("ошибка загрузки моделей: %w", err)
	}
	defer modelsRows.Close()

	var dataRows []Model
	for modelsRows.Next() {
		var m Model
		err := modelsRows.Scan(&m.ID, &m.BrandID, &m.Name)
		if err != nil {
			return nil, fmt.Errorf("ошибка сканирования модели: %v", err)
		}
		dataRows = append(dataRows, m)
	}

	//логи
	if err := modelsRows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при обработке строк: %v", err)
	}

	fmt.Printf("Отладка: загружено моделей — %d\n", len(dataRows))
	return dataRows, nil
}

func getCars(connection *pgx.Conn, ctx context.Context, brandID, modelID, year, price string) ([]Car, error) {
	log.Printf("Поиск автомобилей: brand=%s, model=%s, year=%s, price=%s", brandID, modelID, year, price)
	query := `
        SELECT c.id, b.name AS brand_name, m.name AS model_name, c.year, c.price
        FROM cars c
        JOIN models m ON c.model_id = m.id
        JOIN brands b ON m.brand_id = b.id
        WHERE 1=1
    `

	var args []interface{}
	var conditions []string

	if brandID != "" {
		args = append(args, brandID)
		conditions = append(conditions, fmt.Sprintf("AND b.id = $%d", len(args)))
	}

	if modelID != "" {
		args = append(args, modelID)
		conditions = append(conditions, fmt.Sprintf("AND m.id = $%d", len(args)))
	}

	if year != "" {
		args = append(args, year)
		conditions = append(conditions, fmt.Sprintf("AND c.year = $%d", len(args)))
	}

	if price != "" {
		args = append(args, price)
		conditions = append(conditions, fmt.Sprintf("AND c.price <= $%d", len(args)))
	}

	query += strings.Join(conditions, " ")

	rows, err := connection.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса: %w", err)
	}
	defer rows.Close()

	var cars []Car
	for rows.Next() {
		var c Car
		err := rows.Scan(&c.ID, &c.BrandName, &c.ModelName, &c.Year, &c.Price)
		if err != nil {
			return nil, fmt.Errorf("ошибка сканирования автомобиля: %w", err)
		}
		cars = append(cars, c)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при обработке строк: %w", err)
	}

	return cars, nil
}
