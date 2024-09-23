package repository

import (
	"context"
	//"database/sql"

	"github.com/metalpoch/olt-blueprint/auth/entity"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *userRepository {
	return &userRepository{db}
}

func (repo userRepository) Create(ctx context.Context, user *entity.User) error {
	/*q := `
	    INSERT INTO users (
			id,
			fullname,
			email,
			password,
			change_password,
			is_admin)
	    	VALUES ($1, $2, $3, $4, $5, $6)
	    	RETURNING id;
	    `
		if err := repo.db.QueryRowContext(
			ctx, q, user.ID, user.Fullname, user.Email, user.Password,
			user.ChangePassword, user.IsAdmin,
		).Scan(&user.ID); err != nil {
			return err
		}

		return nil*/
	if err := repo.db.WithContext(ctx).Create(&user).Error; err != nil {
		return err
	}
	return nil
}

func (repo userRepository) GetAll(ctx context.Context) ([]*entity.User, error) {
	users := []*entity.User{}
	/*q := "SELECT * FROM users;"

	rows, err := repo.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		u := new(entity.User)
		rows.Scan(
			&u.ID,
			&u.Fullname,
			&u.Email,
			&u.Password,
			&u.ChangePassword,
			&u.IsAdmin,
			&u.IsDisabled,
			&u.CreatedAt,
			&u.UpdatedAt,
		)
		users = append(users, u)
	}

	return users, nil*/

	if err := repo.db.WithContext(ctx).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil

}

func (repo userRepository) GetUserByID(ctx context.Context, id uint) (*entity.User, error) {
	u := new(entity.User)
	/*row := repo.db.QueryRowContext(ctx, "SELECT * FROM users WHERE id=$1;", id)

	err := row.Scan(
		&u.ID,
		&u.Fullname,
		&u.Email,
		&u.Password,
		&u.ChangePassword,
		&u.IsAdmin,
		&u.IsDisabled,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return u, nil*/
	if err := repo.db.WithContext(ctx).Find(&u, "id=?", id).Error; err != nil {
		return nil, err
	}
	return u, nil
}

func (repo userRepository) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	u := new(entity.User)
	/*row := repo.db.QueryRowContext(ctx, "SELECT * FROM users WHERE email=$1 AND is_disabled=false;", email)

	err := row.Scan(
		&u.ID,
		&u.Fullname,
		&u.Email,
		&u.Password,
		&u.ChangePassword,
		&u.IsAdmin,
		&u.IsDisabled,
		&u.CreatedAt,
		&u.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}*/
	if err := repo.db.WithContext(ctx).Find(&u, "email=?", email).Error; err != nil {
		return nil, err
	}

	return u, nil
}

func (repo userRepository) SoftDelete(ctx context.Context, id uint) error {
	u := new(entity.User)
	/*q := "UPDATE users set is_disabled=true  WHERE id=$1;"

	stmt, err := repo.db.PrepareContext(ctx, q)
	if err != nil {
		return err
	}

	defer stmt.Close()

	if _, err = stmt.ExecContext(ctx, id); err != nil {
		return err
	}*/

	if err := repo.db.Find(&u, "id=?", id).Error; err != nil {
		return err
	}

	if err := repo.db.WithContext(ctx).Delete(&u).Error; err != nil {
		return err
	}
	return nil
}

func (repo userRepository) ChangePassword(ctx context.Context, id uint, password string) error {
	u := new(entity.User)
	/*q := "UPDATE users set change_password=false, password=$1 WHERE id=$2;"

	stmt, err := repo.db.PrepareContext(ctx, q)
	if err != nil {
		return err
	}

	defer stmt.Close()

	if _, err = stmt.ExecContext(ctx, password, id); err != nil {
		return err
	}*/
	if err := repo.db.Find(&u, "id=?", id, password).Error; err != nil {
		return err
	}

	if err := repo.db.Update("Password", password).Error; err != nil {
		return err
	}
	return nil
}
