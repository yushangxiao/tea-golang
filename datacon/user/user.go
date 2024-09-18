package user

import (
	"errors"
	"fmt"
	"tea/sqllite"

	"gorm.io/gorm"
)

type User struct {
	Name       string `json:"name"`
	Password   string `json:"password"`
	IsAdmin    bool   `json:"is_admin"`
	AdminSerct string `json:"adminSerct"`
	gorm.Model
}

type UserManager struct {
	db *gorm.DB
}

func NewUserManager(db *gorm.DB) *UserManager {
	return &UserManager{db: db}
}

func (manager *UserManager) CreateUser(user User) error {
	if user.Name == "" || user.Password == "" {
		return errors.New("必须填写用户名和密码")
	}
	// 检查name是否已经存在
	var count int64
	manager.db.Model(&User{}).Where("name = ?", user.Name).Count(&count)
	if count > 0 {
		return errors.New("用户名已存在")
	}
	result := manager.db.Create(&user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (manager *UserManager) DeleteUser(name string) error {
	var user User
	if err := manager.db.Where("name = ?", name).First(&user).Error; err != nil {
		return err
	}

	result := manager.db.Delete(&user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (manager *UserManager) GetAllUsers() ([]User, error) {
	var users []User
	result := manager.db.Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	// 删除用户密码
	for i := range users {
		users[i].Password = ""
	}
	return users, nil
}

// 登录
func (manager *UserManager) Login(user User) error {
	var u User
	if err := manager.db.Where("name = ? AND password = ?", user.Name, user.Password).First(&u).Error; err != nil {
		return err
	}
	return nil
}

// 检查是否是管理员
func (manager *UserManager) CheckAdmin(name string) (bool, error) {
	var user User
	if err := manager.db.Where("name = ?", name).First(&user).Error; err != nil {
		return false, err
	}
	return user.IsAdmin, nil
}

// 查询用户
func (manager *UserManager) GetUserInfo(name string) (User, error) {
	var user User
	var count int64
	manager.db.Model(&User{}).Where("name = ?", name).Count(&count)
	if count == 0 {
		return user, errors.New("用户不存在")
	}
	if err := manager.db.Where("name = ?", name).First(&user).Error; err != nil {
		return user, err
	}
	// 删除Passwor字段
	user.Password = ""
	return user, nil
}

var MyUserManager *UserManager

func init() {
	db := sqllite.GetDB()
	if !db.Migrator().HasTable(&User{}) {
		db.AutoMigrate(&User{})
		fmt.Println("创建User表")
	}
	MyUserManager = NewUserManager(db)
}
