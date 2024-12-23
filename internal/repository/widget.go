package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/mlvieira/store/internal/models"
)

// widgetRepo handles database operations for widgets.
type widgetRepo struct {
	db *sql.DB
}

// NewWidgetRepository creates a new WidgetRepository
func NewWidgetRepository(db *sql.DB) WidgetRepository {
	return &widgetRepo{db: db}
}

// GetWidgetByID fetches a widget by its ID.
func (r *widgetRepo) GetWidgetByID(ctx context.Context, id int) (models.Widget, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var widget models.Widget

	stmt := `
		SELECT id, name, description, inventory_level, price, 
		       COALESCE(image, '') AS image, is_recurring, plan_id,
			   created_at, updated_at
		FROM widgets
		WHERE id = ?
	`

	row := r.db.QueryRowContext(ctx, stmt, id)
	if err := row.Scan(
		&widget.ID,
		&widget.Name,
		&widget.Description,
		&widget.InventoryLevel,
		&widget.Price,
		&widget.Image,
		&widget.IsRecurring,
		&widget.PlanID,
		&widget.CreatedAt,
		&widget.UpdatedAt,
	); err != nil {
		return widget, err
	}

	return widget, nil
}
