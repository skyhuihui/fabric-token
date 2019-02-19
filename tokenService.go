package blockchain

import (
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"time"
)

func regitserEvent(client *channel.Client, chaincodeID, eventID string) (fab.Registration, <-chan *fab.CCEvent) {

	reg, notifier, err := client.RegisterChaincodeEvent(chaincodeID, eventID)
	if err != nil {
		fmt.Println("注册链码事件失败: %s", err)
	}
	return reg, notifier
}

func eventResult(notifier <-chan *fab.CCEvent, eventID string) error {
	select {
	case ccEvent := <-notifier:
		fmt.Printf("接收到链码事件: %v\n", ccEvent)
	case <-time.After(time.Second * 20):
		return fmt.Errorf("不能根据指定的事件ID接收到相应的链码事件(%s)", eventID)
	}
	return nil
}

// 注册管理员
func (setup *FabricSetup) InitLedger() (string, error) {

	// Prepare arguments
	var args []string
	args = append(args, "initLedger")

	eventID := "tokenInvoke"
	reg, notifier, err := setup.event[1].RegisterChaincodeEvent(setup.ChainCodeID[1], eventID)
	if err != nil {
		return "", err
	}
	defer setup.event[1].Unregister(reg)

	// Create a request (proposal) and send it
	response, err := setup.client[1].Execute(channel.Request{ChaincodeID: setup.ChainCodeID[1], Fcn: args[0], Args: [][]byte{}})
	if err != nil {
		return "", fmt.Errorf("failed to move funds: %v", err)
	}

	// Wait for the result of the submission
	select {
	case ccEvent := <-notifier:
		fmt.Printf("Received CC event: %s\n", ccEvent)
	case <-time.After(time.Second * 20):
		return "", fmt.Errorf("did NOT receive CC event for eventId(%s)", eventID)
	}

	return string(response.Payload), nil
}

// 创建账户
func (setup *FabricSetup) CreateAccount(value string) (string, error) {

	// Prepare arguments
	var args []string
	args = append(args, "createAccount")
	args = append(args, value)

	//eventID := "tokenInvoke"
	//reg, notifier, err := setup.event[1].RegisterChaincodeEvent(setup.ChainCodeID[1], eventID)
	//if err != nil {
	//	return "", err
	//}
	//defer setup.event[1].Unregister(reg)

	// Create a request (proposal) and send it
	response, err := setup.client[1].Execute(channel.Request{ChaincodeID: setup.ChainCodeID[1], Fcn: args[0], Args: [][]byte{[]byte(args[1])}})
	if err != nil {
		return "", fmt.Errorf("failed to move funds: %v", err)
	}

	// Wait for the result of the submission
	//select {
	//case ccEvent := <-notifier:
	//	fmt.Printf("Received CC event: %s\n", ccEvent)
	//case <-time.After(time.Second * 20):
	//	return "", fmt.Errorf("did NOT receive CC event for eventId(%s)", eventID)
	//}

	return string(response.Payload), nil
}

// 创建代币 (1) 代币全称 (2) 代币简称 (3) 代币总量 (4) 代币生成以后持有人 (5) 是否锁仓
func (setup *FabricSetup) InitCurrency(value []string) (string, error) {

	// Prepare arguments
	var args []string
	args = append(args, "initCurrency")
	args = append(args, value[0])
	args = append(args, value[1])
	args = append(args, value[2])
	args = append(args, value[3])
	args = append(args, value[4])

	//eventID := "tokenInvoke"
	//reg, notifier, err := setup.event[1].RegisterChaincodeEvent(setup.ChainCodeID[1], eventID)
	//if err != nil {
	//	return "", err
	//}
	//defer setup.event[1].Unregister(reg)

	// Create a request (proposal) and send it
	response, err := setup.client[1].Execute(channel.Request{ChaincodeID: setup.ChainCodeID[1], Fcn: args[0], Args: [][]byte{[]byte(args[1]), []byte(args[2]), []byte(args[3]), []byte(args[4]), []byte(args[5])}})
	if err != nil {
		return "", fmt.Errorf("failed to move funds: %v", err)
	}

	// Wait for the result of the submission
	//select {
	//case ccEvent := <-notifier:
	//	fmt.Printf("Received CC event: %s\n", ccEvent)
	//case <-time.After(time.Second * 20):
	//	return "", fmt.Errorf("did NOT receive CC event for eventId(%s)", eventID)
	//}

	return string(response.Payload), nil
}

