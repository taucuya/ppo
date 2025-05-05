package user_rep

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	structs "github.com/taucuya/ppo/internal/core/structs"
	rep_struct "github.com/taucuya/ppo/internal/repository/postgres/structs"
	"golang.org/x/crypto/bcrypt"
)

type Repository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (rep *Repository) Create(ctx context.Context, u structs.User) (uuid.UUID, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return uuid.UUID{}, err
	}

	usr := rep_struct.User{
		Id:            u.Id,
		Name:          u.Name,
		Date_of_birth: u.Date_of_birth,
		Mail:          u.Mail,
		Password:      string(hashedPassword),
		Phone:         u.Phone,
		Address:       u.Address,
		Status:        u.Status,
		Role:          u.Role,
	}
	var id uuid.UUID
	err = rep.db.GetContext(ctx, &id,
		`insert into "user" 
		(name, date_of_birth, mail, password, phone, address, status, role) 
		values (:name, :date_of_birth, :mail, :password, :phone, :address, :status, :role) 
		returning id`, usr)

	if err != nil {
		return uuid.UUID{}, err
	}
	fmt.Println(err)
	return id, err
}

func (rep *Repository) GetById(ctx context.Context, id uuid.UUID) (structs.User, error) {
	var u rep_struct.User
	err := rep.db.GetContext(ctx, &u, "select * from \"user\" where id = $1", id)
	if err != nil {
		return structs.User{}, fmt.Errorf("failed to get user: %w", err)
	}
	usr := structs.User{
		Id:            u.Id,
		Name:          u.Name,
		Date_of_birth: u.Date_of_birth,
		Mail:          u.Mail,
		Password:      u.Password,
		Phone:         u.Phone,
		Address:       u.Address,
		Status:        u.Status,
		Role:          u.Role,
	}
	return usr, nil
}

func (rep *Repository) GetByMail(ctx context.Context, mail string) (structs.User, error) {
	var u rep_struct.User
	err := rep.db.GetContext(ctx, &u, "select * from \"user\" where mail = $1", mail)
	if err != nil {
		return structs.User{}, fmt.Errorf("failed to get user by mail: %w", err)
	}
	usr := structs.User{
		Id:            u.Id,
		Name:          u.Name,
		Date_of_birth: u.Date_of_birth,
		Mail:          u.Mail,
		Password:      u.Password,
		Phone:         u.Phone,
		Address:       u.Address,
		Status:        u.Status,
		Role:          u.Role,
	}
	return usr, nil
}

func (rep *Repository) GetByPhone(ctx context.Context, phone string) (structs.User, error) {
	var u rep_struct.User
	err := rep.db.GetContext(ctx, &u, "select * from \"user\" where phone = $1", phone)
	if err != nil {
		return structs.User{}, fmt.Errorf("failed to get user by phone: %w", err)
	}
	usr := structs.User{
		Id:            u.Id,
		Name:          u.Name,
		Date_of_birth: u.Date_of_birth,
		Mail:          u.Mail,
		Password:      u.Password,
		Phone:         u.Phone,
		Address:       u.Address,
		Status:        u.Status,
		Role:          u.Role,
	}
	return usr, nil
}
