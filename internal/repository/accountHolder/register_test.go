package accountHolderRepository

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jetaimejeteveux/simple-banking-app/internal/model"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestAccountHolderRepository_Register(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	assert.NoError(t, err)

	type args struct {
		ctx           context.Context
		accountHolder *model.AccountHolder
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		mockFn  func(args args)
	}{
		{
			name: "SUCCESS",
			args: args{
				ctx: context.Background(),
				accountHolder: &model.AccountHolder{
					FullName:       "John Doe",
					IdentityNumber: "1234567890",
					PhoneNumber:    "08123456789",
					AccountNumber:  "1234-5678-9012",
					Balance:        0,
				},
			},
			wantErr: false,
			mockFn: func(args args) {
				mock.ExpectBegin()
				mock.ExpectQuery(`INSERT INTO "account_holders"(.+) VALUES (.+)`).
					WithArgs(
						sqlmock.AnyArg(), // CreatedAt
						sqlmock.AnyArg(), // UpdatedAt
						sqlmock.AnyArg(), // DeletedAt
						args.accountHolder.FullName,
						args.accountHolder.IdentityNumber,
						args.accountHolder.PhoneNumber,
						args.accountHolder.AccountNumber,
						args.accountHolder.Balance,
					).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
				mock.ExpectCommit()
			},
		},
		{
			name: "ERROR: database error",
			args: args{
				ctx: context.Background(),
				accountHolder: &model.AccountHolder{
					FullName:       "John Doe",
					IdentityNumber: "1234567890",
					PhoneNumber:    "08123456789",
					AccountNumber:  "1234-5678-9012",
					Balance:        0,
				},
			},
			wantErr: true,
			mockFn: func(args args) {
				mock.ExpectBegin()
				mock.ExpectQuery(`INSERT INTO "account_holders"(.+) VALUES (.+)`).
					WithArgs(
						sqlmock.AnyArg(), // CreatedAt
						sqlmock.AnyArg(), // UpdatedAt
						sqlmock.AnyArg(), // DeletedAt
						args.accountHolder.FullName,
						args.accountHolder.IdentityNumber,
						args.accountHolder.PhoneNumber,
						args.accountHolder.AccountNumber,
						args.accountHolder.Balance,
					).WillReturnError(assert.AnError)
				mock.ExpectRollback()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn(tt.args)
			r := &AccountHolderRepository{
				db: gormDB,
			}
			err := r.Register(tt.args.ctx, tt.args.accountHolder)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
