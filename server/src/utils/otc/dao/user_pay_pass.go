package dao

import (
	"common"
	"fmt"
	"github.com/astaxie/beego/orm"
	"time"
	"utils/admin/dao"
	"utils/otc/models"
)

type PayPassValidator struct {
	Token      string `json:"token"`
	SignSalt   string `json:"sign_salt"`
	VerifyStep int    `json:"verify_step"`
}

// user payment password
const (
	UserPayPwdUnknown = iota
	UserPayPwdActive
	UserPayPwdDisable
	UserPayPwdMax
)

//user payment verify step
const (
	UserPayPwdStepOnce = iota
	UserPayPwdStepTwice
)

const (
	UserPayPwdMethodUnknown = iota
	UserPayPwdMethodPass    // set password method by old password
	UserPayPwdMethodSms     // set password method by phone sms
	UserPayPwdMethodMax
)
const (
	UserPayPwdActionUnknown     = iota
	UserPayPwdActionSellConfirm // set password method by old password
)

const (
	MaxPwdVerifyFail = 5   //连续失败5次
	PwdForbidSecond  = 300 //锁定5分钟
)

type UserPayPassDao struct {
	common.BaseDao
}

func NewUserPayPassDao(db string) *UserPayPassDao {
	return &UserPayPassDao{
		BaseDao: common.NewBaseDao(db),
	}
}

var UserPayPassDaoEntity *UserPayPassDao

// 设置密码
func (d *UserPayPassDao) InsertOrUpdatePwd(uid uint64, newPwd string, method int8, verifyStep int8) (data *models.UserPayPassword, err error) {
	data, err = d.QueryPwdByUid(uid)
	if err != nil {
		return
	}

	data = &models.UserPayPassword{
		Uid:        uid,
		Salt:       common.AppSignMgr.GenerateSalt(),
		SignSalt:   common.AppSignMgr.GenerateSalt(),
		Timestamp:  int64(common.NowUint32()),
		Status:     UserPayPwdActive,
		Method:     method,
		VerifyStep: verifyStep,
	}

	data.Password, err = d.generateEncryPwd(uid, newPwd, data.Salt, uint32(data.Timestamp))
	if err != nil {
		common.LogFuncError("generate sign err: %v", err)
		return
	}

	if _, err = d.Orm.InsertOrUpdate(data); err != nil {
		common.LogFuncError("db err: %v", err)
		return
	}

	return
}

// 仅验证密码
func (d *UserPayPassDao) ValidatePassword(uid uint64, pwd string) bool {
	info, err := d.QueryPwdByUid(uid)
	if err != nil {
		return false
	}

	if info.Uid == 0 || info.Status != UserPayPwdActive || info.Password == "" || d.IsLockPwd(uid) {
		return false
	}

	encryPwd, err := d.generateEncryPwd(uid, pwd, info.Salt, uint32(info.Timestamp))
	if err != nil {
		common.LogFuncError("generate sign err: %v", err)
		return false
	}
	return encryPwd == info.Password
}

// 查询支付密码信息
func (d *UserPayPassDao) QueryPwdByUid(uid uint64) (info *models.UserPayPassword, err error) {
	info = &models.UserPayPassword{
		Uid: uid,
	}
	if uid == 0 {
		return info, dao.ErrParam
	}

	if err = d.Orm.Read(info); err != nil {
		if err == orm.ErrNoRows {
			err = nil
			return
		}
		common.LogFuncError("mysql query err: %v", err)
		return info, err
	}

	return info, nil
}

// 查询支付密码状态
func (d *UserPayPassDao) QueryStatusByUid(uid uint64) (int, int, error) {
	if uid == 0 {
		return UserPayPwdUnknown, UserPayPwdStepOnce, dao.ErrParam
	}

	info, err := d.QueryPwdByUid(uid)
	if err != nil {
		common.LogFuncDebug("mysql query err: %v", err)
		return UserPayPwdUnknown, UserPayPwdStepOnce, err
	}

	if info.Status >= UserPayPwdMax || info.Status < UserPayPwdUnknown {
		info.Status = UserPayPwdUnknown
	}

	return int(info.Status), int(info.VerifyStep), nil
}

