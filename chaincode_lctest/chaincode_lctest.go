package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"

	"github.com/hyperledger/fabric/common/util"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

//define lc struct
type LCContent struct {
	LCNO                 string
	ISOURBANKLCISSUE     string
	BRANCHORGNO          string
	SENDMODE             string
	MAILTYPE             string
	MAILENO              string
	RECVBANKSWIFTCODE    string
	RECVBANKENNAME       string
	RECVBANKLARGENO      string
	RECVBANKNAME         string
	RECVBANKADDR         string
	ADVBANKSWIFTCODE     string
	ADVBANKLARGENO       string
	ADVBANKNAME          string
	ADVBANKADDR          string
	ADVBANKNAMEADDR      string
	APPBANKNO            string
	APPNAME              string
	APPADDR              string
	ISSUEDATE            string
	EXPIRYDATE           string
	EXPIRYPLACE          string
	LCCURSIGN            string
	LCAMT                float32
	LCAMTTOLERDOWN       float32
	LCAMTTOLERUP         float32
	LCMAXAMT             float32
	DRAFTINVPCT          float32
	LCFORM               string
	LCAVAILTYPE          string
	APPLICABLERULES      string
	DEFERPAYTYPE         string
	TENORTYPE            string
	DEFERPAYDEADLINE     string
	DEFERPAYDESC         string
	NEGOBANKSAID         string
	NEGOBANKSWIFTCODE    string
	NEGOBANKENADDR       string
	NEGOBANKLARGENO      string
	NEGOBANKCNNAME       string
	NEGOBANKCNADDR       string
	SENDBANKSWIFTCODE    string
	SENDBANKNAMEADDR     string
	BENEFNAME            string
	BENEFADDR            string
	BENEFBANKNAME        string
	BENEFACCTNO          string
	GOODSNAME            string
	PARTTIALSHIPMENT     string
	TRANSSHIPMENT        string
	TRANSPORTNAME        string
	LOADPORTNAME         string
	LATESTSHIPDATE       string
	TRANSPORTMODE        string
	LOADAIRPORTDEST      string
	DISCHAIRPORTDEST     string
	GOODSSERVDESCR       string
	DOCREQURED           string
	OTHERCLAUSES         string
	DEPOSITPCT           float32
	PAYEXPENSE           string
	PRESENTPERIOD        string
	ISTRANSFER           string
	TRANBANKSWIFTCODE    string
	TRANBANKNAMEADDR     string
	ISCONFIRMING         string
	CONFIRMBANKSWIFTCODE string
	CONFBANKNAMEADDR     string
	TRADETYPE            string
	ISINSTALLMENT        string
	ISMAYADD             string
	CONFCHRTAKER         string
	CONFBANKLARGENO      string
	CONFBANKCNNAME       string
	INSTALLMENTDESC      string
	MEMO                 string
}

type LCChaincode struct {
}

//初始化信用证合约的数据
func (t *LCChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("lctest Init() start ........")
	//定义接收数据
	var content string   //接收信用证初始化的内容
	var conStu LCContent //结构体内容
	var err error
	//获取方法名称和参数
	fun_name, fun_params := stub.GetFunctionAndParameters()
	//参数长度必须为4
	if len(fun_params) != 1 {
		return shim.Error("Incorrect number of arguments...")
	}
	//将整个信用证的数据变成内容结构体
	err = json.Unmarshal([]byte(fun_params[0]), &conStu)
	if err != nil {
		return shim.Error("convert json to lc struct error" + err.Error())
	}
	//打印参数
	fmt.Println(conStu)
	//将json数据保存到世界状态
	err = stub.PutState("content", json.Marshal(conStu))
	if err != nil {
		return shim.Error("convert lc struct to json error" + err.Error())
	}
	fmt.Println("lctest Init() finish ........")
	return shim.Success(nil)
}

//智能合约调用的方法
func (t *LCChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("lctest Invoke() start ........")
	//获取方法名称和参数
	fun_name, fun_params := stub.GetFunctionAndParameters()
	if fun_name == "invoke" {
		return t.invoke(stub, fun_params)
	} else if fun_name == "query" {
		return t.query(stub, fun_params)
	}

	return shim.Error("Invalid invoke function name . Expecting [invoke,query]")
}

//触发程序的修改数据
func (t *LCChaincode) invoke(stub shim.ChaincodeStubInterface, params []string) pb.Response {
	fmt.Println("invoke lctest content start ........")
	var err error
	var conStu LCContent
	if len(params) != 1 {
		return shim.Error("Incorrect number of arguments...")
	}
	//取出信用证的内容
	content, err := stub.GetState("content")
	if err != nil {
		return shim.Error("Failed to get state content" + err.Error())
	}
	if content == nil {
		return shim.Error("LC content is null")
	}
	err = json.Unmarshal(content, &conStu)
	if err != nil {
		shim.Error("convert json to lc struct error" + err.Error())
	}
	//定义一个map处理数据
	var dataMap map[string]string
	dataMap = make(map[string]string)
	err = json.Unmarshal([]byte(params[0]), &dataMap)
	if err != nil {
		return shim.Error("params[0] is Error ,el:{\"LCAMT\":\"1000\"}")
	}
	//反射获取到结构体
	immutable := reflect.ValueOf(&conStu).Elem()
	for key, value := range dataMap {
		field := immutable.FieldByName(key)
		if !field.CanSet() {
			return shim.Error("can not update field " + key)
		}
		if field.Kind() == reflect.String {
			field.SetString(value)
		} else if field.Kind() == reflect.Float32 {
			flt, err := strconv.ParseFloat(value, 32)
			if err != nil {
				return shim.Error("covert string to float32 , field :" + key + "fieldval :" + value + err.Error())
			}
			field.SetFloat(flt)
		}
	}

}

