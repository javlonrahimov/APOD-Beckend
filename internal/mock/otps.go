package mock

import (
	"javlonrahimov/apod/internal/data"
	"time"
)

var mockOtp string

var otps = make([]data.Otp, 10)

type otpModelMock struct{}

func NewOtpsMock() data.OtpService {
	otpModelMock := otpModelMock{}
	otpModelMock.New(1, 15*time.Minute)
	return &otpModelMock
}

func (m otpModelMock) New(userId int64, ttl time.Duration) (*data.Otp, error) {
	mockOtp = "123456"
	otp, err := data.GenerateOtp(userId, ttl, &mockOtp)
	if err != nil {
		return nil, err
	}
	_ = m.DeleteAllForUser(userId)
	err = m.Insert(otp)
	return otp, err
}

func (m otpModelMock) Insert(otp *data.Otp) error {
	otps = append(otps, *otp)
	return nil
}

func (m otpModelMock) GetForUser(userId int64) (*data.Otp, error) {

	for i := 0; i < len(otps); i++ {
		if otps[i].UserId == userId {
			if !otps[i].Expiry.After(time.Now()) {
				return nil, data.ErrOtpExpired
			}
			return &otps[i], nil
		}
	}
	return nil, data.ErrRecordNotFound
}

func (m otpModelMock) DeleteAllForUser(userId int64) error {
	for i := 0; i < len(otps); i++ {
		if otps[i].UserId == userId {
			otps = removeOtp(otps, i)
		}
	}
	return nil
}

func removeOtp(slice []data.Otp, s int) []data.Otp {
	return append(slice[:s], slice[s+1:]...)
}
