package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// struct user gorm model
type User struct {
	gorm.Model
	// ID          uint `gorm:"primaryKey"`
	// CreatedAt   time.Time
	// UpdatedAt   time.Time
	// DeletedAt   gorm.DeletedAt `gorm:"index"`
	Name        string `json:"name" form:"name"`
	Email       string `gorm:"unique" json:"email" form:"email"`
	Password    string `json:"password" form:"password"`
	Address     string `json:"address" form:"address"`
	PhoneNumber string `json:"phone_number" form:"phone_number"`
	Role        string `json:"role" form:"role"`
}

/*
TODO 1
buat struct products
id uint
created_at
updated_at
deleted_at
name string
user_id uint FK
description string
*/

type Product struct {
	gorm.Model
	Name string `json:"name" form:"name"`
	Description string `json:"description" form:"description"`
	UserID uint  `json:"user_id" form:"user_id"`
	User   User `gorm:"foreignKey:UserID"`
}

var DB *gorm.DB

// database connection
func InitDB() {
	// declare struct config & variable connectionString

	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading .env file")
	}
	connectionString := os.Getenv("CONNECTION_DB") + "?charset=utf8mb4&parseTime=True&loc=Local"

	DB, err = gorm.Open(mysql.Open(connectionString), &gorm.Config{})

	if err != nil {
		fmt.Println("Error initializing database:", err)
	}
}

// db migration
func InitialMigration() {
	DB.AutoMigrate(&User{})
	DB.AutoMigrate(&Product{})
	/*
		TODO 2
		migrate struct product 
	*/
}

// insert data user
func CreateUserController(c echo.Context) error {
	newUser := User{}
	errBind := c.Bind(&newUser) // mendapatkan data yang dikirim oleh FE melalui request body
	if errBind != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{
			"message": "error bind data. data not valid",
		})
	}

	// simpan ke DB
	tx := DB.Create(&newUser) // proses query insert
	if tx.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]any{
			"message": "error insert data. insert failed",
		})
	}

	return c.JSON(http.StatusCreated, map[string]any{
		"message": "insert success",
	})
}

// read data user
func GetAllUserController(c echo.Context) error {
	var usersData []User
	tx := DB.Find(&usersData) // select * from users;
	if tx.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]any{
			"message": "error read data",
		})
	}
	fmt.Println("users:", usersData)
	return c.JSON(http.StatusOK, map[string]any{
		"message": "success",
		"data":    usersData,
	})
}

func UpdateUserByIdController(c echo.Context) error {
	id := c.Param("user_id")
	idParam, errConv := strconv.Atoi(id)
	if errConv != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{
			"message": "error. id should be number",
		})
	}
	var userData = User{}
	errBind := c.Bind(&userData)
	if errBind != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{
			"message": "error bind data. data not valid",
		})
	}

	tx := DB.Model(&User{}).Where("id = ?", idParam).Updates(userData)
	if tx.Error != nil {
		// fmt.Println("err:", tx.Error)
		return c.JSON(http.StatusInternalServerError, map[string]any{
			"message": "error update " + tx.Error.Error(),
		})
	}

	if tx.RowsAffected == 0 {
		return c.JSON(http.StatusBadRequest, map[string]any{
			"message": "error record not found ",
		})
	}
	return c.JSON(http.StatusOK, map[string]any{
		"message": "success",
	})
}

func DeleteUserController (c echo.Context) error {
	 // mendapatkan id dari parameter
	id := c.Param("user_id")
	idParam, errConv := strconv.Atoi(id)

	// Jika terjadi kesalahan dalam konversi id, kembalikan pesan kesalahan
	if errConv != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{
			"message": "error. id should be number",
		})
	}

	// Jika Tidak Error Maka Menghapus pengguna dari database berdasarkan id
	tx := DB.Where("id = ?", idParam).Delete(&User{})
	if tx.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]any{
			"message": "error delete " + tx.Error.Error(),
		})
	}

	// Jika user tidak ditemukan, kembalikan pesan kesalahan
	if tx.RowsAffected == 0 {
		return c.JSON(http.StatusBadRequest, map[string]any{
			"message": "error record not found ",
		})
	}

	// Jika tidak error penghapusan berhasil, kembalikan pesan sukses
	return c.JSON(http.StatusOK, map[string]any{
		"messages": "success delete user by id",
	})
}

func CreateProductController (c echo.Context) error {
	newProduct := Product{}
	errBind := c.Bind(&newProduct) // mendapatkan data yang dikirim oleh FE melalui request body
	if errBind != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{
			"message": "error bind data. data not valid",
		})
	}

	// simpan ke DB
	tx := DB.Create(&newProduct) // proses query insert
	if tx.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]any{
			"message": "error insert data. insert failed",
		})
	}

	return c.JSON(http.StatusCreated, map[string]any{
		"message": "insert success",
	})
}

