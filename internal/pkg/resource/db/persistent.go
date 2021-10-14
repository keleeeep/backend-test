/*
 * @Author: Adrian Faisal
 * @Date: 14/10/21 13.43
 */

package db

import (
	"context"
	"fmt"
	"github.com/keleeeep/test/internal/pkg/model"
)

func (p *persistent) CreateUser(ctx context.Context, user *model.User) (*model.User, error) {
	query := fmt.Sprintf("INSERT INTO %s (name, phone, password, role) VALUES (?, ?, ?, ?)",
		user.TableName())

	tx, err := p.conn.Begin()
	if err != nil {
		return nil, fmt.Errorf("can't start db transaction: %v", err)
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx,
		query, user.Name, user.Phone, user.Password, user.Role)
	if err != nil {
		return nil, fmt.Errorf("exec query failed: %v", err)
	}

	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("commiting transaction failed: %v", err)
	}

	resp, err := p.FindUser(ctx, user.Name, "name")
	if err != nil {
		return nil, fmt.Errorf("failed to find id: %v", err)
	}

	return resp, nil
}

func (p *persistent) FindUser(ctx context.Context, data, column string) (*model.User, error) {
	m := &model.User{} // struct literal

	query := fmt.Sprintf("SELECT * FROM %s WHERE %s = ?", m.TableName(), column)

	row := p.conn.QueryRowContext(ctx, query, data)

	err := row.Scan(&m.ID, &m.Name, &m.Phone, &m.Password, &m.Role, &m.Timestamp)
	if err != nil {
		return nil, nil
	}

	return m, nil
}
