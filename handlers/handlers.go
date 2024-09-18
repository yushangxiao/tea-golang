package handlers

import (
	"tea/datacon/product"
	"tea/datacon/user"
	"tea/utils"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// var jwtKey = []byte("eastTea")

type Claims struct {
	Username string `json:"username"`
	IsAdmin  bool   `json:"is_admin"`
	jwt.StandardClaims
}

var userManager *user.UserManager
var productManager *product.ProductManager
var jwtKey []byte
var admimSerctList []string
var rootSerct string

func init() {
	userManager = user.MyUserManager
	productManager = product.MyProductManager
	config := utils.MyConfig
	jwtKey = []byte(config.JwtKey)
	admimSerctList = config.AdminSerct
	rootSerct = config.RootSerct
}

func PingHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func CreateUserHandler(c *gin.Context) {
	var user user.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(400, gin.H{"success": false, "error": err.Error()})
		return
	}
	user.IsAdmin = false
	// adminSerct := "eastTea"
	clientAdminSecret := user.AdminSerct
	// fmt.Println(clientAdminSecret)
	if clientAdminSecret != "" {
		// for 判断密钥
		for _, v := range admimSerctList {
			if clientAdminSecret == v {
				user.IsAdmin = true
			}
		}
		if !user.IsAdmin {
			c.JSON(400, gin.H{"success": false, "error": "管理员密钥错误"})
			return
		}
	}
	if err := userManager.CreateUser(user); err != nil {
		c.JSON(500, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"success": true})
}

func GetAllUsersHandler(c *gin.Context) {
	users, err := userManager.GetAllUsers()
	if err != nil {
		c.JSON(500, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"success": true, "data": users})
}

func DeleteUserHandler(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		c.JSON(400, gin.H{"success": false, "error": "name must be provided"})
		return
	}
	if err := userManager.DeleteUser(name); err != nil {
		c.JSON(500, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"success": true})
}

func LoginHandler(c *gin.Context) {
	var user user.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(400, gin.H{"success": false, "error": "数据格式错误"})
		return
	}
	if err := userManager.Login(user); err != nil {
		c.JSON(500, gin.H{"success": false, "error": "用户名或密码错误"})
		return
	}
	// 设置 JWT 令牌的有效期
	expirationTime := time.Now().Add(24 * time.Hour)
	IsAdmin, _ := userManager.CheckAdmin(user.Name)
	// fmt.Println(IsAdmin)
	claims := &Claims{
		Username: user.Name,
		IsAdmin:  IsAdmin,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// 使用密钥生成 JWT 令牌
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// fmt.Println("login", jwtKey)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		c.JSON(500, gin.H{"success": false, "error": "无法生成令牌"})
		return
	}

	// 返回 JWT 令牌给客户端
	c.JSON(200, gin.H{"success": true, "token": tokenString})
}

func UserInfoHandler(c *gin.Context) {
	username := c.MustGet("userInfo").(*utils.Claims).Username
	userInfo, err := userManager.GetUserInfo(username)
	if err != nil {
		c.JSON(500, gin.H{"success": false, "error": "无法获取用户信息"})
		return
	}
	c.JSON(200, gin.H{"success": true, "data": userInfo})
}
func GetProductsHandler(c *gin.Context) {
	products, err := productManager.GetProducts()
	if err != nil {
		c.JSON(500, gin.H{"success": false, "error": "无法获取产品信息"})
		return
	}
	c.JSON(200, gin.H{"success": true, "data": products})
}
func CreateProductHandler(c *gin.Context) {
	var product product.Product
	if err := c.BindJSON(&product); err != nil {
		c.JSON(400, gin.H{"success": false, "error": "数据格式错误"})
		return
	}
	product.IsOnSale = true
	if err := productManager.CreateProduct(product); err != nil {
		c.JSON(500, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"success": true})
}
func DeleteProductHandler(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		c.JSON(400, gin.H{"success": false, "error": "name must be provided"})
		return
	}
	if err := productManager.DeleteProduct(name); err != nil {
		c.JSON(500, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"success": true})
}
func ChangeSaleHandler(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		c.JSON(400, gin.H{"success": false, "error": "name must be provided"})
		return
	}
	if err := productManager.ChangeSale(name); err != nil {
		c.JSON(500, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"success": true})
}

type DelAdminModel struct {
	Name      string `json:"name"`
	RootSerct string `json:"rootSerct"`
}

func DeleteAdminHandler(c *gin.Context) {
	var delAdminModel DelAdminModel
	if err := c.BindJSON(&delAdminModel); err != nil {
		c.JSON(400, gin.H{"success": false, "error": err.Error()})
		return
	}
	if delAdminModel.RootSerct != rootSerct {
		c.JSON(400, gin.H{"success": false, "error": "Root密钥错误"})
		return
	}
	username := c.MustGet("userInfo").(*utils.Claims).Username
	if username == delAdminModel.Name {
		c.JSON(400, gin.H{"success": false, "error": "无法删除自己"})
		return
	}
	if err := userManager.DeleteUser(delAdminModel.Name); err != nil {
		c.JSON(500, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"success": true})
}
