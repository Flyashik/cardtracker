package storage_test

import (
	"github.com/stretchr/testify/assert"
	"server/internal/app/models"
	"server/internal/app/storage"
	"testing"
)

func TestPhoneRepository_Create(t *testing.T) {
	s, teardown := storage.TestStorage(t, dbUrl)
	defer teardown("phones")

	p, err := s.Phone().Create(&models.Phone{
		Manufacturer: "Samsung",
		ModelTag:     "beyond1",
		ModelNumber:  "SM-G973F/DS",
		OsVersion:    "12",
		ApiVersion:   "31",
		Cpu:          "exynos9820",
		Firmware:     "G9773FXXSGHWA1",
		Bootloader:   "G9773FXXSGHWC3",
		SupportedArchs: []string{
			"arm64-v8a",
			"armeabi-v7a",
			"armeabi",
		},
		SimSlots: 1,
		SdSlots:  0,
	})
	assert.NoError(t, err)
	assert.NotNil(t, p)
}

func TestSimRepository_Create(t *testing.T) {
	s, teardown := storage.TestStorage(t, dbUrl)
	defer teardown("phones", "sim_cards")

	p, err := s.Phone().Create(&models.Phone{
		Manufacturer: "Samsung",
		ModelTag:     "beyond1",
		ModelNumber:  "SM-G973F/DS",
		OsVersion:    "12",
		ApiVersion:   "31",
		Cpu:          "exynos9820",
		Firmware:     "G9773FXXSGHWA1",
		Bootloader:   "G9773FXXSGHWC3",
		SupportedArchs: []string{
			"arm64-v8a",
			"armeabi-v7a",
			"armeabi",
		},
	})

	assert.NoError(t, err)
	assert.NotNil(t, p)

	sim, err := s.Sim().Create(&models.SimInfo{
		PhoneNumber: "79889484608",
		Operator:    "MTS",
	}, p)

	assert.NoError(t, err)
	assert.NotNil(t, sim)
}

func TestSdRepository_Create(t *testing.T) {
	s, teardown := storage.TestStorage(t, dbUrl)
	defer teardown("phones", "sd_cards")

	p, err := s.Phone().Create(&models.Phone{
		Manufacturer: "Samsung",
		ModelTag:     "beyond1",
		ModelNumber:  "SM-G973F/DS",
		OsVersion:    "12",
		ApiVersion:   "31",
		Cpu:          "exynos9820",
		Firmware:     "G9773FXXSGHWA1",
		Bootloader:   "G9773FXXSGHWC3",
		SupportedArchs: []string{
			"arm64-v8a",
			"armeabi-v7a",
			"armeabi",
		},
	})

	assert.NoError(t, err)
	assert.NotNil(t, p)

	sd, err := s.SdCard().Create(&models.SdInfo{
		SdManufacturerId: "0x00001b",
		SerialNo:         "0x1a8ed52f",
		TotalSpace:       60874,
		UsedSpace:        16036,
		FreeSpace:        44848,
	}, p)

	assert.NoError(t, err)
	assert.NotNil(t, sd)
}

func TestPhoneRepository_SelectByModelTag(t *testing.T) {
	s, teardown := storage.TestStorage(t, dbUrl)
	defer teardown("phones")

	_, err := s.Phone().SelectByModelTag("beyond1")
	assert.Error(t, err)

	s.Phone().Create(&models.Phone{
		Manufacturer: "Samsung",
		ModelTag:     "beyond1",
		ModelNumber:  "SM-G973F/DS",
		OsVersion:    "12",
		ApiVersion:   "31",
		Cpu:          "exynos9820",
		Firmware:     "G9773FXXSGHWA1",
		Bootloader:   "G9773FXXSGHWC3",
		SupportedArchs: []string{
			"arm64-v8a",
			"armeabi-v7a",
			"armeabi",
		},
		SimSlots: 1,
		SdSlots:  1,
	})

	p, err := s.Phone().SelectByModelTag("beyond1")
	assert.NoError(t, err)
	assert.NotNil(t, p)
}
