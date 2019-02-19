package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type Msg struct {
	Status  bool   `json:"Status"`
	Code    int    `json:"Code"`
	Message string `json:"Message"`
}

type Currency struct {
	Lock        bool               `json:"Lock"`
	TokenName   string             `json:"TokenName"`
	TokenSymbol string             `json:"TokenSymbol"`
	TotalSupply float64            `json:"TotalSupply"`
	User        map[string]float64 `json:"User"`
}

type Token struct {
	Currency map[string]Currency `json:"Currency"`
}

func (token *Token) transfer(_from *Account, _to *Account, _currency string, _value float64) []byte {

	var rev []byte
	if token.Currency[_currency].Lock {
		msg := &Msg{Status: false, Code: 0, Message: "锁仓状态，停止一切转账活动"}
		rev, _ = json.Marshal(msg)
		return rev
	}
	if _from.Frozen {
		msg := &Msg{Status: false, Code: 0, Message: "From 账号冻结"}
		rev, _ = json.Marshal(msg)
		return rev
	}
	if _to.Frozen {
		msg := &Msg{Status: false, Code: 0, Message: "To 账号冻结"}
		rev, _ = json.Marshal(msg)
		return rev
	}
	if !token.isCurrency(_currency) {
		msg := &Msg{Status: false, Code: 0, Message: "货币符号不存在"}
		rev, _ = json.Marshal(msg)
		return rev
	}
	if _from.BalanceOf[_currency] >= _value {
		_from.BalanceOf[_currency] -= _value
		_to.BalanceOf[_currency] += _value
		cur := token.Currency[_currency]
		cur.User[_from.Name] -= _value
		if cur.User[_to.Name] == 0 {
			cur.User[_to.Name] = _value
		} else {
			cur.User[_to.Name] += _value
		}
		token.Currency[_currency] = cur
		msg := &Msg{Status: true, Code: 0, Message: "转账成功"}
		rev, _ = json.Marshal(msg)
		return rev
	} else {
		msg := &Msg{Status: false, Code: 0, Message: "余额不足"}
		rev, _ = json.Marshal(msg)
		return rev
	}

}
func (token *Token) initialSupply(_name string, _symbol string, _supply float64, _account *Account, lock bool) []byte {
	if _, ok := token.Currency[_symbol]; ok {
		msg := &Msg{Status: false, Code: 0, Message: "代币已经存在"}
		rev, _ := json.Marshal(msg)
		return rev
	}

	if _account.BalanceOf[_symbol] > 0 {
		msg := &Msg{Status: false, Code: 0, Message: "账号中存在代币"}
		rev, _ := json.Marshal(msg)
		return rev
	} else {
		user := make(map[string]float64)
		user[_account.Name] = _supply
		token.Currency[_symbol] = Currency{TokenName: _name, TokenSymbol: _symbol, TotalSupply: _supply, Lock: lock, User: user}
		_account.BalanceOf[_symbol] = _supply

		msg := &Msg{Status: true, Code: 0, Message: "代币初始化成功"}
		rev, _ := json.Marshal(msg)
		return rev
	}

}

func (token *Token) mint(_currency string, _amount float64, _account *Account) []byte {
	if !token.isCurrency(_currency) {
		msg := &Msg{Status: false, Code: 0, Message: "货币符号不存在"}
		rev, _ := json.Marshal(msg)
		return rev
	}
	cur := token.Currency[_currency]
	cur.TotalSupply += _amount
	cur.User[_account.Name] += _amount
	token.Currency[_currency] = cur
	_account.BalanceOf[_currency] += _amount

	msg := &Msg{Status: true, Code: 0, Message: "代币增发成功"}
	rev, _ := json.Marshal(msg)
	return rev

}
func (token *Token) burn(_currency string, _amount float64, _account *Account) []byte {
	if !token.isCurrency(_currency) {
		msg := &Msg{Status: false, Code: 0, Message: "货币符号不存在"}
		rev, _ := json.Marshal(msg)
		return rev
	}
	if _account.BalanceOf[_currency] >= _amount {
		cur := token.Currency[_currency]
		cur.TotalSupply -= _amount
		cur.User[_account.Name] -= _amount
		token.Currency[_currency] = cur
		_account.BalanceOf[_currency] -= _amount

		msg := &Msg{Status: false, Code: 0, Message: "代币回收成功"}
		rev, _ := json.Marshal(msg)
		return rev
	} else {
		msg := &Msg{Status: false, Code: 0, Message: "代币回收失败，回收额度不足"}
		rev, _ := json.Marshal(msg)
		return rev
	}

}
func (token *Token) isCurrency(_currency string) bool {
	if _, ok := token.Currency[_currency]; ok {
		return true
	} else {
		return false
	}
}
func (token *Token) setLock(_currency string, _look bool) bool {
	cur := token.Currency[_currency]
	cur.Lock = _look
	token.Currency[_currency] = cur
	return token.Currency[_currency].Lock
}

