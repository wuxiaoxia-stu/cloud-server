package model

import (
	"aiyun_cloud_srv/app/controller/exam/model/exam_answer"
	"aiyun_cloud_srv/app/controller/exam/model/exam_user"
	"aiyun_cloud_srv/app/model/auth_client"
	"aiyun_cloud_srv/app/model/auth_leader_bind"
	"aiyun_cloud_srv/app/model/auth_leader_key"
	"aiyun_cloud_srv/app/model/auth_licence"
	"aiyun_cloud_srv/app/model/auth_licence_num"
	"aiyun_cloud_srv/app/model/auth_ukey"
	"aiyun_cloud_srv/app/model/hospital"
	"aiyun_cloud_srv/app/model/kl_feature"
	"aiyun_cloud_srv/app/model/kl_feature_atlas"
	"aiyun_cloud_srv/app/model/kl_syndrome"
	"aiyun_cloud_srv/app/model/kl_syndrome_feature"
	"aiyun_cloud_srv/app/model/kl_version"
	"aiyun_cloud_srv/app/model/order"
	"aiyun_cloud_srv/app/model/order_device"
	"aiyun_cloud_srv/app/model/order_product"
	"aiyun_cloud_srv/app/model/order_review"
	"aiyun_cloud_srv/app/model/order_upgrade"
	"aiyun_cloud_srv/app/model/sys_admin"
	"aiyun_cloud_srv/app/model/sys_role"
	"aiyun_cloud_srv/app/model/version"
	"fmt"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
	"os"
	"reflect"
	"strings"
)

type Table struct {
	TableName  string `orm:"table_name"`
	ColumnName string `orm:"column_name"`
	UdtName    string `orm:"udt_name"`
}

type Unique struct {
	Tablename string `orm:"tablename"`
	Indexname string `orm:"indexname"`
}

//分页查询
type PageReqParams struct {
	KeyWord   string `p:"keyword"`
	Status    int    `p:"status" default:"-1"`
	Page      int    `p:"page" default:"1"`
	PageSize  int    `p:"page_size" default:"10"`
	Order     string `p:"order" default:"id"`
	Sort      string `p:"sort" default:"DESC"`
	StartTime string `p:"start_time"`
	EndTime   string `p:"end_time"`
}

//状态设置
type SetStatusParams struct {
	Id     int `json:"id" p:"id" v:"required#参数错误"`
	Status int `json:"status" p:"status" v:"required|in:0,1#参数错误|参数错误"`
}

//批量操作提交信息绑定
type Ids struct {
	Ids []int `json:"ids"`
}

func DbInit() {
	//ImportData()
	// 自动建表
	if err := dbAutoMigrate(
		&hospital.Entity{},
		&sys_admin.Entity{},
		&sys_role.Entity{},
		&auth_client.Entity{},
		&auth_licence.Entity{},
		&auth_licence_num.Entity{},
		&auth_ukey.Entity{},
		&auth_leader_key.Entity{},
		&auth_leader_bind.Entity{},
		&order.Entity{},
		&order_device.Entity{},
		&order_product.Entity{},
		&order_review.Entity{},
		&order_upgrade.Entity{},
		&version.Entity{},
		&kl_feature.Entity{},
		&kl_feature_atlas.Entity{},
		&kl_syndrome.Entity{},
		&kl_syndrome_feature.Entity{},
		&kl_version.Entity{},
		&exam_user.Entity{},
		&exam_answer.Entity{},
	); err != nil {
		panic("数据库初始化失败")
	}

	if err := initAdminData(); err != nil {
		panic("Admin表数据初始化失败")
	}
}

//导入数据
func ImportData() (err error) {
	file, err := os.Open("./data/public_v1.sql")
	if err != nil {
		g.Log().Error(err)
		return
	}
	defer file.Close()

	fileinfo, err := file.Stat()
	if err != nil {
		g.Log().Error(err)
		return
	}

	filesize := fileinfo.Size()
	buffer := make([]byte, filesize)

	_, err = file.Read(buffer)
	if err != nil {
		g.Log().Error(err)
		return
	}

	sql := string(buffer)

	_, err = g.DB().Exec(sql)
	if err != nil {
		g.Log().Error("数据导入失败", err)
		return
	} else {
		g.Log().Info("数据导入成功")
	}

	return
}

