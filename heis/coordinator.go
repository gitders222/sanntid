package main

import (
	"./fsm"
	heis "./heisdriver" //"./simulator/client"
	"encoding/json"
	"fmt"
	_ "os"
	"time"
)

//elevtype heis.ElevType = heis.ET_Simulation

const (
	MAX_NUM_ELEVS = 10
	N_FLOORS      = heis.N_FLOORS
	UP            = heis.BUTTON_CALL_UP
	DOWN          = heis.BUTTON_CALL_DOWN
	COMMAND       = heis.BUTTON_COMMAND
)

type GlobalOrderStruct struct {
	Available         [2][N_FLOORS]bool                  //'json:"Available"'
	Taken             [2][N_FLOORS]bool                  //'json:"Taken"'
	Timestamps        [2][N_FLOORS]uint                  //'json:"Timestamps"'
	Clock             uint                               //'json:"Clock"'
	Scores            [MAX_NUM_ELEVS][2][N_FLOORS]int    //'json:"Scores"'
	LocalOrdersBackup [MAX_NUM_ELEVS]fsm.LocalOrderState //'json:"LocalOrdersBackup"'
} //

var (
	GlobalOrders            GlobalOrderStruct
	//unverified_GlobalOrders GlobalOrderStruct
	LocalOrders             fsm.LocalOrderState

	online         bool
	localElevIndex int
	activeElevs    []int

	orderTimestamp int
)

/*********************************
Testing for network encoding
*********************************/
var str string

/*********************************
Main
*********************************/

func main() {

	// Init stuff
	/************************/
	localElevIndex = 0
	online = false
	activeElevs = make([]int, 0, MAX_NUM_ELEVS)

	// Sets all Taken values to -1. Not o because 0 can be a elevIndex
	for ordertype := UP; ordertype <= DOWN; ordertype++ {
		for floor := 0; floor < N_FLOORS; floor++ {
			GlobalOrders.Taken[ordertype][floor] = false
			unverified_GlobalOrders.Taken[ordertype][floor] = false
		}
	}
	/************************/

	orderChan := make(chan heis.Order, 5)
	eventChan := make(chan heis.Event, 5)
	fsmChan := make(chan fsm.LocalOrderState)

	/*******TEST DECODING ENCODING PACKET*********
	GlobalPacketENC := EncodeGlobalPacket()
	//fmt.Println(string(GlobalPacketENC))
	_ = GlobalPacketENC

	GlobalPacketDEC, err := DecodeGlobalPacket(GlobalPacketENC)
	fmt.Println("Test PacketDEC: ", GlobalPacketDEC.Taken)
	_ = err
	****************************************/

	/*****ADD*****
	networkChan := make(chan string)
	**************/
	heis.ElevInit()
	go heis.Poller(orderChan, eventChan)
	go fsm.Fsm(eventChan, fsmChan)
	/*****ADD*****
	go networkMonitor(network)

	**************/
	//timestamp := time.Now()
	for {
		// online, localElevIndex,  = networkStatus()
		select {
		case newOrder := <-orderChan:
			fmt.Println("Coord new order:", newOrder)
			
			if !online {
				//update local orders
				updateLocalOrders(newOrder)

				time.Sleep(time.Millisecond * 5)
				fsmChan <- LocalOrders
			}else{
				//update global avaliable
				updateGlobalOrders(newOrder)
				//save in buffer
			}
		case completedOrders := <-fsmChan:
			fmt.Println("Coord completed order")
			if !online{
				//update local complete
				updateLocalComplete()
				updateLocalState(completedOrders)

			}
			else{
				//update global complete
				updateGlobalComplete()
				//save in bufferout on network

			}

		/*********ADD*********
		case msg := <-networkChan:
			//decode msg
			//update global all

			GlobalOrders, err := DecodeGlobalPacket(msg)
			_ = err
			
			
			//if all scores done && you best score
				//update global taken
				//timestamp
				//delete all global scores for order
				
				//update local orders from global
				
			//else 
				updateLocalState()
				//set your score
				//update global scores
				
				//send msg	
			
	
		/*********************


			/******MAYBE********
			MSG = decode(msg)
			// replace with handle_msg()

			if MSG.header == "orders"{
				unverified_GlobalState = mergeOrders(LocalState, newGlobalState)
				//unverified_GlobalState = scoreOrders(unverified_GlobalState)
				msg = encode(unverified_GlobalState)
				networkChan<- msg
				status := <-networkchan
				if status == "SUCCESS" {
					GlobalState = unverified_GlobalState
				} else {
					//troubleshoot network
					online = false
				}
				updateLights()

			}
			updateLights()
			/**********************/

		default:
			//timeout error handeling
				//if global taken too long
					//move from taken to avaliable
				// 
			for ordertype := UP; ordertype <= DOWN; ordertype++ {
				for floor := 0; floor < N_FLOORS; floor++ {
					if time.Since(GlobalOrder.Timestamps[ordertype][floor]) > time.Second*10{
						GlobalOrder.Available[ordertype][floor] = true
						GlobalOrder.Taken[ordertype][floor] = false
					}

					if time.Since(LocalOrder.Timestamps[ordertype][floor]) > time.Second*10{
						GlobalOrder.Available[ordertype][floor] = true
						GlobalOrder.Taken[ordertype][floor] = false
						LocalOrders.Pending[ordertype][floor] = false

					}
				}
			}


			/********MAYBE********
			// COULD ALSO BE DONE THROUGH THE CHANNEL
			online, localElevIndex, activeElevs = network.getNetworkState()
			**********************/
			/********REMOVE*******
			//getNextOrder()
			if time.Since(timestamp) > time.Second*100 {
				//online = false
				// testing getNextOrder(). Works!
				fmt.Println("timer, global avalable: ", GlobalOrders.Available)
				fmt.Println("timer, global taken: ", GlobalOrders.Taken)
				fmt.Println(getNextOrder())
				timestamp = time.Now()
			}
			*********************/
			// set nextOrder?
			continue
		}
	}
}

