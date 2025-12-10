package services

import (
	"fmt"
	"log/slog"
	"strings"

	"challenge/pkg/models"
)

type DistributionService struct {
	distributors map[string]*models.Distributor
	locations    map[string]models.Location
	logger       *slog.Logger
}

func NewDistributionService(logger *slog.Logger) *DistributionService {
	return &DistributionService{
		distributors: make(map[string]*models.Distributor),
		locations:    make(map[string]models.Location),
		logger:       logger,
	}
}

func (s *DistributionService) SetLocations(locations map[string]models.Location) {
	s.locations = locations
}

func (s *DistributionService) SetDistributors(distributors map[string]*models.Distributor) {
	s.distributors = distributors
}

func (s *DistributionService) GetDistributors() map[string]*models.Distributor {
	return s.distributors
}

func (s *DistributionService) AddDistributor(name string, parentName string) error {
	s.logger.Debug("attempting to add distributor",
		slog.String("name", name),
		slog.String("parent", parentName))

	if _, exists := s.distributors[name]; exists {
		s.logger.Warn("distributor already exists",
			slog.String("name", name))
		return fmt.Errorf("distributor %s already exists", name)
	}

	if parentName != "" {
		if _, exists := s.distributors[parentName]; !exists {
			s.logger.Warn("parent distributor not found",
				slog.String("parent", parentName))
			return fmt.Errorf("parent distributor %s not found", parentName)
		}
	}

	dist := &models.Distributor{
		Name:        name,
		Parent:      parentName,
		Permissions: []models.Permission{},
	}

	s.distributors[name] = dist

	s.logger.Info("distributor added successfully",
		slog.String("name", name),
		slog.String("parent", parentName))

	return nil
}

func (s *DistributionService) AddPermission(distName string, isInclude bool, locationStr string) error {
	s.logger.Debug("attempting to add permission",
		slog.String("distributor", distName),
		slog.Bool("is_include", isInclude),
		slog.String("location", locationStr))

	dist, exists := s.distributors[distName]
	if !exists {
		s.logger.Warn("distributor not found",
			slog.String("distributor", distName))
		return fmt.Errorf("distributor %s not found", distName)
	}

	loc := s.parseLocation(locationStr)
	perm := models.Permission{
		IsInclude: isInclude,
		Location:  loc,
	}

	dist.Permissions = append(dist.Permissions, perm)

	permType := "EXCLUDE"
	if isInclude {
		permType = "INCLUDE"
	}

	s.logger.Info("permission added successfully",
		slog.String("distributor", distName),
		slog.String("permission_type", permType),
		slog.String("location", locationStr),
		slog.String("city", loc.City),
		slog.String("state", loc.State),
		slog.String("country", loc.Country))

	return nil
}

func (s *DistributionService) CanDistribute(distName string, locationStr string) (bool, error) {
	s.logger.Debug("checking distribution permission",
		slog.String("distributor", distName),
		slog.String("location", locationStr))

	dist, exists := s.distributors[distName]
	if !exists {
		s.logger.Warn("distributor not found for check",
			slog.String("distributor", distName))
		return false, fmt.Errorf("distributor %s not found", distName)
	}

	targetLoc := s.parseLocation(locationStr)
	result := s.checkPermission(dist, targetLoc)

	s.logger.Info("permission check completed",
		slog.String("distributor", distName),
		slog.String("location", locationStr),
		slog.Bool("can_distribute", result))

	return result, nil
}

func (s *DistributionService) GetDistributor(name string) (*models.Distributor, error) {
	dist, exists := s.distributors[name]
	if !exists {
		return nil, fmt.Errorf("distributor %s not found", name)
	}
	return dist, nil
}

// parseLocation parses a location string like "HYDERABAD-TELANGANA-INDIA"
func (s *DistributionService) parseLocation(locStr string) models.Location {
	parts := strings.Split(locStr, "-")

	if len(parts) == 3 {
		return models.Location{City: parts[0], State: parts[1], Country: parts[2]}
	} else if len(parts) == 2 {
		return models.Location{State: parts[0], Country: parts[1]}
	} else if len(parts) == 1 {
		return models.Location{Country: parts[0]}
	}

	return models.Location{}
}

func (s *DistributionService) checkPermission(dist *models.Distributor, target models.Location) bool {
	allPermissions := s.collectPermissions(dist)

	s.logger.Debug("checking permissions",
		slog.String("distributor", dist.Name),
		slog.Int("total_permissions", len(allPermissions)),
		slog.String("target_city", target.City),
		slog.String("target_state", target.State),
		slog.String("target_country", target.Country))

	included := false
	excluded := false

	for _, perm := range allPermissions {
		if s.matchesLocation(target, perm.Location) {
			if perm.IsInclude {
				included = true
			} else {
				excluded = true
			}
		}
	}

	for _, perm := range allPermissions {
		if !perm.IsInclude && s.matchesLocation(target, perm.Location) {
			if s.isMoreSpecific(perm.Location, target) {
				return false
			}
		}
	}

	return included && !excluded
}

// collectPermissions collects all permissions from distributor chain
func (s *DistributionService) collectPermissions(dist *models.Distributor) []models.Permission {
	var perms []models.Permission

	if dist.Parent != "" {
		if parent, exists := s.distributors[dist.Parent]; exists {
			perms = append(perms, s.collectPermissions(parent)...)
		}
	}

	perms = append(perms, dist.Permissions...)

	return perms
}

// matchesLocation checks if target location matches permission location
func (s *DistributionService) matchesLocation(target models.Location, perm models.Location) bool {
	if perm.IsCountryLevel() {
		return target.Country == perm.Country
	}

	if perm.IsStateLevel() {
		return target.State == perm.State && target.Country == perm.Country
	}

	return target.City == perm.City && target.State == perm.State && target.Country == perm.Country
}

// isMoreSpecific checks if location is more or equally specific
func (s *DistributionService) isMoreSpecific(loc models.Location, target models.Location) bool {
	cityMatch := loc.City == "" || loc.City == target.City
	stateMatch := loc.State == "" || loc.State == target.State
	countryMatch := loc.Country == target.Country

	return cityMatch && stateMatch && countryMatch
}
