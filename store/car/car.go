package car

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/MarNawar/carZone/models"
	"github.com/google/uuid"
)

type Store struct {
	db *sql.DB
}

func New(db *sql.DB) Store {
	return Store{db: db}
}

func (s Store) GetCarById(ctx context.Context, id string) (models.Car, error) {
	var car models.Car
	query := `SELECT c.id, c.name, c.year, c.brand, c.fuel_type, c.engine_id, c.price, c.created_at, c.updated_at, e.id, e.displacement, e.no_of_cylinders, e.car_range FROM car c JOIN engine e ON c.engine_id = e.id WHERE c.id = $1`

	row := s.db.QueryRowContext(ctx, query, id)

	err := row.Scan(
		&car.ID,
		&car.Name,
		&car.Year,
		&car.Brand,
		&car.FuelType,
		&car.Engine.EngineID,
		&car.Price,
		&car.CreatedAt,
		&car.UpdatedAt,
		&car.Engine.EngineID,
		&car.Engine.Displacement,
		&car.Engine.NoOfCylinders,
		&car.Engine.CarRange,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return car, nil
		}
		return car, err
	}
	return car, nil
}

func (s Store) GetCarByBrand(ctx context.Context, brand string, isEngine bool) ([]models.Car, error) {
	var cars []models.Car
	var query string

	// Build query based on isEngine flag
	if isEngine {
		query = `
			SELECT 
				c.id, c.name, c.year, c.brand, c.fuel_type, c.engine_id, c.price, c.created_at, c.updated_at,
				e.id, e.displacement, e.no_of_cylinders, e.car_range
			FROM car c
			LEFT JOIN engine e ON c.engine_id = e.id
			WHERE c.brand = $1
		`
	} else {
		query = `
			SELECT 
				id, name, year, brand, fuel_type, engine_id, price, created_at, updated_at
			FROM car
			WHERE brand = $1
		`
	}

	// Execute query
	rows, err := s.db.QueryContext(ctx, query, brand)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch cars: %w", err)
	}
	defer rows.Close()

	// Iterate through result rows
	for rows.Next() {
		var car models.Car
		if isEngine {
			err := rows.Scan(
				&car.ID,
				&car.Name,
				&car.Year,
				&car.Brand,
				&car.FuelType,
				&car.Engine.EngineID,
				&car.Price,
				&car.CreatedAt,
				&car.UpdatedAt,
				&car.Engine.EngineID,
				&car.Engine.Displacement,
				&car.Engine.NoOfCylinders,
				&car.Engine.CarRange,
			)
			if err != nil {
				return nil, fmt.Errorf("failed to scan car with engine: %w", err)
			}
		} else {
			err := rows.Scan(
				&car.ID,
				&car.Name,
				&car.Year,
				&car.Brand,
				&car.FuelType,
				&car.Engine.EngineID,
				&car.Price,
				&car.CreatedAt,
				&car.UpdatedAt,
			)
			if err != nil {
				return nil, fmt.Errorf("failed to scan car: %w", err)
			}
		}

		// Append car to the result list
		cars = append(cars, car)
	}

	// Check for any errors during iteration
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return cars, nil
}


func (s Store) CreateCar(ctx context.Context, carReq *models.CarRequest) (models.Car, error) {
	var createdCar models.Car

	// Validate engine existence
	var engineExists bool
	err := s.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM engine WHERE id = $1)", carReq.Engine.EngineID).Scan(&engineExists)
	if err != nil {
		return createdCar, fmt.Errorf("failed to verify engine existence: %w", err)
	}
	if !engineExists {
		return createdCar, fmt.Errorf("engine with ID %s does not exist", carReq.Engine.EngineID)
	}

	// Prepare car data
	carID := uuid.New()
	currentTime := time.Now()

	newCar := models.Car{
		ID:        carID,
		Name:      carReq.Name,
		Year:      carReq.Year,
		Brand:     carReq.Brand,
		FuelType:  carReq.FuelType,
		Engine:    carReq.Engine,
		Price:     carReq.Price,
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
	}

	// Begin transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return createdCar, fmt.Errorf("failed to begin transaction: %w", err)
	}
	// Ensure proper rollback on error
	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	// Insert car into database
	query := `
		INSERT INTO car (id, name, year, brand, fuel_type, engine_id, price, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) 
		RETURNING id, name, year, brand, fuel_type, engine_id, price, created_at, updated_at
	`

	err = tx.QueryRowContext(
		ctx,
		query,
		newCar.ID,
		newCar.Name,
		newCar.Year,
		newCar.Brand,
		newCar.FuelType,
		newCar.Engine.EngineID,
		newCar.Price,
		newCar.CreatedAt,
		newCar.UpdatedAt,
	).Scan(
		&createdCar.ID,
		&createdCar.Name,
		&createdCar.Year,
		&createdCar.Brand,
		&createdCar.FuelType,
		&createdCar.Engine.EngineID,
		&createdCar.Price,
		&createdCar.CreatedAt,
		&createdCar.UpdatedAt,
	)

	if err != nil {
		return createdCar, fmt.Errorf("failed to create car: %w", err)
	}

	return createdCar, nil
}


