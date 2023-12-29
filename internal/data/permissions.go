package data

import (
	"context"
	"database/sql"
	"time"
)

type UserPermissions []string

func (p UserPermissions) Include(code string) bool {
	for i := range p {
		if code == p[i] {
			return true
		}
	}
	return false
}

type PermissionModel struct {
	DB *sql.DB
}

type Permissions interface {
	GetAllForUser(userID int64) (UserPermissions, error)
}

// GetAllForUser returns all permission codes for a specific user in a Permissions slice.
func (m PermissionModel) GetAllForUser(userID int64) (UserPermissions, error) {
	query := `
		select permissions.code
		from permissions
		inner join users_permissions on users_permissions.permission_id = permissions.id
		inner join users on users_permissions.user_id = users.id
		where users.id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions UserPermissions
	for rows.Next() {
		var permission string

		if err := rows.Scan(&permission); err != nil {
			return nil, err
		}

		permissions = append(permissions, permission)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return permissions, nil
}