/*********************************
Functions
Needed:
updateGlobalState()

*********************************/

func addNewOrder(order heis.Order) {
	ordertype := order.OrderType
	floor := order.Floor
	switch ordertype {
	case UP:
		if online {
			unverified_GlobalOrders.Available[ordertype][floor] = true
		} else {
			LocalOrders.Pending[ordertype][floor] = true
			LocalOrders.Completed[ordertype][floor] = false
		}
	case DOWN:
		if online {
			unverified_GlobalOrders.Available[ordertype][floor] = true
		} else {
			LocalOrders.Pending[ordertype][floor] = true
			LocalOrders.Completed[ordertype][floor] = false
		}
	case COMMAND:
		LocalOrders.Pending[ordertype][floor] = true
		LocalOrders.Completed[ordertype][floor] = false
	default:
		fmt.Println("Invalid OrderType in addNewOrder()")
	}
}

func updateLocalState(newLocalState fsm.LocalOrderState) {
	LocalOrders.PrevFloor = newLocalState.PrevFloor
	LocalOrders.Direction = newLocalState.Direction
	/*if online {
		LocalOrders.Completed = newLocalState.Completed
	} else {
		LocalOrders.Completed = {{0}}
	}
	for ordertype := UP; ordertype <= COMMAND; ordertype++ {
		for floor := 0; floor < N_FLOORS; floor++ {
			LocalOrders.Pending[ordertype][floor] = LocalOrders.Pending[ordertype][floor] && !(LocalOrders.Completed[ordertype][floor])
		}
	}*/
}

/**************sverre lør***************************/
func updateLocalOrders(order heis.Order){
	ordertype := order.OrderType
	floor := order.Floor
	stamp := time.Now()
	switch ordertype {

	case UP:
		LocalOrders.Pending[ordertype][floor] = true
		LocalOrders.Completed[ordertype][floor] = false
		LocalOrders.Timestamps[ordertype][floor] = stamp
	case DOWN:
	
		LocalOrders.Pending[ordertype][floor] = true
		LocalOrders.Completed[ordertype][floor] = false
		LocalOrders.Timestamps[ordertype][floor] = stamp


	case COMMAND:
		LocalOrders.Pending[ordertype][floor] = true
		LocalOrders.Completed[ordertype][floor] = false
		LocalOrders.Timestamps[ordertype][floor] = stamp
		
	default:
		fmt.Println("Invalid OrderType in updatelocalorders()")
	}

}

func updateLocalComplete(order heis.Order){
	LocalOrders.Completed = newLocalState.Completed
	for ordertype := UP; ordertype <= COMMAND; ordertype++ {
		for floor := 0; floor < N_FLOORS; floor++ {
			LocalOrders.Pending[ordertype][floor] = LocalOrders.Pending[ordertype][floor] && !(LocalOrders.Completed[ordertype][floor])
		}
	}
}

func updateGlobalOrders(){
	stamp := time.Now()
	for ordertype := UP; ordertype <= DOWN; ordertype++ {
		for floor := 0; floor < N_FLOORS; floor++ {
			if !GlobalOrders.Available[ordertype][floor] && LocalOrders.Pending[ordertype][floor]{
				!GlobalOrders.Available[ordertype][floor] = LocalOrders.Pending[ordertype][floor]
				GlobalOrders.Timestamps[ordertype][floor] = stamp 
			} 
		}
	}	
}

func updateGlobalComplete(){
	for ordertype := UP; ordertype <= DOWN; ordertype++ {
		for floor := 0; floor < N_FLOORS; floor++ {
			if !GlobalOrders.Completed[ordertype][floor] && LocalOrders.Completed[ordertype][floor]{
				!GlobalOrders.Completed[ordertype][floor] && LocalOrders.Completed[ordertype][floor]
			} 
		}
	}	
}