//锁仓某个代币 (1) 代币简称 (2) 是否锁仓 (3) 操作人
func (setup *FabricSetup) SetLock(value []string) (string, error) {

	// Prepare arguments
	var args []string
	args = append(args, "setLock")
	args = append(args, value[0])
	args = append(args, value[1])
	args = append(args, value[2])

	//eventID := "tokenInvoke"
	//reg, notifier, err := setup.event[1].RegisterChaincodeEvent(setup.ChainCodeID[1], eventID)
	//if err != nil {
	//	return "", err
	//}
	//defer setup.event[1].Unregister(reg)

	// Create a request (proposal) and send it
	response, err := setup.client[1].Execute(channel.Request{ChaincodeID: setup.ChainCodeID[1], Fcn: args[0], Args: [][]byte{[]byte(args[1]), []byte(args[2]), []byte(args[3])}})
	if err != nil {
		return "", fmt.Errorf("failed to move funds: %v", err)
	}

	// Wait for the result of the submission
	//select {
	//case ccEvent := <-notifier:
	//	fmt.Printf("Received CC event: %s\n", ccEvent)
	//case <-time.After(time.Second * 20):
	//	return "", fmt.Errorf("did NOT receive CC event for eventId(%s)", eventID)
	//}

	return string(response.Payload), nil
}

//转账 (1) 发送人(2) 接收人(3) 代币名(4)发送代币量
func (setup *FabricSetup) TransferToken(value []string) (string, error) {

	// Prepare arguments
	var args []string
	args = append(args, "transferToken")
	args = append(args, value[0])
	args = append(args, value[1])
	args = append(args, value[2])
	args = append(args, value[3])

	//eventID := "tokenInvoke"
	//reg, notifier, err := setup.event[1].RegisterChaincodeEvent(setup.ChainCodeID[1], eventID)
	//if err != nil {
	//	return "", err
	//}
	//defer setup.event[1].Unregister(reg)

	// Create a request (proposal) and send it
	response, err := setup.client[1].Execute(channel.Request{ChaincodeID: setup.ChainCodeID[1], Fcn: args[0], Args: [][]byte{[]byte(args[1]), []byte(args[2]), []byte(args[3]), []byte(args[4])}})
	if err != nil {
		return "", fmt.Errorf("failed to move funds: %v", err)
	}

	// Wait for the result of the submission
	//select {
	//case ccEvent := <-notifier:
	//	fmt.Printf("Received CC event: %s\n", ccEvent)
	//case <-time.After(time.Second * 20):
	//	return "", fmt.Errorf("did NOT receive CC event for eventId(%s)", eventID)
	//}

	return string(response.Payload), nil
}

//冻结账户 (1) 要冻结的账户 (2) 是否冻结 (3) 操作人
func (setup *FabricSetup) FrozenAccount(value []string) (string, error) {

	// Prepare arguments
	var args []string
	args = append(args, "frozenAccount")
	args = append(args, value[0])
	args = append(args, value[1])
	args = append(args, value[2])

	//eventID := "tokenInvoke"
	//reg, notifier, err := setup.event[1].RegisterChaincodeEvent(setup.ChainCodeID[1], eventID)
	//if err != nil {
	//	return "", err
	//}
	//defer setup.event[1].Unregister(reg)

	// Create a request (proposal) and send it
	response, err := setup.client[1].Execute(channel.Request{ChaincodeID: setup.ChainCodeID[1], Fcn: args[0], Args: [][]byte{[]byte(args[1]), []byte(args[2]), []byte(args[3])}})
	if err != nil {
		return "", fmt.Errorf("failed to move funds: %v", err)
	}

	// Wait for the result of the submission
	//select {
	//case ccEvent := <-notifier:
	//	fmt.Printf("Received CC event: %s\n", ccEvent)
	//case <-time.After(time.Second * 20):
	//	return "", fmt.Errorf("did NOT receive CC event for eventId(%s)", eventID)
	//}

	return string(response.Payload), nil
}

//代币增发 (1)代币名称(2)增发数量(3)操作人，也是代币增发接收人
func (setup *FabricSetup) MintToken(value []string) (string, error) {

	// Prepare arguments
	var args []string
	args = append(args, "mintToken")
	args = append(args, value[0])
	args = append(args, value[1])
	args = append(args, value[2])

	//eventID := "tokenInvoke"
	//reg, notifier, err := setup.event[1].RegisterChaincodeEvent(setup.ChainCodeID[1], eventID)
	//if err != nil {
	//	return "", err
	//}
	//defer setup.event[1].Unregister(reg)

	// Create a request (proposal) and send it
	response, err := setup.client[1].Execute(channel.Request{ChaincodeID: setup.ChainCodeID[1], Fcn: args[0], Args: [][]byte{[]byte(args[1]), []byte(args[2]), []byte(args[3])}})
	if err != nil {
		return "", fmt.Errorf("failed to move funds: %v", err)
	}

	// Wait for the result of the submission
	//select {
	//case ccEvent := <-notifier:
	//	fmt.Printf("Received CC event: %s\n", ccEvent)
	//case <-time.After(time.Second * 20):
	//	return "", fmt.Errorf("did NOT receive CC event for eventId(%s)", eventID)
	//}

	return string(response.Payload), nil
}

