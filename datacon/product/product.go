package product

import (
	"errors"
	"fmt"
	"tea/sqllite"

	"gorm.io/gorm"
)

type Product struct {
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	IsOnSale bool    `json:"is_on_sale"`
	gorm.Model
}
type ProductManager struct {
	db *gorm.DB
}

func NewProductManager(db *gorm.DB) *ProductManager {
	return &ProductManager{db: db}
}
func (manager *ProductManager) CreateProduct(product Product) error {
	if product.Name == "" || product.Price <= 0 {
		return errors.New("请检查产品名称和价格")
	}
	// 检查name是否已经存在
	var count int64
	manager.db.Model(&Product{}).Where("name = ?", product.Name).Count(&count)
	if count > 0 {
		return errors.New("产品已存在")
	}
	result := manager.db.Create(&product)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (manager *ProductManager) DeleteProduct(name string) error {
	var product Product
	if err := manager.db.Where("name = ?", name).First(&product).Error; err != nil {
		return err
	}
	if product.IsOnSale {
		return errors.New("On sale product can not be deleted")
	}
	result := manager.db.Delete(&product)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (manager *ProductManager) GetProducts() ([]Product, error) {
	var products []Product
	result := manager.db.Find(&products)
	if result.Error != nil {
		return nil, result.Error
	}
	// 排序，优先展示在售商品
	saleList := make([]Product, 0)
	noSaleList := make([]Product, 0)
	for _, product := range products {
		if product.IsOnSale {
			saleList = append(saleList, product)
		} else {
			noSaleList = append(noSaleList, product)
		}
	}
	products = append(saleList, noSaleList...)
	// 增加数量，测试
	// products = append(products, noSaleList...)
	// products = append(products, noSaleList...)
	// products = append(products, noSaleList...)
	return products, nil
}
func (manager *ProductManager) ChangeSale(name string) error {
	var product Product
	if err := manager.db.Where("name = ?", name).First(&product).Error; err != nil {
		return err
	}
	product.IsOnSale = !product.IsOnSale
	result := manager.db.Save(&product)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

var MyProductManager *ProductManager

func init() {
	db := sqllite.GetDB()
	if !db.Migrator().HasTable(&Product{}) {
		db.AutoMigrate(Product{})
		fmt.Println("创建Product表")
	}
	MyProductManager = NewProductManager(db)
}
