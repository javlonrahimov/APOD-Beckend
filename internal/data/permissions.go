package data

import (
	"context"
	"database/sql"
	"time"

	"github.com/lib/pq"
)

type Permissions []string

type PermissonService interface {
	GetAllForUser(userId int64) (Permissions, error)
	AddForUser(userId int64, codes ...string) error
}

var (
	ApodsRead  = "apods:read"
	ApodsWrite = "apods:write"
)

func (p Permissions) Include(code string) bool {
	for i := range p {
		if code == p[i] {
			return true
		}
	}
	return false
}

type permissionModel struct {
	db *sql.DB
}

func NewPermissonModel(db *sql.DB) PermissonService {
	return &permissionModel{db: db}
}

func (m permissionModel) GetAllForUser(userId int64) (Permissions, error) {
	query := `
	SELECT permissions.code
	FROM permissions
	INNER JOIN users_permissions ON users_permissions.permission_id = permissions.id
	INNER JOIN users ON users_permissions.user_id = users.id
	WHERE users.id = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	rows, err := m.db.QueryContext(ctx, query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var permissions Permissions
	for rows.Next() {
		var permission string
		err := rows.Scan(&permission)
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, permission)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return permissions, nil
}

func (m permissionModel) AddForUser(userId int64, codes ...string) error {
	query := `
	INSERT INTO users_permissions
	SELECT $1, permissions.id FROM permissions WHERE permissions.code = ANY($2)`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := m.db.ExecContext(ctx, query, userId, pq.Array(codes))
	return err
}
