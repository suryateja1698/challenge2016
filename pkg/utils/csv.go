package utils

import (
	"encoding/csv"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"

	"challenge/pkg/models"
)

type CSVUtil struct {
	logger *slog.Logger
}

func NewCSVUtil(logger *slog.Logger) *CSVUtil {
	return &CSVUtil{
		logger: logger,
	}
}

// LoadCities loads cities from a CSV file and returns a map of locations
func (r *CSVUtil) LoadCities(filename string) (map[string]models.Location, error) {
	r.logger.Info("loading cities from CSV",
		slog.String("filename", filename))

	file, err := os.Open(filename)
	if err != nil {
		r.logger.Error("failed to open cities file",
			slog.String("filename", filename),
			slog.String("error", err.Error()))
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// Skip header
	_, err = reader.Read()
	if err != nil {
		r.logger.Error("failed to read CSV header",
			slog.String("error", err.Error()))
		return nil, fmt.Errorf("error reading header: %v", err)
	}

	locations := make(map[string]models.Location)
	count := 0

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			r.logger.Error("failed to read CSV record",
				slog.String("error", err.Error()))
			return nil, fmt.Errorf("error reading record: %v", err)
		}

		if len(record) >= 3 {
			city := strings.TrimSpace(record[0])
			state := strings.TrimSpace(record[1])
			country := strings.TrimSpace(record[2])

			// Store city-state-country combination
			loc := models.Location{
				City:    city,
				State:   state,
				Country: country,
			}
			key := fmt.Sprintf("%s-%s-%s", city, state, country)
			locations[key] = loc

			// Store state-country combination
			stateKey := fmt.Sprintf("%s-%s", state, country)
			if _, exists := locations[stateKey]; !exists {
				locations[stateKey] = models.Location{
					State:   state,
					Country: country,
				}
			}

			// Store country
			if _, exists := locations[country]; !exists {
				locations[country] = models.Location{
					Country: country,
				}
			}

			count++
		}
	}

	r.logger.Info("successfully loaded cities",
		slog.Int("total_locations", count),
		slog.Int("unique_keys", len(locations)))

	return locations, nil
}