type Account struct {
	Name      string             `json:"Name"`
	Frozen    bool               `json:"Frozen"`
	BalanceOf map[string]float64 `json:"BalanceOf"`
}

func (account *Account) balance(_currency string) map[string]float64 {
	bal := map[string]float64{_currency: account.BalanceOf[_currency]}
	return bal
}

func (account *Account) balanceAll() map[string]float64 {
	return account.BalanceOf
}

// -----------
const TokenKey = "Token"
const Admin = "Admin"

// Define the Smart Contract structure
type SmartContract struct {
}

func (s *SmartContract) Init(stub shim.ChaincodeStubInterface) pb.Response {

	token := &Token{Currency: map[string]Currency{}}

	tokenAsBytes, err := json.Marshal(token)
	err = stub.PutState(TokenKey, tokenAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	} else {
		fmt.Printf("Init Token %s \n", string(tokenAsBytes))
	}
	err = stub.SetEvent("tokenInvoke", []byte{})
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

func (s *SmartContract) Query(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	if function == "balance" {
		return s.balance(stub, args)
	} else if function == "balanceAll" {
		return s.balanceAll(stub, args)
	} else if function == "showAccount" {
		return s.showAccount(stub, args)
	}
	return shim.Error("Invalid Smart Contract function name.")
}

func (s *SmartContract) Invoke(stub shim.ChaincodeStubInterface) pb.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := stub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger appropriately
	//CORE_PEER_ADDRESS=peer:7052 CORE_CHAINCODE_ID_NAME=token:0 ./token
	//peer chaincode install -p chaincodedev/chaincode/token -n token -v 0
	//peer chaincode instantiate -n token -v 0 -c '{"Args":[]}' -C myc
	if function == "initLedger" {
		// 注册管理员账户
		//peer chaincode invoke -C myc -n token -c '{"function":"initLedger","Args":[]}'
		return s.initLedger(stub, args)
	} else if function == "createAccount" {
		//创建账户 参数 ： (1)账户名
		//peer chaincode invoke -C myc -n token -c '{"function":"createAccount","Args":["123"]}'
		return s.createAccount(stub, args)
	} else if function == "initCurrency" {
		//创建代币 (1) 代币全称 (2) 代币简称 (3) 代币总量 (4) 代币生成以后持有人 (5) 是否锁仓
		// peer chaincode invoke -C myc -n token -c '{"function":"initCurrency","Args":["Netkiller Token","NKC","1000000","skyhuihui","false"]}'
		return s.initCurrency(stub, args)
	} else if function == "setLock" {
		//锁仓某个代币 (1) 代币简称 (2) 是否锁仓 (3) 操作人
		//peer chaincode invoke -C myc -n token -c '{"function":"setLock","Args":["NKC","true","skyhuihui"]}'
		return s.setLock(stub, args)
	} else if function == "transferToken" {
		//转账 (1) 发送人(2) 接收人(3) 代币名(4)发送代币量
		//peer chaincode invoke -C myc -n token -c '{"function":"transferToken","Args":["skyhuihui","123","ada","12.584"]}'
		return s.transferToken(stub, args)
	} else if function == "frozenAccount" {
		//冻结账户 (1) 要冻结的账户 (2) 是否冻结 (3) 操作人
		//peer chaincode invoke -C myc -n token -c '{"function":"frozenAccount","Args":["netkiller","true","skyhuihui"]}'
		return s.frozenAccount(stub, args)
	} else if function == "mintToken" {
		//代币增发 (1)代币名称(2)增发数量(3)操作人，也是代币增发接收人
		//peer chaincode invoke -C myc -n token -c '{"function":"mintToken","Args":["NKC","5000","skyhuihui"]}'
		return s.mintToken(stub, args)
	} else if function == "burnToken" {
		//代币销毁 (1)代币名称(2)回收数量(3)回收的账户（回收谁的代币）(4)操作人
		//peer chaincode invoke -C myc -n token -c '{"function":"burnToken","Args":["NKC","5000","123","skyhuihui"]}'
		return s.burnToken(stub, args)
	} else if function == "balance" {
		//查询指定账户指定代币 (1)查询账户 （2） 代币名称
		//peer chaincode invoke -C myc -n token -c '{"function":"balance","Args":["skyhuihui","NKC"]}'
		return s.balance(stub, args)
	} else if function == "balanceAll" {
		//查询某个用户所有资金 (1)账户名
		//peer chaincode invoke -C myc -n token -c '{"function":"balanceAll","Args":["skyhuihui"]}'
		return s.balanceAll(stub, args)
	} else if function == "showAccount" {
		//查看某个账户(1)账户名
		//peer chaincode invoke -C myc -n token -c '{"function":"showAccount","Args":["skyhuihui"]}'
		return s.showAccount(stub, args)
	} else if function == "showToken" {
		//查看所有代币
		//peer chaincode invoke -C myc -n token -c '{"function":"showToken","Args":[]}'
		return s.showToken(stub, args)
	} else if function == "showTokenUser" {
		//查看代币的所有持有用户 （1）代币名
		//peer chaincode invoke -C myc -n token -c '{"function":"showTokenUser","Args":["ada"]}'
		return s.showTokenUser(stub, args)
	}

	return shim.Error("Invalid Smart Contract function name.")
}

