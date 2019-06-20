package controllers

import (
	"admin/controllers/errcode"
	"common"
	"encoding/json"
	"fmt"
	"time"
	"utils/admin/dao"
	"utils/admin/models"
)

type MenuController struct {
	BaseController
}

type addReqMsg struct {
	Pid        uint64 `json:"pid"` //父菜单id，0代表这个菜单是一级菜单
	Level      int32  `json:"level"`
	Name       string `json:"name"`
	Path       string `json:"path"`
	Icon       string `json:"icon"`
	HideInMenu bool   `json:"hide_in_menu"`
	Component  string `json:"component"`
	OrderId    uint32 `json:"order_id"`
}

//添加接口
func (c *MenuController) Add() {
	//_, errCode := c.CheckPermission()
	//if errCode != controllers.ERROR_CODE_SUCCESS {
	//	c.ErrorResponseAndLog(OPActionDeleteMonthDividendCfg, errCode, string(c.Ctx.Input.RequestBody))
	//	return
	//}
	reqMsg := &addReqMsg{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &reqMsg)
	if err != nil {
		common.LogFuncDebug("json decode: %s \nfailed : %v", string(c.Ctx.Input.RequestBody), err)
		c.ErrorResponseAndLog(OPActionAddMenu, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	menuConf := new(models.MenuConf)
	menuConf.PId = reqMsg.Pid

	menuConf.Level = reqMsg.Level

	menuConf.Name = reqMsg.Name
	menuConf.Path = reqMsg.Path
	menuConf.Icon = reqMsg.Icon
	menuConf.HideInMenu = reqMsg.HideInMenu
	menuConf.Component = reqMsg.Component
	menuConf.OrderId = reqMsg.OrderId
	now := time.Now().Unix()
	menuConf.CTime = now
	menuConf.UTime = now
	err = dao.MenuConfDaoEntity.Add(menuConf)
	if err != nil {
		c.ErrorResponseAndLog(OPActionAddMenu, controllers.ERROR_CODE_DB, string(c.Ctx.Input.RequestBody))
		return
	}
	c.SuccessResponseWithoutData()
}

//删除接口
type delMsg struct {
	Id uint64 `json:"id"`
}

func (c *MenuController) Del() {
	//_, errCode := c.CheckPermission()
	//if errCode != controllers.ERROR_CODE_SUCCESS {
	//	c.ErrorResponseAndLog(OPActionDeleteMonthDividendCfg, errCode, string(c.Ctx.Input.RequestBody))
	//	return
	//}
	id, err := c.GetUint64("id")
	if err != nil {
		common.LogFuncDebug("json decode: %s \nfailed : %v", string(c.Ctx.Input.RequestBody), err)
		c.ErrorResponseAndLog(OPActionDelMenu, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}
	err = dao.MenuConfDaoEntity.Delete(id)
	if err != nil {
		c.ErrorResponseAndLog(OPActionDelMenu, controllers.ERROR_CODE_DB, string(c.Ctx.Input.RequestBody))
		return
	}
	c.SuccessResponseWithoutData()
}

type UpdateMsg struct {
	Id uint64 `json:"id"`
	addReqMsg
}

//修改接口,前端传一个菜单列表过来 传什么改什么
func (c *MenuController) UpdateMenu() {
	//_, errCode := c.CheckPermission()
	//if errCode != controllers.ERROR_CODE_SUCCESS {
	//	c.ErrorResponseAndLog(OPActionDeleteMonthDividendCfg, errCode, string(c.Ctx.Input.RequestBody))
	//	return
	//}
	msg := &UpdateMsg{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &msg)
	if err != nil {
		common.LogFuncDebug("json decode: %s \nfailed : %v", string(c.Ctx.Input.RequestBody), err)
		c.ErrorResponseAndLog(OPActionUpdateMenu, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}
	now := time.Now().Unix()
	menuConf := models.MenuConf{}
	menuConf.Id = msg.Id
	menuConf.OrderId = msg.OrderId
	menuConf.Component = msg.Component
	menuConf.HideInMenu = msg.HideInMenu
	menuConf.Icon = msg.Icon
	menuConf.Path = msg.Path
	menuConf.Level = msg.Level
	menuConf.Name = msg.Name
	menuConf.PId = msg.Pid
	menuConf.UTime = now

	err = dao.MenuConfDaoEntity.Update(menuConf)
	if err != nil {
		c.ErrorResponseAndLog(OPActionUpdateMenu, controllers.ERROR_CODE_DB, string(c.Ctx.Input.RequestBody))
		return
	}
	c.SuccessResponseWithoutData()
}

var oneLevel int32 = 1
var firstOrderId = 1

//获取菜单列表
func (c *MenuController) GetAllMenu() {
	//_, errCode := c.CheckPermission()
	//if errCode != controllers.ERROR_CODE_SUCCESS {
	//	c.ErrorResponseAndLog(OPActionDeleteMonthDividendCfg, errCode, string(c.Ctx.Input.RequestBody))
	//	return
	//}
	menus, err := dao.MenuConfDaoEntity.FindAll()
	if err != nil {
		c.ErrorResponseAndLog(OPActionGetAllMenu, controllers.ERROR_CODE_DB, string(c.Ctx.Input.RequestBody))
		return
	}

	oneLevelMenus := make([]*models.MenuConf, 0)
	towLevelMenuMap := make(map[uint64][]*models.MenuConf, 0)
	for _, menu := range menus {
		if menu.Level == oneLevel {
			oneLevelMenus = append(oneLevelMenus, menu)
			continue
		}
		if _, ok := towLevelMenuMap[menu.PId]; !ok {
			towLevelMenuMap[menu.PId] = make([]*models.MenuConf, 0)
		}
		towLevelMenuMap[menu.PId] = append(towLevelMenuMap[menu.PId], menu)
	}
	resultMenu := make([]*models.MenuConf, 0, len(menus))
	for _, oneLevelMenu := range oneLevelMenus {
		resultMenu = append(resultMenu, oneLevelMenu)
		if _, ok := towLevelMenuMap[oneLevelMenu.Id]; !ok {
			continue
		}
		childrenMenu := towLevelMenuMap[oneLevelMenu.Id]
		resultMenu = append(resultMenu, childrenMenu...)
	}

	res := map[string]interface{}{
		KEY_LIST: resultMenu,
	}
	res[KEY_META] = map[string]interface{}{
		KEY_LIMIT: 100,
		KEY_PAGE:  1,
		"total":   len(resultMenu),
	}
	c.SuccessResponse(res)
}

//查询接口,按照约定的格式返回所有菜单
type MenuMsg struct {
	AdminId uint64 `json:"admin_id"`
}
type oneLevelMenuMsg struct {
	*towLevelMenuMsg
	Routes []*towLevelMenuMsg `json:"routes"`
}
type towLevelMenuMsg struct {
	Id         uint64 `json:"-"`
	Pid        uint64 `json:"-"`
	HideInMenu bool   `json:"hideInMenu"`
	Path       string `json:"path"`
	Name       string `json:"name"`
	Component  string `json:"component,omitempty"`
	Icon       string `json:"icon,omitempty"`
}

func (c *MenuController) GetAllAccessMenus() {
	//_, errCode := c.CheckPermission()
	//if errCode != controllers.ERROR_CODE_SUCCESS {
	//	c.ErrorResponseAndLog(OPActionDeleteMonthDividendCfg, errCode, string(c.Ctx.Input.RequestBody))
	//	return
	//}

	routes := make([]*oneLevelMenuMsg, 0)
	res := map[string]interface{}{
		"routes": routes,
	}

	menus, err := dao.MenuConfDaoEntity.FindAll()
	if err != nil {
		c.ErrorResponseAndLog(OPActionGetAccessMenus, controllers.ERROR_CODE_DB, string(c.Ctx.Input.RequestBody))
		return
	}

	oneLevelMenuMap := make(map[uint64]*oneLevelMenuMsg, 0)
	towLevelMenuMap := make(map[uint64]*towLevelMenuMsg, 0)
	towLevelMenus := make([]*towLevelMenuMsg, 0)
	for _, menu := range menus {
		if menu.Level == 1 {
			msg := new(oneLevelMenuMsg)
			msg.towLevelMenuMsg = InitMenuMsg(menu)
			msg.Routes = make([]*towLevelMenuMsg, 0)

			oneLevelMenuMap[menu.Id] = msg
			routes = append(routes, msg)
		} else if menu.Level == 2 {
			msg := InitMenuMsg(menu)
			towLevelMenuMap[menu.Id] = msg
			towLevelMenus = append(towLevelMenus, msg)
		}
	}
	for _, towLevelMenu := range towLevelMenus {
		if _, ok := oneLevelMenuMap[towLevelMenu.Pid]; !ok {
			continue
		}
		oneLevelMenuMap[towLevelMenu.Pid].Routes = append(oneLevelMenuMap[towLevelMenu.Pid].Routes, towLevelMenu)
	}
	res["routes"] = routes

	c.SuccessResponse(res)
}

//获取角色Id可见的所有菜单
//要传admin_id进来，然后根据admin_id去role_admin表里查到对应的role_id,再根据role_id去查询有权限的菜单
func (c *MenuController) GetAccessMenus() {
	//_, errCode := c.CheckPermission()
	//if errCode != controllers.ERROR_CODE_SUCCESS {
	//	c.ErrorResponseAndLog(OPActionDeleteMonthDividendCfg, errCode, string(c.Ctx.Input.RequestBody))
	//	return
	//}

	adminId, err := c.GetUint64("admin_id")
	if err != nil {
		common.LogFuncDebug("json decode: %s \nfailed : %v", string(c.Ctx.Input.RequestBody), err)
		c.ErrorResponseAndLog(OPActionGetAccessMenus, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}
	routes := make([]*oneLevelMenuMsg, 0)
	res := map[string]interface{}{
		"routes": routes,
	}
	//查找这个role_id可以获取的所有菜单数据
	roleIds, err := dao.RoleAdminDaoEntity.QueryRoleIdByAdminId(adminId)
	if err != nil {
		c.ErrorResponseAndLog(OPActionGetAccessMenus, controllers.ERROR_CODE_DB, string(c.Ctx.Input.RequestBody))
		return
	}
	menuAccesses, err := dao.MenuAccessDaoEntity.FindByRoleId(roleIds)
	if len(menuAccesses) == 0 {
		c.SuccessResponse(res)
		return
	}

	//去除重复的菜单id
	menuIdMap := make(map[uint64]uint64, 0)
	menuIds := make([]uint64, 0, len(menuAccesses))
	towLevelMenus := make([]*towLevelMenuMsg, 0)
	for _, menuAccess := range menuAccesses {
		if _, ok := menuIdMap[menuAccess.MenuId]; ok {
			continue
		}
		menuIdMap[menuAccess.MenuId] = menuAccess.MenuId
		menuIds = append(menuIds, menuAccess.MenuId)
	}

	menus, err := dao.MenuConfDaoEntity.FindByMenuIds(menuIds)
	if err != nil {
		c.ErrorResponseAndLog(OPActionGetAccessMenus, controllers.ERROR_CODE_DB, string(c.Ctx.Input.RequestBody))
		return
	}

	oneLevelMenuMap := make(map[uint64]*oneLevelMenuMsg, 0)
	towLevelMenuMap := make(map[uint64]*towLevelMenuMsg, 0)
	for _, menu := range menus {
		if menu.Level == 1 {
			msg := new(oneLevelMenuMsg)
			msg.towLevelMenuMsg = InitMenuMsg(menu)
			msg.Routes = make([]*towLevelMenuMsg, 0)

			oneLevelMenuMap[menu.Id] = msg

			routes = append(routes, msg)
		} else if menu.Level == 2 {
			msg := InitMenuMsg(menu)
			towLevelMenuMap[menu.Id] = msg
			towLevelMenus = append(towLevelMenus, msg)
		}
	}
	for _, towLevelMenu := range towLevelMenus {
		if _, ok := oneLevelMenuMap[towLevelMenu.Pid]; !ok {
			continue
		}
		oneLevelMenuMap[towLevelMenu.Pid].Routes = append(oneLevelMenuMap[towLevelMenu.Pid].Routes, towLevelMenu)
	}

	res["routes"] = routes

	c.SuccessResponse(res)
}
func InitMenuMsg(menu *models.MenuConf) (msg *towLevelMenuMsg) {
	msg = new(towLevelMenuMsg)
	msg.Id = menu.Id
	msg.Pid = menu.PId
	msg.HideInMenu = menu.HideInMenu
	msg.Path = menu.Path
	msg.Name = menu.Name
	msg.Component = menu.Component
	msg.Icon = menu.Icon
	return
}

type MenuAccessMsg struct {
	RoleId  uint64   `json:"role_id"`
	MenusId []uint64 `json:"menus_id"`
}

//菜单权限管理
func (c *MenuController) UpdateMenuAccess() {
	//_, errCode := c.CheckPermission()
	//if errCode != controllers.ERROR_CODE_SUCCESS {
	//	c.ErrorResponseAndLog(OPActionDeleteMonthDividendCfg, errCode, string(c.Ctx.Input.RequestBody))
	//	return
	//}
	msg := &MenuAccessMsg{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &msg)
	if err != nil {
		common.LogFuncDebug("json decode: %s \nfailed : %v", string(c.Ctx.Input.RequestBody), err)
		c.ErrorResponseAndLog(OPActionUpdateMenuAccess, controllers.ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	err = dao.MenuAccessDaoEntity.Delete(msg.RoleId)
	if err != nil {
		c.ErrorResponseAndLog(OPActionUpdateMenuAccess, controllers.ERROR_CODE_DB, string(c.Ctx.Input.RequestBody))
		return
	}
	menuAccess := make([]*models.MenuAccess, 0)
	for _, menuId := range msg.MenusId {
		data := new(models.MenuAccess)
		data.RoleId = msg.RoleId
		data.MenuId = menuId
		menuAccess = append(menuAccess, data)
	}
	fmt.Println("menuAccess ", menuAccess)
	err = dao.MenuAccessDaoEntity.Insert(menuAccess)
	if err != nil {
		c.ErrorResponseAndLog(OPActionUpdateMenuAccess, controllers.ERROR_CODE_DB, string(c.Ctx.Input.RequestBody))
		return
	}
	c.SuccessResponseWithoutData()
}