//初始化数据
func initAdminData() (err error) {
	tx, err := g.DB().Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
		}
	}()

	var admin *sys_admin.Entity
	err = tx.Model(sys_admin.Table).Where("username", "admin").Scan(&admin)
	if err != nil {
		g.Log().Error(err)
		tx.Rollback()
		return
	}

	if admin == nil {
		_, err = tx.Model(sys_admin.Table).Data(g.Map{
			"role_id":  1,
			"username": "admin",
			"password": "315193380d2fedffa677b7fe236fadb5",
			"salt":     "sinb",
			"email":    "778774780@qq.com",
			"status":   1,
		}).Insert()

		if err != nil {
			g.Log().Error("初始化sys_admin表数据失败", err)
			tx.Rollback()
			return
		} else {
			g.Log().Info("初始化sys_admin表数据成功")
		}
	}

	var role *sys_role.Entity
	err = tx.Model(sys_role.Table).Where("id", 1).Scan(&role)
	if err != nil {
		g.Log().Error(err)
		tx.Rollback()
		return
	}

	if role == nil {
		_, err = tx.Model(sys_role.Table).Data(g.Map{
			"id":     1,
			"name":   "超级管理员",
			"rule":   "[\"auth\",\"order/list\",\"order/list\",\"order/create\",\"order/details\",\"order/upgrade\",\"order/review\",\"order/manager_review\",\"order/deploy\",\"order/revoke\",\"order/delete\",\"order/device_details\",\"licence/list\",\"licence/list\",\"licence/set-status\",\"licence/delete\",\"kl\",\"kl-syndrome/tree\",\"kl-syndrome/tree\",\"kl-syndrome/add\",\"kl-syndrome/edit\",\"kl-syndrome/set-status\",\"kl-syndrome/delete\",\"kl-feature/tree\",\"kl-feature/tree\",\"kl-feature/add\",\"kl-feature/edit\",\"kl-feature/set-status\",\"kl-feature/delete\",\"kl-version/list\",\"kl-version/list\",\"kl-version/add\",\"kl-version/download\",\"kl-version/delete\",\"init\",\"ukey/list\",\"ukey/list\",\"ukey/add\",\"ukey/set-status\",\"ukey/delete\",\"leader-key/list\",\"leader-key/list\",\"leader-key/add\",\"leader-key/set-status\",\"leader-key/delete\",\"device\",\"server/list\",\"client/list\",\"version/list\",\"version/list\",\"version/add\",\"version/edit\",\"version/set-status\",\"version/delete\",\"sys\",\"sys-admin/list\",\"sys-admin/list\",\"sys-admin/add\",\"sys-admin/edit\",\"sys-admin/set-status\",\"sys-admin/delete\",\"sys-role/list\",\"sys-role/list\",\"sys-role/add\",\"sys-role/edit\",\"sys-role/set-status\",\"sys-role/delete\",\"hospital/list\",\"hospital/list\",\"hospital/add\",\"hospital/edit\",\"hospital/set-status\",\"hospital/delete\"]",
			"status": 1,
		}).Insert()

		if err != nil {
			g.Log().Error("初始化sys_role表数据失败", err)
			tx.Rollback()
			return
		} else {
			g.Log().Info("初始化sys_role表数据成功")
		}
	}

	return tx.Commit()
}

