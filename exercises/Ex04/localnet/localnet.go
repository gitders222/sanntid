package localnet

import (
	"net"
	"strings"
	"errors"
	"sort"
)

var (
	localIP string
	broadcastIP string
	KnownIPs [] string
)

func GetLocalIP() (string, error) {
	if localIP == "" {
		conn, err := net.DialTCP("tcp4", nil, &net.TCPAddr{IP: []byte{8, 8, 8, 8}, Port: 53})
		if err != nil {
			return "", err
		}
		defer conn.Close()
		localIP = strings.Split(conn.LocalAddr().String(), ":")[0]
	}
	return localIP, nil
}

func GetBroadcastIP() (string, error) {
	if broadcastIP == "" {
		if localIP == "" {
			GetLocalIP()
		}
		temp := strings.Split(localIP, ".")
		broadcastIP = temp[0]+"."+temp[1]+"."+temp[2]+".255"
	}
	return broadcastIP, nil
}

func GetKnownIPs()([]string, error){
	return KnownIPs, nil
}

func GetNumberOfNodes()(int){
	return len(KnownIPs) +1
}

func AddNewNodeIP(newIP string)(error){
	for i:=0; i<len(KnownIPs); i++ {
		if KnownIPs[i] == newIP {
			return errors.New("IP is already in list")
		} 
	}
	KnownIPs = append(KnownIPs, newIP)
	sort.Strings(KnownIPs)
	return nil
}

func RemoveNodeIP( newIP string)(error){
	for i:=0; i<len(KnownIPs); i++ {
		if KnownIPs[i] == newIP {
			KnownIPs = append(KnownIPs[:i], KnownIPs[i+1:]...)
			return nil
		} 
	}
	return errors.New("IP to delete is not in list")
}

func GetNextNodeIP()(string){
	if GetNumberOfNodes() == 2 {
		return KnownIPs[0]
	} 
	// smallest member
	if localIP < KnownIPs[0] {
		return KnownIPs[0]
	}
	// somewhere inbetween
	for i:=0;i<len(KnownIPs)-1;i++ {
		if localIP == KnownIPs[i] {
			return "BadIPlist" //shouldn't happen
		} else if localIP > KnownIPs[i] && localIP < KnownIPs[i+1] {
			return KnownIPs[i+1]
		} 
	}
	// reached end of list, wrap around
	return KnownIPs[0]
}
func IsStartNode()(bool){
	if localIP < KnownIPs[0] {
		return true
	}
	return false
}

