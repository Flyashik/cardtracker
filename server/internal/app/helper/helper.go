package helper

import "server/internal/app/models"

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
