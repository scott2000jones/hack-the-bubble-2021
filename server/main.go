package main

import (
    "fmt" 
    "net"  
	"strconv"
)

const (
	screenWidth  = 1000
	screenHeight = 600
)

type SpritePos struct {
	x, y int;
}

func sendResponse(conn *net.UDPConn, addr *net.UDPAddr, luiPos SpritePos, mioPos SpritePos) {
	msg := fmt.Sprintf("%d,%d,%d,%d,|", luiPos.x, luiPos.y, mioPos.x, mioPos.y)
    _,err := conn.WriteToUDP([]byte(msg), addr)
    if err != nil {
        fmt.Printf("Couldn't send response %v", err)
    }
}


func main() {
	// port 22069 is player 1
	// port 22420 is player 2

	var luiIP string
	var mioIP string
	var luiPort int
	var mioPort int
	// var luiAddr *net.UDPAddr
	// var mioAddr *net.UDPAddr

	luiPos := SpritePos{x: 200, y: 200}
	mioPos := SpritePos{x: 300, y: 200}

	fmt.Println(strconv.Itoa(luiPos.x) + "  " + strconv.Itoa(mioPos.x))

    p := make([]byte, 2048)
    addr := net.UDPAddr{
        Port: 22068,
        IP: net.ParseIP("138.251.29.189"),
    }
    ser, err := net.ListenUDP("udp", &addr)
    if err != nil {
        fmt.Printf("Some error %v\n", err)
        return
    }
    for {
        _,remoteaddr,err := ser.ReadFromUDP(p)
		if err !=  nil {
            fmt.Printf("Some error  %v", err)
            continue
        }
        // fmt.Printf("Read a message from %v %s \n", remoteaddr, p)
		fmt.Println("read: " + string(p) + "----------------------------")
		fmt.Println(string(p)[:4])
		if string(p)[:2] == "up" {
			// move up
			fmt.Println("up")
			// fmt.Printf("%s\n", (remoteaddr.IP))
			// fmt.Println(strconv.Itoa(remoteaddr.Port))
			// fmt.Println(remoteaddr.Zone)
			updatePlayerPos(0, 1, remoteaddr, &luiPos, &mioPos, luiIP, luiPort)
		} else if string(p)[:4] == "down" {
			// move down
			fmt.Println("down")
			updatePlayerPos(0, -1, remoteaddr, &luiPos, &mioPos, luiIP, luiPort)
		} else if string(p)[:4] == "left" {
			// move left
			fmt.Println("left")
			updatePlayerPos(-1, 0, remoteaddr, &luiPos, &mioPos, luiIP, luiPort)
		} else if string(p)[:5] == "right" {
			// move right
			fmt.Println("right")
			updatePlayerPos(1, 0, remoteaddr, &luiPos, &mioPos, luiIP, luiPort)
		} else if string(p)[:4] == "init" {
			if string(p)[5:8] == "lui" {
				luiIP = fmt.Sprintf("%s", remoteaddr.IP)
				luiPort = remoteaddr.Port
			} else if string(p)[5:8] == "mio" {
				mioIP = fmt.Sprintf("%s", remoteaddr.IP)
				mioPort = remoteaddr.Port
			}
			fmt.Println(luiIP)
			fmt.Println(mioIP)
			fmt.Println(mioPort)
			// fmt.Printf("%s\n", remoteaddr.IP)
			// fmt.Println(strconv.Itoa(remoteaddr.Port))
		}
		fmt.Println(strconv.Itoa(luiPos.x) + "  " + strconv.Itoa(luiPos.y) + " || " + strconv.Itoa(mioPos.x) + "  " + strconv.Itoa(mioPos.y))
		

		if luiIP != "" {
			raddr := net.UDPAddr{IP: net.ParseIP(luiIP), Port: luiPort}
			go sendResponse(ser, &raddr, luiPos, mioPos)
		}

		if mioIP != "" {
			raddr := net.UDPAddr{IP: net.ParseIP(mioIP), Port: mioPort}
			// fmt.Println((&raddr).Port)
			go sendResponse(ser, &raddr, luiPos, mioPos)
		}

    }
}

func updatePlayerPos(xdiff int, ydiff int, remoteaddr *net.UDPAddr, luiPos *SpritePos, mioPos *SpritePos, luiIP string, luiPort int) {
	if fmt.Sprintf("%s", remoteaddr.IP) == luiIP && remoteaddr.Port == luiPort {
		luiPos.x += xdiff
		luiPos.y -= ydiff
		if luiPos.x < 0 { luiPos.x = 0 }
		if luiPos.y < 0 { luiPos.y = 0 }
		if luiPos.x > screenWidth { luiPos.x = screenWidth }
		if luiPos.y > screenHeight { luiPos.y = screenHeight }
	} else {
		mioPos.x += xdiff
		mioPos.y -= ydiff
		if mioPos.x < 0 { mioPos.x = 0 }
		if mioPos.y < 0 { mioPos.y = 0 }
		if mioPos.x > screenWidth { mioPos.x = screenWidth }
		if mioPos.y > screenHeight { mioPos.y = screenHeight }
	}
}