func GetAllProductController (c echo.Context) error {
	var productsData []Product
	tx := DB.Preload("User").Find(&productsData) // select * from products;
	if tx.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]any{
			"message": "error read data",
		})
	}
	fmt.Println("users:", productsData)
	return c.JSON(http.StatusOK, map[string]any{
		"message": "success",
		"data":    productsData,
	})
}

func GetProductByIdController (c echo.Context) error {
	// mendapatkan id dari parameter
	id := c.Param("product_id")
	idParam, errConv := strconv.Atoi(id)

	// Jika terjadi kesalahan dalam konversi id, kembalikan pesan kesalahan
	if errConv != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{
			"message": "error. id should be number",
		})
	}

	// Jika Tidak Erorr Maka Lanjut Mencari produk berdasarkan id
	var productData Product
	tx := DB.Preload("User").First(&productData, idParam)
	if tx.Error !=  nil {
		return c.JSON(http.StatusInternalServerError, map[string]any{
            "message": "error read data",
        })
	}
	// Jika produk tidak ditemukan Tampilkan Error
	if tx.RowsAffected == 0 {
        return c.JSON(http.StatusNotFound, map[string]any{
            "message": "product not found",
        })
    }
	// Jika produk ditemukan Tampilkan Success
	 return c.JSON(http.StatusOK, map[string]any{
        "message": "success",
        "data":    productData,
    })
}

func UpdateProductByIdController (c echo.Context) error {
	// mendapatkan id dari parameter
	id := c.Param("product_id")
	idParam, errConv := strconv.Atoi(id)

	// Jika terjadi kesalahan dalam konversi id, kembalikan pesan kesalahan
	if errConv != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{
			"message": "error. id should be number",
		})
	}


	var productData = Product{}
	errBind := c.Bind(&productData)
	if errBind != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{
			"message": "error bind data. data not valid",
		})
	}

	tx := DB.Model(&Product{}).Where("id = ?", idParam).Updates(productData)
	if tx.Error != nil {
		// fmt.Println("err:", tx.Error)
		return c.JSON(http.StatusInternalServerError, map[string]any{
			"message": "error update " + tx.Error.Error(),
		})
	}

	if tx.RowsAffected == 0 {
		return c.JSON(http.StatusBadRequest, map[string]any{
			"message": "error record not found ",
		})
	}
	return c.JSON(http.StatusOK, map[string]any{
		"message": "success",
	})
}

func DeleteProductController (c echo.Context) error {
	 // mendapatkan id dari parameter
	id := c.Param("product_id")
	idParam, errConv := strconv.Atoi(id)

	// Jika terjadi kesalahan dalam konversi id, kembalikan pesan kesalahan
	if errConv != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{
			"message": "error. id should be number",
		})
	}

	// Jika Tidak Error Maka Menghapus pengguna dari database berdasarkan id
	tx := DB.Where("id = ?", idParam).Delete(&Product{})
	if tx.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]any{
			"message": "error delete " + tx.Error.Error(),
		})
	}

	// Jika user tidak ditemukan, kembalikan pesan kesalahan
	if tx.RowsAffected == 0 {
		return c.JSON(http.StatusBadRequest, map[string]any{
			"message": "error record not found ",
		})
	}

	// Jika tidak error penghapusan berhasil, kembalikan pesan sukses
	return c.JSON(http.StatusOK, map[string]any{
		"messages": "success delete product by id",
	})
}

func main() {
	fmt.Println("running")
	InitDB()
	InitialMigration()

	// create a new echo instance
	e := echo.New()
	// define routes/ endpoint
	e.POST("/users", CreateUserController)
	e.GET("/users", GetAllUserController)
	e.PUT("/users/:user_id", UpdateUserByIdController)

	e.DELETE("/users/:user_id", DeleteUserController)
	e.POST("/products", CreateProductController)
	e.GET("/products", GetAllProductController)
	e.GET("/products/:product_id", GetProductByIdController)
	e.PUT("/products/:product_id", UpdateProductByIdController)
	e.DELETE("/products/:product_id", DeleteProductController)
	/*
	TODO 3
	tambahkan endpoint untuk:
		DELETE /users/:user_id
		POST /products
		GET /products
		GET /products/:product_id
		PUT /products/:product_id
		DELETE /products/:product_id
	*/

	//start server and port
	e.Logger.Fatal(e.Start(":8000"))
}