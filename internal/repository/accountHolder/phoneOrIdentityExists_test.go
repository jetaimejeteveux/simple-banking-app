package accountHolderRepository

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestAccountHolderRepository_IsPhoneOrIdentityExist(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	assert.NoError(t, err)

	type args struct {
		ctx            context.Context
		phoneNumber    string
		identityNumber string
	}
	tests := []struct {
		name          string
		args          args
		expectedExist bool
		wantErr       bool
		mockFn        func(args args)
	}{
		{
			name: "SUCCESS: exists",
			args: args{
				ctx:            context.Background(),
				phoneNumber:    "08123456789",
				identityNumber: "1234567890",
			},
			expectedExist: true,
			wantErr:       false,
			mockFn: func(args args) {
				rows := sqlmock.NewRows([]string{"count"}).AddRow(1)

				mock.ExpectQuery(`SELECT count\(\*\) FROM "account_holders" WHERE identity_number = \$1 OR phone_number = \$2 AND "account_holders"."deleted_at" IS NULL`).
					WithArgs(args.phoneNumber, args.identityNumber).
					WillReturnRows(rows)
			},
		},
		{
			name: "SUCCESS: does not exist",
			args: args{
				ctx:            context.Background(),
				phoneNumber:    "08123456789",
				identityNumber: "1234567890",
			},
			expectedExist: false,
			wantErr:       false,
			mockFn: func(args args) {
				rows := sqlmock.NewRows([]string{"count"}).AddRow(0)

				mock.ExpectQuery(`SELECT count\(\*\) FROM "account_holders" WHERE identity_number = \$1 OR phone_number = \$2 AND "account_holders"."deleted_at" IS NULL`).
					WithArgs(args.phoneNumber, args.identityNumber).
					WillReturnRows(rows)
			},
		},
		{
			name: "ERROR: database error",
			args: args{
				ctx:            context.Background(),
				phoneNumber:    "08123456789",
				identityNumber: "1234567890",
			},
			expectedExist: false,
			wantErr:       true,
			mockFn: func(args args) {
				mock.ExpectQuery(`SELECT count\(\*\) FROM "account_holders" WHERE identity_number = \$1 OR phone_number = \$2 AND "account_holders"."deleted_at" IS NULL`).
					WithArgs(args.phoneNumber, args.identityNumber).
					WillReturnError(assert.AnError)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn(tt.args)
			r := &AccountHolderRepository{
				db: gormDB,
			}

			exists, err := r.IsPhoneOrIdentityExist(tt.args.ctx, tt.args.phoneNumber, tt.args.identityNumber)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedExist, exists)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