func (s Store) UpdateCar(ctx context.Context, id string, carReq *models.CarRequest) (models.Car, error) {
	var updatedCar models.Car

	// Fetch existing car to validate ID
	var exists bool
	err := s.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM car WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return updatedCar, fmt.Errorf("failed to check car existence: %w", err)
	}
	if !exists {
		return updatedCar, fmt.Errorf("car with ID %s does not exist", id)
	}

	// Start building the dynamic query
	var queryBuilder strings.Builder
	queryBuilder.WriteString("UPDATE car SET ")
	var args []interface{}
	argID := 1

	if carReq.Name != "" {
		queryBuilder.WriteString(fmt.Sprintf("name = $%d, ", argID))
		args = append(args, carReq.Name)
		argID++
	}
	if carReq.Year != "" {
		queryBuilder.WriteString(fmt.Sprintf("year = $%d, ", argID))
		args = append(args, carReq.Year)
		argID++
	}
	if carReq.Brand != "" {
		queryBuilder.WriteString(fmt.Sprintf("brand = $%d, ", argID))
		args = append(args, carReq.Brand)
		argID++
	}
	if carReq.FuelType != "" {
		queryBuilder.WriteString(fmt.Sprintf("fuel_type = $%d, ", argID))
		args = append(args, carReq.FuelType)
		argID++
	}
	if carReq.Engine != (models.Engine{}) {
		queryBuilder.WriteString(fmt.Sprintf("engine_id = $%d, ", argID))
		args = append(args, carReq.Engine.EngineID)
		argID++
	}
	if carReq.Price != 0.0 {
		queryBuilder.WriteString(fmt.Sprintf("price = $%d, ", argID))
		args = append(args, carReq.Price)
		argID++
	}

	// Always update the updated_at field
	queryBuilder.WriteString(fmt.Sprintf("updated_at = $%d ", argID))
	args = append(args, time.Now())
	argID++

	// Add the WHERE clause
	queryBuilder.WriteString(fmt.Sprintf("WHERE id = $%d RETURNING id, name, year, brand, fuel_type, engine_id, price, created_at, updated_at", argID))
	args = append(args, id)

	// Begin transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return updatedCar, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	// Execute the query
	err = tx.QueryRowContext(ctx, queryBuilder.String(), args...).
		Scan(
			&updatedCar.ID,
			&updatedCar.Name,
			&updatedCar.Year,
			&updatedCar.Brand,
			&updatedCar.FuelType,
			&updatedCar.Engine.EngineID,
			&updatedCar.Price,
			&updatedCar.CreatedAt,
			&updatedCar.UpdatedAt,
		)

	if err != nil {
		return updatedCar, fmt.Errorf("failed to update car: %w", err)
	}

	return updatedCar, nil
}


func (s Store) DeleteCar(ctx context.Context, id string) (models.Car, error) {
	var deletedCar models.Car

	// Begin transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return deletedCar, fmt.Errorf("failed to begin transaction: %w", err)
	}
	// Ensure rollback on error
	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	// Delete and return the car details
	query := `
		DELETE FROM car 
		WHERE id = $1 
		RETURNING id, name, year, brand, fuel_type, engine_id, price, created_at, updated_at
	`
	err = tx.QueryRowContext(ctx, query, id).Scan(
		&deletedCar.ID,
		&deletedCar.Name,
		&deletedCar.Year,
		&deletedCar.Brand,
		&deletedCar.FuelType,
		&deletedCar.Engine.EngineID,
		&deletedCar.Price,
		&deletedCar.CreatedAt,
		&deletedCar.UpdatedAt,
	)

	// Handle error when no rows are affected
	if err == sql.ErrNoRows {
		return deletedCar, fmt.Errorf("car with ID %s does not exist", id)
	} else if err != nil {
		return deletedCar, fmt.Errorf("failed to delete car: %w", err)
	}

	return deletedCar, nil
}

