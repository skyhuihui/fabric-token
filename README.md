# fabric-token
在联盟链上构建token

测试命令：

//CORE_PEER_ADDRESS=peer:7052 CORE_CHAINCODE_ID_NAME=token:0 ./token
	//peer chaincode install -p chaincodedev/chaincode/token -n token -v 0
	//peer chaincode instantiate -n token -v 0 -c '{"Args":[]}' -C myc
		// 注册管理员账户
		//peer chaincode invoke -C myc -n token -c '{"function":"initLedger","Args":[]}'
		//创建账户 参数 ： (1)账户名
		//peer chaincode invoke -C myc -n token -c '{"function":"createAccount","Args":["123"]}'
		//创建代币 (1) 代币全称 (2) 代币简称 (3) 代币总量 (4) 代币生成以后持有人 (5) 是否锁仓
		// peer chaincode invoke -C myc -n token -c '{"function":"initCurrency","Args":["Netkiller Token","NKC","1000000","skyhuihui","false"]}'
		//锁仓某个代币 (1) 代币简称 (2) 是否锁仓 (3) 操作人
		//peer chaincode invoke -C myc -n token -c '{"function":"setLock","Args":["NKC","true","skyhuihui"]}'
		//转账 (1) 发送人(2) 接收人(3) 代币名(4)发送代币量
		//peer chaincode invoke -C myc -n token -c '{"function":"transferToken","Args":["skyhuihui","123","ada","12.584"]}'
		//冻结账户 (1) 要冻结的账户 (2) 是否冻结 (3) 操作人
		//peer chaincode invoke -C myc -n token -c '{"function":"frozenAccount","Args":["netkiller","true","skyhuihui"]}'
		//代币增发 (1)代币名称(2)增发数量(3)操作人，也是代币增发接收人
		//peer chaincode invoke -C myc -n token -c '{"function":"mintToken","Args":["NKC","5000","skyhuihui"]}'
		//代币销毁 (1)代币名称(2)回收数量(3)回收的账户（回收谁的代币）(4)操作人
		//peer chaincode invoke -C myc -n token -c '{"function":"burnToken","Args":["NKC","5000","123","skyhuihui"]}'
		//查询指定账户指定代币 (1)查询账户 （2） 代币名称
		//peer chaincode invoke -C myc -n token -c '{"function":"balance","Args":["skyhuihui","NKC"]}'
		//查询某个用户所有资金 (1)账户名
		//peer chaincode invoke -C myc -n token -c '{"function":"balanceAll","Args":["skyhuihui"]}'
		//查看某个账户(1)账户名
		//peer chaincode invoke -C myc -n token -c '{"function":"showAccount","Args":["skyhuihui"]}'
		//查看所有代币
		//peer chaincode invoke -C myc -n token -c '{"function":"showToken","Args":[]}'
		//查看代币的所有持有用户
		//peer chaincode invoke -C myc -n token -c '{"function":"showTokenUser","Args":["ada"]}'
