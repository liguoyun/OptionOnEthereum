package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"io/ioutil"
	"net/http"
	"strconv"
)

type User struct {
	ID      string   `json:"id,omitemty"`
	Name    string   `json:"fullname,omitempty"`
	Address *Address `json:"address,omitempty"`
}

type Address struct {
	PrivateKey    string `json:"privatekey,omitempty"`
	PublicAddress string `json:"publicaddress,omitempty"`
}

type TranferBody struct {
	FromId string `json:"fromid,omitempty"`
	ToId   string `json:"toid,omitempty"`
	Amount string `json:"amount,omitempty"`
}
type TranferRep struct {
	TxId string `json:"txid,omitempty"`
}
type StockShare struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"stock,omitempty"`
}

var users []User

func GetUserById(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	for _, item := range users {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(users)
}

func GetUsers(w http.ResponseWriter, req *http.Request) {
	json.NewEncoder(w).Encode(users)
}

func Hello(response http.ResponseWriter, request *http.Request) {
	tmpl, err := template.ParseFiles("./user.tpl")
	if err != nil {
		fmt.Println("Error happened..")
	}
	tmpl.Execute(response, "区块链股权转让演示")
}

func LoadUsers(jsonFile string) error {
	buf, err := ioutil.ReadFile(jsonFile) //read all into buf
	if err != nil {
		fmt.Println("error:", err)
	}
	err = json.Unmarshal(buf, &users)
	if err != nil {
		fmt.Println("error:", err)
	}
	return nil
}

func LoadStocks(jsonFile string) []StockShare {
	buf, err := ioutil.ReadFile(jsonFile) //read all into buf
	if err != nil {
		fmt.Println("error:", err)
	}
	var ss []StockShare
	err = json.Unmarshal(buf, &ss)
	if err != nil {
		fmt.Println("error:", err)
	}
	return ss
}

func GetStockShares(w http.ResponseWriter, req *http.Request) {
	fmt.Println("从区块链中查询最新股权数据")

	var reqBody TranferBody
	_ = json.NewDecoder(req.Body).Decode(&reqBody)

	var ss []StockShare

	//区块链获取数据
	contractor := NewStockTokenWapper(remote_api_address, eth_contract_address)
	for _, person := range users {
		addr := person.Address.PublicAddress
		banl, err := contractor.GetBalanceOfAddr(addr)
		if err == nil {
			ss = append(ss, StockShare{Name: person.Name, Value: banl})
		}
	}
	//测试数据
	if len(ss) == 0 {
		fmt.Println("区块链网络出现问题，返回最近一次成功的股权数据")
		ss = LoadStocks("stocks.json")
	}

	body, err := json.Marshal(ss)
	if err != nil {
		panic(err.Error())
	}
	ioutil.WriteFile("stocks.json", body, 0666) //写入文件(字节数组)
	w.Header().Set("Content-Length", strconv.Itoa(len(body)))
	w.Write(body) //发送数据
}

func GetPrivateKeyFromList(Id string) (string, string) {
	var result string
	var name string
	for _, p := range users {
		if p.ID == Id {
			result = p.Address.PrivateKey
			name = p.Name
		}
	}
	return result, name
}

func GetPubAddressFromList(Id string) (string, string) {
	var result string
	var name string
	for _, p := range users {
		if p.ID == Id {
			result = p.Address.PublicAddress
			name = p.Name
		}
	}
	return result, name
}

func PostTransfer(w http.ResponseWriter, req *http.Request) {
	var reqBody TranferBody
	_ = json.NewDecoder(req.Body).Decode(&reqBody)
	var rep TranferRep

	if reqBody.FromId != "" && reqBody.ToId != "" {
		priKey, fromN := GetPrivateKeyFromList(reqBody.FromId)
		pubAddr, toN := GetPubAddressFromList(reqBody.ToId)
		amount, _ := strconv.ParseInt(reqBody.Amount, 10, 64)
		fmt.Printf("从  %s  向 %s 转账 %d \n", fromN, toN, amount)
		contractor := NewStockTokenWapper(remote_api_address, eth_contract_address)
		rep.TxId, _ = contractor.TransferToken(priKey, pubAddr, amount)
	}
	body, err := json.Marshal(rep)
	if err != nil {
		fmt.Println("生成JSON出错", err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

func main() {
	LoadUsers("./users.json")
	router := mux.NewRouter()
	//设置Static目录
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	router.HandleFunc("/", Hello).Methods("GET")
	router.HandleFunc("/api/transfer", PostTransfer).Methods("POST")
	router.HandleFunc("/api/stockshares", GetStockShares).Methods("GET")
	router.HandleFunc("/api/users", GetUsers).Methods("GET")
	router.HandleFunc("/api/people/{id}", GetUserById).Methods("GET")

	http.ListenAndServe(":8080", router)
}
