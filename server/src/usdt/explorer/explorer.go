package explorer

import (
	"common"
	"fmt"
	"github.com/astaxie/beego/httplib"
)

// for api.omniexplorer.info restful api usage.

const (
	OMNI_EXPLORER_URL = "https://api.omniwallet.org"

	PROTO_GET_BALANCES                       = "/v2/address/addr/"
	PROTO_GET_ADDR_DETAILS                   = "/v1/address/addr/details/"
	PROTO_ARMORY_GET_UNSIGNED                = "/v1/armory/getunsigned"
	PROTO_ARMORY_GET_RAWTRANSACTION          = "/v1/armory/getrawtransaction"
	PROTO_DECODE                             = "/v1/decode/"
	PROTO_OMNIDEX_DESIGNATING_CURRENCIES     = "/v1/omnidex/designatingcurrencies"
	PROTO_PROPERTIES_GET_HISTORY             = "/v1/properties/gethistory/"
	PROTO_PROPERTIES_LIST_BY_OWNER           = "/v1/properties/listbyowner"
	PROTO_PROPERTIES_LIST_ACTIVE_CROWD_SALES = "/v1/properties/listactivecrowdsales"
	PROTO_PROPERTIES_LIST_BY_ECOSYSTEM       = "/v1/properties/listbyecosystem"
	PROTO_PROPERTIES_LIST                    = "/v1/properties/list"
	PROTO_SEARCH                             = "/v1/search"
	PROTO_TRANSACTION_ADDRESS                = "/v1/transaction/address"
	PROTO_TRANSACTION_PUSH                   = "/v1/transaction/pushtx/"
	PROTO_TRANSACTION_TX					 = "/v1/transaction/tx/"
)

type ExplorerInterface interface {
	Balances(addrs []string) (resp map[string]interface{}, err error)
	AddrDetails(addr string) (resp map[string]interface{}, err error)
	GetArmoryGetUnsigned(unsignedHex, pubKey string) (resp map[string]interface{}, err error)
	GetRawTransaction(armoryTransaction string) (resp map[string]interface{}, err error)
	DecodeRaw(hex string) (resp map[string]interface{}, err error)
	GetOmnidexDesignatingCurrencies(ecosystem int) (resp map[string]interface{}, err error)
	PropertyGetHistory(page int) (resp map[string]interface{}, err error)
	PropertiesListByOwner(addresses string) (resp map[string]interface{}, err error)
	PropertiesListActiveCrowdSales(ecosystem int) (resp map[string]interface{}, err error)
	PropertiesListByEcosystem(ecosystem int) (resp map[string]interface{}, err error)
	PropertiesList(ecosystem int) (resp map[string]interface{}, err error)
	Search(query string) (resp map[string]interface{}, err error)
	TransactionAddress(addr string, page int) (resp map[string]interface{}, err error)
	TransactionPush(signedTransaction string) (resp map[string]interface{}, err error)
	Transaction(hash string) (resp map[string]interface{}, err error)
}

//
type Explorer struct {
	req *httplib.BeegoHTTPRequest
}

func NewExplorer() *Explorer {
	return &Explorer{}
}

// Return the balance information for multiple addresses
// Data :
// addr : address
func (e *Explorer) Balances(addrs []string) (resp map[string]interface{}, err error) {
	e.initRequest(PROTO_GET_BALANCES)
	for _, addr := range addrs {
		e.param("addr", addr)
	}

	return e.do()
}

// Return the balance information and transaction history list for a given address
// Data :
// addr : address
func (e *Explorer) AddrDetails(addr string) (resp map[string]interface{}, err error) {
	e.initRequest(PROTO_GET_ADDR_DETAILS)
	e.param("addr", addr)
	return e.do()
}

// Returns the Armory encoded version of an unsigned transaction for use with Armory offline transactions.
// Data:
// unsigned_hex : raw bitcoin hex formatted tx to be converted
// pubkey : pubkey of the sending address
func (e *Explorer) GetArmoryGetUnsigned(unsignedHex, pubKey string) (resp map[string]interface{}, err error) {
	e.initRequest(PROTO_ARMORY_GET_UNSIGNED)
	e.param("unsigned_hex", unsignedHex)
	e.param("pubkey", pubKey)
	return e.do()
}

//Decodes and returns the raw hex and signed status from an armory transaction.
// Data:
// armory_tx : armory transaction in text format
func (e *Explorer) GetRawTransaction(armoryTransaction string) (resp map[string]interface{}, err error) {
	e.initRequest(PROTO_ARMORY_GET_RAWTRANSACTION)
	e.param("armory_tx", armoryTransaction)
	return e.do()
}

// Decodes raw hex returning Omni and Bitcoin transaction information
// Data:
// hex :
func (e *Explorer) DecodeRaw(hex string) (resp map[string]interface{}, err error) {
	e.initRequest(PROTO_DECODE)
	e.param("hex", hex)
	return e.do()
}

// Return a list of currently active/available base currencies the omnidex has open orders against.
// Data:
// ecosystem : 1 for main / production ecosystem or 2 for test/development ecosystem
func (e *Explorer) GetOmnidexDesignatingCurrencies(ecosystem int) (resp map[string]interface{}, err error) {
	e.initRequest(PROTO_OMNIDEX_DESIGNATING_CURRENCIES)
	e.param("ecosystem", fmt.Sprint(ecosystem))
	return e.do()
}