//代币销毁 (1)代币名称(2)回收数量(3)回收的账户（回收谁的代币）(4)操作人
func (setup *FabricSetup) BurnToken(value []string) (string, error) {

	// Prepare arguments
	var args []string
	args = append(args, "burnToken")
	args = append(args, value[0])
	args = append(args, value[1])
	args = append(args, value[2])
	args = append(args, value[3])

	//eventID := "tokenInvoke"
	//reg, notifier, err := setup.event[1].RegisterChaincodeEvent(setup.ChainCodeID[1], eventID)
	//if err != nil {
	//	return "", err
	//}
	//defer setup.event[1].Unregister(reg)

	// Create a request (proposal) and send it
	response, err := setup.client[1].Execute(channel.Request{ChaincodeID: setup.ChainCodeID[1], Fcn: args[0], Args: [][]byte{[]byte(args[1]), []byte(args[2]), []byte(args[3]), []byte(args[4])}})
	if err != nil {
		return "", fmt.Errorf("failed to move funds: %v", err)
	}

	// Wait for the result of the submission
	//select {
	//case ccEvent := <-notifier:
	//	fmt.Printf("Received CC event: %s\n", ccEvent)
	//case <-time.After(time.Second * 20):
	//	return "", fmt.Errorf("did NOT receive CC event for eventId(%s)", eventID)
	//}

	return string(response.Payload), nil
}

//查询指定账户指定代币 (1)查询账户 （2） 代币名称
func (setup *FabricSetup) Balance(value []string) (string, error) {

	// Prepare arguments
	var args []string
	args = append(args, "balance")
	args = append(args, value[0])
	args = append(args, value[1])

	// Create a request (proposal) and send it
	response, err := setup.client[1].Execute(channel.Request{ChaincodeID: setup.ChainCodeID[1], Fcn: args[0], Args: [][]byte{[]byte(args[1]), []byte(args[2])}})
	if err != nil {
		return "", fmt.Errorf("failed to move funds: %v", err)
	}
	return string(response.Payload), nil
}

//查询某个用户所有资金 (1)账户名
func (setup *FabricSetup) BalanceAll(value string) (string, error) {

	// Prepare arguments
	var args []string
	args = append(args, "balanceAll")
	args = append(args, value)

	// Create a request (proposal) and send it
	response, err := setup.client[1].Execute(channel.Request{ChaincodeID: setup.ChainCodeID[1], Fcn: args[0], Args: [][]byte{[]byte(args[1])}})
	if err != nil {
		return "", fmt.Errorf("failed to move funds: %v", err)
	}
	return string(response.Payload), nil
}

//查看某个账户(1)账户名
func (setup *FabricSetup) ShowAccount(value string) (string, error) {

	// Prepare arguments
	var args []string
	args = append(args, "showAccount")
	args = append(args, value)

	response, err := setup.client[1].Execute(channel.Request{ChaincodeID: setup.ChainCodeID[1], Fcn: args[0], Args: [][]byte{[]byte(args[1])}})
	if err != nil {
		return "", fmt.Errorf("failed to move funds: %v", err)
	}
	return string(response.Payload), nil
}

//查看所有代币
func (setup *FabricSetup) ShowToken() (string, error) {

	// Prepare arguments
	var args []string
	args = append(args, "showToken")

	response, err := setup.client[1].Execute(channel.Request{ChaincodeID: setup.ChainCodeID[1], Fcn: args[0], Args: [][]byte{}})
	if err != nil {
		return "", fmt.Errorf("failed to move funds: %v", err)
	}
	return string(response.Payload), nil
}

//查看代币的所有持有用户 (1)代币名
func (setup *FabricSetup) ShowTokenUser(value string) (string, error) {

	// Prepare arguments
	var args []string
	args = append(args, "showTokenUser")
	args = append(args, value)

	// Create a request (proposal) and send it
	response, err := setup.client[1].Execute(channel.Request{ChaincodeID: setup.ChainCodeID[1], Fcn: args[0], Args: [][]byte{[]byte(args[1])}})
	if err != nil {
		return "", fmt.Errorf("failed to move funds: %v", err)
	}

	return string(response.Payload), nil
}
