package controllers

import (
	. "admin/controllers/errcode"
	"common"
	"common/systoolbox"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/orm"
	"github.com/mailru/easyjson"
	admindao "utils/admin/dao"
	adminmodels "utils/admin/models"
	utilscommon "utils/common"
)

type TaskController struct {
	BaseController
}

const (
	adminSvr = "admin"
)

// Query query task
func (c *TaskController) Query() {
	//_, errCode := c.CheckPermission()
	//if errCode != ERROR_CODE_SUCCESS {
	//	c.ErrorResponseAndLog(OPActionTaskRead, errCode, string(c.Ctx.Input.RequestBody))
	//	return
	//}

	type Msg struct {
		AppName string `json:"app_name"`
		Page    int    `json:"page"`
		PerPage int    `json:"per_page"`
	}
	msg := Msg{}

	msg.AppName = c.GetString(KeyAppName)

	msg.Page, _ = c.GetInt(KEY_PAGE)
	msg.PerPage, _ = c.GetInt(KEY_LIMIT)
	if msg.Page == 0 {
		msg.Page = 1
	}
	if msg.PerPage == 0 {
		msg.PerPage = DEFAULT_PER_PAGE
	}

	meta, list, err := admindao.TaskDaoEntity.Query(msg.AppName, msg.Page, msg.PerPage)
	if err != nil {
		c.ErrorResponseAndLog(OPActionTaskRead, ERROR_CODE_DB, string(c.Ctx.Input.RequestBody))
		return
	}

	res := map[string]interface{}{}
	res["list"] = list
	res["meta"] = meta

	c.SuccessResponseAndLog(OPActionTaskRead, fmt.Sprint(msg), res)
}

