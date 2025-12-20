package evidenceDemo

import (
	"testing"

	"github.com/bubicn/bubichain-v4-sdk-go/src/model"
	"github.com/bubicn/bubichain-v4-sdk-go/src/sdk"
)

var (
	Url                      = "https://dev-openapi.bubi.cn/api/v1/fa978237669b4ebc858ce2e0bb2ec472"
	evidenceSenderPrivateKey = "privbyWmLKhtHY2WaTNjgYRwDL37w9WcNS1t36ZZyNxBt9uyzK1XqWJV"
)

func TestEvidenceStore(t *testing.T) {
	// 初始化SDK
	var sdkObj sdk.Sdk
	var initReq model.SDKInitRequest
	initReq.SetUrl(Url)
	initRes := sdkObj.Init(initReq)
	if initRes.ErrorCode != 0 {
		t.Errorf("SDK初始化失败: %s", initRes.ErrorDesc)
		return
	}

	// 准备存证请求
	var storeReq model.EvidenceStoreRequest
	storeReq.SetKey("test_key_001")                  // 设置存证的key
	storeReq.SetContent("这是一条测试存证内容")                // 设置存证的内容
	storeReq.SetPrivateKey(evidenceSenderPrivateKey) // 使用一个有效的私钥示例

	// 调用存证接口
	storeRes := sdkObj.Evidence.Store(storeReq)
	if storeRes.ErrorCode != 0 {
		t.Errorf("存证失败: 错误代码 %d, 错误描述: %s", storeRes.ErrorCode, storeRes.ErrorDesc)
		return
	}

	t.Logf("存证成功，Hash: %s", storeRes.Result.Hash)
}

func TestEvidenceGet(t *testing.T) {
	// 初始化SDK
	var sdkObj sdk.Sdk
	var initReq model.SDKInitRequest
	initReq.SetUrl(Url) // 替换为实际的API地址
	initRes := sdkObj.Init(initReq)
	if initRes.ErrorCode != 0 {
		t.Errorf("SDK初始化失败: %s", initRes.ErrorDesc)
		return
	}

	// 准备查询请求
	var getReq model.EvidenceGetRequest
	// 使用一个示例hash，实际使用时应该使用真实的存证hash
	getReq.SetHash("4d2cad7fb6bf6b622918c90b312aeecbd4e1bf7d775bde20502180a175858c53")

	// 调用查询接口
	getRes := sdkObj.Evidence.Get(getReq)
	if getRes.ErrorCode != 0 {
		t.Errorf("查询存证失败: 错误代码 %d, 错误描述: %s", getRes.ErrorCode, getRes.ErrorDesc)
		return
	}

	t.Logf("查询存证成功:")
	t.Logf("  Hash: %s", getRes.Result.Hash)
	t.Logf("  区块序号: %d", getRes.Result.LedgerSeq)
	t.Logf("  关闭时间: %d", getRes.Result.CloseTime)
	t.Logf("  内容Key: %s", getRes.Result.Content.Key)
	t.Logf("  内容: %s", getRes.Result.Content.Content)
	t.Logf("  时间戳: %d", getRes.Result.Content.Timestamp)
	t.Logf("  签名: %s", getRes.Result.Content.Signature)
	t.Logf("  公钥: %s", getRes.Result.Content.PublicKey)
}