//自动更新表结构，添加新字段
func dbAutoMigrate(models ...interface{}) (err error) {

	if len(models) <= 0 {
		return
	}

	for _, model := range models {
		t := reflect.TypeOf(model)

		if t.Kind() != reflect.Ptr || t.Elem().Kind() != reflect.Struct {
			err = gerror.New("模型参数应为结构体指针")
			g.Log().Error(err)
			return err
		}

		model_name := t.String()
		model_arr := strings.Split(strings.TrimLeft(model_name, "*"), ".")
		if len(model_arr) != 2 {
			err = gerror.New("模型参数错误")
			g.Log().Error(err)
			return err
		}

		//判断表是否存在
		//select * from pg_tables where tablename = 'lpm_sys_role'
		//select * from information_schema.TABLES where TABLE_NAME = 'lpm_sys_role';

		table_prefix := g.Cfg().GetString("database.prefix")
		full_table_name := table_prefix + model_arr[0]

		table := Table{}
		err = g.DB().GetScan(&table, "select * from information_schema.TABLES where TABLE_NAME = ?", full_table_name)
		if err != nil {
			if err.Error() == "sql: no rows in result set" {
			} else {
				g.Log().Error(err)
				return
			}
		}

		if table.TableName == "" {
			//CREATE TABLE "public"."lpm_sys_role" ();
			_, err = g.DB().Exec("CREATE TABLE \"" + full_table_name + "\" ()")
			if err != nil {
				g.Log().Error(err)
				return
			} else {
				g.Log().Infof("%s表创建成功", full_table_name)
			}
		}

		v := reflect.ValueOf(model).Elem()

		//查询表字段
		//select * from information_schema.COLUMNS where table_name = 'lpm_sys_admin'
		columns := []Table{}
		err = g.DB().GetScan(&columns, "select * from information_schema.COLUMNS where table_name = ?", full_table_name)
		if err != nil {
			if err.Error() == "sql: no rows in result set" {
			} else {
				g.Log().Error(err)
				return
			}
		}

		add_colums_sql := []string{}       //字段信息
		alter_colums_comment := []string{} //备注信息
		table_comment := ""
		unique := []string{}
		for i := 0; i < v.NumField(); i++ {

			tagOrmInfo := v.Type().Field(i).Tag.Get("orm")
			if tagOrmInfo == "" {
				break
			}

			tag_orms := strings.Split(tagOrmInfo, ",")
			if len(tag_orms) <= 0 {
				err = gerror.New("")
				g.Log().Error(err)
				return err
			}

			column := tag_orms[0]
			column_type := v.Field(i).Type().String()
			size := "4"
			default_val := ""
			is_primary := false

			not_null := ""
			for _, v2 := range tag_orms {
				if strings.Contains(v2, "primary") {
					is_primary = true
				}

				if strings.Contains(v2, "table_comment:") {
					table_comment = strings.TrimLeft(v2, "table_comment:")
				}

				if strings.Contains(v2, "unique") {
					unique = append(unique, column)
				}

				if strings.Contains(v2, "size:") {
					size = strings.TrimLeft(v2, "size:")
				}

				if strings.Contains(strings.ToUpper(v2), "NOT NULL") {
					not_null = "NOT NULL"
				}

				if strings.Contains(v2, "default:") {
					default_val = strings.TrimLeft(v2, "default:")
				}

				if strings.Contains(v2, "comment:") && !strings.Contains(v2, "table_comment:") {
					comment := strings.TrimLeft(v2, "comment:")
					if comment != "" {
						alter_colums_comment = append(alter_colums_comment, fmt.Sprintf("COMMENT ON COLUMN %s.%s IS %s", full_table_name, column, comment))
					}
				}
			}

			column_exist := false
			for _, v2 := range columns {
				if column == v2.ColumnName {
					column_exist = true
					break
				}
			}

			//如果字段你不存在，就创建字段， 预定义字段类型
			// int => int(size)
			// string =>  varchar(size)
			// gtime.Time =>  timestamptz(6)
			if !column_exist {
				t := "int"
				if column_type == "string" {
					t = "varchar"
				} else if strings.Contains(column_type, "gtime.Time") {
					t = "timestamptz"
					size = "6"
				}

				if size == "text" {
					t = "text"
				}

				sql := "ADD COLUMN "
				if t == "int" {
					sql += column + " " + t + size
				} else if t == "text" {
					sql += column + " " + t
				} else {
					sql += column + " " + t + "(" + size + ")"
				}

				if is_primary {
					sql += " NOT NULL GENERATED BY DEFAULT AS IDENTITY ( INCREMENT 1 MINVALUE  1 START 1 CACHE 1)"
				}

				if not_null != "" {
					sql += " NOT NULL"
				}

				if default_val != "" {
					sql += " DEFAULT " + default_val
				}

				add_colums_sql = append(add_colums_sql, sql)
			}

		}

		//ALTER TABLE lpm_sys_role
		//	ADD COLUMN id int4 NOT NULL GENERATED BY DEFAULT AS IDENTITY ( INCREMENT 1 MINVALUE  1 START 1 CACHE 1),
		//	ADD COLUMN name VARCHAR(50) NOT NULL DEFAULT '',
		//	ADD COLUMN create_at timestamptz(6)
		if len(add_colums_sql) > 0 {
			_, err = g.DB().Exec("ALTER TABLE " + full_table_name + " " + strings.Join(add_colums_sql, ","))
			if err != nil {
				g.Log().Error(err)
				return
			}
		}

		//修改行注释
		//comment on column table_name.column_name is '名称';
		if len(alter_colums_comment) > 0 {
			_, err = g.DB().Exec(strings.Join(alter_colums_comment, ";"))
			if err != nil {
				g.Log().Error(err)
				return
			}
		}
		//修改表注释
		//comment on table table_name is '表名称';
		if table_comment != "" {
			_, err = g.DB().Exec(fmt.Sprintf("COMMENT ON TABLE %s IS %s;", full_table_name, table_comment))
			if err != nil {
				g.Log().Error(err)
				return
			}
		}

		// 添加唯一索引
		//alter table lpm_auth_licence ADD CONSTRAINT unique_author_number UNIQUE(author_number)
		if len(unique) > 0 {
			//查询索引是否存在,不存在则创建
			//select * from pg_indexes where tablename = 'lpm_auth_licence';
			unique_data := []Unique{}
			err = g.DB().GetScan(&unique_data, "select * from pg_indexes where tablename = ?", full_table_name)
			if err != nil {
				g.Log().Error(err)
				return
			}
			for _, v := range unique {
				index_exist := false
				for _, v2 := range unique_data {
					if v2.Indexname == "unique_"+v {
						index_exist = true
					}
				}

				if !index_exist {
					_, err = g.DB().Exec(fmt.Sprintf("alter table %s ADD CONSTRAINT unique_%s UNIQUE(%s);", full_table_name, v, v))
					if err != nil {
						g.Log().Error(err)
						return
					}
				}

			}
		}
	}

	return
}
