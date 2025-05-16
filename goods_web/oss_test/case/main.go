package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"hash"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aliyun/credentials-go/credentials"
)

// 定义全局变量
var (
	region     string
	bucketName string
	product    = "oss"
)

// PolicyToken 结构体用于存储生成的表单数据
type PolicyToken struct {
	Policy           string `json:"policy"`
	SecurityToken    string `json:"security_token"`
	SignatureVersion string `json:"x_oss_signature_version"`
	Credential       string `json:"x_oss_credential"`
	Date             string `json:"x_oss_date"`
	Signature        string `json:"signature"`
	Host             string `json:"host"`
	Dir              string `json:"dir"`
	Callback         string `json:"callback"`
}

type CallbackParam struct {
	CallbackUrl      string `json:"callbackUrl"`
	CallbackBody     string `json:"callbackBody"`
	CallbackBodyType string `json:"callbackBodyType"`
}

func main() {
	// 定义默认的IP和端口字符串
	strIPPort := ":8080"
	if len(os.Args) == 3 {
		strIPPort = fmt.Sprintf("%s:%s", os.Args[1], os.Args[2])
	} else if len(os.Args) != 1 {
		fmt.Println("Usage   : go run test1.go                ")
		fmt.Println("Usage   : go run test1.go ip port        ")
		fmt.Println("")
		os.Exit(0)
	}
	// 打印服务器运行的地址和端口
	fmt.Printf("server is running on %s \n", strIPPort)
	// 注册处理根路径请求的函数
	http.HandleFunc("/", handlerRequest)
	// 注册处理获取签名请求的函数
	http.HandleFunc("/get_post_signature_for_oss_upload", handleGetPostSignature)
	// 启动HTTP服务器
	err := http.ListenAndServe(strIPPort, nil)
	if err != nil {
		strError := fmt.Sprintf("http.ListenAndServe failed : %s \n", err.Error())
		panic(strError)
	}
}

// handlerRequest 函数处理根路径请求
func handlerRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		http.ServeFile(w, r, "templates/index.html")
		return
	}
	http.NotFound(w, r)
}

// handleGetPostSignature 函数处理获取签名请求
func handleGetPostSignature(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		response := getPolicyToken()
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*") // 允许跨域
		w.Write([]byte(response))
		return
	}
	http.NotFound(w, r)
}

// getPolicyToken 函数生成 OSS 上传所需的签名和凭证
func getPolicyToken() string {
	// 设置bucket所处地域
	region = "cn-hangzhou"
	// 设置bucket名称
	bucketName = "examplebucket"
	// 设置 OSS 上传地址
	host := fmt.Sprintf("https://%s.oss-%s.aliyuncs.com", bucketName, region)
	// 设置上传目录
	dir := "user-dir"
	// callbackUrl为 上传回调服务器的URL，请将下面的IP和Port配置为您自己的真实信息。
	callbackUrl := "http://oss-demo.aliyuncs.com:23450/callback"

	config := new(credentials.Config).
		SetType("ram_role_arn").
		SetAccessKeyId(os.Getenv("OSS_ACCESS_KEY_ID")).
		SetAccessKeySecret(os.Getenv("OSS_ACCESS_KEY_SECRET")).
		SetRoleArn(os.Getenv("OSS_STS_ROLE_ARN")).
		SetRoleSessionName("Role_Session_Name").
		SetPolicy("").
		SetRoleSessionExpiration(3600)

	// 根据配置创建凭证提供器
	provider, err := credentials.NewCredential(config)
	if err != nil {
		log.Fatalf("NewCredential fail, err:%v", err)
	}

	// 从凭证提供器获取凭证
	cred, err := provider.GetCredential()
	if err != nil {
		log.Fatalf("GetCredential fail, err:%v", err)
	}

	// 构建policy
	utcTime := time.Now().UTC()
	date := utcTime.Format("20060102")
	expiration := utcTime.Add(1 * time.Hour)
	policyMap := map[string]any{
		"expiration": expiration.Format("2006-01-02T15:04:05.000Z"),
		"conditions": []any{
			map[string]string{"bucket": bucketName},
			map[string]string{"x-oss-signature-version": "OSS4-HMAC-SHA256"},
			map[string]string{"x-oss-credential": fmt.Sprintf("%v/%v/%v/%v/aliyun_v4_request", *cred.AccessKeyId, date, region, product)},
			map[string]string{"x-oss-date": utcTime.Format("20060102T150405Z")},
			map[string]string{"x-oss-security-token": *cred.SecurityToken},
		},
	}

	// 将policy转换为 JSON 格式
	policy, err := json.Marshal(policyMap)
	if err != nil {
		log.Fatalf("json.Marshal fail, err:%v", err)
	}

	// 构造待签名字符串（StringToSign）
	stringToSign := base64.StdEncoding.EncodeToString([]byte(policy))

	hmacHash := func() hash.Hash { return sha256.New() }
	// 构建signing key
	signingKey := "aliyun_v4" + *cred.AccessKeySecret
	h1 := hmac.New(hmacHash, []byte(signingKey))
	io.WriteString(h1, date)
	h1Key := h1.Sum(nil)

	h2 := hmac.New(hmacHash, h1Key)
	io.WriteString(h2, region)
	h2Key := h2.Sum(nil)

	h3 := hmac.New(hmacHash, h2Key)
	io.WriteString(h3, product)
	h3Key := h3.Sum(nil)

	h4 := hmac.New(hmacHash, h3Key)
	io.WriteString(h4, "aliyun_v4_request")
	h4Key := h4.Sum(nil)

	// 生成签名
	h := hmac.New(hmacHash, h4Key)
	io.WriteString(h, stringToSign)
	signature := hex.EncodeToString(h.Sum(nil))

	var callbackParam CallbackParam
	callbackParam.CallbackUrl = callbackUrl
	callbackParam.CallbackBody = "filename=${object}&size=${size}&mimeType=${mimeType}&height=${imageInfo.height}&width=${imageInfo.width}"
	callbackParam.CallbackBodyType = "application/x-www-form-urlencoded"
	callback_str, err := json.Marshal(callbackParam)
	if err != nil {
		fmt.Println("callback json err:", err)
	}
	callbackBase64 := base64.StdEncoding.EncodeToString(callback_str)
	// 构建返回给前端的表单
	policyToken := PolicyToken{
		Policy:           stringToSign,
		SecurityToken:    *cred.SecurityToken,
		SignatureVersion: "OSS4-HMAC-SHA256",
		Credential:       fmt.Sprintf("%v/%v/%v/%v/aliyun_v4_request", *cred.AccessKeyId, date, region, product),
		Date:             utcTime.UTC().Format("20060102T150405Z"),
		Signature:        signature,
		Host:             host,           // 返回 OSS 上传地址
		Dir:              dir,            // 返回上传目录
		Callback:         callbackBase64, // 返回上传回调参数
	}

	response, err := json.Marshal(policyToken)
	if err != nil {
		fmt.Println("json err:", err)
	}
	return string(response)
}