func (c *TaskController) QuerySingle() {
	//_, errCode := c.CheckPermission()
	//if errCode != ERROR_CODE_SUCCESS {
	//	c.ErrorResponseAndLog(OPActionTaskRead, errCode, string(c.Ctx.Input.RequestBody))
	//	return
	//}

	appName := c.GetString(KeyAppName)
	if len(appName) <= 0 {
		c.ErrorResponseAndLog(OPActionTaskRead, ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	var (
		regionId, serverId int64
		err                error
	)
	regionId, err = c.GetInt64(KeyRegionId)
	if err != nil {
		c.ErrorResponseAndLog(OPActionTaskRead, ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	serverId, err = c.GetInt64(KeyServerId)
	if err != nil {
		c.ErrorResponseAndLog(OPActionTaskRead, ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	detail := systoolbox.TaskDetail{}
	if appName == adminSvr { //get from local
		detail = systoolbox.TaskMgr.Detail()
	} else { //get from remote
		// send rpc request
		resp, err := common.RabbitMQRpcSend(appName, regionId, serverId, &common.RabbitMQRpcMsg{
			Cmd: utilscommon.RabbitMQRpcMsgCmdTaskQuerySingle,
		}, 10000)
		if err != nil {
			common.LogFuncError("%v", err)
			c.ErrorResponseAndLog(OPActionTaskRead, ERROR_CODE_RABBIT_MQ_PUBLISH_FAIL, string(c.Ctx.Input.RequestBody))
			return
		}

		if ERROR_CODE(resp.Code) != ERROR_CODE_SUCCESS {
			common.LogFuncError("rpc response code %d", resp.Code)
			c.ErrorResponseAndLog(OPActionTaskRead, ERROR_CODE_RABBIT_MQ_PUBLISH_FAIL, string(c.Ctx.Input.RequestBody))
			return
		}

		err = easyjson.Unmarshal(resp.Body, &detail)
		if err != nil {
			common.LogFuncError("%v", err)
			c.ErrorResponseAndLog(OPActionTaskRead, ERROR_CODE_DECODE_FAIL, string(c.Ctx.Input.RequestBody))
			return
		}
	}

	c.SuccessResponseAndLog(OPActionTaskRead, string(c.Ctx.Input.RequestBody), detail)
}

func (c *TaskController) QueryResult() {
	//_, errCode := c.CheckPermission()
	//if errCode != ERROR_CODE_SUCCESS {
	//	c.ErrorResponseAndLog(OPActionTaskRead, errCode, string(c.Ctx.Input.RequestBody))
	//	return
	//}

	var (
		appName            string
		regionId, serverId int64
		page, perPage      int
		err                error
	)

	appName = c.GetString(KeyAppName)

	regionId, _ = c.GetInt64(KeyRegionId, admindao.RegionIdImpossible)
	serverId, _ = c.GetInt64(KeyServerId, admindao.ServerIdImpossible)

	page, _ = c.GetInt(KEY_PAGE)
	perPage, _ = c.GetInt(KEY_LIMIT)
	if page == 0 {
		page = 1
	}
	if perPage == 0 {
		perPage = DEFAULT_PER_PAGE
	}

	meta, list, err := admindao.TaskResultDaoEntity.Query(appName, regionId, serverId, page, perPage)
	if err != nil {
		c.ErrorResponseAndLog(OPActionTaskRead, ERROR_CODE_DB, string(c.Ctx.Input.RequestBody))
		return
	}

	res := map[string]interface{}{}
	res["list"] = list
	res["meta"] = meta

	c.SuccessResponseAndLog(OPActionTaskRead, string(c.Ctx.Input.RequestBody), res)
}

// Add add task
func (c *TaskController) Add() {
	//_, errCode := c.CheckPermission()
	//if errCode != ERROR_CODE_SUCCESS {
	//	c.ErrorResponseAndLog(OPActionTaskWrite, errCode, string(c.Ctx.Input.RequestBody))
	//	return
	//}

	type Msg struct {
		Name     string `json:"name"`
		AppName  string `json:"app_name"`
		Alia     string `json:"alia"`
		Spec     string `json:"spec"`
		FuncName string `json:"func_name"`
		Desc     string `json:"desc"`
		Status   uint8  `json:"status"`
	}
	msg := Msg{}

	err := json.Unmarshal(c.Ctx.Input.RequestBody, &msg)
	if err != nil {
		common.LogFuncDebug("json decode: %s \nfailed : %v", string(c.Ctx.Input.RequestBody), err)
		c.ErrorResponseAndLog(OPActionTaskWrite, ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	t := &adminmodels.Task{
		Name:     msg.Name,
		Alia:     msg.Alia,
		AppName:  msg.AppName,
		Spec:     msg.Spec,
		FuncName: msg.FuncName,
		Desc:     msg.Desc,
		Status:   msg.Status,
	}

	err = admindao.TaskDaoEntity.Insert(t)
	if err != nil {
		c.ErrorResponseAndLog(OPActionTaskWrite, ERROR_CODE_DB, fmt.Sprint(msg))
		return
	}

	c.SuccessResponseAndLog(OPActionTaskWrite, fmt.Sprint(msg), t)
}

// Edit edit task
func (c *TaskController) Edit() {
	//_, errCode := c.CheckPermission()
	//if errCode != ERROR_CODE_SUCCESS {
	//	c.ErrorResponseAndLog(OPActionTaskWrite, errCode, string(c.Ctx.Input.RequestBody))
	//	return
	//}

	type Msg struct {
		Id       uint32 `json:"id"`
		Name     string `json:"name"`
		AppName  string `json:"app_name"`
		Alia     string `json:"alia"`
		Spec     string `json:"spec"`
		FuncName string `json:"func_name"`
		Desc     string `json:"desc"`
		Status   uint8  `json:"status"`
	}
	msg := Msg{}

	err := json.Unmarshal(c.Ctx.Input.RequestBody, &msg)
	if err != nil {
		common.LogFuncDebug("json decode: %s \nfailed : %v", string(c.Ctx.Input.RequestBody), err)
		c.ErrorResponseAndLog(OPActionTaskWrite, ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	t := &adminmodels.Task{
		Id:       msg.Id,
		Name:     msg.Name,
		Alia:     msg.Alia,
		AppName:  msg.AppName,
		Spec:     msg.Spec,
		FuncName: msg.FuncName,
		Desc:     msg.Desc,
		Status:   msg.Status,
	}
	err = admindao.TaskDaoEntity.Update(t)
	if err != nil {
		c.ErrorResponseAndLog(OPActionTaskWrite, ERROR_CODE_DB, fmt.Sprint(msg))
		return
	}

	c.SuccessResponseAndLog(OPActionTaskWrite, fmt.Sprint(msg), t)
}

// Edit delete task
func (c *TaskController) Delete() {
	//_, errCode := c.CheckPermission()
	//if errCode != ERROR_CODE_SUCCESS {
	//	c.ErrorResponseAndLog(OPActionTaskDelete, errCode, string(c.Ctx.Input.RequestBody))
	//	return
	//}

	type Msg struct {
		Id uint32 `json:"id"`
	}
	msg := Msg{}

	err := json.Unmarshal(c.Ctx.Input.RequestBody, &msg)
	if err != nil {
		common.LogFuncDebug("json decode: %s \nfailed : %v", string(c.Ctx.Input.RequestBody), err)
		c.ErrorResponseAndLog(OPActionTaskDelete, ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	//// query task
	//var task adminmodels.Task
	//task, err = admindao.TaskDaoEntity.QueryById(msg.Id)
	//if err != nil {
	//	c.ErrorResponseAndLog(OPActionTaskDelete, ERROR_CODE_DB, string(c.Ctx.Input.RequestBody))
	//	return
	//}

	// tasks on server nodes will delete by the other way.

	//send to servers
	//taskMsgList := systoolbox.TaskMsgList{
	//	List: []systoolbox.TaskMsg{
	//		{
	//			Name: task.Name,
	//		},
	//	},
	//}

	//var buf []byte
	//buf, err = easyjson.Marshal(&taskMsgList)
	//if err != nil {
	//	common.LogFuncError("%v", err)
	//	return
	//}

	//err = common.RabbiteMQPublish(taskBusinessName(task.AppName),
	//	common.RabbitMQRoutingKeyTaskDelete,
	//	buf,
	//)
	//if err != nil {
	//	c.ErrorResponseAndLog(OPActionTaskDelete, ERROR_CODE_RABBIT_MQ_PUBLISH_FAIL, fmt.Sprint(msg))
	//	return
	//}

	// delete task data
	err = admindao.TaskDaoEntity.Remove(msg.Id)
	if err != nil {
		c.ErrorResponseAndLog(OPActionTaskDelete, ERROR_CODE_DB, fmt.Sprint(msg))
		return
	}

	c.SuccessResponseAndLog(OPActionTaskDelete, fmt.Sprint(msg), msg)
}

// Distribute distribute task to servers
func (c *TaskController) Distribute() {
	//_, errCode := c.CheckPermission()
	//if errCode != ERROR_CODE_SUCCESS {
	//	c.ErrorResponseAndLog(OPActionTaskWrite, errCode, string(c.Ctx.Input.RequestBody))
	//	return
	//}

	type Msg struct {
		Id uint32 `json:"id"`
	}

	msg := Msg{}

	err := json.Unmarshal(c.Ctx.Input.RequestBody, &msg)
	if err != nil {
		common.LogFuncError("json decode: %s \nfailed : %v", string(c.Ctx.Input.RequestBody), err)
		c.ErrorResponseAndLog(OPActionTaskWrite, ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	taskMsgList, appName, errCode := c.getTaskMsgListById(msg.Id)
	if errCode != ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionTaskWrite, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	if appName == adminSvr {
		if taskList, ok := systoolbox.CheckTaskFunc(taskMsgList, true, FuncContainer, taskResultStoreLocal); ok {
			if len(taskList) > 0 {
				systoolbox.TaskMgr.AddTaskList(taskList)
			}
		}
	} else {
		buf, appName, errCode := c.getTaskMsgListBuf(taskMsgList)
		if errCode != ERROR_CODE_SUCCESS {
			c.ErrorResponseAndLog(OPActionTaskWrite, errCode, fmt.Sprint(msg))
			return
		}

		err = common.RabbitMQPublish(taskBusinessName(appName),
			common.RabbitMQRoutingKeyTaskDistribute,
			buf)
		if err != nil {
			c.ErrorResponseAndLog(OPActionTaskWrite, ERROR_CODE_RABBIT_MQ_PUBLISH_FAIL, fmt.Sprint(msg))
			return
		}
	}
	c.SuccessResponseAndLog(OPActionTaskWrite, fmt.Sprint(msg), msg)
}

// Distribute distribute task to single server
func (c *TaskController) DistributeSingle() {
	//_, errCode := c.CheckPermission()
	//if errCode != ERROR_CODE_SUCCESS {
	//	c.ErrorResponseAndLog(OPActionTaskWrite, errCode, string(c.Ctx.Input.RequestBody))
	//	return
	//}

	type Msg struct {
		Id       uint32 `json:"id"`
		RegionId int64  `json:"region_id"`
		ServerId int64  `json:"server_id"`
	}

	msg := Msg{}

	err := json.Unmarshal(c.Ctx.Input.RequestBody, &msg)
	if err != nil {
		common.LogFuncDebug("json decode: %s \nfailed : %v", string(c.Ctx.Input.RequestBody), err)
		c.ErrorResponseAndLog(OPActionTaskWrite, ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	taskMsgList, appName, errCode := c.getTaskMsgListById(msg.Id)
	if errCode != ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionTaskWrite, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	if appName == adminSvr {
		if taskList, ok := systoolbox.CheckTaskFunc(taskMsgList, true, FuncContainer, taskResultStoreLocal); ok {
			if len(taskList) > 0 {
				systoolbox.TaskMgr.AddTaskList(taskList)
			}
		}
	} else {
		buf, appName, errCode := c.getTaskMsgListBuf(taskMsgList)
		if errCode != ERROR_CODE_SUCCESS {
			c.ErrorResponseAndLog(OPActionTaskWrite, errCode, string(c.Ctx.Input.RequestBody))
			return
		}

		err = common.RabbitMQPublish(svrTaskBusinessName(appName, msg.RegionId, msg.ServerId),
			common.RabbitMQRoutingKeyTaskDistribute,
			buf)
		if err != nil {
			c.ErrorResponseAndLog(OPActionTaskWrite, ERROR_CODE_RABBIT_MQ_PUBLISH_FAIL, fmt.Sprint(msg))
			return
		}
	}
	c.SuccessResponseAndLog(OPActionTaskWrite, fmt.Sprint(msg), msg)
}

// Switch switch task switch of servers
func (c *TaskController) Switch() {
	//_, errCode := c.CheckPermission()
	//if errCode != ERROR_CODE_SUCCESS {
	//	c.ErrorResponseAndLog(OPActionTaskWrite, errCode, string(c.Ctx.Input.RequestBody))
	//	return
	//}

	res, errCode := c.trySwitch()
	if errCode != ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionTaskWrite, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	c.SuccessResponseAndLog(OPActionTaskWrite, string(c.Ctx.Input.RequestBody), res)
}

// SwitchSingle switch task switch of specific server
func (c *TaskController) SwitchSingle() {
	//_, errCode := c.CheckPermission()
	//if errCode != ERROR_CODE_SUCCESS {
	//	c.ErrorResponseAndLog(OPActionTaskWrite, errCode, string(c.Ctx.Input.RequestBody))
	//	return
	//}

	res, errCode := c.trySwitchSingle()
	if errCode != ERROR_CODE_SUCCESS {
		c.ErrorResponseAndLog(OPActionTaskWrite, errCode, string(c.Ctx.Input.RequestBody))
		return
	}

	c.SuccessResponseAndLog(OPActionTaskWrite, string(c.Ctx.Input.RequestBody), res)
}

// RunSingle run single server's task
func (c *TaskController) Run() {
	//_, errCode := c.CheckPermission()
	//if errCode != ERROR_CODE_SUCCESS {
	//	c.ErrorResponseAndLog(OPActionTaskWrite, errCode, string(c.Ctx.Input.RequestBody))
	//	return
	//}

	type Msg struct {
		Name     string `json:"name"`
		AppName  string `json:"app_name"`
		RegionId int64  `json:"region_id"`
		ServerId int64  `json:"server_id"`
	}

	msg := Msg{}

	err := json.Unmarshal(c.Ctx.Input.RequestBody, &msg)
	if err != nil {
		common.LogFuncError("json decode: %s \nfailed : %v", string(c.Ctx.Input.RequestBody), err)
		c.ErrorResponseAndLog(OPActionTaskWrite, ERROR_CODE_PARAMS_ERROR, string(c.Ctx.Input.RequestBody))
		return
	}

	if adminSvr == msg.AppName {
		taskMsgList, errCode := c.getTaskMsgListByAppNameAndName(msg.AppName, msg.Name)
		if errCode != ERROR_CODE_SUCCESS {
			c.ErrorResponseAndLog(OPActionTaskWrite, errCode, string(c.Ctx.Input.RequestBody))
			return
		}

		if taskList, ok := systoolbox.CheckTaskFunc(taskMsgList, false, FuncContainer, nil); ok {
			if len(taskList) > 0 {
				systoolbox.TaskMgr.RunTaskList(taskList)
			}
		}

	} else {
		buf, _, errCode := c.getTaskMsgListBufByAppNameAndName(msg.AppName, msg.Name)
		if errCode != ERROR_CODE_SUCCESS {
			c.ErrorResponseAndLog(OPActionTaskWrite, errCode, fmt.Sprint(msg))
			return
		}

		err = common.RabbitMQPublish(svrTaskBusinessName(msg.AppName, msg.RegionId, msg.ServerId),
			common.RabbitMQRoutingKeyTaskRun,
			buf)
		if err != nil {
			c.ErrorResponseAndLog(OPActionTaskWrite, ERROR_CODE_RABBIT_MQ_PUBLISH_FAIL, fmt.Sprint(msg))
			return
		}
	}
	c.SuccessResponseAndLog(OPActionTaskWrite, fmt.Sprint(msg), "")

}

// get task rabbitmq business name
func taskBusinessName(appName string) string {
	return fmt.Sprintf("%s.%s", common.RabbitMQBusinessNameTask, appName)
}

// get task rabbitmq business name of application
func svrTaskBusinessName(appName string, regionId, serverId int64) string {
	return fmt.Sprintf("%s.%s.%d.%d", common.RabbitMQBusinessNameTask, appName, regionId, serverId)
}

func (c *TaskController) getTaskMsgListByAppNameAndName(appName, name string) (taskMsgList systoolbox.TaskMsgList, errCode ERROR_CODE) {
	errCode = ERROR_CODE_SUCCESS

	task, err := admindao.TaskDaoEntity.QueryByAppNameAndName(appName, name)
	if err != nil {
		common.LogFuncError("%v", err)
		if err == orm.ErrNoRows {
			errCode = ERROR_CODE_TASK_NO_EXIST_ERROR
		} else {
			errCode = ERROR_CODE_DB
		}
		return
	}

	taskMsgList = systoolbox.TaskMsgList{
		List: []systoolbox.TaskMsg{
			{
				Id:       task.Id,
				Name:     task.Name,
				Spec:     task.Spec,
				FuncName: task.FuncName,
			},
		},
	}

	return
}

func (c *TaskController) getTaskMsgListBufByAppNameAndName(appName, name string) (buf []byte, appNameOut string, errCode ERROR_CODE) {
	errCode = ERROR_CODE_SUCCESS

	var taskMsgList systoolbox.TaskMsgList
	taskMsgList, errCode = c.getTaskMsgListByAppNameAndName(appName, name)
	if errCode != ERROR_CODE_SUCCESS {
		return
	}

	var err error
	buf, err = easyjson.Marshal(&taskMsgList)
	if err != nil {
		common.LogFuncError("%v", err)
		errCode = ERROR_CODE_ENCODE_FAIL
		return
	}

	appNameOut = appName

	return
}

func (c *TaskController) getTaskMsgListById(id uint32) (taskMsgList systoolbox.TaskMsgList, appName string, errCode ERROR_CODE) {
	errCode = ERROR_CODE_SUCCESS

	task, err := admindao.TaskDaoEntity.QueryById(id)
	if err != nil {
		common.LogFuncError("%v", err)
		if err == orm.ErrNoRows {
			errCode = ERROR_CODE_TASK_NO_EXIST_ERROR
		} else {
			errCode = ERROR_CODE_DB
		}
		return
	}

	taskMsgList = systoolbox.TaskMsgList{
		List: []systoolbox.TaskMsg{
			{
				Id:       task.Id,
				Name:     task.Name,
				Spec:     task.Spec,
				FuncName: task.FuncName,
			},
		},
	}

	appName = task.AppName

	return
}

func (c *TaskController) getTaskMsgListBufById(id uint32) (buf []byte, appName string, errCode ERROR_CODE) {
	errCode = ERROR_CODE_SUCCESS
	var taskMsgList systoolbox.TaskMsgList
	taskMsgList, appName, errCode = c.getTaskMsgListById(id)
	if errCode != ERROR_CODE_SUCCESS {
		return
	}

	var err error
	buf, err = easyjson.Marshal(&taskMsgList)
	if err != nil {
		common.LogFuncError("%v", err)
		errCode = ERROR_CODE_ENCODE_FAIL
		return
	}

	return
}

func (c *TaskController) getTaskMsgListBuf(taskMsgList systoolbox.TaskMsgList) (buf []byte, appName string, errCode ERROR_CODE) {
	errCode = ERROR_CODE_SUCCESS

	var err error
	buf, err = easyjson.Marshal(&taskMsgList)
	if err != nil {
		common.LogFuncError("%v", err)
		errCode = ERROR_CODE_ENCODE_FAIL
		return
	}

	return
}

func (c *TaskController) getTaskStatusMsgListByIdList(idList []uint32, status uint8, isCheckStatus, isCheckAppName bool) (taskMsgList systoolbox.TaskMsgList, appName string, errCode ERROR_CODE) {
	errCode = ERROR_CODE_SUCCESS
	tasks, err := admindao.TaskDaoEntity.QueryByIdList(idList)
	if err != nil {
		common.LogFuncError("%v", err)
		if err == orm.ErrNoRows || len(tasks) <= 0 {
			errCode = ERROR_CODE_TASK_NO_EXIST_ERROR
		} else {
			errCode = ERROR_CODE_DB
		}
		return
	}

	for _, task := range tasks {
		if isCheckStatus && task.Status == status {
			common.LogFuncError("task %d status is %d already", task.Id, task.Status)
			errCode = ERROR_CODE_TASK_STATUS_ERROR
			return
		}
		if isCheckAppName && appName != "" && appName != task.AppName {
			common.LogFuncError("task %d 's app name %s is different from %s", task.Id, task.AppName, appName)
			errCode = ERROR_CODE_TASK_APPNAME_DIFFERENT_ERROR
			return
		}
		appName = task.AppName
		taskMsgList.List = append(taskMsgList.List, systoolbox.TaskMsg{
			Name:   task.Name,
			Switch: status,
		})
	}

	return
}

func (c *TaskController) getTaskStatusMsgListBufByIdList(idList []uint32, status uint8, isCheckStatus, isCheckAppName bool) (buf []byte, appName string, errCode ERROR_CODE) {
	errCode = ERROR_CODE_SUCCESS

	var taskMsgList systoolbox.TaskMsgList
	taskMsgList, appName, errCode = c.getTaskStatusMsgListByIdList(idList, status, isCheckStatus, isCheckAppName)
	if errCode != ERROR_CODE_SUCCESS {
		return
	}

	var err error
	buf, err = easyjson.Marshal(&taskMsgList)
	if err != nil {
		common.LogFuncError("%v", err)
		errCode = ERROR_CODE_ENCODE_FAIL
		return
	}

	return
}

func (c *TaskController) getTaskStatusMsgListByNameList(appName string, nameList []string, status uint8, isCheckStatus, isCheckAppName bool) (taskMsgList systoolbox.TaskMsgList, errCode ERROR_CODE) {
	errCode = ERROR_CODE_SUCCESS
	tasks, err := admindao.TaskDaoEntity.QueryByAppNameAndNameList(appName, nameList)
	if err != nil {
		common.LogFuncError("%v", err)
		if err == orm.ErrNoRows || len(tasks) <= 0 {
			errCode = ERROR_CODE_TASK_NO_EXIST_ERROR
		} else {
			errCode = ERROR_CODE_DB
		}
		return
	}

	for _, task := range tasks {
		if isCheckStatus && task.Status == status {
			common.LogFuncError("task %d status is %d already", task.Id, task.Status)
			errCode = ERROR_CODE_TASK_STATUS_ERROR
			return
		}
		if isCheckAppName && appName != "" && appName != task.AppName {
			common.LogFuncError("task %d 's app name %s is different from %s", task.Id, task.AppName, appName)
			errCode = ERROR_CODE_TASK_APPNAME_DIFFERENT_ERROR
			return
		}

		taskMsgList.List = append(taskMsgList.List, systoolbox.TaskMsg{
			Id:     task.Id,
			Name:   task.Name,
			Switch: status,
		})
	}

	if len(taskMsgList.List) <= 0 {
		errCode = ERROR_CODE_TASK_NO_EXIST_ERROR
		return
	}

	return
}

func (c *TaskController) getTaskStatusMsgListBufByNameList(appName string, nameList []string, status uint8, isCheckStatus, isCheckAppName bool) (buf []byte, errCode ERROR_CODE) {
	errCode = ERROR_CODE_SUCCESS

	taskMsgList, errCode := c.getTaskStatusMsgListByNameList(appName, nameList, status, isCheckStatus, isCheckAppName)
	if errCode != ERROR_CODE_SUCCESS {
		return
	}

	var err error
	buf, err = easyjson.Marshal(&taskMsgList)
	if err != nil {
		common.LogFuncError("%v", err)
		errCode = ERROR_CODE_ENCODE_FAIL
		return
	}

	common.LogFuncDebug("buf : %v", buf)

	return
}

func (c *TaskController) trySwitch() (res interface{}, errCode ERROR_CODE) {
	errCode = ERROR_CODE_SUCCESS

	type Msg struct {
		IdList []uint32 `json:"id_list"`
		Switch uint8    `json:"switch"`
	}

	msg := Msg{}

	err := json.Unmarshal(c.Ctx.Input.RequestBody, &msg)
	if err != nil {
		common.LogFuncDebug("json decode: %s \nfailed : %v", string(c.Ctx.Input.RequestBody), err)
		errCode = ERROR_CODE_PARAMS_ERROR
		return
	}

	err = admindao.TaskDaoEntity.UpdateStatusMulti(msg.IdList, msg.Switch)
	if err != nil {
		errCode = ERROR_CODE_DB
		return
	}

	res = msg

	return
}

func (c *TaskController) trySwitchSingle() (res interface{}, errCode ERROR_CODE) {
	errCode = ERROR_CODE_SUCCESS
	//_, errCode = c.CheckPermission()
	//if errCode != ERROR_CODE_SUCCESS {
	//	c.ErrorResponseAndLog(OPActionTaskWrite, errCode, string(c.Ctx.Input.RequestBody))
	//	return
	//}

	type Msg struct {
		NameList []string `json:"name_list"`
		Status   uint8    `json:"status"`
		AppName  string   `json:"app_name"`
		RegionId int64    `json:"region_id"`
		ServerId int64    `json:"server_id"`
	}

	msg := Msg{}

	err := json.Unmarshal(c.Ctx.Input.RequestBody, &msg)
	if err != nil {
		common.LogFuncDebug("json decode: %s \nfailed : %v", string(c.Ctx.Input.RequestBody), err)
		errCode = ERROR_CODE_PARAMS_ERROR
		return
	}

	if msg.AppName == adminSvr { //switch from local
		var taskMgList systoolbox.TaskMsgList
		taskMgList, errCode = c.getTaskStatusMsgListByNameList(msg.AppName, msg.NameList, msg.Status, false, true)
		if errCode != ERROR_CODE_SUCCESS {
			return
		}
		if taskList, ok := systoolbox.CheckTaskFunc(taskMgList, false, FuncContainer, nil); ok {
			if len(taskList) > 0 {
				systoolbox.TaskMgr.SwitchTaskList(taskList)
			}
		}
	} else { //switch from remote
		var (
			buf []byte
		)
		buf, errCode = c.getTaskStatusMsgListBufByNameList(msg.AppName, msg.NameList, msg.Status, false, true)
		if errCode != ERROR_CODE_SUCCESS {
			return
		}

		err = common.RabbitMQPublish(svrTaskBusinessName(msg.AppName, msg.RegionId, msg.ServerId),
			common.RabbitMQRoutingKeyTaskSwitch,
			buf)
		if err != nil {
			errCode = ERROR_CODE_RABBIT_MQ_PUBLISH_FAIL
			return
		}
	}

	res = msg

	return
}

// stroe task result locally
func taskResultStoreLocal(name string, code int32, detail string, beginTime, endTime uint32) {
	admindao.TaskResultDaoEntity.Add(adminmodels.TaskResult{
		AppName:   adminSvr,
		RegionId:  0,
		ServerId:  0,
		Name:      name,
		Code:      code,
		Detail:    detail,
		BeginTime: beginTime,
		EndTime:   endTime,
		Ctime:     endTime,
	})
}
