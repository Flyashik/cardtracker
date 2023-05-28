package helper

import (
	"encoding/csv"
	"io"
	"os"
	"path/filepath"
	"server/internal/app/models"
	"strings"
	"golang.org/x/crypto/bcrypt"
	"github.com/dgrijalva/jwt-go"
)

const DataPath = "data"

func IsEmptySdSlot(sd models.SdInfo) bool {
	if sd.SerialNo == "" && sd.SdManufacturerId == "" && sd.FreeSpace == 0 && sd.UsedSpace == 0 && sd.TotalSpace == 0 {
		return true
	}

	return false
}

func IsEmptySimSlot(sim models.SimInfo) bool {
	if sim.PhoneNumber == "" && sim.Operator == "" {
		return true
	}

	return false
}

func ConvertModelTagToMarketingName(modelTag string) (string, error) {
	f, err := os.Open(filepath.Join(DataPath, "supported_devices.csv"))
	if err != nil {
		return "", err
	}
	defer f.Close()

	reader := csv.NewReader(f)
	reader.Comma = ','

	_, err = reader.Read()
	if err != nil {
		return "", err
	}

	for {
		rec, err := reader.Read()
		if err == io.EOF {
			return modelTag, nil
		}
		if strings.ToLower(rec[2]) == strings.ToLower(modelTag) {
			return rec[1], nil
		}
	}
}

func ConvertManufacturerIdToCompanyName(manufacturerId string) (string, error) {
	f, err := os.Open(filepath.Join(DataPath, "sd_cards.csv"))
	if err != nil {
		return "", err
	}
	defer f.Close()

	reader := csv.NewReader(f)
	reader.Comma = ','

	_, err = reader.Read()
	if err != nil {
		return "", err
	}

	for {
		rec, err := reader.Read()
		if err == io.EOF {
			return manufacturerId, nil
		}
		if strings.ToLower(rec[1]) == strings.ToLower(manufacturerId) {
			return rec[0], nil
		}
	}
}

func CompareHashPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateHashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func ParseToken(tokenString string) (claims *models.Claims, err error) {
	token, err := jwt.ParseWithClaims(tokenString, &models.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("my_secret_key"), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*models.Claims)

	if !ok {
		return nil, err
	}

	return claims, nil
}
