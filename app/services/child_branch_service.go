package services

import (
	"errors"
	"time"

	"github.com/followCode/djjs-event-reporting-backend/app/models"
	"github.com/followCode/djjs-event-reporting-backend/config"
)

// CreateChildBranch creates a new child branch
func CreateChildBranch(childBranch *models.ChildBranch) error {
	childBranch.CreatedOn = time.Now()
	if err := config.DB.Create(childBranch).Error; err != nil {
		return err
	}
	return nil
}

// GetAllChildBranches fetches all child branches
func GetAllChildBranches() ([]models.ChildBranch, error) {
	var childBranches []models.ChildBranch
	if err := config.DB.
		Preload("ParentBranch").
		Preload("Country").
		Preload("State").
		Preload("District").
		Preload("City").
		Preload("Infrastructures").
		Preload("Members").
		Order("id DESC").
		Find(&childBranches).Error; err != nil {
		return nil, err
	}
	return childBranches, nil
}

// GetChildBranch fetches a child branch by ID
func GetChildBranch(childBranchID uint) (*models.ChildBranch, error) {
	var childBranch models.ChildBranch
	if err := config.DB.
		Preload("ParentBranch").
		Preload("Country").
		Preload("State").
		Preload("District").
		Preload("City").
		Preload("Infrastructures").
		Preload("Members").
		First(&childBranch, childBranchID).Error; err != nil {
		return nil, errors.New("child branch not found")
	}
	return &childBranch, nil
}

// GetChildBranchesByParent fetches all child branches of a parent branch
func GetChildBranchesByParent(parentBranchID uint) ([]models.ChildBranch, error) {
	var childBranches []models.ChildBranch
	if err := config.DB.
		Where("parent_branch_id = ?", parentBranchID).
		Preload("ParentBranch").
		Preload("Country").
		Preload("State").
		Preload("District").
		Preload("City").
		Preload("Infrastructures").
		Preload("Members").
		Order("id DESC").
		Find(&childBranches).Error; err != nil {
		return nil, err
	}
	return childBranches, nil
}

// UpdateChildBranch updates a child branch
func UpdateChildBranch(childBranchID uint, updatedData map[string]interface{}) error {
	var childBranch models.ChildBranch
	if err := config.DB.First(&childBranch, childBranchID).Error; err != nil {
		return errors.New("child branch not found")
	}

	// Validate parent_branch_id if being updated
	if parentID, ok := updatedData["parent_branch_id"]; ok {
		var parentIDVal uint
		switch v := parentID.(type) {
		case float64:
			parentIDVal = uint(v)
		case uint:
			parentIDVal = v
		case int:
			parentIDVal = uint(v)
		}
		if parentIDVal > 0 {
			var parentBranch models.Branch
			if err := config.DB.First(&parentBranch, parentIDVal).Error; err != nil {
				return errors.New("invalid parent_branch_id")
			}
		}
	}

	// Validate location IDs if being updated
	if countryID, ok := updatedData["country_id"]; ok && countryID != nil {
		var countryIDVal uint
		switch v := countryID.(type) {
		case float64:
			countryIDVal = uint(v)
		case uint:
			countryIDVal = v
		case int:
			countryIDVal = uint(v)
		}
		if countryIDVal > 0 {
			var country models.Country
			if err := config.DB.First(&country, countryIDVal).Error; err != nil {
				return errors.New("invalid country_id")
			}
		}
	}

	if stateID, ok := updatedData["state_id"]; ok && stateID != nil {
		var stateIDVal uint
		switch v := stateID.(type) {
		case float64:
			stateIDVal = uint(v)
		case uint:
			stateIDVal = v
		case int:
			stateIDVal = uint(v)
		}
		if stateIDVal > 0 {
			var state models.State
			if err := config.DB.First(&state, stateIDVal).Error; err != nil {
				return errors.New("invalid state_id")
			}
		}
	}

	now := time.Now()
	updatedData["updated_on"] = &now

	if err := config.DB.Model(&childBranch).Updates(updatedData).Error; err != nil {
		return err
	}
	return nil
}

// DeleteChildBranch deletes a child branch by ID
func DeleteChildBranch(childBranchID uint) error {
	if err := config.DB.Delete(&models.ChildBranch{}, childBranchID).Error; err != nil {
		return err
	}
	return nil
}

// *************************************** Child Branch Infrastructure ****************************************************** //

// CreateChildBranchInfrastructure creates a new child branch infrastructure record
func CreateChildBranchInfrastructure(infra *models.ChildBranchInfrastructure) error {
	infra.CreatedOn = time.Now()
	if err := config.DB.Create(infra).Error; err != nil {
		return err
	}
	return nil
}

// GetInfrastructureByChildBranch fetches infrastructure records by child branch ID
func GetInfrastructureByChildBranch(childBranchID uint) ([]models.ChildBranchInfrastructure, error) {
	var infra []models.ChildBranchInfrastructure
	if err := config.DB.Where("child_branch_id = ?", childBranchID).Preload("ChildBranch").Find(&infra).Error; err != nil {
		return nil, err
	}
	return infra, nil
}

// UpdateChildBranchInfrastructure updates a child branch infrastructure record
func UpdateChildBranchInfrastructure(id uint, updatedData map[string]interface{}) error {
	var infra models.ChildBranchInfrastructure
	if err := config.DB.First(&infra, id).Error; err != nil {
		return errors.New("infrastructure not found")
	}

	now := time.Now()
	updatedData["updated_on"] = &now

	if err := config.DB.Model(&infra).Updates(updatedData).Error; err != nil {
		return err
	}
	return nil
}

// DeleteChildBranchInfrastructure deletes a child branch infrastructure record
func DeleteChildBranchInfrastructure(id uint) error {
	if err := config.DB.Delete(&models.ChildBranchInfrastructure{}, id).Error; err != nil {
		return err
	}
	return nil
}

// *************************************** Child Branch Member ****************************************************** //

// CreateChildBranchMember creates a new child branch member
func CreateChildBranchMember(member *models.ChildBranchMember) error {
	member.CreatedOn = time.Now()
	if err := config.DB.Create(member).Error; err != nil {
		return err
	}
	return nil
}

// GetMembersByChildBranch fetches all members of a child branch
func GetMembersByChildBranch(childBranchID uint) ([]models.ChildBranchMember, error) {
	var members []models.ChildBranchMember
	if err := config.DB.Where("child_branch_id = ?", childBranchID).Preload("ChildBranch").Find(&members).Error; err != nil {
		return nil, err
	}
	return members, nil
}

// UpdateChildBranchMember updates a child branch member
func UpdateChildBranchMember(memberID uint, updatedData map[string]interface{}) error {
	var member models.ChildBranchMember
	if err := config.DB.First(&member, memberID).Error; err != nil {
		return errors.New("member not found")
	}

	now := time.Now()
	updatedData["updated_on"] = &now

	if err := config.DB.Model(&member).Updates(updatedData).Error; err != nil {
		return err
	}
	return nil
}

// DeleteChildBranchMember deletes a child branch member
func DeleteChildBranchMember(memberID uint) error {
	if err := config.DB.Delete(&models.ChildBranchMember{}, memberID).Error; err != nil {
		return err
	}
	return nil
}


