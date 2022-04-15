package data

import (
	"context"
	"crypto/rand"
	"database/sql"
	"errors"
	"io"
	"javlonrahimov/apod/internal/validator"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type OtpService interface {
	New(userId int64, ttl time.Duration) (*Otp, error)
	Insert(otp *Otp) error
	GetForUser(userId int64) (*Otp, error)
	DeleteAllForUser(userId int64) error
}

type Otp struct {
	Plaintext string
	Hash      []byte
	UserId    int64
	Expiry    time.Time
}

type otpModel struct {
	db *sql.DB
}

var table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}

func NewOtpModel(db *sql.DB) OtpService {
	return &otpModel{db: db}
}

func generateOtp(userId int64, ttl time.Duration) (*Otp, error) {

	otp := &Otp{
		UserId: userId,
		Expiry: time.Now().Add(ttl),
	}

	otp.Plaintext = encodeToString(6)

	hash, err := bcrypt.GenerateFromPassword([]byte(otp.Plaintext), 12)
	if err != nil {
		return nil, err
	}

	otp.Hash = hash

	return otp, nil
}

func encodeToString(max int) string {
	b := make([]byte, max)
	n, err := io.ReadAtLeast(rand.Reader, b, max)
	if n != max {
		panic(err)
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	return string(b)
}

func ValidateOtpPlaintext(v *validator.Validator, otpPlaintext string) {
	v.Check(otpPlaintext != "", "otp", "must be provided")
	v.Check(len(otpPlaintext) == 6, "otp", "must be 6 bytes long")
}

func (m otpModel) New(userId int64, ttl time.Duration) (*Otp, error) {
	otp, err := generateOtp(userId, ttl)
	if err != nil {
		return nil, err
	}
	_ = m.DeleteAllForUser(userId)
	err = m.Insert(otp)
	return otp, err
}

func (m otpModel) Insert(otp *Otp) error {
	query := `
	INSERT INTO otps (hash, user_id, expiry)
	VALUES ($1, $2, $3)`

	args := []interface{}{otp.Hash, otp.UserId, otp.Expiry}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.db.ExecContext(ctx, query, args...)
	return err
}

func (m otpModel) GetForUser(userId int64) (*Otp, error) {
	query := `
	SELECT hash, user_id, expiry FROM otps
	WHERE user_id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var otp Otp

	err := m.db.QueryRowContext(ctx, query, userId).Scan(
		&otp.Hash,
		&otp.UserId,
		&otp.Expiry,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	if !otp.Expiry.After(time.Now()) {
		return nil, ErrOtpExpired
	}
	return &otp, nil
}

func (m otpModel) DeleteAllForUser(userId int64) error {
	query := `
	DELETE FROM otps
	WHERE user_id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.db.ExecContext(ctx, query, userId)
	return err
}