func scoresDone(){
	for ordertype := UP; ordertype <= DOWN; ordertype++ {
		for floor := 0; floor < N_FLOORS; floor++ {
			if !GlobalOrders.Completed[ordertype][floor] && LocalOrders.Completed[ordertype][floor]{
				!GlobalOrders.Completed[ordertype][floor] && LocalOrders.Completed[ordertype][floor]
			} 
		}
	}

	GlobalOrders.Taken[][][]
	Scores            [MAX_NUM_ELEVS][2][N_FLOORS]int    //'json:"Scores"'


}
/*************************************************/

func updateGlobalState() {
	//merge unvertified and global orders
	GlobalOrders = unverified_GlobalOrders
	setLights()

}

func setLights() {

	for b := 0; b < heis.N_BUTTONS-1; b++ {
		for floor := 0; floor < N_FLOORS; floor++ {
			if GlobalOrders.Available[b][floor] || GlobalOrders.Taken[b][floor] {
				heis.ElevSetButtonLamp(b, floor, 1)
			} else {
				heis.ElevSetButtonLamp(b, floor, 0)
			}

		}
	}
}

// getNextOrder() will be replaced by updateLocalOrders() AND mergeOrders()()
func getNextOrder() (heis.ElevButtonType, int, bool) {
	//IMPROVEMENTS:
	//Collect and choose best option instead of taking first one

	// If elev score is best in table AND not 0, it claims the order
	for ordertype := UP; ordertype <= DOWN; ordertype++ {
		for floor := 0; floor < N_FLOORS; floor++ {
			if GlobalOrders.Available[ordertype][floor] {
				if isBestScore(ordertype, floor) {
					GlobalOrders.Available[ordertype][floor] = false
					GlobalOrders.Taken[ordertype][floor] = true
					GlobalOrders.Timestamps[ordertype][floor] = GlobalOrders.Clock
					return ordertype, floor, true
				}
			}
		}
	}
	return -1, -1, false
}

func scoreOrders() {
	for ordertype := UP; ordertype <= DOWN; ordertype++ {
		for floor := 0; floor < N_FLOORS; floor++ {
			//INSERT COST FUNC HERE
			//simple:
			//GlobalOrders.Available[ordertype][floor]

			GlobalOrders.Scores[localElevIndex][ordertype][floor] = 10 //LocalOrders.Pending[ordertype][floor]*10;
		}
	}
}

func isBestScore(ordertype heis.ElevButtonType, floor int) bool {
	// returns false if it finds a better competitor, else returns true.
	// If this is somehow called when the network is down, ignore globalorders.
	if len(activeElevs) == 0 || !online {
		return true
	}
	for _, extElevIndex := range activeElevs {
		if GlobalOrders.Scores[localElevIndex][ordertype][floor] < GlobalOrders.Scores[extElevIndex][ordertype][floor] {
			return false
		}
	}
	return true
}



/********REMOVE**********************************************************************/
// LocalOrders are either updated through the fsmChan or updateLocalOrders()
func completeOrder(order heis.Order) {
	LocalOrders.Pending[order.OrderType][order.Floor] = false
	LocalOrders.Completed[order.OrderType][order.Floor] = true
	// Turn off light?
	// GlobalOrders.Available[i] = false
	// GlobalOrders.Taken[i] = -1
	// stop and open doors
}

/********REMOVE**********************************************************************/
func sendOrdersToNetwork() bool {
	// Replaced with an encode/decode function pair and netChan
	
	msg, _ := json.Marshal(GlobalOrders)
	str = string(msg)
	return true
}

/********REMOVE**********************************************************************/
func recvOrdersFromNetwork(network chan []byte) bool {
	// Replaced with an encode/decode function pair and netChan
	/*
		temp := make(GlobalOrderStruct)
		select {
		case msg := <-network:
			err = json.Unmarshal([]byte(str), &temp)
			if err == nil {
				Orders = mergeOrders(temp)
				return true
			}
		default:
			return false
		}
	*/
	_ = json.Unmarshal([]byte(str), &GlobalOrders)
	fmt.Println(GlobalOrders)
	return true
}

/********REMOVE**********************************************************************/

func mergeOrders(newGlobalOrders GlobalOrderStruct) {
	GlobalOrders = newGlobalOrders
	//^uint(0) gives the maximum value for uint
	GlobalOrders.Clock = (GlobalOrders.Clock + 1) % (^uint(0))
	// If an order is listed as taken, but this elev has completed it, the order is removed from globalorders
	for ordertype := UP; ordertype <= DOWN; ordertype++ {
		for floor := 0; floor < N_FLOORS; floor++ {
			if newGlobalOrders.Taken[ordertype][floor] && LocalOrders.Completed[ordertype][floor] {
				GlobalOrders.Taken[ordertype][floor] = false
				LocalOrders.Completed[ordertype][floor] = false
			}
		}
	}
}

func EncodeGlobalPacket() (b []byte) {
	GlobalPacketD, err := json.Marshal(GlobalOrders)
	_ = err
	return GlobalPacketD
}

func DecodeGlobalPacket(JsonPacket []byte) (PacketDEC GlobalOrderStruct, err error) {

	var GlobalPacketDEC GlobalOrderStruct
	err = json.Unmarshal(JsonPacket, &GlobalPacketDEC)
	return GlobalPacketDEC, err
}
