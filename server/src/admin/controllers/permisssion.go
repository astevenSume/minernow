package controllers

import (
	"admin/controllers/errcode"
	"common"
	"encoding/json"
	"fmt"
	"strconv"
	"utils/admin/dao"
	"utils/admin/models"
)

type PermissionController struct {
	BaseController
}

func (c *PermissionController) GetPageInfo() (uint64, int, int, error) {
	//请求参数
	var id uint64
	strId := c.GetString(KEY_ID)
	if len(strId) > 0 {
		var err error
		id, err = strconv.ParseUint(strId, 10, 64)
		if err != nil {
			return 0, 0, 0, err
		}
	}
	page, err := c.GetInt(KEY_PAGE, 1)
	if err != nil {
		common.LogFuncError("err:%v", err)
		return 0, 0, 0, err
	}
	perPage, err := c.GetInt(KEY_LIMIT, DEFAULT_PER_PAGE)
	if err != nil {
		common.LogFuncError("err:%v", err)
		return 0, 0, 0, err
	}

	return id, page, perPage, nil
}

//获取所有权限
func (c *PermissionController) HandleGetPermissions() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionReadPermission, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	id, page, perPage, err := c.GetPageInfo()
	if err != nil {
		c.ErrorResponseAndLog(OPActionReadPermission, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	//数据获取
	var input string
	if id == 0 {
		//分页查询
		res := map[string]interface{}{}
		input = fmt.Sprintf("{\"id\":%v,\"page\":%v,\"per_page\":%v}", id, page, perPage)
		permissions, pageInf, err := dao.PermissonDaoEntity.QueryPagePermission(page, perPage)
		if err != nil {
			c.ErrorResponseAndLog(OPActionReadPermission, controllers.ERROR_CODE_DB, input)
			return
		}
		res["meta"] = pageInf
		res["list"] = permissions
		c.SuccessResponseAndLog(OPActionReadPermission, input, res)
	} else {
		input = fmt.Sprintf("{\"id\":%v}", id)
		permission := models.Permission{Id: uint64(id)}
		err := dao.PermissonDaoEntity.QueryPermission(&permission, models.COLUMN_Permission_Id)
		if err != nil {
			c.ErrorResponseAndLog(OPActionReadPermission, controllers.ERROR_CODE_DB, input)
			return
		}
		c.SuccessResponseAndLog(OPActionReadPermission, input, permission)
	}

}

//创建权限
func (c *PermissionController) HandleCreatePermissions() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionAddPermission, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	type msg struct {
		Slug string `json:"slug"`
		Desc string `json:"desc"`
	}
	req := &msg{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &req)
	if err != nil {
		common.LogFuncDebug("json decode: %s \nfailed : %v", string(c.Ctx.Input.RequestBody), err)
		c.ErrorResponseAndLog(OPActionAddPermission, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	if req.Slug == "" || req.Desc == "" {
		c.ErrorResponseAndLog(OPActionAddPermission, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	//插入数据
	permission := new(models.Permission)
	permission.Slug = req.Slug
	permission.Desc = req.Desc
	permission.Ctime = common.NowInt64MS()
	permission.Utime = common.NowInt64MS()
	id, err := dao.PermissonDaoEntity.CreatePermissionBySlug(permission)
	if err != nil {
		c.ErrorResponseAndLog(OPActionAddPermission, controllers.ERROR_CODE_DB, string(c.Ctx.Input.RequestBody))
		return
	}
	permission.Id = id

	//返回c查询结果
	c.SuccessResponseAndLog(OPActionAddPermission, string(c.Ctx.Input.RequestBody), permission)
}

//更新权限
func (c *PermissionController) HandleUpdatePermissions() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionEditPermission, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	type msg struct {
		Id   int    `json:"id"`
		Slug string `json:"slug"`
		Desc string `json:"desc"`
	}
	req := &msg{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &req)
	if err != nil {
		common.LogFuncDebug("json decode: %s \nfailed : %v", string(c.Ctx.Input.RequestBody), err)
		c.ErrorResponseAndLog(OPActionEditPermission, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	if req.Slug == "" || req.Desc == "" {
		c.ErrorResponseAndLog(OPActionEditPermission, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	permission := models.Permission{
		Id:    uint64(req.Id),
		Slug:  req.Slug,
		Desc:  req.Desc,
		Utime: common.NowInt64MS(),
	}
	err = dao.PermissonDaoEntity.UpdatePermissionBySlug(&permission, models.COLUMN_Permission_Slug, models.COLUMN_Permission_Desc, models.COLUMN_Permission_Utime)
	if err != nil {
		c.ErrorResponseAndLog(OPActionEditPermission, controllers.ERROR_CODE_DB, string(c.Ctx.Input.RequestBody))
		return
	}

	//返回c查询结果
	c.SuccessResponseAndLog(OPActionEditPermission, string(c.Ctx.Input.RequestBody), permission)
}

//删除权限
func (c *PermissionController) HandleDeletePermissions() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionDelPermission, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	id, err := c.GetInt("id", 0)
	if err != nil || id == 0 {
		c.ErrorResponseAndLog(OPActionDelPermission, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}
	input := fmt.Sprintf("{\"id\":%v}", id)

	err = dao.PermissonDaoEntity.DelPermission(uint64(id))
	if err != nil {
		c.ErrorResponseAndLog(OPActionDelPermission, controllers.ERROR_CODE_DB, input)
		return
	}

	//返回c查询结果
	data := map[string]interface{}{
		"id": id,
	}
	c.SuccessResponseAndLog(OPActionDelPermission, input, data)
}
