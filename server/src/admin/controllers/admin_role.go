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

type AdminRoleController struct {
	BaseController
}

func (c *AdminRoleController) GetPageInfo() (uint64, int, int, error) {
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

//获取角色信息
func (c *AdminRoleController) HandleGetRoles() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionReadRole, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	id, page, perPage, err := c.GetPageInfo()
	if err != nil {
		c.ErrorResponseAndLog(OPActionReadRole, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	var input string
	if id == 0 {
		res := map[string]interface{}{}
		//分页数据
		input = fmt.Sprintf("{\"id\":%v,\"page\":%v,\"per_page\":%v}", id, page, perPage)
		roles, pageInf, err := dao.RoleDaoEntity.QueryPageRole(page, perPage)
		if err != nil {
			c.ErrorResponseAndLog(OPActionReadRole, controllers.ERROR_CODE_DB, input)
			return
		}
		res["meta"] = pageInf
		res["list"] = roles
		c.SuccessResponseAndLog(OPActionReadRole, input, res)
	} else {
		input = fmt.Sprintf("{\"page\":%v,\"per_page\":%v}", page, perPage)
		role := &models.Role{Id: uint64(id)}
		info, err := dao.RoleDaoEntity.QueryRole(role, models.COLUMN_Role_Id)
		if err != nil {
			c.ErrorResponseAndLog(OPActionReadRole, controllers.ERROR_CODE_DB, input)
			return
		}
		c.SuccessResponseAndLog(OPActionReadRole, input, info)
	}
}

//创建角色
func (c *AdminRoleController) HandleCreateRole() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionAddRole, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	type msg struct {
		Name        string   `json:"name"`
		Desc        string   `json:"desc"`
		Permissions []string `json:"permissions"`
	}
	req := &msg{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &req)
	if err != nil {
		common.LogFuncDebug("json decode: %s \nfailed : %v", string(c.Ctx.Input.RequestBody), err)
		c.ErrorResponseAndLog(OPActionAddRole, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	if req.Name == "" || req.Desc == "" {
		c.ErrorResponseAndLog(OPActionAddRole, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	//创建角色
	role := &models.Role{
		Name:  req.Name,
		Desc:  req.Desc,
		Ctime: common.NowInt64MS(),
		Utime: common.NowInt64MS(),
	}
	permissions, err := dao.RoleDaoEntity.CreateRole(role, req.Permissions)
	if err != nil {
		c.ErrorResponseAndLog(OPActionAddRole, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	roleInfo := dao.GetBaseRoleInfo(*role)
	roleInfo.Permissions = permissions

	//返回查询结果
	c.SuccessResponseAndLog(OPActionAddRole, string(c.Ctx.Input.RequestBody), roleInfo)
}

//更新角色
func (c *AdminRoleController) HandleUpdateRole() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionEditRole, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	type msg struct {
		Id          int      `json:"id"`
		Name        string   `json:"name"`
		Desc        string   `json:"desc"`
		Permissions []string `json:"permissions"`
	}

	req := &msg{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &req)
	if err != nil {
		common.LogFuncDebug("json decode: %s \nfailed : %v", string(c.Ctx.Input.RequestBody), err)
		c.ErrorResponseAndLog(OPActionEditRole, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	if req.Name == "" || req.Desc == "" {
		c.ErrorResponseAndLog(OPActionEditRole, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	role := &models.Role{
		Id:    uint64(req.Id),
		Name:  req.Name,
		Desc:  req.Desc,
		Utime: common.NowInt64MS(),
	}
	permissions, err := dao.RoleDaoEntity.UpdateRole(role, req.Permissions, models.COLUMN_Role_Name, models.COLUMN_Role_Desc, models.COLUMN_Role_Utime)
	if err != nil {
		c.ErrorResponseAndLog(OPActionEditRole, controllers.ERROR_CODE_DB, string(c.Ctx.Input.RequestBody))
		return
	}
	roleInfo := dao.GetBaseRoleInfo(*role)
	roleInfo.Permissions = permissions

	//返回c查询结果
	c.SuccessResponseAndLog(OPActionEditRole, string(c.Ctx.Input.RequestBody), roleInfo)
}

//删除角色
func (c *AdminRoleController) HandleDeleteRole() {
	c.setOPAction(OPActionDelRole)
	c.setRequestData(string(c.Ctx.Input.RequestBody))
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionDelRole, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	id, err := c.GetInt("id", 0)
	if err != nil || id == 0 {
		c.ErrorResponse(controllers.ERROR_CODE_PARAMS_ERROR)
		return
	}
	input := fmt.Sprintf("{\"id:\":%d}", id)

	err = dao.RoleDaoEntity.DelRole(uint64(id))
	if err != nil {
		c.ErrorResponseAndLog(OPActionDelRole, controllers.ERROR_CODE_DB, input)
		return
	}
	res := map[string]interface{}{}
	res["data"] = map[string]interface{}{
		"id": id,
	}

	//返回c查询结果
	c.SuccessResponseAndLog(OPActionDelRole, input, res)
}

//获取角色权限
func (c *AdminRoleController) HandleGetRolePermissions() {
	c.setOPAction(OPActionReadRolePermissions)
	c.setRequestData(string(c.Ctx.Input.RequestBody))
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionReadRolePermissions, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	id, err := c.GetInt("id", 0)
	if err != nil || id <= 0 {
		c.ErrorResponse(controllers.ERROR_CODE_PARAMS_ERROR)
		return
	}
	input := fmt.Sprintf("{\"id:\":%d}", id)

	//获取数据
	permissions, err := dao.RolePermissionDaoEntity.GetRolePermissions(uint64(id))
	if err != nil {
		c.ErrorResponseAndLog(OPActionReadRolePermissions, controllers.ERROR_CODE_DB, input)
		return
	}

	//返回c查询结果
	res := map[string]interface{}{}
	res["list"] = permissions
	c.SuccessResponseAndLog(OPActionReadRolePermissions, input, res)
}

//添加角色权限
func (c *AdminRoleController) HandleAddRolePermissions() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionEditRolePermissions, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	type msg struct {
		Id          uint64   `json:"id"`
		Permissions []string `json:"permissions"`
	}

	req := &msg{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &req)
	if err != nil || req.Id == 0 {
		common.LogFuncDebug("json decode: %s \nfailed : %v", string(c.Ctx.Input.RequestBody), err)
		c.ErrorResponseAndLog(OPActionEditRolePermissions, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	permissions, err := dao.RolePermissionDaoEntity.AddRolePermissions(req.Id, req.Permissions)
	if err != nil {
		c.ErrorResponseAndLog(OPActionEditRolePermissions, controllers.ERROR_CODE_DB, string(c.Ctx.Input.RequestBody))
		return
	}

	//返回查询结果
	res := map[string]interface{}{}
	res["list"] = permissions
	c.SuccessResponseAndLog(OPActionEditRolePermissions, string(c.Ctx.Input.RequestBody), res)
}

//删除角色权限
func (c *AdminRoleController) HandleDeleteRolePermissions() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionDelRolePermissions, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	type msg struct {
		Id          uint64   `json:"id"`
		Permissions []string `json:"permissions"`
	}
	req := &msg{}

	err := json.Unmarshal(c.Ctx.Input.RequestBody, &req)
	if err != nil || req.Id == 0 {
		common.LogFuncDebug("json decode: %s \nfailed : %v", string(c.Ctx.Input.RequestBody), err)
		c.ErrorResponseAndLog(OPActionDelRolePermissions, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	permissions, err := dao.RolePermissionDaoEntity.DelRolePermissions(req.Id, req.Permissions)
	if err != nil {
		c.ErrorResponseAndLog(OPActionDelRolePermissions, controllers.ERROR_CODE_DB, string(c.Ctx.Input.RequestBody))
		return
	}

	//返回查询结果
	res := map[string]interface{}{}
	res["list"] = permissions
	c.SuccessResponseAndLog(OPActionDelRolePermissions, string(c.Ctx.Input.RequestBody), res)
}

//获取角色所有成员
func (c *AdminRoleController) HandleGetRoleMembers() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionReadRoleMember, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	id, page, perPage, err := c.GetPageInfo()
	if err != nil {
		c.ErrorResponseAndLog(OPActionReadRoleMember, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}
	input := fmt.Sprintf("{\"id\":%d,\"page\":%d,\"per_page\":%d", id, page, perPage)

	members, pageInf, err := dao.AdminUserDaoEntity.QueryPageRoleAdminUser(uint64(id), page, perPage)
	if err != nil {
		c.ErrorResponseAndLog(OPActionReadRoleMember, controllers.ERROR_CODE_DB, input)
		return
	}

	//返回查询结果
	res := map[string]interface{}{}
	res["list"] = members
	res["meta"] = pageInf
	c.SuccessResponseAndLog(OPActionReadRoleMember, input, res)
}
