package constants

const (
	IdentityNumberExistsError   = "nik sudah terdaftar"
	PhoneNumberExistsError      = "no_hp sudah terdaftar"
	AccountNumberExistsError    = "no_rekening sudah terdaftar"
	AccountNotFoundError        = "no_rekening tidak ditemukan"
	PhoneNumberNotFoundError    = "no_hp tidak ditemukan"
	IdentityNumberNotFoundError = "nik tidak ditemukan"
	InsufficientBalanceError    = "saldo tidak mencukupi"
	IdentityNumberCheckError    = "Terjadi kesalahan saat memeriksa nomor identitas. Silakan coba lagi nanti."
	PhoneNumberCheckError       = "Terjadi kesalahan saat memeriksa nomor telepon. Silakan coba lagi nanti."
	AccountNumberCheckError     = "Terjadi kesalahan saat memeriksa nomor rekening. Silakan coba lagi nanti."
	RegisterAccountError        = "Terjadi kesalahan saat mendaftar nomor rekening. Silakan coba lagi nanti."
	FetchAccountHolderError     = "Terjadi kesalahan saat mengambil data pemilik rekening. Silakan coba lagi nanti."
	DepositError                = "Terjadi kesalahan saat melakukan setoran. Silakan coba lagi nanti."
	WithdrawError               = "Terjadi kesalahan saat melakukan penarikan. Silakan coba lagi nanti."
	InvalidRequestError         = "Format data tidak valid. Pastikan semua field sudah benar."
	MissingFieldError           = "Pastikan semua field yang diperlukan sudah diisi."
	InvalidAccountNumberError   = "Nomor rekening tidak valid. Silakan periksa nomor rekening yang Anda masukkan."
)