//触发程序的查询数据
func (t *LCChaincode) query(stub shim.ChaincodeStubInterface, params []string) pb.Response {
	fmt.Println("query lctest content start ........")
	//定义变量
	var err error
	var conStu LCContent
	//取出信用证的内容
	content, err := stub.GetState("content")
	if err != nil {
		return shim.Error("Failed to get state content" + err.Error())
	}
	if content == nil {
		return shim.Error("LC content is null")
	}
	fmt.Println("query lctest content finish ........")
	return shim.Success(string(content))
}
func main() {
	var content LCContent
	//userJson := "{\"username\":\"system\",\"password\":123456}"

	contentJson := "{\"LCNO\":\"LC00001\",\"ISOURBANKLCISSUE\":\"Y\",\"BRANCHORGNO\":\"172001\",\"SENDMODE\":\"YJ\",\"MAILTYPE\":\"EMS\",\"MAILENO\":\"EMS001\",\"RECVBANKSWIFTCODE\":\"CITIUS33\",\"RECVBANKENNAME\":\"HUA QI BANK\",\"RECVBANKLARGENO\":\"QSHH1001\",\"RECVBANKNAME\":\"HUA QI BANK\",\"RECVBANKADDR\":\"HUA QI ADDR\",\"ADVBANKSWIFTCODE\":\"ICBCCNBK\",\"ADVBANKLARGENO\":\"DEHH001\",\"ADVBANKNAME\":\"TZHMC\",\"ADVBANKADDR\":\"TZHDZ\",\"ADVBANKNAMEADDR\":\"TZHMCDZ\",\"APPBANKNO\":\"SQRZH\",\"APPNAME\":\"SQRMC\",\"APPADDR\":\"SQRDZ\",\"ISSUEDATE\":\"20170523\",\"EXPIRYDATE\":\"20170623\",\"EXPIRYPLACE\":\"DQDD\",\"LCCURSIGN\":\"USD\",\"LCAMT\":4000000.0,\"LCAMTTOLERDOWN\":2.0,\"LCAMTTOLERUP\":2.0,\"LCMAXAMT\":4000000.0,\"DRAFTINVPCT\":2.0,\"LCFORM\":\"JQ\",\"LCAVAILTYPE\":\"YQCD\",\"APPLICABLERULES\":\"123\",\"DEFERPAYTYPE\":\"ZXZF\",\"TENORTYPE\":\"123\",\"DEFERPAYDEADLINE\":\"20170823\",\"DEFERPAYDESC\":\"333333333333\",\"NEGOBANKSAID\":\"11\",\"NEGOBANKSWIFTCODE\":\"ICBCCBBK\",\"NEGOBANKENADDR\":\"YFHYWMCDZ\",\"NEGOBANKLARGENO\":\"QSHH\",\"NEGOBANKCNNAME\":\"ZWMCDZ\",\"NEGOBANKCNADDR\":\"ZWMCDZ\",\"SENDBANKSWIFTCODE\":\"CITIUS33\",\"SENDBANKNAMEADDR\":\"FBHTWMCDZ\",\"BENEFNAME\":\"SYRMC\",\"BENEFADDR\":\"SYRDZ\",\"BENEFBANKNAME\":\"SYYHMC\",\"BENEFACCTNO\":\"NO123456\",\"GOODSNAME\":\"HWMC\",\"PARTTIALSHIPMENT\":\"N\",\"TRANSSHIPMENT\":\"N\",\"TRANSPORTNAME\":\"JHDD\",\"LOADPORTNAME\":\"SHDD\",\"LATESTSHIPDATE\":\"20170523\",\"TRANSPORTMODE\":\"CY\",\"LOADAIRPORTDEST\":\"LY\",\"DISCHAIRPORTDEST\":\"DL\",\"GOODSSERVDESCR\":\"2345\",\"DOCREQURED\":\"Y\",\"OTHERCLAUSES\":\"123456\",\"DEPOSITPCT\":3.0,\"PAYEXPENSE\":\"WF\",\"PRESENTPERIOD\":\"20\",\"ISTRANSFER\":\"N\",\"TRANBANKSWIFTCODE\":\"CITIUS33\",\"TRANBANKNAMEADDR\":\"ZRHDZ\",\"ISCONFIRMING\":\"N\",\"CONFIRMBANKSWIFTCODE\":\"CITIUS33\",\"CONFBANKNAMEADDR\":\"WERTYU\",\"TRADETYPE\":\"1\",\"ISINSTALLMENT\":\"Q\",\"ISMAYADD\":\"Q\",\"CONFCHRTAKER\":\"Q\",\"CONFBANKLARGENO\":\"123456\",\"CONFBANKCNNAME\":\"DSFDF\",\"INSTALLMENTDESC\":\"QWEQWEQWESDSADAD\",\"MEMO\":\"BEIZHU\"}"

	err := json.Unmarshal([]byte(contentJson), &content)
	if err != nil {
		fmt.Println("error ", err)
	}
	fmt.Printf(contentJson)
	fmt.Println(content) //打印结果:map[password:123456 username:system]

	var key string
	key = "LCNO"
	immutable := reflect.ValueOf(&content).Elem()
	field := immutable.FieldByName(key)
	fmt.Println(field.Kind())
	fmt.Println(content.LCNO)
	field1 := immutable.FieldByName("ABC")
	if err != nil {
		fmt.Println("999999")
	}
	if !field1.CanSet() {
		fmt.Println("000000")
	}

	immutable.FieldByName(key).SetString("ABCD1234")
	fmt.Println(content.LCNO)
	fmt.Println("%s", content.LCAMT)
}
