package common

import "github.com/google/uuid"

func ConvertToUUIDs(strSlice []string) ([]uuid.UUID, error) {
	var uuidSlice []uuid.UUID

	for _, str := range strSlice {
		uuidVal, err := uuid.Parse(str)
		if err != nil {
			return nil, err
		}

		uuidSlice = append(uuidSlice, uuidVal)
	}

	return uuidSlice, nil
}