// 更新二次验证
func (d *UserPayPassDao) UpdateVerifyStep(uid uint64, verifyStep int8) error {
	if uid == 0 || (verifyStep != UserPayPwdStepOnce && verifyStep != UserPayPwdStepTwice) {
		return dao.ErrParam
	}
	info := &models.UserPayPassword{
		Uid:        uid,
		VerifyStep: verifyStep,
	}

	_, err := d.Orm.Update(info, models.COLUMN_UserPayPassword_VerifyStep)
	if err != nil {
		common.LogFuncDebug("mysql query err: %v", err)
		return err
	}

	return nil
}

// 校验redis缓存中的支付密码二次认证的token
// validator 为唯一验证器
// 验证完成后需显示删除token ： d.DelCacheToken(uid uint64)
func (d *UserPayPassDao) ValidateCacheToken(uid uint64, token string, clearCache bool) (isValid bool) {
	if token == "" || uid == 0 {
		return false
	}
	var key = d.generateCacheKey(uid)
	tmp, err := common.RedisManger.Get(key).Result()
	if err != nil {
		return false
	}

	if clearCache {
		d.DelCacheToken(uid)
	}
	return token == tmp
}

// 删除token
func (d *UserPayPassDao) DelCacheToken(uid uint64) {
	common.RedisManger.Del(d.generateCacheKey(uid))
}

// 验证密码并返回新的签名的盐，并返回唯一验证器：
// 验证器作为二次token验证的唯一标识
func (d *UserPayPassDao) ValidateByPassword(uid uint64, pwd string, expiration time.Duration) (isValid bool, passValidator PayPassValidator) {
	info, err := d.QueryPwdByUid(uid)
	if err != nil {
		return
	}
	if info.Uid == 0 || info.Status != UserPayPwdActive || info.Password == "" {
		return
	}

	// 验证密码
	encryPwd, err := d.generateEncryPwd(uid, pwd, info.Salt, uint32(info.Timestamp))
	if err != nil {
		common.LogFuncError("generate sign err: %v", err)
		return
	}
	if encryPwd != info.Password {
		return
	}
	passValidator.VerifyStep = int(info.VerifyStep)
	isValid = true

	// 生成临时token
	passValidator.Token = common.AppSignMgr.GenerateSalt()
	err = common.RedisManger.Set(d.generateCacheKey(uid), passValidator.Token, expiration).Err()
	if err != nil {
		isValid = false
		return
	}

	//重置签名的盐
	passValidator.SignSalt, err = d.ResetSignSalt(uid)
	if err != nil {
		common.LogFuncError(fmt.Sprintf("reset payment sign salt err:%v", err))
		passValidator.SignSalt = info.SignSalt
	}
	return
}

// 验证签名并返回新的签名的盐，并返回唯一验证器：
// 验证器作为二次token验证的唯一标识
func (d *UserPayPassDao) ValidateBySign(uid uint64, sign string, timestamp uint32, expiration time.Duration) (isValid bool, passValidator PayPassValidator) {

	info, err := d.QueryPwdByUid(uid)
	if err != nil {
		return
	}
	if info.Uid == 0 || info.Status != UserPayPwdActive || info.SignSalt == "" || d.IsLockPwd(uid) {
		return
	}

	src := map[string]string{
		"uid": fmt.Sprintf("%d", uid),
	}
	if signature, err := common.AppSignMgr.GenerateMSign(src, timestamp, info.SignSalt); err != nil || signature != sign {
		return
	}

	isValid = true
	passValidator.VerifyStep = int(info.VerifyStep)

	// 生成临时token
	passValidator.Token = common.AppSignMgr.GenerateSalt()
	err = common.RedisManger.Set(d.generateCacheKey(uid), passValidator.Token, expiration).Err()
	if err != nil {
		isValid = false
		return
	}

	// 设置新的盐， 设置失败返回原始盐
	passValidator.SignSalt, err = d.ResetSignSalt(uid)
	if err != nil {
		common.LogFuncError(fmt.Sprintf("reset payment sign salt err:%v", err))
		passValidator.SignSalt = info.SignSalt
	}

	return
}

