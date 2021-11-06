package main

import (
    "fmt"
    "net"
    "bufio"
	"time"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	screenWidth  = 1000
	screenHeight = 800
)


func main() {
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

	plr, err := newPlayer(renderer)
	if err != nil {
		fmt.Println("creating player:", err)
		return
	}

	c := make(chan string)
	go UDPLoop(c)
	testcon, _ := net.Dial("udp", "127.0.0.1:22068")
	
	for {
		msg := <- c
		fmt.Println(msg)

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				return
			}
		}

		renderer.SetDrawColor(255, 255, 255, 255)
		renderer.Clear()

		plr.draw(renderer)

		renderer.Present()
	}
	testcon.Close()
}

func channelTest(c chan<- string) {
	for {
		c <- "every two seconds"
		time.Sleep(time.Second * 2)
	}
}


func UDPLoop(c chan<- string) {
    p :=  make([]byte, 2048)
    conn, err := net.Dial("udp", "127.0.0.1:22068")
    if err != nil {
        fmt.Printf("Some error %v", err)
        return
    }
	for {
		fmt.Fprintf(conn, "Hi UDP Server, How are you doing?")
		_, err = bufio.NewReader(conn).Read(p)
		if err == nil {
			fmt.Printf("%s\n", p)
		} else {
			fmt.Printf("Some error %v\n", err)
		}
		time.Sleep(time.Second * 1)
	}
    
    conn.Close()
}