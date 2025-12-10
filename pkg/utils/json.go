package utils

import (
	"encoding/json"
	"log/slog"
	"os"

	"challenge/pkg/models"
)

type JSONUtil struct {
	logger *slog.Logger
}

func NewJSONUtil(logger *slog.Logger) *JSONUtil {
	return &JSONUtil{
		logger: logger,
	}
}

type StateData struct {
	Distributors map[string]*models.Distributor `json:"distributors"`
}

// Save saves the distributor state to a JSON file
func (r *JSONUtil) Save(filename string, distributors map[string]*models.Distributor) error {
	r.logger.Debug("saving state to file",
		slog.String("filename", filename))

	state := StateData{
		Distributors: distributors,
	}

	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		r.logger.Error("failed to marshal state",
			slog.String("error", err.Error()))
		return err
	}

	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		r.logger.Error("failed to write state file",
			slog.String("filename", filename),
			slog.String("error", err.Error()))
		return err
	}

	r.logger.Info("state saved successfully",
		slog.String("filename", filename),
		slog.Int("distributors", len(distributors)))

	return nil
}

// Load loads the distributor state from a JSON file
func (r *JSONUtil) Load(filename string) (map[string]*models.Distributor, error) {
	r.logger.Debug("loading state from file",
		slog.String("filename", filename))

	data, err := os.ReadFile(filename)
	if err != nil {
		r.logger.Error("failed to read state file",
			slog.String("filename", filename),
			slog.String("error", err.Error()))
		return nil, err
	}

	var state StateData
	err = json.Unmarshal(data, &state)
	if err != nil {
		r.logger.Error("failed to unmarshal state",
			slog.String("error", err.Error()))
		return nil, err
	}

	r.logger.Info("state loaded successfully",
		slog.String("filename", filename),
		slog.Int("distributors", len(state.Distributors)))

	return state.Distributors, nil
}

func (r *JSONUtil) Exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}
