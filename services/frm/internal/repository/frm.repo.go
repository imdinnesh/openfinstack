package fraud

import (
	"github.com/imdinnesh/openfinstack/services/frm/models"
	"gorm.io/gorm"
)

type RuleRepository interface {
	Create(rule *models.Rule) error
	FindAll() ([]models.Rule, error)
	FindByID(id uint) (*models.Rule, error)
	Update(rule *models.Rule) error
	Delete(id uint) error
}

type ruleRepository struct {
	db *gorm.DB
}

func NewRuleRepository(db *gorm.DB) RuleRepository {
	return &ruleRepository{db}
}

func (r *ruleRepository) Create(rule *models.Rule) error {
	return r.db.Create(rule).Error
}

func (r *ruleRepository) FindAll() ([]models.Rule, error) {
	var rules []models.Rule
	err := r.db.Find(&rules).Error
	return rules, err
}

func (r *ruleRepository) FindByID(id uint) (*models.Rule, error) {
	var rule models.Rule
	err := r.db.First(&rule, id).Error
	return &rule, err
}

func (r *ruleRepository) Update(rule *models.Rule) error {
	return r.db.Save(rule).Error
}

func (r *ruleRepository) Delete(id uint) error {
	return r.db.Delete(&models.Rule{}, id).Error
}
