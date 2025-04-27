package accountHolderRepository

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jetaimejeteveux/simple-banking-app/internal/model"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestAccountHolderRepository_GetByAccountNumber(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	assert.NoError(t, err)

	type args struct {
		ctx           context.Context
		accountNumber string
	}
	tests := []struct {
		name         string
		args         args
		expectedUser *model.AccountHolder
		wantErr      bool
		mockFn       func(args args)
	}{
		{
			name: "SUCCESS",
			args: args{
				ctx:           context.Background(),
				accountNumber: "1234-5678-9012",
			},
			expectedUser: &model.AccountHolder{
				FullName:       "John Doe",
				IdentityNumber: "1234567890",
				PhoneNumber:    "08123456789",
				AccountNumber:  "1234-5678-9012",
				Balance:        1000.0,
			},
			wantErr: false,
			mockFn: func(args args) {
				createdAt := time.Now()
				updatedAt := time.Now()

				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "full_name", "identity_number", "phone_number", "account_number", "balance"}).
					AddRow(1, createdAt, updatedAt, nil, "John Doe", "1234567890", "08123456789", "1234-5678-9012", 1000.0)

				mock.ExpectQuery(`SELECT \* FROM "account_holders" WHERE account_number = \$1 AND "account_holders"."deleted_at" IS NULL ORDER BY "account_holders"."id" LIMIT \$2`).
					WithArgs(args.accountNumber, 1).
					WillReturnRows(rows)
			},
		},
		{
			name: "ERROR: account not found",
			args: args{
				ctx:           context.Background(),
				accountNumber: "NONEXISTENT",
			},
			expectedUser: nil,
			wantErr:      true,
			mockFn: func(args args) {
				mock.ExpectQuery(`SELECT \* FROM "account_holders" WHERE account_number = \$1 AND "account_holders"."deleted_at" IS NULL ORDER BY "account_holders"."id" LIMIT \$2`).
					WithArgs(args.accountNumber, 1).
					WillReturnError(gorm.ErrRecordNotFound)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn(tt.args)
			r := &AccountHolderRepository{
				db: gormDB,
			}
			accountHolder, err := r.GetByAccountNumber(tt.args.ctx, tt.args.accountNumber)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, accountHolder)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedUser.FullName, accountHolder.FullName)
				assert.Equal(t, tt.expectedUser.IdentityNumber, accountHolder.IdentityNumber)
				assert.Equal(t, tt.expectedUser.PhoneNumber, accountHolder.PhoneNumber)
				assert.Equal(t, tt.expectedUser.AccountNumber, accountHolder.AccountNumber)
				assert.Equal(t, tt.expectedUser.Balance, accountHolder.Balance)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
