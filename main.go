package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"gameServers/config"
	"gameServers/utils"
	"github.com/buger/jsonparser"
	"github.com/go-redis/redis"
	"github.com/go-resty/resty/v2"
	"golang.org/x/net/websocket"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)
const max_room_num = 2
var roomInfo = make(map[string][]string)
var userguessInfo = make(map[string]UserGuessInfo)
var gameroomInfo  = make(map[string]GameInfo)
var Addr = config.Config.GetString("Address")
var PrivateKey = config.Config.GetString("PrivateKey")
var ViewKey = config.Config.GetString("ViewKey")
var ApiUrl = config.Config.GetString("ApiUrl")
var Contract = config.Config.GetString("Contract")

const Query = "https://vm.aleo.org/api"
const Broadcast = "https://vm.aleo.org/api/testnet3/transaction/broadcast"
const Prefix = "at"
const FEE = "100000"

var (
	JSON          = websocket.JSON              // codec for JSON
	Message       = websocket.Message           // codec for string, []byte
	ActiveClients = make(map[string]ClientConn) // map containing clients  //在线websocket列表
	User          = make(map[string]string)
)

type ClientConn struct {
	websocket *websocket.Conn
}

type OutPut struct {
	Owner string
	Gates string
	Game_id string
	Player_a string
	Player_b string
}

type UserMsg struct {
	Room string
	Cmd string
	User string
	AvatarUrl string
	Content string
	Uuid string
	HandNum string
	GuessNum string
}

type UserInfo struct {
	User string
	AvatarUrl string
	Uuid string
}

type ReplyMsg struct {
	Room string
	Cmd string
	Data string
	Uuid string
}

type GuessResult struct {
	Room string
	Uuid string
	Result string
}

type UserGuessInfo  struct {
	Uuid string
	HandNum int
	WinCount int
	LoseCount int
	Win      string    //0-lost, 1-win, 2-draw
	Tid      string    //Transaction id

}

type GameInfo struct {
	Addr string
	Gates string
	Game_id string
	Player_a string
	Player_b string
}

// GetLatestFeeRecord 获取fee transition outputs 0 value
func GetLatestFeeRecord() (string, error) {
	id, err := utils.GetId()
	if err != nil {
		return "", err
	}

	client := resty.New()
	resp, err := client.R().
		EnableTrace().
		Get(ApiUrl + id)
	if err != nil {
		log.Println("获取最新record数据 error: ", err)
		return "", nil
	}

	ciphertext, err := jsonparser.GetString(resp.Body(), "fee", "transition", "outputs", "[0]", "value")
	if err != nil {
		log.Println("获取最新record数据 value error: ", err)
		return "", err
	}

	record, err := DecryptCiphertext(ciphertext)
	if err != nil {
		log.Println("获取最新record数据 进行解密 error: ", err)
		return "", err
	}
	log.Println("获取最新record数据: ", record)
	return record, nil
}

func DecryptCiphertext(ciphertext string) (string, error) {
	cmd := "snarkos"

	args := []string{
		"developer",
		"decrypt",
		"--ciphertext",
		ciphertext,
		"--view-key",
		ViewKey}
	log.Println("args: ", args)
	record, err := utils.ExecCmdWithTimeout(60, cmd, args...)
	if err != nil {
		log.Println("DecryptCiphertext err:", err, record)
		log.Println("result:", record)
		return "", err
	}
	log.Println("DecryptCiphertext record:", strings.TrimSpace(record))
	return strings.TrimSpace(record), nil
}

// GetExecOutputValue  Get execution transitions value
func GetExecOutputValue(id string) (string, error) {
	client := resty.New()
	resp, err := client.R().
		EnableTrace().
		Get(ApiUrl + id)
	if err != nil {
		log.Println("send http request, get output->value data error: ", err)
		return "", nil
	}
	cipherText, err := jsonparser.GetString(resp.Body(), "execution", "transitions", "[0]", "outputs", "[0]", "value")
	if err != nil {
		log.Println("get json output->value data error: ", err)
		return "", err
	}
	log.Println("get output->value data cipherText: ", cipherText)
	return cipherText, nil
}