// Returns list of transactions (up to 10 per page) relevant to queried Property ID.
// Returned transaction types include:
// Creation Tx, Change issuer txs, Grant Txs, Revoke Txs, Crowdsale Participation Txs, Close Crowdsale earlier tx
func (e *Explorer) PropertyGetHistory(page int) (resp map[string]interface{}, err error) {
	e.initRequest(PROTO_PROPERTIES_GET_HISTORY)
	e.param("page", fmt.Sprint(page))
	return e.do()
}

// Return list of properties created by a queried address.
// Data ：
// addresses ：
func (e *Explorer) PropertiesListByOwner(addresses string) (resp map[string]interface{}, err error) {
	e.initRequest(PROTO_PROPERTIES_LIST_BY_OWNER)
	e.param("addresses", addresses)
	return e.do()
}

// Returns list of currently active crowdsales.
// Data:
// ecosystem : 1 for production/main ecosystem. 2 for test/dev ecosystem
func (e *Explorer) PropertiesListActiveCrowdSales(ecosystem int) (resp map[string]interface{}, err error) {
	e.initRequest(PROTO_PROPERTIES_LIST_ACTIVE_CROWD_SALES)
	e.param("ecosystem", fmt.Sprint(ecosystem))
	return e.do()
}

// Returns list of created properties filtered by ecosystem.
// Data:
// ecosystem : 1 for production/main ecosystem. 2 for test/dev ecosystem
func (e *Explorer) PropertiesListByEcosystem(ecosystem int) (resp map[string]interface{}, err error) {
	e.initRequest(PROTO_PROPERTIES_LIST_BY_ECOSYSTEM)
	e.param("ecosystem", fmt.Sprint(ecosystem))
	return e.do()
}

// Returns list of all created properties.
func (e *Explorer) PropertiesList(ecosystem int) (resp map[string]interface{}, err error) {
	e.initRequest(PROTO_PROPERTIES_LIST)
	return e.do()
}

// Search by transaction id, address or property id.
// Data:
// query : text string of either Transaction ID, Address, or property id to search for
func (e *Explorer) Search(query string) (resp map[string]interface{}, err error) {
	e.initRequest(PROTO_SEARCH)
	e.param("query", query)
	return e.do()
}

// Returns list of transactions for queried address.
// Data:
// addr : address to query
// page : cycle through available response pages (10 txs per page)
func (e *Explorer) TransactionAddress(addr string, page int) (resp map[string]interface{}, err error) {
	e.initRequest(PROTO_TRANSACTION_ADDRESS)
	e.param("addr", addr)
	e.param("page", fmt.Sprint(page))
	return e.do()
}

// Broadcast a signed transaction to the network.
// Data:
// signedTransaction : signed hex to broadcast
func (e *Explorer) TransactionPush(signedTransaction string) (resp map[string]interface{}, err error) {
	e.initRequest(PROTO_TRANSACTION_PUSH)
	e.param("signedTransaction", signedTransaction)
	return e.do()
	////req := httplib.Post("https://api.omniwallet.org/v1/transaction/pushtx/")
	////req = req.Header("Content-Type", "application/x-www-form-urlencoded")
	//req := e.req.Param("signedTransaction", signedTransaction)
	////var resp map[string]interface{}
	//err = req.ToJSON(&resp)
	//if err != nil {
	//	common.LogFuncError("%v", err)
	//	return
	//}

	//return
}

func (e *Explorer) TransactionTX(txid string) (resp map[string]interface{}, err error) {
	e.initRequestGet(PROTO_TRANSACTION_TX + txid)
	return e.do()
	////req := httplib.Post("https://api.omniwallet.org/v1/transaction/pushtx/")
	////req = req.Header("Content-Type", "application/x-www-form-urlencoded")
	//req := e.req.Param("signedTransaction", signedTransaction)
	////var resp map[string]interface{}
	//err = req.ToJSON(&resp)
	//if err != nil {
	//	common.LogFuncError("%v", err)
	//	return
	//}

	//return
}

// Returns transaction details of a queried transaction hash.
func (e *Explorer) Transaction(hash string) (resp map[string]interface{}, err error) {
	e.initRequest(PROTO_TRANSACTION_ADDRESS + "/" + hash)
	return e.do()
}

func (e *Explorer) initRequest(proto string) {
	e.req = httplib.Post(OMNI_EXPLORER_URL + proto)
	e.req.Header("Content-Type", "application/x-www-form-urlencoded")
	common.LogFuncDebug("url : %s", OMNI_EXPLORER_URL+proto)
	return
}

func (e *Explorer) initRequestGet(proto string) {
	e.req = httplib.Get(OMNI_EXPLORER_URL + proto)
	e.req.Header("Content-Type", "application/x-www-form-urlencoded")
	common.LogFuncDebug("get url : %s", OMNI_EXPLORER_URL+proto)
	return
}

func (e *Explorer) do() (resp map[string]interface{}, err error) {
	resp = make(map[string]interface{}, 0)
	err = e.req.ToJSON(&resp)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	return
}

func (e *Explorer) param(key, value string) {
	e.req.Param(key, value)
}
