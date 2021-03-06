package main

import (
    "fmt" 
    "net"  
	"math/rand"
	"time"
	"math"
)

const (
	playerHeight = 162
	playerWidth = 100
	screenWidth  = 1000
	screenHeight = 600
	enemyCount = 6
)

type SpritePos struct {
	x, y int;
}

func sendResponse(conn *net.UDPConn, addr *net.UDPAddr, luiPos SpritePos, mioPos SpritePos, enemyPositions [enemyCount]SpritePos, isEnemyDead [enemyCount]int, luiScore int, mioScore int) {
	msg := fmt.Sprintf("%d,%d,%d,%d", luiPos.x, luiPos.y, mioPos.x, mioPos.y)
	for i := 0; i < enemyCount; i++  {
		msg = fmt.Sprintf("%s,%d,%d", msg, enemyPositions[i].x, enemyPositions[i].y)
	}
	for i := 0; i < enemyCount; i++  {
		msg = fmt.Sprintf("%s,%d", msg, isEnemyDead[i])
	}
	msg = fmt.Sprintf("%s|%d|%d|", msg, luiScore, mioScore)
	msg = fmt.Sprintf("%s::::::", msg)
	fmt.Println(msg)
    _,err := conn.WriteToUDP([]byte(msg), addr)
    if err != nil {
        fmt.Printf("Couldn't send response %v", err)
    }
}


func main() {
	// port 22069 is player 1
	// port 22420 is player 2
	var startTime int64

	var luiIP string
	var mioIP string
	var luiPort int
	var mioPort int

	luiPos := SpritePos{x: 100, y: 500}
	mioPos := SpritePos{x: 700, y: 200}

	luiScore := 0
	mioScore := 0

	var enemyPositions [enemyCount]SpritePos
	for i := 0; i < enemyCount; i++ {
		enemyPositions[i] = SpritePos{x: 400, y: 300}
	}
	var isEnemyDead [enemyCount]int 
	for i := 0; i < enemyCount; i++ {
		isEnemyDead[i] = 1
	}
	
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
		if string(p)[:2] == "up" {
			// move up
			updatePlayerPos(0, 1, remoteaddr, &luiPos, &mioPos, luiIP, luiPort)
			for i := 0; i < enemyCount; i++ {
				updateEnemyPos(&enemyPositions[i])
			}
		} else if string(p)[:4] == "down" {
			// move down
			updatePlayerPos(0, -1, remoteaddr, &luiPos, &mioPos, luiIP, luiPort)
			for i := 0; i < enemyCount; i++ {
				updateEnemyPos(&enemyPositions[i])
			}
		} else if string(p)[:4] == "left" {
			// move left
			updatePlayerPos(-1, 0, remoteaddr, &luiPos, &mioPos, luiIP, luiPort)
			for i := 0; i < enemyCount; i++ {
				updateEnemyPos(&enemyPositions[i])
			}
		} else if string(p)[:5] == "right" {
			// move right
			updatePlayerPos(1, 0, remoteaddr, &luiPos, &mioPos, luiIP, luiPort)
			for i := 0; i < enemyCount; i++ {
				updateEnemyPos(&enemyPositions[i])
			}
		} else if string(p)[:4] == "init" {
			if string(p)[5:8] == "lui" {
				luiIP = fmt.Sprintf("%s", remoteaddr.IP)
				luiPort = remoteaddr.Port
				startTime = time.Now().Unix()
			} else if string(p)[5:8] == "mio" {
				mioIP = fmt.Sprintf("%s", remoteaddr.IP)
				mioPort = remoteaddr.Port
				startTime = time.Now().Unix()
			}
		}
		luiScore, mioScore = checkCollisions(enemyPositions, &isEnemyDead, luiPos, mioPos, startTime, luiScore, mioScore)

		if luiIP != "" {
			raddr := net.UDPAddr{IP: net.ParseIP(luiIP), Port: luiPort}
			go sendResponse(ser, &raddr, luiPos, mioPos, enemyPositions, isEnemyDead, luiScore, mioScore)
		}

		if mioIP != "" {
			raddr := net.UDPAddr{IP: net.ParseIP(mioIP), Port: mioPort}
			go sendResponse(ser, &raddr, luiPos, mioPos, enemyPositions, isEnemyDead, luiScore, mioScore)
		}
    }
}

func checkCollisions(enemyPositions [enemyCount]SpritePos, isEnemyDead *[enemyCount]int, luiPos SpritePos, mioPos SpritePos, startTime int64, luiScore int, mioScore int) (newLui int, newMio int) {
	if time.Now().Unix() < startTime + 10 {
		return
	}
	for i := 0; i < enemyCount; i++ {
		if isEnemyDead[i] == 1 {
			luiDiffx := math.Abs(float64(luiPos.x - enemyPositions[i].x))
			luiDiffy := math.Abs(float64(luiPos.y - enemyPositions[i].y))
			if (luiDiffx < playerWidth/2 && luiDiffy < playerHeight/2) {
				fmt.Println("COLLISION")
				isEnemyDead[i] = 0
				luiScore += 1
			}

			mioDiffx := math.Abs(float64(mioPos.x - enemyPositions[i].x))
			mioDiffy := math.Abs(float64(mioPos.y - enemyPositions[i].y))
			if (mioDiffx < playerWidth/2 && mioDiffy < playerHeight/2) {
				fmt.Println("COLLISION")
				isEnemyDead[i] = 0
				mioScore += 1
			}
			
		}
	}
	return luiScore, mioScore
}

func updateEnemyPos(e *SpritePos) {
	s1 := rand.NewSource(time.Now().UnixNano())
    r1 := rand.New(s1)

	direction := r1.Intn(4)
	coinFlip := r1.Intn(2)
	coinFlip = coinFlip * 3
	if direction == 0 {
		// up 
		e.y -= (1 * coinFlip)

	} else if direction == 1 {
		// down 
		e.y += (1 * coinFlip)
	} else if direction == 2 {
		// down 
		e.x += (1 * coinFlip)
	} else if direction == 3 {
		// down 
		e.x -= (1 * coinFlip)
	}		
	if e.y < 0 { e.y = 0 }
	if e.y > screenHeight { e.y = screenHeight }
	if e.x < 0 { e.x = 0 }
	if e.x > screenHeight { e.x = screenHeight }
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