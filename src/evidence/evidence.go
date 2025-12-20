// evidence
package evidence

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/bubicn/bubichain-v4-sdk-go/src/common"
	"github.com/bubicn/bubichain-v4-sdk-go/src/crypto/keypair"
	"github.com/bubicn/bubichain-v4-sdk-go/src/crypto/signature"
	"github.com/bubicn/bubichain-v4-sdk-go/src/exception"
	"github.com/bubicn/bubichain-v4-sdk-go/src/model"
)

type EvidenceOperation struct {
	Url string
}

// Store 存证接口
// 封装了签名过程，用户只需提供 key、content 和 privateKey
func (evidence *EvidenceOperation) Store(reqData model.EvidenceStoreRequest) model.EvidenceStoreResponse {
	var resData model.EvidenceStoreResponse

	// 验证请求参数
	key := reqData.GetKey()
	if key == "" {
		SDKRes := exception.GetSDKRes(exception.INVALID_DATAKEY_ERROR)
		resData.ErrorCode = SDKRes.ErrorCode
		resData.ErrorDesc = SDKRes.ErrorDesc
		return resData
	}
	if len(key) > 1024 {
		SDKRes := exception.GetSDKRes(exception.INVALID_DATAKEY_ERROR)
		resData.ErrorCode = SDKRes.ErrorCode
		resData.ErrorDesc = SDKRes.ErrorDesc
		return resData
	}

	content := reqData.GetContent()
	if content == "" {
		SDKRes := exception.GetSDKRes(exception.SYSTEM_ERROR)
		resData.ErrorCode = SDKRes.ErrorCode
		resData.ErrorDesc = "Content cannot be empty"
		return resData
	}
	if len(content) > 256*1024 {
		SDKRes := exception.GetSDKRes(exception.INVALID_DATAVALUE_ERROR)
		resData.ErrorCode = SDKRes.ErrorCode
		resData.ErrorDesc = SDKRes.ErrorDesc
		return resData
	}

	if reqData.GetPrivateKey() == "" {
		SDKRes := exception.GetSDKRes(exception.PRIVATEKEY_NULL_ERROR)
		resData.ErrorCode = SDKRes.ErrorCode
		resData.ErrorDesc = "PrivateKey cannot be empty"
		return resData
	}

	// 生成时间戳（毫秒级）
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)

	// 从私钥获取公钥
	publicKey, err := keypair.GetEncPublicKey(reqData.GetPrivateKey())
	if err != nil {
		SDKRes := exception.GetSDKRes(exception.PRIVATEKEY_ONE_ERROR)
		resData.ErrorCode = SDKRes.ErrorCode
		resData.ErrorDesc = "Invalid private key: " + err.Error()
		return resData
	}

	// 生成签名数据：content + timestamp
	signData := reqData.GetContent() + strconv.FormatInt(timestamp, 10)
	// 使用私钥对签名数据进行签名
	signatureResult, err := signature.Sign(reqData.GetPrivateKey(), []byte(signData))
	if err != nil {
		SDKRes := exception.GetSDKRes(exception.SIGN_ERROR)
		resData.ErrorCode = SDKRes.ErrorCode
		resData.ErrorDesc = "Signature failed: " + err.Error()
		return resData
	}

	// 构建请求体
	requestBody := map[string]interface{}{
		"key":       reqData.GetKey(),
		"content":   reqData.GetContent(),
		"timestamp": timestamp,
		"signature": signatureResult, // 已经是十六进制字符串
		"publicKey": publicKey,
	}

	// 转换为JSON
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		SDKRes := exception.GetSDKRes(exception.SYSTEM_ERROR)
		resData.ErrorCode = SDKRes.ErrorCode
		resData.ErrorDesc = "JSON marshal failed: " + err.Error()
		return resData
	}

	// 发送POST请求到存证接口
	post := "/evidence/store"
	response, SDKRes := common.PostRequest(evidence.Url, post, jsonBody)
	if SDKRes.ErrorCode != 0 {
		resData.ErrorCode = SDKRes.ErrorCode
		resData.ErrorDesc = SDKRes.ErrorDesc
		return resData
	}
	defer response.Body.Close()

	// 解析响应
	if response.StatusCode == 200 {
		decoder := json.NewDecoder(response.Body)
		err := decoder.Decode(&resData)
		if err != nil {
			SDKRes := exception.GetSDKRes(exception.SYSTEM_ERROR)
			resData.ErrorCode = SDKRes.ErrorCode
			resData.ErrorDesc = "Response decode failed: " + err.Error()
			return resData
		}

		// 检查API返回的错误代码
		if resData.ErrorCode != 0 {
			return resData
		}
	} else {
		// 尝试解析API返回的错误响应
		var errorResp struct {
			ErrorCode int    `json:"error_code"`
			ErrorDesc string `json:"error_desc"`
		}

		decoder := json.NewDecoder(response.Body)
		err := decoder.Decode(&errorResp)
		if err == nil && errorResp.ErrorCode != 0 {
			// API返回了结构化的错误
			resData.ErrorCode = errorResp.ErrorCode
			resData.ErrorDesc = errorResp.ErrorDesc
			return resData
		}

		// 如果不是结构化错误，读取响应内容以便调试
		bodyBytes := make([]byte, 1024)
		n, _ := response.Body.Read(bodyBytes)
		bodyStr := string(bodyBytes[:n])

		SDKRes := exception.GetSDKRes(exception.CONNECTNETWORK_ERROR)
		resData.ErrorCode = SDKRes.ErrorCode
		resData.ErrorDesc = fmt.Sprintf("HTTP %d: %s", response.StatusCode, bodyStr)
		return resData
	}

	return resData
}

