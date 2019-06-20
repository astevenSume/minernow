package usdt

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	"github.com/golang/mock/gomock"
	"testing"
	"usdt/dao/mock"
	. "usdt/error"
	"utils/usdt/dao"
	"utils/usdt/models"
)

//func initTest(t *testing.T) {
//	var (
//		dbName = "otc_test"
//		usr    = "otc"
//		pwd    = "otc"
//		port   = 3306
//		host   = "localhost"
//	)
//	err := orm.RegisterDriver("mysql", orm.DRMySQL)
//	if err != nil {
//		t.Fatalf("RegisterDriver failed : %v", err)
//	}
//
//	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4", usr, pwd, host, port, dbName)
//	err = orm.RegisterDataBase(dbName, "mysql", dsn, 10, 20)
//	if err != nil {
//		t.Fatalf("RegisterDataBase mysql for %s failed : %v", dsn, err)
//	}
//
//	db, err := orm.GetDB(dbName)
//	if err != nil {
//		t.Fatalf("orm.GetDB(\"mysql\") failed : %v", err)
//	}
//
//	if db == nil {
//		t.Fatalf("orm.GetDB(\"mysql\") is nil")
//	}
//
//	db.SetConnMaxLifetime(time.Second * time.Duration(10))
//
//	DbOrm = orm.NewOrm()
//	_ = DbOrm.Using(dbName)
//
//	err = models.ModelsInit()
//	if err != nil {
//		t.Fatalf("%v", err)
//	}
//}

func TestGenerateDbTransations_Single(t *testing.T) {

	var (
		uid             = uint64(2)
		addr, addr1     = "1DXDso8dNotaWxKu9MDt48ZNNgqgxJ1ggC", ""
		newLastestTx    string
		currPage, pages int
		dbTxs           []models.UsdtOnchainTransaction
		txId            = "2138bbe91afc4d7d19ed37614cb61a1d33fd09abc9027589b58e1fbb023b912e"
	)

	UsdtConfig.ConfirmationLimit = 6

	//init
	//initTest(t)

	ctl := gomock.NewController(t)
	defer ctl.Finish()

	mockOnchainData := mock_dao.NewMockOnchainDataDaoInterface(ctl)
	dao.OnchainDataDaoEntity = mockOnchainData

	// ============== Transaction format invalid
	resp := map[string]interface{}{
		"address":      addr,
		"pages":        "1", //need int 1
		"current_page": 0,
		"transactions": []Transaction{
			{
				Amount:           "10.00000000",
				Block:            568095,
				Blockhash:        "0000000000000000001daa5737c6364d3b93c8b985a0ba54dfd52b4d8cadc81b",
				Blocktime:        1553155216,
				Confirmations:    446,
				Divisible:        true,
				Fee:              "0.00002000",
				Ismine:           false,
				Positioninblock:  1859,
				Propertyid:       31,
				Propertyname:     "TetherUS",
				ReferenceAddress: "1GPa2CEu887nocR2MVzc993q925xHMRmpJ",
				Sendingaddress:   "1DXDso8dNotaWxKu9MDt48ZNNgqgxJ1ggC",
				Txid:             "2138bbe91afc4d7d19ed37614cb61a1d33fd09abc9027589b58e1fbb023b912e",
				Type:             "Simple Send",
				TypeInt:          0,
				Valid:            true,
				Version:          0,
			},
		},
	}
	_, _, _, _, errCode := generateDbTransations(uid, resp, &newLastestTx)
	if errCode != ERROR_CODE_CHAIN_TRANSACTION_ERROR {
		t.Fatalf("error code should be ERROR_CODE_CHAIN_TRANSACTION_ERROR")
	}

	// get lastestTx error.
	resp["pages"] = 1
	mockOnchainData.EXPECT().GetLastestTransaction(addr).Return("", fmt.Errorf("some error"))
	_, _, _, _, errCode = generateDbTransations(uid, resp, nil)
	if errCode != ERROR_CODE_DB {
		t.Errorf("error code should be ERROR_CODE_DB but %d", errCode)
	}

	//no lastestTx exists.
	mockOnchainData.EXPECT().GetLastestTransaction(addr).Return("", orm.ErrNoRows)
	_, _, _, _, errCode = generateDbTransations(uid, resp, &newLastestTx)
	if errCode != ERROR_CODE_SUCCESS {
		t.Fatalf("error code %d", errCode)
	}

	//generate succeed, tx data and newLastestTx return
	mockOnchainData.EXPECT().GetLastestTransaction(addr).Return("", orm.ErrNoRows)
	addr1, currPage, pages, dbTxs, errCode = generateDbTransations(uid, resp, &newLastestTx)
	if errCode != ERROR_CODE_SUCCESS {
		t.Errorf("error code should be ERROR_CODE_SUCCESS but %d", errCode)
	}
	if currPage != 0 {
		t.Errorf("currPage should be 0 but %d", currPage)
	}
	if pages != 1 {
		t.Errorf("pages should be 1 but %d", pages)
	}
	if len(dbTxs) != 1 {
		t.Errorf("len(dbTxs) should be 1 but %d", len(dbTxs))
	}
	if addr1 != addr {
		t.Errorf("addr1 should be %s but %s", addr, addr1)
	}

	//lastestTx exsits and just equals to the txId input.
	mockOnchainData.EXPECT().GetLastestTransaction(addr).Return(txId, nil)
	addr1, currPage, pages, dbTxs, errCode = generateDbTransations(uid, resp, &newLastestTx)
	if len(dbTxs) != 0 {
		t.Errorf("len(dbTxs) should be 0 but %d", len(dbTxs))
	}

	//lastestTx exsits but no equals to the txId input.
	mockOnchainData.EXPECT().GetLastestTransaction(addr).Return("some tx", nil)
	addr1, currPage, pages, dbTxs, errCode = generateDbTransations(uid, resp, &newLastestTx)
	if len(dbTxs) != 1 {
		t.Errorf("len(dbTxs) should be 0 but %d", len(dbTxs))
	}
	if newLastestTx != txId {
		t.Errorf("newLastestTx should be %s but %s", txId, newLastestTx)
	}

	//tx confirmation small than 6
	mockOnchainData.EXPECT().GetLastestTransaction(addr).Return("", nil)
	resp["transactions"].([]Transaction)[0].Confirmations = 5
	addr1, currPage, pages, dbTxs, errCode = generateDbTransations(uid, resp, &newLastestTx)
	if len(dbTxs) != 0 {
		t.Errorf("len(dbTxs) should be 0 but %d", len(dbTxs))
	}

	//tx valid = false
	mockOnchainData.EXPECT().GetLastestTransaction(addr).Return("", nil)
	resp["transactions"].([]Transaction)[0].Confirmations = 7
	resp["transactions"].([]Transaction)[0].Valid = false
	addr1, currPage, pages, dbTxs, errCode = generateDbTransations(uid, resp, &newLastestTx)
	if len(dbTxs) != 0 {
		t.Errorf("len(dbTxs) should be 0 but %d", len(dbTxs))
	}
}
