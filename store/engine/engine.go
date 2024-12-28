package engine

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/MarNawar/carZone/models"
	"github.com/google/uuid"
)

type EngineStore struct {
	db *sql.DB
}

func New(db *sql.DB) *EngineStore {
	return &EngineStore{db: db}
}

func (e EngineStore) EngineById(ctx context.Context, id string) (models.Engine, error) {
	var engine models.Engine

	query := `
		SELECT 
			engine_id, displacement, no_of_cylinders, car_range
		FROM engine
		WHERE engine_id = $1
	`

	// Use QueryRowContext for single row retrieval
	err := e.db.QueryRowContext(ctx, query, id).Scan(
		&engine.EngineID,
		&engine.Displacement,
		&engine.NoOfCylinders,
		&engine.CarRange,
	)

	// Handle errors
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return engine, fmt.Errorf("engine with ID %s does not exist", id)
		}
		return engine, fmt.Errorf("failed to fetch engine: %w", err)
	}

	return engine, nil
}

func (e EngineStore) CreateEngine(ctx context.Context, engineReq *models.EngineRequest) (models.Engine, error) {
	var createdEngine models.Engine
	engineID := uuid.New()

	newEngine := models.Engine{
		EngineID:      engineID,
		Displacement:  engineReq.Displacement,
		NoOfCylinders: engineReq.NoOfCylinders,
		CarRange:      engineReq.CarRange,
	}

	// Begin transaction
	tx, err := e.db.BeginTx(ctx, nil)
	if err != nil {
		return newEngine, fmt.Errorf("failed to begin transaction: %w", err)
	}
	// Ensure proper rollback on error
	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	// Insert engine into database
	query := `
			INSERT INTO engine (engine_id, displacement, noOfCylinders, carRange) 
			VALUES ($1, $2, $3, $4) 
			RETURNING engine_id, displacement, noOfCylinders, carRange
		`

	err = tx.QueryRowContext(
		ctx,
		query,
		newEngine.EngineID,
		newEngine.Displacement,
		newEngine.NoOfCylinders,
		newEngine.CarRange,
	).Scan(
		&createdEngine.EngineID,
		&createdEngine.Displacement,
		&createdEngine.NoOfCylinders,
		&createdEngine.CarRange,
	)

	if err != nil {
		return newEngine, fmt.Errorf("failed to create engine:  %w", err)
	}

	return createdEngine, nil
}

func (e EngineStore) EngineUpdate(ctx context.Context, id string, engineReq *models.EngineRequest) (models.Engine, error) {
	var updatedEngine models.Engine

	// Fetch existing car to validate ID
	var exists bool
	err := e.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM engine WHERE engine_id = $1)", id).Scan(&exists)
	if err != nil {
		return updatedEngine, fmt.Errorf("failed to check engine existence: %w", err)
	}
	if !exists {
		return updatedEngine, fmt.Errorf("engine with ID %s does not exist", id)
	}

	// Start building the dynamic query
	var queryBuilder strings.Builder
	queryBuilder.WriteString("UPDATE engine SET ")
	var args []interface{}
	argID := 1

	if engineReq.Displacement != 0 {
		queryBuilder.WriteString(fmt.Sprintf("displacement = $%d, ", argID))
		args = append(args, engineReq.Displacement)
		argID++
	}
	if engineReq.NoOfCylinders != 0 {
		queryBuilder.WriteString(fmt.Sprintf("noOfCylinders = $%d, ", argID))
		args = append(args, engineReq.NoOfCylinders)
		argID++
	}
	if engineReq.CarRange != 0 {
		queryBuilder.WriteString(fmt.Sprintf("carRange = $%d, ", argID))
		args = append(args, engineReq.CarRange)
		argID++
	}

	// Add the WHERE clause
	queryBuilder.WriteString(fmt.Sprintf(" WHERE engine_id = $%d RETURNING engine_id, displacement, noOfCylinders, carRange", argID))
	args = append(args, id)

	// Begin transaction
	tx, err := e.db.BeginTx(ctx, nil)
	if err != nil {
		return updatedEngine, fmt.Errorf("failed to begin transaction: %w", err)
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
			&updatedEngine.EngineID,
			&updatedEngine.Displacement,
			&updatedEngine.NoOfCylinders,
			&updatedEngine.CarRange,
		)

	if err != nil {
		return updatedEngine, fmt.Errorf("failed to update engine: %w", err)
	}

	return updatedEngine, nil
}

func (e EngineStore) EngineDelete(ctx context.Context, id string) (models.Engine, error) {
	var deletedEngine models.Engine

	// Begin transaction
	tx, err := e.db.BeginTx(ctx, nil)
	if err != nil {
		return deletedEngine, fmt.Errorf("failed to begin transaction: %w", err)
	}
	// Ensure rollback on error
	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	// Delete and return the engine details
	query := `
		DELETE FROM engine 
		WHERE engine_id = $1 
		RETURNING engine_id, displacement, noOfCylinders, carRange
	`
	err = tx.QueryRowContext(ctx, query, id).Scan(
		&deletedEngine.EngineID,
		&deletedEngine.Displacement,
		&deletedEngine.NoOfCylinders,
		&deletedEngine.CarRange,
	)

	// Handle error when no rows are affected
	if err == sql.ErrNoRows {
		return deletedEngine, fmt.Errorf("engine with ID %s does not exist", id)
	} else if err != nil {
		return deletedEngine, fmt.Errorf("failed to delete engine: %w", err)
	}

	return deletedEngine, nil
}