func create_game(record string, room string, uuid string) (string, error) {
	id := ""

	addr := Addr
	game_id := room + "field"
	player := uuid + "field"
	args := []string{
		"developer",
		"execute",
		Contract,
		"create_game",
		addr,
		game_id,
		player,
		"--private-key",
		PrivateKey,
		"--query",
		Query,
		"--broadcast",
		Broadcast,
		"--fee",
		FEE,
		"--record",
		record}
	log.Println("Execute contract create_game data args:", args)

	cmd := "snarkos"
	result, err := utils.ExecCmdWithTimeout(60*5, cmd, args...)
	if err != nil {
		log.Println("Execute contract create_game data err:", err)
		log.Println("Execute contract create_game data result:", result)
		return id, err
	}
	log.Println("Execute contract create_game data result:", result)

	split := strings.Split(strings.TrimSpace(result), "\n")
	id = split[len(split)-1]
	if !strings.HasPrefix(id, Prefix) {
		log.Println("Execution error. Incorrect data was obtained, id: ", id)
		return "", errors.New("Execution error. Incorrect data was obtained")
	}

	log.Println("Execute contract submission data id: ", id)

	err = utils.SetId(id)
	if err != nil {
		log.Println("utils.SetId err:", err)
		return id, err
	}

	gameInfo := GameInfo{}
	gameInfo.Addr = Addr
	gameInfo.Gates = "0u64"
	gameInfo.Game_id = game_id
	gameInfo.Player_a = player
	gameInfo.Player_b = "0field"
	gameroomInfo[room] = gameInfo


	return id, nil
}

func join_game(record string, room string, uuid string, info GameInfo)(string, error){
	id := ""

	player := uuid + "field"
	params := `{addr: {{addr}}, gates: {{gates}}, game_id: {{game_id}}, player_a: {{player_a}}, player_b: {{player_b}}}`
	params = strings.Replace(params, "{{addr}}",info.Addr, -1)
	params = strings.Replace(params, "{{gates}}", info.Gates, -1)
	params = strings.Replace(params, "{{game_id}}", info.Game_id, -1)
	params = strings.Replace(params, "{{player_a}}", info.Player_a, -1)
	params = strings.Replace(params, "{{player_b}}", info.Player_b, -1)
	args := []string{
		"developer",
		"execute",
		Contract,
		"join_game",
		params,
		player,
		"--private-key",
		PrivateKey,
		"--query",
		Query,
		"--broadcast",
		Broadcast,
		"--fee",
		FEE,
		"--record",
		record}
	log.Println("Execute contract join_game data args:", args)

	cmd := "snarkos"
	result, err := utils.ExecCmdWithTimeout(60*5, cmd, args...)
	if err != nil {
		log.Println("Execute contract join_game data err:", err)
		log.Println("Execute contract join_game data result:", result)
		return id, err
	}
	log.Println("Execute contract join_game data result:", result)

	split := strings.Split(strings.TrimSpace(result), "\n")
	id = split[len(split)-1]
	if !strings.HasPrefix(id, Prefix) {
		log.Println("Execution error. Incorrect data was obtained, id: ", id)
		return "", errors.New("Execution error. Incorrect data was obtained")
	}

	log.Println("Execute contract join_game data id: ", id)

	err = utils.SetId(id)
	if err != nil {
		log.Println("utils.SetId err:", err)
		return id, err
	}

	gameInfo := GameInfo{}
	gameInfo.Addr = info.Addr
	gameInfo.Gates = info.Gates
	gameInfo.Game_id = info.Game_id
	gameInfo.Player_a = info.Player_a
	gameInfo.Player_b = player
	gameroomInfo[room] = gameInfo

	return id, nil
}


func start_game(record string, room string, uuid string, info GameInfo, handNum string)(string, error){
	id := ""

	player := uuid + "field"
	player_id := room + "field"
	hand_num := handNum + "u8"
	params := `{addr: {{addr}}, gates: {{gates}}, game_id: {{game_id}}, player_a: {{player_a}}, player_b: {{player_b}}}`
	params = strings.Replace(params, "{{addr}}",info.Addr, -1)
	params = strings.Replace(params, "{{gates}}", info.Gates, -1)
	params = strings.Replace(params, "{{game_id}}", info.Game_id, -1)
	params = strings.Replace(params, "{{player_a}}", info.Player_a, -1)
	params = strings.Replace(params, "{{player_b}}", info.Player_b, -1)
	args := []string{
		"developer",
		"execute",
		Contract,
		"start_game",
		params,
		player_id,
		player,
		hand_num,
		"--private-key",
		PrivateKey,
		"--query",
		Query,
		"--broadcast",
		Broadcast,
		"--fee",
		FEE,
		"--record",
		record}
	log.Println("Execute contract start_game data args:", args)

	cmd := "snarkos"
	result, err := utils.ExecCmdWithTimeout(60*5, cmd, args...)
	if err != nil {
		log.Println("Execute contract start_game data err:", err)
		log.Println("Execute contract start_game data result:", result)
		return id, err
	}
	log.Println("Execute contract start_game data result:", result)

	split := strings.Split(strings.TrimSpace(result), "\n")
	id = split[len(split)-1]
	if !strings.HasPrefix(id, Prefix) {
		log.Println("Execution error. Incorrect data was obtained, id: ", id)
		return "", errors.New("Execution error. Incorrect data was obtained")
	}

	log.Println("Execute contract start_game data id: ", id)

	err = utils.SetId(id)
	if err != nil {
		log.Println("utils.SetId err:", err)
		return id, err
	}


	return id, nil
}

