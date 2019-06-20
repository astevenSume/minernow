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

type AdminUserController struct {
	BaseController
}

func (c *AdminUserController) GetPageInfo() (uint64, int, int, error) {
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

//获取管理员
func (c *AdminUserController) HandleGetAdmins() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionReadAdmin, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	id, page, perPage, err := c.GetPageInfo()
	if err != nil {
		common.LogFuncInfo("err:%v", err)
		c.ErrorResponseAndLog(OPActionReadAdmin, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	roleId, err := c.GetInt("role_id", 0)
	if err != nil {
		c.ErrorResponseAndLog(OPActionReadAdmin, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}
	status, err := c.GetInt(KEY_STATUS, -1)
	if err != nil {
		common.LogFuncError("param err:%v", err)
		c.ErrorResponseAndLog(OPActionReadAdmin, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	res := map[string]interface{}{}
	var input string
	//获取数据
	if id == 0 {
		input = fmt.Sprintf("{\"role_id\":%v,\"status\":%v,\"page\":%v,\"limit\":%v}", roleId, status, page, perPage)
		members, pageInf, err := dao.AdminUserDaoEntity.QueryPageAdminUser(uint64(roleId), page, perPage, int8(status))
		if err != nil {
			c.ErrorResponseAndLog(OPActionReadAdmin, controllers.ERROR_CODE_DB, input)
			return
		}
		res["meta"] = pageInf
		res["list"] = members
	} else {
		input = fmt.Sprintf("{\"id\":%d}", id)
		q, err := dao.GetAdminInfo(uint64(id))
		if err != nil {
			c.ErrorResponseAndLog(OPActionReadAdmin, controllers.ERROR_CODE_DB, input)
			return
		}
		c.SuccessResponseAndLog(OPActionReadAdmin, input, dao.TransAdminMember(q))
		return
	}

	c.SuccessResponseAndLog(OPActionReadAdmin, input, res)
}

//创建管理员
func (c *AdminUserController) HandleCreateAdmins() {
	adminUser, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionAddAdmin, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	type msg struct {
		Name         string   `json:"name"`
		Email        string   `json:"email"`
		RoleIds      []uint64 `json:"role_ids"`
		WhitelistIps string   `json:"whitelist_ips"`
	}
	req := &msg{}

	err := json.Unmarshal(c.Ctx.Input.RequestBody, &req)
	if err != nil {
		common.LogFuncDebug("json decode: %s \nfailed : %v", string(c.Ctx.Input.RequestBody), err)
		c.ErrorResponseAndLog(OPActionAddAdmin, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	if req.Name == "" || req.Email == "" || len(req.RoleIds) == 0 {
		c.ErrorResponseAndLog(OPActionAddAdmin, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	newAdminUser := &models.AdminUser{
		Email:        req.Email,
		Name:         req.Name,
		WhitelistIps: req.WhitelistIps,
	}

	member, err := dao.AdminUserDaoEntity.CreateAdminUser(newAdminUser, req.RoleIds, adminUser.AdminUser.Name)
	if err != nil {
		c.ErrorResponseAndLog(OPActionAddAdmin, controllers.ERROR_CODE_DB, string(c.Ctx.Input.RequestBody))
		return
	}

	//返回查询结果
	c.SuccessResponseAndLog(OPActionAddAdmin, string(c.Ctx.Input.RequestBody), member)
}

//更新管理员
func (c *AdminUserController) HandleUpdateAdmins() {
	adminUser, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionEditAdmin, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	type msg struct {
		Id           uint64   `json:"id"`
		Name         string   `json:"name"`
		Email        string   `json:"email"`
		RoleIds      []uint64 `json:"role_ids"`
		Status       int8     `json:"status"`
		WhitelistIps string   `json:"whitelist_ips"`
	}
	req := &msg{}

	err := json.Unmarshal(c.Ctx.Input.RequestBody, &req)
	if err != nil {
		common.LogFuncDebug("json decode: %s \nfailed : %v", string(c.Ctx.Input.RequestBody), err)
		c.ErrorResponseAndLog(OPActionEditAdmin, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	updateUser := &models.AdminUser{
		Id:           req.Id,
		Name:         req.Name,
		Email:        req.Email,
		WhitelistIps: req.WhitelistIps,
		Status:       req.Status,
	}

	member, err := dao.AdminUserDaoEntity.UpdateAdminUser(updateUser, req.RoleIds, adminUser.AdminUser.Name)
	if err != nil {
		c.ErrorResponseAndLog(OPActionEditAdmin, controllers.ERROR_CODE_DB, string(c.Ctx.Input.RequestBody))
		return
	}

	//返回查询结果
	c.SuccessResponseAndLog(OPActionEditAdmin, string(c.Ctx.Input.RequestBody), member)
}

//删除管理员
func (c *AdminUserController) HandleDelAdmins() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionDelAdmin, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	id, err := c.GetInt("id", 0)
	if err != nil || id == 0 {
		c.ErrorResponseAndLog(OPActionDelAdmin, controllers.ERROR_CODE_PARAMS_ERROR, "")
		return
	}

	delAdminUser := &models.AdminUser{
		Id:     uint64(id),
		Status: dao.AdminUserStatusDeleted,
		Dtime:  common.NowInt64MS(),
	}

	err = dao.AdminUserDaoEntity.UpdateAdminUserStatus(uint64(id), delAdminUser, models.COLUMN_AdminUser_Status, models.COLUMN_AdminUser_Dtime)
	if err != nil {
		c.ErrorResponseAndLog(OPActionDelAdmin, controllers.ERROR_CODE_DB, string(c.Ctx.Input.RequestBody))
		return
	}

	//返回查询结果
	data := map[string]interface{}{
		"id": id,
	}
	c.SuccessResponseAndLog(OPActionDelAdmin, string(c.Ctx.Input.RequestBody), data)
}

//禁用管理员
func (c *AdminUserController) HandleSuspendAdmins() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionAddSuspendAdmin, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	id, err := c.GetInt("id", 0)
	if err != nil || id == 0 {
		c.ErrorResponseAndLog(OPActionAddSuspendAdmin, controllers.ERROR_CODE_PARAMS_ERROR, "")
		return
	}
	input := fmt.Sprintf("{\"id\":%d}", id)

	suspendAdminUser := &models.AdminUser{
		Id:     uint64(id),
		Status: dao.AdminUserStatusSuspended,
		Utime:  common.NowInt64MS(),
	}
	err = dao.AdminUserDaoEntity.UpdateAdminUserStatus(uint64(id), suspendAdminUser, models.COLUMN_AdminUser_Status, models.COLUMN_AdminUser_Utime)
	if err != nil {
		c.ErrorResponseAndLog(OPActionAddSuspendAdmin, controllers.ERROR_CODE_DB, input)
		return
	}

	//返回查询结果
	data := map[string]interface{}{
		"id": id,
	}
	c.SuccessResponseAndLog(OPActionAddSuspendAdmin, input, data)
}

//恢复管理员
func (c *AdminUserController) HandleRestoreAdmins() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionAddRestoreAdmin, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	id, err := c.GetInt("id", 0)
	if err != nil || id == 0 {
		c.ErrorResponseAndLog(OPActionAddRestoreAdmin, controllers.ERROR_CODE_PARAMS_ERROR, "")
		return
	}
	input := fmt.Sprintf("{\"id\":%d}", id)

	rsAdminUser := &models.AdminUser{
		Id:     uint64(id),
		Status: dao.AdminUserStatusActive,
		Utime:  common.NowInt64MS(),
	}
	err = dao.AdminUserDaoEntity.UpdateAdminUserStatus(uint64(id), rsAdminUser, models.COLUMN_AdminUser_Status, models.COLUMN_AdminUser_Utime)
	if err != nil {
		c.ErrorResponseAndLog(OPActionAddRestoreAdmin, controllers.ERROR_CODE_DB, input)
		return
	}

	//返回查询结果
	data := map[string]interface{}{
		"id": id,
	}
	c.SuccessResponseAndLog(OPActionAddRestoreAdmin, input, data)
}

//获取管理员角色
func (c *AdminUserController) HandleGetAdminsRoles() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionReadAdminsRole, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	id, page, perPage, err := c.GetPageInfo()
	if err != nil || id == 0 {
		c.ErrorResponseAndLog(OPActionReadAdminsRole, controllers.ERROR_CODE_PARAMS_ERROR, "")
		return
	}
	input := fmt.Sprintf("{\"id\":%d,\"page\":%d,\"limit\":%d", id, page, perPage)

	roles, pageInf, err := dao.AdminUserDaoEntity.QueryPageAdminRole(uint64(id), page, perPage)
	if err != nil {
		c.ErrorResponseAndLog(OPActionReadAdminsRole, controllers.ERROR_CODE_DB, input)
		return
	}

	//返回c查询结果
	res := map[string]interface{}{}
	res["list"] = roles
	res["meta"] = pageInf
	c.SuccessResponseAndLog(OPActionReadAdminsRole, input, res)
}

//添加管理员角色
func (c *AdminUserController) HandleAddAdminsRoles() {
	adminUser, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionAddAdminsRole, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	type msg struct {
		Id     uint64 `json:"id"`
		RoleId uint64 `json:"role_id"`
	}
	req := &msg{}

	err := json.Unmarshal(c.Ctx.Input.RequestBody, &req)
	if err != nil || req.Id == 0 {
		common.LogFuncDebug("json decode: %s \nfailed : %v", string(c.Ctx.Input.RequestBody), err)
		c.ErrorResponseAndLog(OPActionAddAdminsRole, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	addAdminUser := &models.AdminUser{Id: uint64(req.Id)}
	err = dao.AdminUserDaoEntity.AddAdminUserRoles(addAdminUser, req.RoleId, adminUser.AdminUser.Name, models.COLUMN_AdminUser_Id)
	if err != nil {
		c.ErrorResponseAndLog(OPActionAddAdminsRole, controllers.ERROR_CODE_DB, string(c.Ctx.Input.RequestBody))
		return
	}

	//返回查询结果
	data := map[string]interface{}{
		"role_id":    req.RoleId,
		"admin_id":   req.Id,
		"granted_by": adminUser.AdminUser.Name,
		"granted_at": common.NowInt64MS(),
	}
	c.SuccessResponseAndLog(OPActionAddAdminsRole, string(c.Ctx.Input.RequestBody), data)
}

//删除管理员角色
func (c *AdminUserController) HandleDelAdminsRoles() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionDelAdminsRole, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	type msg struct {
		Id     uint64 `json:"id"`
		RoleId int    `json:"role_id"`
	}
	req := &msg{}

	err := json.Unmarshal(c.Ctx.Input.RequestBody, &req)
	if err != nil || req.Id == 0 {
		common.LogFuncDebug("json decode: %s \nfailed : %v", string(c.Ctx.Input.RequestBody), err)
		c.ErrorResponseAndLog(OPActionDelAdminsRole, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	err = dao.RoleAdminDaoEntity.DelRolesAdmin(uint64(req.RoleId), uint64(req.Id))
	if err != nil {
		c.ErrorResponseAndLog(OPActionDelAdminsRole, controllers.ERROR_CODE_DB, string(c.Ctx.Input.RequestBody))
		return
	}

	//返回查询结果
	data := map[string]interface{}{
		"role_id":  req.RoleId,
		"admin_id": req.Id,
	}
	c.SuccessResponseAndLog(OPActionDelAdminsRole, string(c.Ctx.Input.RequestBody), data)
}

//解除管理员google验证
func (c *AdminUserController) UnBind() {
	_, errCode := c.CheckPermission()
	if errCode != controllers.ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionAddAdminsUnBind, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	id, err := c.GetUint64("id", 0)
	if err != nil || id == 0 {
		c.ErrorResponseAndLog(OPActionAddAdminsUnBind, controllers.ERROR_CODE_PARAMS_ERROR, "")
		return
	}
	input := fmt.Sprintf("{\"id\":%d}", id)

	err = dao.AdminUserDaoEntity.UpdateGoogleAuthBind(id, false)
	if err != nil {
		c.ErrorResponseAndLog(OPActionAddAdminsUnBind, controllers.ERROR_CODE_DB, input)
		return
	}

	c.SuccessResponseWithoutDataAndLog(OPActionAddAdminsUnBind, input)
}
