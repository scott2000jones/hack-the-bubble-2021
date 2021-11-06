package main

import (
    "fmt"
    "net"
    "bufio"
	// "time"
	"github.com/veandco/go-sdl2/sdl"
	"strconv"
	"os"
	"strings"
)

const (
	screenWidth  = 1000
	screenHeight = 600
	enemyCount = 6
)

type SpritePos struct {
	x int;
	y int;
}

func main() {
	// Get args
    args := os.Args[1:]


	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		fmt.Println("initializing SDL:", err)
		return
	}

	window, err := sdl.CreateWindow(
		"Hack the Bubble 2021",
		sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		screenWidth, screenHeight,
		sdl.WINDOW_OPENGL)
	if err != nil {
		fmt.Println("initializing window:", err)
		return
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		fmt.Println("initializing renderer:", err)
		return
	}
	defer renderer.Destroy()

	plrLui, err := newPlayer(renderer, "sprites/gide_lui.bmp", 200, 200)
	if err != nil {
		fmt.Println("creating player lui:", err)
		return
	}
	plrMio, err := newPlayer(renderer, "sprites/gide_mio.bmp", 201, 202)
	if err != nil {
		fmt.Println("creating player lui2:", err)
		return
	}
	plrbg, err := newPlayer(renderer, "sprites/bg.bmp", 0, 0)
	if err != nil {
		fmt.Println("creating player bg:", err)
		return
	}

	var enemyPositions [enemyCount]player
	for i := 0; i < enemyCount; i++ {
		temp, err := newPlayer(renderer, "sprites/ok.bmp", 200, 200)
		if err != nil {
			fmt.Println("creating player enemy:", err)
			return
		}
		enemyPositions[i] = temp
	}
	var isEnemyDead [enemyCount]int
	for i := 0; i < enemyCount; i++ {
		isEnemyDead[i] = 1
	}

	laddr, err := net.ResolveUDPAddr("udp", args[1] + ":" + args[2])
	// argPort, _ := strconv.Atoi(args[2])
	raddr := net.UDPAddr{IP: net.ParseIP(args[3]), Port: 22068}
	conn, err := net.DialUDP("udp", laddr, &raddr)
	if err != nil {
        fmt.Printf("Some error %v", err)
        return
    }
	c := make(chan string)
	go UDPLoop(c, conn)
	if args[0] != "lui" {
		SendUDP("init,mio", conn)
	} else {
		SendUDP("init,lui", conn)
	}

	for {
		select {
		case msg := <-c:
			fmt.Println("received message", msg)
			rmsg := strings.Split(msg, ",")
			// fmt.Println(rmsg[0])
			// floatNum, errF := strconv.ParseFloat(rmsg[0], 64)
			// if errF != nil {
			// 	fmt.Printf("Parse error %v\n", errF)
			// }
			// fmt.Println(floatNum)
			plrLui.x, _ = strconv.ParseFloat(rmsg[0], 64)
			plrLui.y, _ = strconv.ParseFloat(rmsg[1], 64)
			plrMio.x, _ = strconv.ParseFloat(rmsg[2], 64)
			plrMio.y, _ = strconv.ParseFloat(rmsg[3], 64)
			for i := 0; i < enemyCount; i++ {
				enemyPositions[i].x, _ = strconv.ParseFloat(rmsg[4+i], 64)
				enemyPositions[i].y, _ = strconv.ParseFloat(rmsg[5+i], 64)
			}
			for i := 0; i < enemyCount; i++ {
				isEnemyDead[i], _ = strconv.Atoi(rmsg[16+i])
			}
		default:
			// fmt.Println("no message received")
		}

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				return
			}
		}

		renderer.SetDrawColor(255, 255, 255, 255)
		renderer.Clear()
		plrbg.drawbg(renderer)

		for i := 0; i < enemyCount; i++ {
			if isEnemyDead[i] == 1 {
				enemyPositions[i].draw(renderer)
			}
		}

		plrLui.draw(renderer)
		plrMio.draw(renderer)

		// fmt.Println(plrLui.x) 
		// fmt.Println(plrLui.y) 
		// fmt.Println(plrMio.x)
		// fmt.Println(plrMio.y)

		if args[0] == "mio" {
			plrMio.update(conn)
			// plrLui2.update(conn)
		} else if args[0] == "lui" {
			plrLui.update(conn)
		}
		

		renderer.Present()
	}
	conn.Close()
}


func UDPLoop(c chan<- string, conn net.Conn) {
    p :=  make([]byte, 2048)
	for {
		// fmt.Fprintf(conn, "Hi UDP Server, How are you doing?")
		_, err := bufio.NewReader(conn).Read(p)
		if err == nil {
			// fmt.Printf("%s\n", p)
			c <- string(p)
		} else {
			fmt.Printf("Some error %v\n", err)
		}
		// time.Sleep(time.Second * 1)
	}
}

func SendUDP(msg string, conn net.Conn) {
	fmt.Fprintf(conn, msg)
}