func finish_game(record string, id1 string, id2 string)(string, error){

	id := ""
	ciphertext1, err := GetExecOutputValue(id1)
	if err != nil {
		return id1, err
	}

	// 获取计算健康数据的入参
	value1, err := DecryptCiphertext(ciphertext1)
	if err != nil {
		return id1, err
	}

	ciphertext2, err := GetExecOutputValue(id2)
	if err != nil {
		return id1, err
	}

	// 获取计算健康数据的入参
	value2, err := DecryptCiphertext(ciphertext2)
	if err != nil {
		return id1, err
	}

	args := []string{
		"developer",
		"execute",
		Contract,
		"finish_game",
		value1,
		value2,
		"--private-key",
		PrivateKey,
		"--query",
		Query,
		"--broadcast",
		Broadcast,
		"--fee",
		FEE,
		"--record",
		record}
	log.Println("Execute contract finish_game data args:", args)

	cmd := "snarkos"
	result, err := utils.ExecCmdWithTimeout(60*5, cmd, args...)
	if err != nil {
		log.Println("Execute contract finish_game data err:", err)
		log.Println("Execute contract finish_game data result:", result)
		return id, err
	}
	log.Println("Execute contract finish_game data result:", result)

	split := strings.Split(strings.TrimSpace(result), "\n")
	id = split[len(split)-1]
	if !strings.HasPrefix(id, Prefix) {
		log.Println("Execution error. Incorrect data was obtained, id: ", id)
		return "", errors.New("Execution error. Incorrect data was obtained")
	}

	log.Println("Execute contract finish_game data id: ", id)

	err = utils.SetId(id)
	if err != nil {
		log.Println("utils.SetId err:", err)
		return id, err
	}

	time.Sleep(time.Duration(1000*18) * time.Millisecond)

	ciphertext3, err := GetExecOutputValue(id)
	if err != nil {
		return id, err
	}

	value3, err := DecryptCiphertext(ciphertext3)
	if err != nil {
		return id1, err
	}

	log.Println("Execute contract finish_game data value3: ", value3)
	records := strings.Split(value3, "\n")
	isWin:=""
	for i := 0; i < len(records); i++ {
		if strings.HasPrefix(strings.TrimSpace(records[i]), "player:") {
			temp := strings.TrimSpace(strings.Replace(strings.Split(strings.TrimSpace(records[i]), ":")[1], "field.private,", "", -1))
			fmt.Println("temp1:", temp)
		} else if strings.HasPrefix(strings.TrimSpace(records[i]), "is_winner:") {
			isWin = strings.TrimSpace(strings.Replace(strings.Split(strings.TrimSpace(records[i]), ":")[1], "u8.private,", "", -1))
		}
	}

	return isWin, nil
}

func echoHandler(ws *websocket.Conn) {
	var err error
	var userMsg UserMsg

	for {

		var data []byte
		if err = websocket.Message.Receive(ws, &data); err != nil {
			//fmt.Println("can't receive:", err)
			break
		}

		fmt.Println("data:", string(data))
		err = json.Unmarshal(data, &userMsg)
		fmt.Println(userMsg)

		go wsHandler(ws,userMsg)

	}

}