func (s *SmartContract) createAccount(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	key := args[0]
	name := args[0]
	existAsBytes, err := stub.GetState(key)
	fmt.Printf("GetState(%s) %s \n", key, string(existAsBytes))
	if string(existAsBytes) != "" {
		fmt.Println("Failed to create account, Duplicate key.")
		return shim.Error("Failed to create account, Duplicate key.")
	}

	account := Account{
		Name:      name,
		Frozen:    false,
		BalanceOf: map[string]float64{}}

	accountAsBytes, _ := json.Marshal(account)
	err = stub.PutState(key, accountAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Printf("createAccount %s \n", string(accountAsBytes))

	err = stub.SetEvent("tokenInvoke", []byte{})
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(accountAsBytes)
}
func (s *SmartContract) initLedger(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	key := "skyhuihui"
	name := "skyhuihui"
	existAsBytes, err := stub.GetState(Admin)
	fmt.Printf("GetState(%s) %s \n", key, string(existAsBytes))
	if string(existAsBytes) != "" {
		fmt.Println("Failed to create account, Duplicate key.")
		return shim.Error("Failed to create account, Duplicate key.")
	}

	account := Account{
		Name:      name,
		Frozen:    false,
		BalanceOf: map[string]float64{}}

	accountAsBytes, _ := json.Marshal(account)
	err = stub.PutState(key, accountAsBytes)
	nameByte, _ := json.Marshal(name)
	err = stub.PutState(Admin, nameByte)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.SetEvent("tokenInvoke", []byte{})
	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Printf("createAccount %s \n", string(accountAsBytes))

	return shim.Success(accountAsBytes)
}

func (s *SmartContract) showToken(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	tokenAsBytes, err := stub.GetState(TokenKey)
	if err != nil {
		return shim.Error(err.Error())
	} else {
		fmt.Printf("GetState(%s)) %s \n", TokenKey, string(tokenAsBytes))
	}
	return shim.Success(tokenAsBytes)
}

func (s *SmartContract) showTokenUser(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	_token := args[0]
	token := Token{}
	existAsBytes, err := stub.GetState(TokenKey)
	if err != nil {
		return shim.Error(err.Error())
	} else {
		fmt.Printf("GetState(%s)) %s \n", TokenKey, string(existAsBytes))
	}
	json.Unmarshal(existAsBytes, &token)
	reToekn, err := json.Marshal(token.Currency[_token])
	if err != nil {
		return shim.Error(err.Error())
	} else {
		fmt.Printf("Account balance %s \n", string(reToekn))
	}
	return shim.Success(reToekn)
}

func (s *SmartContract) initCurrency(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}
	admin, err := stub.GetState(Admin)
	if admin == nil {
		return shim.Error("The administrator account is empty")
	}
	account, _ := json.Marshal(args[3])
	if string(account) != string(admin) {
		return shim.Error("Current account is not an admin account")
	}
	_name := args[0]
	_symbol := args[1]
	_supply, _ := strconv.ParseFloat(args[2], 64)
	_account := args[3]
	lock := args[4]

	coinbaseAsBytes, err := stub.GetState(_account)
	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Printf("Coinbase before %s \n", string(coinbaseAsBytes))

	coinbase := &Account{}

	json.Unmarshal(coinbaseAsBytes, &coinbase)

	token := Token{}
	existAsBytes, err := stub.GetState(TokenKey)
	if err != nil {
		return shim.Error(err.Error())
	} else {
		fmt.Printf("GetState(%s)) %s \n", TokenKey, string(existAsBytes))
	}
	json.Unmarshal(existAsBytes, &token)
	var blog bool
	if lock == "false" {
		blog = false
	} else {
		blog = true
	}
	result := token.initialSupply(_name, _symbol, _supply, coinbase, blog)
	tokenAsBytes, _ := json.Marshal(token)
	err = stub.PutState(TokenKey, tokenAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	} else {
		fmt.Printf("Init Token %s \n", string(tokenAsBytes))
	}

	coinbaseAsBytes, _ = json.Marshal(coinbase)
	err = stub.PutState(_account, coinbaseAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Printf("Coinbase after %s \n", string(coinbaseAsBytes))

	err = stub.SetEvent("tokenInvoke", []byte{})
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(result)
}

func (s *SmartContract) transferToken(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}
	_from := args[0]
	_to := args[1]
	_currency := args[2]
	_amount, _ := strconv.ParseFloat(args[3], 32)

	if _amount <= 0 {
		return shim.Error("Incorrect number of amount")
	}

	fromAsBytes, err := stub.GetState(_from)
	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Printf("fromAccount %s \n", string(fromAsBytes))
	fromAccount := &Account{}
	json.Unmarshal(fromAsBytes, &fromAccount)

	toAsBytes, err := stub.GetState(_to)
	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Printf("toAccount %s \n", string(toAsBytes))
	toAccount := &Account{}
	json.Unmarshal(toAsBytes, &toAccount)

	tokenAsBytes, err := stub.GetState(TokenKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Printf("Token %s \n", string(toAsBytes))
	token := Token{Currency: map[string]Currency{}}
	json.Unmarshal(tokenAsBytes, &token)

	result := token.transfer(fromAccount, toAccount, _currency, _amount)
	fmt.Printf("Result %s \n", string(result))

	tokenAsBytes, err = json.Marshal(token)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.PutState(TokenKey, tokenAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Printf("Token after %s \n", string(tokenAsBytes))

	fromAsBytes, err = json.Marshal(fromAccount)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.PutState(_from, fromAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	} else {
		fmt.Printf("fromAccount %s \n", string(fromAsBytes))
	}

	toAsBytes, err = json.Marshal(toAccount)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.PutState(_to, toAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	} else {
		fmt.Printf("toAccount %s \n", string(toAsBytes))
	}

	err = stub.SetEvent("tokenInvoke", []byte{})
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(result)
}
func (s *SmartContract) mintToken(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}
	admin, err := stub.GetState(Admin)
	if admin == nil {
		return shim.Error("The administrator account is empty")
	}
	account, _ := json.Marshal(args[2])
	if string(account) != string(admin) {
		return shim.Error("Current account is not an admin account")
	}

	_currency := args[0]
	_amount, _ := strconv.ParseFloat(args[1], 32)
	_account := args[2]

	coinbaseAsBytes, err := stub.GetState(_account)
	if err != nil {
		return shim.Error(err.Error())
	} else {
		fmt.Printf("Coinbase before %s \n", string(coinbaseAsBytes))
	}

	coinbase := &Account{}
	json.Unmarshal(coinbaseAsBytes, &coinbase)

	tokenAsBytes, err := stub.GetState(TokenKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Printf("Token before %s \n", string(tokenAsBytes))

	token := Token{}

	json.Unmarshal(tokenAsBytes, &token)

	result := token.mint(_currency, _amount, coinbase)

	tokenAsBytes, err = json.Marshal(token)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.PutState(TokenKey, tokenAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Printf("Token after %s \n", string(tokenAsBytes))

	coinbaseAsBytes, _ = json.Marshal(coinbase)
	err = stub.PutState(_account, coinbaseAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	} else {
		fmt.Printf("Coinbase after %s \n", string(coinbaseAsBytes))
	}

	fmt.Printf("mintToken %s \n", string(tokenAsBytes))

	err = stub.SetEvent("tokenInvoke", []byte{})
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(result)
}

func (s *SmartContract) burnToken(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}
	admin, err := stub.GetState(Admin)
	if admin == nil {
		return shim.Error("The administrator account is empty")
	}
	account, _ := json.Marshal(args[3])
	if string(account) != string(admin) {
		return shim.Error("Current account is not an admin account")
	}

	_currency := args[0]
	_amount, _ := strconv.ParseFloat(args[1], 32)
	_account := args[2]

	coinbaseAsBytes, err := stub.GetState(_account)
	if err != nil {
		return shim.Error(err.Error())
	} else {
		fmt.Printf("Coinbase before %s \n", string(coinbaseAsBytes))
	}

	coinbase := &Account{}
	json.Unmarshal(coinbaseAsBytes, &coinbase)

	tokenAsBytes, err := stub.GetState(TokenKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Printf("Token before %s \n", string(tokenAsBytes))

	token := Token{}

	json.Unmarshal(tokenAsBytes, &token)

	result := token.burn(_currency, _amount, coinbase)

	tokenAsBytes, err = json.Marshal(token)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.PutState(TokenKey, tokenAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Printf("Token after %s \n", string(tokenAsBytes))

	coinbaseAsBytes, _ = json.Marshal(coinbase)
	err = stub.PutState(_account, coinbaseAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	} else {
		fmt.Printf("Coinbase after %s \n", string(coinbaseAsBytes))
	}

	fmt.Printf("mintToken %s \n", string(tokenAsBytes))
	err = stub.SetEvent("tokenInvoke", []byte{})
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(result)
}

func (s *SmartContract) setLock(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}
	admin, err := stub.GetState(Admin)
	if admin == nil {
		return shim.Error("The administrator account is empty")
	}
	account, _ := json.Marshal(args[2])
	if string(account) != string(admin) {
		return shim.Error("Current account is not an admin account")
	}
	_currency := args[0]
	_look := args[1]

	tokenAsBytes, err := stub.GetState(TokenKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	// fmt.Printf("setLock - begin %s \n", string(tokenAsBytes))

	token := Token{}

	json.Unmarshal(tokenAsBytes, &token)

	if _look == "true" {
		token.setLock(_currency, true)
	} else {
		token.setLock(_currency, false)
	}

	tokenAsBytes, err = json.Marshal(token)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.PutState(TokenKey, tokenAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Printf("setLock - end %s \n", string(tokenAsBytes))
	err = stub.SetEvent("tokenInvoke", []byte{})
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}
func (s *SmartContract) frozenAccount(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}
	admin, err := stub.GetState(Admin)
	if admin == nil {
		return shim.Error("The administrator account is empty")
	}
	acc, _ := json.Marshal(args[2])
	if string(acc) != string(admin) {
		return shim.Error("Current account is not an admin account")
	}

	_account := args[0]
	_status := args[1]

	accountAsBytes, err := stub.GetState(_account)
	if err != nil {
		return shim.Error(err.Error())
	}
	// fmt.Printf("setLock - begin %s \n", string(tokenAsBytes))

	account := Account{}

	json.Unmarshal(accountAsBytes, &account)

	var status bool
	if _status == "true" {
		status = true
	} else {
		status = false
	}

	account.Frozen = status

	accountAsBytes, err = json.Marshal(account)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.PutState(_account, accountAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	} else {
		fmt.Printf("frozenAccount - end %s \n", string(accountAsBytes))
	}
	err = stub.SetEvent("tokenInvoke", []byte{})
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

func (s *SmartContract) showAccount(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	_account := args[0]

	accountAsBytes, err := stub.GetState(_account)
	if err != nil {
		return shim.Error(err.Error())
	} else {
		fmt.Printf("Account balance %s \n", string(accountAsBytes))
	}
	return shim.Success(accountAsBytes)
}

func (s *SmartContract) balance(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	_account := args[0]
	_currency := args[1]

	accountAsBytes, err := stub.GetState(_account)
	if err != nil {
		return shim.Error(err.Error())
	} else {
		fmt.Printf("Account balance %s \n", string(accountAsBytes))
	}

	account := Account{}
	json.Unmarshal(accountAsBytes, &account)
	result := account.balance(_currency)

	resultAsBytes, _ := json.Marshal(result)
	fmt.Printf("%s balance is %s \n", _account, string(resultAsBytes))

	return shim.Success(resultAsBytes)
}

func (s *SmartContract) balanceAll(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	_account := args[0]

	accountAsBytes, err := stub.GetState(_account)
	if err != nil {
		return shim.Error(err.Error())
	} else {
		fmt.Printf("Account balance %s \n", string(accountAsBytes))
	}

	account := Account{}
	json.Unmarshal(accountAsBytes, &account)
	result := account.balanceAll()
	resultAsBytes, _ := json.Marshal(result)
	fmt.Printf("%s balance is %s \n", _account, string(resultAsBytes))

	return shim.Success(resultAsBytes)
}

// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
