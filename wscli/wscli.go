package wscli

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"slices"
	"time"

	"github.com/kachaje/sacco-schema/utils"

	"github.com/gorilla/websocket"
)

func clearScreen() {
	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}

	cmd.Stdout = os.Stdout
	cmd.Run()
}

func Main() {
	var port int64 = 8080
	var phoneNumber string = "1234567890"
	var silentMode bool

	flag.Int64Var(&port, "p", port, "server port")
	flag.StringVar(&phoneNumber, "n", phoneNumber, "phone number")
	flag.BoolVar(&silentMode, "s", silentMode, "silent mode")

	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := utils.WaitForPort("localhost", fmt.Sprint(port), 30*time.Second, 2*time.Second, false)
	if err != nil {
		log.Fatal(err)
	}

	client, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://localhost:%d/ws?phoneNumber=%s", port, phoneNumber), nil)
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(os.Stdin)

	for {
		select {
		case <-ctx.Done():
			return
		default:
			clearScreen()

			fmt.Println("")

			_, message, err := client.ReadMessage()
			if err != nil {
				log.Fatal(err)
				return
			}

			if slices.Contains([]string{
				"Thank you for using our service",
				"Zikomo potidalila",
			}, string(message)) {
				return
			} else if !silentMode {
				fmt.Println(string(message))
			}

			scanner.Scan()

			input := scanner.Text()

			err = client.WriteMessage(websocket.TextMessage, []byte(input))
			if err != nil {
				log.Fatal(err)
				return
			}
		}
	}
}
