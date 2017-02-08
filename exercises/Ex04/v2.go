package main 


import (
	"fmt"
	"net"
	"bufio"
	"time"
	"sort"
	"./localip"
	"./iferror"
)

type elevPacket struct {
	
}

func main(){
	UDPport := ":20022"
	// TCPport := ":20023"
	localIP, err := localip.Get()
	// broadcastIP, _ := localip.GetBroadcast()
	checkError(err, "Retrieving local IP", iferror.Ignore)
	passcode := "svekonrules"
	UDPmsg := passcode+"\n"+localIP+"\n"
	IPList := make([] string, 0, 20)

	UDPReceiveChan := make(chan string)
	UDPTransmitEnable := make(chan bool)
	//TCPReceiveChan := make(chan elevPacket)
	//TCPTransmitChan := make(chan elevPacket)
	
	go UDPReceiver(UDPReceiveChan, passcode, UDPport)
	go UDPBroadcaster(UDPTransmitEnable, UDPmsg, localIP, UDPport)
	UDPTransmitEnable <- true
	for {	
		newIP := <-UDPReceiveChan
		if dontKnowIP(newIP, IPList){
			IPList = append(IPList, newIP)
			sort.Strings(IPList)
		} else {
			// try to reconnect with failing node?
		}
		
		
		fmt.Println(IPList)
		//fmt.Println("main",<-UDPReceiveChan)
		//UDPTransmitEnable <- false
	}
}

func UDPReceiver(UDPReceiveChan chan string, passcode string, port string){
	localIP, _ := localip.Get()	
	addr, _ := net.ResolveUDPAddr("udp", port)
	socket, err := net.ListenUDP("udp", addr)
	checkError(err, "Setting up UDP listener", iferror.Quit)

	reader := bufio.NewReader(socket)
	
	for {
		code, err := reader.ReadString('\n')
		//checkError(err, "UDP datagram received", iferror.Ignore) // very frequent
		if code == (passcode + "\n") {
			msg, _ := reader.ReadString('\n')
			// ignore computer's own messages
			if msg != (localIP + "\n") ||true{	
				UDPReceiveChan <- msg[:len(msg)-1]
			}		
		} else {
			err = nil
			for err == nil {
				_, err = reader.ReadString('\n')	
			}	
		}
	}

}

func UDPBroadcaster(TransmitEnable chan bool, msg string, localBroadcastIP string, UDPport string){
		
	address, err := net.ResolveUDPAddr("udp",localBroadcastIP+UDPport)
	conn, err := net.DialUDP("udp", nil, address)
	checkError(err, "Initialising UDP broadcast", iferror.Ignore)
	
	broadcastOn := false
	for {
		// If off, wait for enable so that resources aren't taken.
		if !broadcastOn {
			broadcastOn = <-TransmitEnable		
		} else {
		// If on, check channel for updates and transmit every second.
			select {
				case broadcastOn = <-TransmitEnable:
					continue
				default:
					_, err = conn.Write( []byte(msg) )
					checkError(err, "Broadcasting IP", iferror.Ignore)
					time.Sleep(time.Second)
			}
		}
	}
}

func TCPReceiver(ReceiveChan chan elevPacket, previousNode *net.Conn){
	
}

func TCPTransmitter(ReceiveChan chan elevPacket, nextNode *net.Conn){}

func dontKnowIP(IP string, IPlist []string)(bool){
	for i:=0; i<len(IPlist); i++ {
		if IPlist[i] == IP {
			return false
		} 
	}
	return true
}

func checkError(err error, msg string, f iferror.Action){
	if err != nil {
		fmt.Println(msg, "... ERROR")
		fmt.Println(err)
		f()
	} else {
		fmt.Println(msg, "... Done")
	}
}

