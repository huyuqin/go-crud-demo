package main

import (
	"fmt"
	"strconv"
	"time"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type List struct{
	gorm.Model
	Name string	`gorm:"type:varchar(20);not null" json:"name" binding:"required"`
	State string `gorm:"type:varchar(20);not null" json:"state" binding:"required"`
	Phone string `gorm:"type:varchar(20);not null" json:"phone" binding:"required"`
	Email string `gorm:"type:varchar(40);not null" json:"email" binding:"required"`
	Address string `gorm:"type:varchar(200);not null" json:"address" binding:"required"`
}

func main(){
	dsn:="root:123456@tcp(127.0.0.1:3306)/demo2?charset=utf8&parseTime=true"
	db,err:=gorm.Open(mysql.Open(dsn),&gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err!=nil{
		panic(err)
	}
	fmt.Println(db)

	sqlDB,err:=db.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(10*time.Second)
	db.AutoMigrate(&List{})

	r:=gin.Default()

	r.POST("/user/add",func(c *gin.Context){
		var data List
		err:=c.ShouldBindJSON(&data)

		if err!=nil{
			c.JSON(200,gin.H{
				"msg":"添加失败",
				"data":gin.H{},
				"code":400,
			})
		}else{
			db.Create(&data)
			c.JSON(200,gin.H{
				"msg":"添加成功",
				"data":data,
				"code":200,
			})
		}
	})
	
	// 找到对应的ID所对应的条目
	// 判断ID存在
	// 数据库中删除
	// 返回ID没有找到
	r.DELETE("/user/delete/:id",func(c *gin.Context) {
		var data []List
		id:=c.Param("id")
		db.Where("id=?",id).Find(&data)
		if len(data)==0{
			c.JSON(200,gin.H{
				"msg":"id error",
				"code":400,
			})
		}else{
			db.Delete(&data)
			c.JSON(200,gin.H{
				"msg":"delete success",
				"code":200,
			})
		}
	})

	r.GET("/user/list/:name",func(c *gin.Context) {
		name:=c.Param("name")
		var dataList []List
		db.Where("name=?",name).Find(&dataList)
		if len(dataList)==0{
			c.JSON(200,gin.H{
				"msg":"没找到数据",
				"code":400,
			})
		}else{
			c.JSON(200,gin.H{
				"msg":"找到了",
				"code":200,
				"data":dataList,
			})
		}
	})

	r.PUT("/user/update/:id",func(c *gin.Context){
		var data List
		id:=c.Param("id")
		db.Select("id").Where("id=?",id).Find(&data)
		if data.ID==0{
			c.JSON(200,gin.H{
				"msg":"没找到数据",
				"code":400,
			})
		}else{
			err:=c.ShouldBindJSON(&data)
			if err!=nil{
				c.JSON(200,gin.H{
					"msg":"修改失败",
					"code":400,
				})
			}else{
				db.Where("id=?",id).Updates(&data)
				c.JSON(200,gin.H{
					"msg":"修改成功",
					"code":200,
				})
			}
		}
	})
	
	r.GET("/user/list",func(c *gin.Context) {
		var data []List
		var total int64
		pageNum,_:=strconv.Atoi(c.Query("pageNum"))
		pageSize,_:=strconv.Atoi(c.Query("pageSize"))
		if pageNum==0{
			pageNum= -1
		}
		if pageSize==0{
			pageSize=-1
		}
		offsetVal:=(pageNum-1)*pageSize
		
		if pageNum==-1&&pageSize==-1{
			offsetVal=-1
		}
		db.Model(data).Count(&total).Limit(pageSize).Offset(offsetVal).Find(&data)
		if len(data)==0{
			c.JSON(200,gin.H{
				"msg":"没找到数据",
				"code":400,
				"data":gin.H{},
			})
		}else{
			c.JSON(200,gin.H{
				"msg":"找到数据",
				"code":200,
				"data":gin.H{
					"list":data,
					"total":total,
					"pageNum":pageNum,
					"pageSize":pageSize,
				},
			})
		}
	})
	
	
	
	r.Run()



}