package main

import (
    "fmt"
    "net"
    "bufio"
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

	plrLui, err := newPlayer(renderer, "sprites/gide_lui.bmp", 100, 500)
	if err != nil {
		fmt.Println("creating player lui:", err)
		return
	}
	plrMio, err := newPlayer(renderer, "sprites/gide_mio.bmp", 700, 200)
	if err != nil {
		fmt.Println("creating player lui2:", err)
		return
	}
	plrbg, err := newPlayer(renderer, "sprites/bg.bmp", 0, 0)
	if err != nil {
		fmt.Println("creating player bg:", err)
		return
	}

	luiWinScreen, err := newPlayer(renderer, "sprites/lui_win.bmp", 0, 0)
	if err != nil {
		fmt.Println("creating player bg:", err)
		return
	}
	mioWinScreen, err := newPlayer(renderer, "sprites/mio_win.bmp", 0, 0)
	if err != nil {
		fmt.Println("creating player bg:", err)
		return
	}
	drawScreen, err := newPlayer(renderer, "sprites/draw.bmp", 0, 0)
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

	var luiScoreImages [enemyCount+1]player
	for i := 0; i <= enemyCount; i++ {
		temp, err := newPlayer(renderer, "sprites/" + strconv.Itoa(i) + ".bmp", 200, 200)
		if err != nil {
			fmt.Println("creating player enemy:", err)
			return
		}
		luiScoreImages[i] = temp
	}

	M, err := newPlayer(renderer, "sprites/M.bmp", 300, 100)
	if err != nil {
		fmt.Println("creating player bg:", err)
		return
	}
	L, err := newPlayer(renderer, "sprites/L.bmp", 200, 100)
	if err != nil {
		fmt.Println("creating player bg:", err)
		return
	}

	var mioScoreImages [enemyCount+1]player
	for i := 0; i <= enemyCount; i++ {
		temp, err := newPlayer(renderer, "sprites/" + strconv.Itoa(i) + ".bmp", 300, 200)
		if err != nil {
			fmt.Println("creating player enemy:", err)
			return
		}
		mioScoreImages[i] = temp
	}

	var luiScore int
	var mioScore int

	laddr, err := net.ResolveUDPAddr("udp", args[1] + ":" + args[2])
	argPort, _ := strconv.Atoi(args[4])
	raddr := net.UDPAddr{IP: net.ParseIP(args[3]), Port: argPort}
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
			rmsg := strings.Split(msg, ",")
			plrLui.x, _ = strconv.ParseFloat(rmsg[0], 64)
			plrLui.y, _ = strconv.ParseFloat(rmsg[1], 64)
			plrMio.x, _ = strconv.ParseFloat(rmsg[2], 64)
			plrMio.y, _ = strconv.ParseFloat(rmsg[3], 64)
			for i := 0; i < enemyCount; i++ {
				newx, _ := strconv.ParseFloat(rmsg[4+i], 64)
				enemyPositions[i].x  = newx
				newy, _ := strconv.ParseFloat(rmsg[5+i], 64)
				enemyPositions[i].y = newy
			}
			for i := 0; i < enemyCount; i++ {
				newv, _ := strconv.Atoi(rmsg[16+i])
				if newv == 0 && isEnemyDead[i] == 1 {
					// pop 
					fmt.Println("collision !")
				}
				isEnemyDead[i] = newv 
			}
			luiScore, _ = strconv.Atoi(strings.Split(msg, "|")[1])
			mioScore, _ = strconv.Atoi(strings.Split(msg, "|")[2])
		default:
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

		plrLui.draw(renderer)
		plrMio.draw(renderer)

		if luiScore + mioScore == 6 {
			if luiScore > mioScore {
				luiWinScreen.drawbg(renderer)
			} else if mioScore > luiScore {
				mioWinScreen.drawbg(renderer)
			} else {
				// draw
				drawScreen.drawbg(renderer)
			}
			if args[0] == "mio" {
				plrMio.update(conn)
			} else if args[0] == "lui" {
				plrLui.update(conn)
			}
		} else {
			for i := 0; i < enemyCount; i++ {
				if isEnemyDead[i] == 1 {
					enemyPositions[i].draw(renderer)
				}
			}
	
			
	
			M.draw(renderer)
			L.draw(renderer)
			luiScoreImages[luiScore].draw(renderer)
			mioScoreImages[mioScore].draw(renderer)
	
			if args[0] == "mio" {
				plrMio.update(conn)
			} else if args[0] == "lui" {
				plrLui.update(conn)
			}
			
		}
		renderer.Present()

		
	}
	conn.Close()
}


func UDPLoop(c chan<- string, conn net.Conn) {
    p :=  make([]byte, 2048)
	for {
		_, err := bufio.NewReader(conn).Read(p)
		if err == nil {
			c <- string(p)
		} else {
			fmt.Printf("Some error %v\n", err)
		}
	}
}

func SendUDP(msg string, conn net.Conn) {
	fmt.Fprintf(conn, msg)
}