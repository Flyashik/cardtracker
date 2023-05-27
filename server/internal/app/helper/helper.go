package helper

import (
	"encoding/csv"
	"io"
	"os"
	"path/filepath"
	"server/internal/app/models"
	"strings"
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
		if err != nil {
			return "", err
		}
		if err == io.EOF {
			return "", nil
		}
		if rec[1] == manufacturerId {
			return rec[0], nil
		}
	}
}