func (d *UserPayPassDao) ValidatePasswordAndRecord(uid uint64, pwd string, expiration time.Duration) (isValid bool, passValidator PayPassValidator) {
	isValid, passValidator = d.ValidateByPassword(uid, pwd, expiration)
	d.VerifyPayPassResult(uid, isValid)
	return
}

// 重置签名
func (d *UserPayPassDao) ResetSignSalt(uid uint64) (signSalt string, err error) {
	data := &models.UserPayPassword{
		Uid:      uid,
		SignSalt: common.AppSignMgr.GenerateSalt(),
	}

	if _, err = d.Orm.Update(data, models.COLUMN_UserPayPassword_SignSalt); err != nil {
		common.LogFuncError("db err: %v", err)
		return
	}

	signSalt = data.SignSalt

	return
}

// 生成支付验证token的key
func (d *UserPayPassDao) generateCacheKey(uid uint64) (pass string) {
	pass = fmt.Sprintf("user_pay_pass_%d", uid)
	return
}

// 生成支付验证失败的key
func (d *UserPayPassDao) payPwdRedisKey(uid uint64) string {
	return fmt.Sprintf("pay_pwd_fail_%v", uid)
}

// 生成锁定密码的key
func (d *UserPayPassDao) payPwdLockRedisKey(uid uint64) string {
	return fmt.Sprintf("pay_pwd_lock_%v", uid)
}

// 生成支付密码加密串
func (d *UserPayPassDao) generateEncryPwd(uid uint64, pwd string, salt string, timestamp uint32) (pass string, err error) {
	src := map[string]string{
		"u": fmt.Sprintf("%d", uid),
		"p": pwd,
		"t": fmt.Sprintf("%d", timestamp),
	}

	var tmp string
	tmp, err = common.AppSignMgr.GenerateMSign(src, timestamp, salt)
	if err != nil {
		common.LogFuncError("generate sign err: %v", err)
		return
	}

	pass, err = common.AppSignMgr.GenerateSign(tmp, salt)
	if err != nil {
		common.LogFuncError("generate sign err: %v", err)
		return
	}
	return
}

func (d *UserPayPassDao) VerifyPayPassResult(uid uint64, ok bool) int {
	var num int
	var err error
	key := d.payPwdRedisKey(uid)
	if ok {
		num = 0
	} else {
		num, err = common.RedisManger.Get(key).Int()
		if err != nil {
			common.LogFuncError("PayPassResult:%v", err)
		}
		num = num + 1
	}

	if num >= MaxPwdVerifyFail {
		num = num % MaxPwdVerifyFail
		keyLock := d.payPwdLockRedisKey(uid)
		lockTime := time.Now().Unix() + PwdForbidSecond
		err = common.RedisManger.Set(keyLock, lockTime, PwdForbidSecond*time.Second).Err()
		if err != nil {
			common.LogFuncError("PayPassResult:%v", err)
		}
	}

	now := time.Now().Unix()
	_, end := common.TodayTimeRange()
	err = common.RedisManger.Set(key, num, time.Duration(end-now)*time.Second).Err()
	if err != nil {
		common.LogFuncError("PayPassResult:%v", err)
	}
	return num
}

func (d *UserPayPassDao) IsLockPwd(uid uint64) bool {
	keyLock := d.payPwdLockRedisKey(uid)
	now := time.Now().Unix()
	lockTime, _ := common.RedisManger.Get(keyLock).Int64()
	if lockTime > now {
		return true
	}

	return false
}

func (d *UserPayPassDao) GetLockInfo(uid uint64) (bool, int, int) {
	key := d.payPwdRedisKey(uid)
	num, err := common.RedisManger.Get(key).Int()
	if err != nil {
		common.LogFuncDebug("PayPassResult:%v", err)
	}

	keyLock := d.payPwdLockRedisKey(uid)
	now := time.Now().Unix()
	lockTime, _ := common.RedisManger.Get(keyLock).Int64()
	if lockTime > now {
		return true, num, int(lockTime - now)
	}

	return false, num, 0
}
