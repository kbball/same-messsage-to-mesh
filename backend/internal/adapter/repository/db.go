package repository

import (
	"database/sql"
	"errors"

	"github.com/kbball/same-message-to-mesh/backend/internal/domain"
)

func mapNotFound(err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return domain.ErrNotFound
	}
	return err
}