func wsHandler(ws *websocket.Conn,userMsg UserMsg) {
	sockCli := ClientConn{ws}
	var err error


	redisClient := redis.NewClient(&redis.Options{
		Addr:     "172.18.10.6:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})


	//登录
	if userMsg.Cmd == "login" {
		fmt.Println("login")
		fmt.Println("usermsg:", (userMsg))
		//判断房间人数是否已满
		userId, _ :=roomInfo[userMsg.Room]
		if len(userId) < max_room_num{

			if len(userId)==0 {
				//First player,create_game
				// 获取最新record数据
				record, err1 := GetLatestFeeRecord()
				if err1 != nil {
					log.Println("err: ", err1)
					return
				}
				create_game(record, userMsg.Room, userMsg.Uuid)
			}else{
				//second player,join_game
				record, err1 := GetLatestFeeRecord()
				if err1 != nil {
					log.Println("err: ", err1)
					return
				}
				join_game(record, userMsg.Room, userMsg.Uuid, gameroomInfo[userMsg.Room])
			}

			//socket用户列表新增当前用户websocket连接
			ActiveClients[userMsg.Uuid] = sockCli

			var rm ReplyMsg
			rm.Room = userMsg.Room
			rm.Cmd = "loginsucess"

			var onlineList []string
			onlineList = append(onlineList, userMsg.Uuid)
			onlineListStr,_ := json.Marshal(onlineList)
			fmt.Println("uuid:", userMsg.Uuid)
			rm.Uuid = string(onlineListStr)

			sendMsg,err2 := json.Marshal(rm)
			sendMsgStr := string(sendMsg)
			fmt.Println("sendmsg:", sendMsgStr)
			if err2 == nil {
				if err = websocket.Message.Send(ws, sendMsgStr); err != nil {
					log.Println("Could not send UsersList to ", userMsg.Uuid, err.Error())
				}
			}

			//初始化用户
			uuid := roomInfo[userMsg.Room]
			uuid = append(uuid, userMsg.Uuid)
			roomInfo[userMsg.Room] = uuid
			initOnlineMsg(redisClient,userMsg)

		}else {

			//The number of people is full
			var rm ReplyMsg
			rm.Room = userMsg.Room
			rm.Cmd = "loginFailed"
			rm.Data = "The number of people is full"

			sendMsg,err2 := json.Marshal(rm)
			sendMsgStr := string(sendMsg)
			fmt.Println(sendMsgStr)
			if err2 != nil {

			} else {
				if err = websocket.Message.Send(ws, sendMsgStr); err != nil {
					log.Println("Could not send UsersList to ", userMsg.User, err.Error())
				}
			}

		}

	} else if userMsg.Cmd == "logout" {
		fmt.Println("logout")

		//socket用户列表删除该用户websocket连接
		delete(ActiveClients,userMsg.Uuid)
		//从redis房间set集合内删除该用户uuid
		redisClient.SRem("ROOM:"+userMsg.Room,userMsg.Uuid)

		//初始化用户
		initOnlineMsg(redisClient,userMsg)


		//出拳
	} else if userMsg.Cmd == "guess" {

		fmt.Println("guess")
		fmt.Println(userMsg.HandNum)

		record, err1 := GetLatestFeeRecord()
		if err1 != nil {
			log.Println("err: ", err1)
			return
		}
		id, _ := start_game(record, userMsg.Room, userMsg.Uuid, gameroomInfo[userMsg.Room], userMsg.HandNum)

		myHandNum,_ := strconv.Atoi(userMsg.HandNum)

		guessInfo := UserGuessInfo{}
		guessInfo.HandNum = myHandNum
		guessInfo.Uuid = userMsg.Uuid
		guessInfo.Tid = id
		userguessInfo[userMsg.Uuid] = guessInfo


		online := roomInfo[userMsg.Room]
		var onlineList []string
		fmt.Println("get online success")
		i := 0
		//循环取在线用户
		if len(online) != 0 {
			for _, na := range online {
				fmt.Println("online uuid:", na)
				if na != "" {
					onlineList = append(onlineList, na)
					tmpGuessInfo := userguessInfo[na]
					handnum := tmpGuessInfo.HandNum
					if handnum > 0 {
						i++
					}
				}
			}
		}

		onlineListStr,_ := json.Marshal(onlineList)
		//房间内所有人都已提交，则计算最后结果
		if i == len(online) && i == max_room_num {
			fmt.Println("return game result:")

			tmpGuessInfo1 := userguessInfo[online[0]]
			tmpGuessInfo2 := userguessInfo[online[1]]
			time.Sleep(time.Duration(1000*18) * time.Millisecond)

			//var record string
			record, err1 := GetLatestFeeRecord()
			if err1 != nil {
				log.Println("err1: ", err1)
				time.Sleep(time.Duration(1000*5) * time.Millisecond)
				record, err1 = GetLatestFeeRecord()
				if err1 != nil {
					log.Println("err1: ", err1)
					return
				}
			}

			isWin, gameerr:=finish_game(record, tmpGuessInfo1.Tid, tmpGuessInfo2.Tid)
			fmt.Println("return game result isWin:", isWin, "gameerr:", gameerr)

			if isWin == "0" {
				tmpGuessInfo2.Win = "2"
				tmpGuessInfo1.Win = "2"
			}else if isWin == "1" {
				tmpGuessInfo2.Win = "1"
				tmpGuessInfo1.Win = "0"
				tmpGuessInfo2.WinCount+=1
				tmpGuessInfo1.LoseCount+=1
			}else if isWin == "2"{
				tmpGuessInfo1.Win = "1"
				tmpGuessInfo2.Win = "0"
				tmpGuessInfo1.WinCount+=1
				tmpGuessInfo2.LoseCount+=1
			}

			var resultList []string
			data1, _ := json.Marshal(tmpGuessInfo1)
			data2, _ := json.Marshal(tmpGuessInfo2)
			resultList = append(resultList, string(data1))
			resultList = append(resultList, string(data2))
			resultListStr,_ := json.Marshal(resultList)
			fmt.Println("resultList", string(resultListStr))

			tmpGuessInfo1.HandNum = 0
			tmpGuessInfo2.HandNum = 0
			userguessInfo[online[0]] = tmpGuessInfo1
			userguessInfo[online[1]] = tmpGuessInfo2
			for _, na := range online {
				if na != "" {

					//guessInfo.HandNum = 0
					//userguessInfo[na] = guessInfo
					var rm ReplyMsg
					rm.Room = userMsg.Room
					rm.Cmd = "result"
					rm.Data = string(resultListStr)
					rm.Uuid = string(onlineListStr)

					sendMsg,_ := json.Marshal(rm)
					sendMsgStr := string(sendMsg)
					fmt.Println("result msg", sendMsgStr)
					if err = websocket.Message.Send(ActiveClients[na].websocket, sendMsgStr); err != nil {
						log.Println("Could not send UsersList to ", "", err.Error())
					}

				}
			}
		}

		//发消息
	} else {

	}
}

//房间成员初始化,有人加入或者退出都要重新初始化，相当于聊天室的在线用户列表的维护
func initOnlineMsg(redisClient *redis.Client,userMsg UserMsg) {

	var err error
	var online []string
	if userId, ok :=roomInfo[userMsg.Room]; ok{
		online = userId
	}
	var onlineList []string

	//循环取在线用户个人信息
	if len(online) != 0 {
		for _, na := range online {
			if na != "" {
				onlineList = append(onlineList, na)
			}
		}
	}
	fmt.Println("get online success")
	//生成在线用户信息json串
	//c, err := json.Marshal(onlineList)

	onlineListStr,err2 := json.Marshal(onlineList)

	var rm ReplyMsg
	rm.Room = userMsg.Room
	rm.Cmd = "init"
	rm.Uuid = string(onlineListStr)

	sendMsg,err2 := json.Marshal(rm)
	sendMsgStr := string(sendMsg)
	fmt.Println("sendMsgStr1111:", sendMsgStr)
	if err2 != nil {

	} else {
		//给所有用户发初始化消息
		if len(online) != 0 {
			for _, na := range online {
				if na != "" {
					if err = websocket.Message.Send(ActiveClients[na].websocket, sendMsgStr); err != nil {
						log.Println("Could not send UsersList to ", "", err.Error())
					}
				}
			}
		}
		//若房间人数满，发送就绪消息
		if len(online) >= max_room_num {
			fmt.Println("full")
			var rm ReplyMsg
			rm.Room = userMsg.Room
			rm.Cmd = "full"
			rm.Uuid = string(onlineListStr)

			sendMsg,_ := json.Marshal(rm)
			sendMsgStr := string(sendMsg)
			fmt.Println("full sendmsg", sendMsgStr)
			for _, na := range online {
				if na != "" {
					if err = websocket.Message.Send(ActiveClients[na].websocket, sendMsgStr); err != nil {
						log.Println("Could not send UsersList to ", "", err.Error())
					}
				}
			}
		}
	}

}


func main() {
	fmt.Println("test games start")
	fmt.Println("addr:", Addr)
	http.Handle("/echo", websocket.Handler(echoHandler))
	http.Handle("/", http.FileServer(http.Dir(".")))

	err := http.ListenAndServe(":8929", nil)
	if err != nil {
		log.Println("ListenAndServe err:", err)
	}

}