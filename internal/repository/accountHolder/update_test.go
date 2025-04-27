package accountHolderRepository

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestAccountHolderRepository_UpdateBalance(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	assert.NoError(t, err)

	type args struct {
		ctx           context.Context
		accountNumber string
		balance       float64
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		mockFn  func(args args)
	}{
		{
			name: "SUCCESS: update balance",
			args: args{
				ctx:           context.Background(),
				accountNumber: "1234-5678-9012",
				balance:       1000.50,
			},
			wantErr: false,
			mockFn: func(args args) {
				mock.ExpectBegin()
				mock.ExpectExec(`UPDATE "account_holders" SET "balance"=\$1,"updated_at"=\$2 WHERE account_number = \$3 AND "account_holders"."deleted_at" IS NULL`).
					WithArgs(args.balance, sqlmock.AnyArg(), args.accountNumber).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
		},
		{
			name: "ERROR: database error",
			args: args{
				ctx:           context.Background(),
				accountNumber: "1234-5678-9012",
				balance:       1000.50,
			},
			wantErr: true,
			mockFn: func(args args) {
				mock.ExpectBegin()
				mock.ExpectExec(`UPDATE "account_holders" SET "balance"=\$1,"updated_at"=\$2 WHERE account_number = \$3 AND "account_holders"."deleted_at" IS NULL`).
					WithArgs(args.balance, sqlmock.AnyArg(), args.accountNumber).
					WillReturnError(assert.AnError)
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

			err := r.UpdateBalance(tt.args.ctx, tt.args.accountNumber, tt.args.balance)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