// Get 查询存证
func (evidence *EvidenceOperation) Get(reqData model.EvidenceGetRequest) model.EvidenceGetResponse {
	var resData model.EvidenceGetResponse

	// 验证请求参数
	if reqData.GetHash() == "" {
		SDKRes := exception.GetSDKRes(exception.INVALID_HASH_ERROR)
		resData.ErrorCode = SDKRes.ErrorCode
		resData.ErrorDesc = "Hash cannot be empty"
		return resData
	}

	// 构建请求URL
	get := "/evidence/" + reqData.GetHash()

	// 发送GET请求
	response, SDKRes := common.GetRequest(evidence.Url, get, "")
	if SDKRes.ErrorCode != 0 {
		resData.ErrorCode = SDKRes.ErrorCode
		resData.ErrorDesc = SDKRes.ErrorDesc
		return resData
	}
	defer response.Body.Close()

	// 解析响应
	if response.StatusCode == 200 {
		decoder := json.NewDecoder(response.Body)
		err := decoder.Decode(&resData)
		if err != nil {
			SDKRes := exception.GetSDKRes(exception.SYSTEM_ERROR)
			resData.ErrorCode = SDKRes.ErrorCode
			resData.ErrorDesc = "Response decode failed: " + err.Error()
			return resData
		}

		// 检查API返回的错误代码
		if resData.ErrorCode != 0 {
			return resData
		}
	} else {
		// 尝试解析API返回的错误响应
		var errorResp struct {
			ErrorCode int    `json:"error_code"`
			ErrorDesc string `json:"error_desc"`
		}

		decoder := json.NewDecoder(response.Body)
		err := decoder.Decode(&errorResp)
		if err == nil && errorResp.ErrorCode != 0 {
			// API返回了结构化的错误
			resData.ErrorCode = errorResp.ErrorCode
			resData.ErrorDesc = errorResp.ErrorDesc
			return resData
		}

		// 如果不是结构化错误，读取响应内容以便调试
		bodyBytes := make([]byte, 1024)
		n, _ := response.Body.Read(bodyBytes)
		bodyStr := string(bodyBytes[:n])

		SDKRes := exception.GetSDKRes(exception.CONNECTNETWORK_ERROR)
		resData.ErrorCode = SDKRes.ErrorCode
		resData.ErrorDesc = fmt.Sprintf("HTTP %d: %s", response.StatusCode, bodyStr)
		return resData
	}

	return resData
}
