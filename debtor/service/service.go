package debtorservice

import (
	"context"
	"errors"
	"fmt"
	"log"

	"microsrv/model"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" // Mysql driver
)

// Service interface
type Service interface {
	Health() bool
	CreateDebtor(ctx context.Context, d model.Debtor) (model.Debtor, error)
	GetDebtor(ctx context.Context, id uint32) (model.Debtor, error)
	GetAll(ctx context.Context, p model.Pagination) (model.DebtorsResponse, error)
	Save(ctx context.Context, debtor model.Debtor, id uint) (model.Debtor, error)
	Delete(ctx context.Context, id uint) error
}

var (
	ErrInconsistentIDs = errors.New("inconsistent IDs")
	ErrAlreadyExists   = errors.New("already exists")
	ErrNotFound        = errors.New("not found")
)

type databaseStore struct{ db *gorm.DB }

// NewDB func
func NewDB(DSN string) Service {
	db, err := gorm.Open("mysql", DSN)
	if err != nil {
		fmt.Println(err.Error())
		panic("failed to connect database")
	}
	//defer db.Close()
	db = db.Set("gorm:table_options", "ENGINE=InnoDB")
	db.Set("sql_mode", "")
	return &databaseStore{db: db}
}

// Health implementation of the Service.
func (ds *databaseStore) Health() bool {
	return ds.db.DB().Ping() == nil
}

func (ds *databaseStore) CreateDebtor(ctx context.Context, d model.Debtor) (model.Debtor, error) {
	debtor := model.Debtor{}
	err := ds.db.Create(&d).Error
	if err != nil {
		return debtor, err
	}
	err = ds.db.
		Preload("Arbitration").
		Preload("BankDetails").
		First(&debtor, d.ID).
		Error
	if err != nil {
		return debtor, err
	}
	return debtor, nil
}

func (ds *databaseStore) GetDebtor(ctx context.Context, id uint32) (model.Debtor, error) {
	debtor := model.Debtor{}
	err := ds.db.
		Preload("Biddings").
		Preload("Arbitration").
		Preload("BankDetails").
		First(&debtor, id).
		Error
	if err != nil {
		return debtor, err
	}
	return debtor, nil
}

func (ds *databaseStore) GetAll(ctx context.Context, p model.Pagination) (model.DebtorsResponse, error) {
	if p.Limit == 0 {
		p.Limit = 20
	}
	res := model.DebtorsResponse{}
	debtors := model.Debtors{}
	count := 0
	q := ds.db
	if p.Name != "" {
		q = q.Where("name COLLATE UTF8_GENERAL_CI LIKE ?", fmt.Sprintf("%%%s%%", p.Name))
	}
	q1 := q.Table("debtors")
	err := q1.Where("deleted_at IS NULL").Count(&count).Error
	if err != nil {
		return res, err
	}
	res.Count = uint(count)
	if p.Sort == "" {
		q = q.Order("name")
	} else {
		q = q.Order("name desc")
	}
	err = q.
		Select("id, name, inn, ogrn, address, arbitration_id, case_no, decision_date, bankruptcy_manager_id").
		Preload("Biddings").
		Preload("BankruptcyManager").
		Preload("Arbitration").
		Limit(p.Limit).
		Offset(p.From).
		Find(&debtors).
		Error
	if err != nil {
		return res, err
	}
	res.Debtors = debtors
	return res, nil
}

// Save func
func (ds *databaseStore) Save(ctx context.Context, debtor model.Debtor, id uint) (model.Debtor, error) {
	dbtr := model.Debtor{}
	err := ds.db.
		Preload("BankDetails").
		Preload("Arbitration").
		Preload("Biddings").
		Find(&dbtr, id).
		Error
	if err != nil {
		log.Println(err)
		return dbtr, err
	}
	debtor.ID = dbtr.ID
	err = ds.db.Model(&dbtr).Updates(debtor).Error
	if err != nil {
		return dbtr, err
	}
	err = ds.db.
		Preload("BankDetails").
		Preload("Arbitration").
		Preload("Biddings").
		Find(&dbtr, id).
		Error
	if err != nil {
		return dbtr, err
	}
	return dbtr, nil
}

// Delete func
func (ds *databaseStore) Delete(ctx context.Context, id uint) error {
	return ds.db.Delete(model.Debtor{}, id).Error
}
