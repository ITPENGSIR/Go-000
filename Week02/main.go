/**
* @Author: 彭光豪
* @Date: 12/1/20 10:45 PM
 */
package main

import (
	"Backend/src/utils/logger"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

var Db *sqlx.DB

func init() {
	database, err := sqlx.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/user")
	if err != nil {
		fmt.Println("open mysql failed,", err)
		return
	}
	Db = database
}

//自定义错误结构体
type ErrorInfo struct {
	Code 	int
	Msg		string
}
//定义数据结构体
type UserInfo struct {
	Id 		int
	Name 	string
}
func (err *ErrorInfo)Error()string{
	return err.Msg
}

var (
	ErrNotFound=&ErrorInfo{Code:100001,Msg:"user not found"}
	ErrServer  =&ErrorInfo{Code:100002,Msg:"server  err"}
)

//dao层，获取用户信息
func GetUserInfoById(id int)([]UserInfo,error){
	var usr []UserInfo
	que:="select id, name from user where id = ? "
	err := Db.Select(&usr, que, id)
	err = sql.ErrNoRows
	if err == sql.ErrNoRows {
		return nil,errors.Wrap(ErrNotFound,fmt.Sprintf("GetUserInfoById id= %d is not found",id))
	}
	if err != nil {
		return nil,errors.Wrap(ErrServer,fmt.Sprintf("GetUserInfoById err sql=%s err=%v ",que,err))
	}
	return usr,err
}

//server

func GetUserServer(id int)([]UserInfo,error){
	return GetUserInfoById(id)
}

func main(){
	data, err := GetUserServer(100)
	if errors.Is(err, ErrNotFound) {
		logger.Error(err)
		//返回对应业务码
	}
	fmt.Println(data,err)
}