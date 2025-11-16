package domain

import (
	"fmt"
	"net"
	"os"

	"github.com/net12labs/cirm/dali/context/cmd"
	"github.com/net12labs/cirm/dali/data"
	domain_context "github.com/net12labs/cirm/dali/domain/context"
	"github.com/net12labs/cirm/dali/rtm"
	webserver "github.com/net12labs/cirm/mali/web-server"
	"github.com/net12labs/cirm/service-daemon/site"
)

type dom struct {
	name           string
	Site           *site.Site
	WebServer      *domain_context.WebServer
	socketListener net.Listener
}

var Main = &dom{}

func (d *dom) Name() string {
	return d.name
}

// we should also be able to bubble from here up

func (d *dom) Init(name string) *dom {
	d.name = name
	d.runtimeInit()

	d.dataInit()
	d.Site = site.NewSite()

	d.Site.Execute = func(cmd *cmd.Cmd) {
		fmt.Println("Executing command via Site:", cmd)
		if cmd.ExitCode == -1 {
			cmd.ExitCode = 1
			cmd.ErrorMsg = "No handler implemented"
		}
	}

	d.startSocket()

	d.WebServer = webserver.NewWebServer()
	d.Site.WebServer = d.WebServer

	d.WebServer.AddRoute("/", func(req *webserver.Request) {
		if req.Path.Path == "/" {
			req.RedirectToUrl("/home")
			return
		}
		req.WriteResponse404()
	})

	d.Site.Init()

	return d
}

func (d *dom) Start() {

	if err := d.Site.Start(); err != nil {
		fmt.Println("Failed to start service:", err)
		rtm.Runtime.Exit(1)
	}

	if err := d.WebServer.Start(); err != nil {
		fmt.Printf("Failed to start web server: %v\n", err)
		rtm.Runtime.Exit(1)
	}

}

func (d *dom) dataInit() {
	dbPath := rtm.Etc.Get("main_db_path").String()
	rtm.Do.InitFsPath(dbPath)
	mainDb := data.Ops.CreateDb("main", dbPath)

	if err := mainDb.Init(); err != nil {
		fmt.Println("Failed to initialize database:", err)
		rtm.Runtime.Exit(1)
	}
}

func (d *dom) runtimeInit() {

	pid := rtm.Pid
	pid.Pid.PidFilePath = rtm.Etc.Get("pid_file_path").String()
	if err := pid.Handle_ExitOnDuplicate(); err != nil {
		fmt.Println(err)
		rtm.Runtime.Exit(1)
	}
	rtm.Runtime.OnPanic.AddListener(func(err any) {
		pid.Handle_CleanupOnExit()
	})
	rtm.Runtime.OnExit.AddListener(func(code any) {
		pid.Handle_CleanupOnExit()
		fmt.Println("Exited with code", code)
	})

}

func (d *dom) startSocket() {
	socketPath := rtm.Etc.Get("socket_path").String()
	rtm.Do.InitFsPath(socketPath)

	// Remove existing socket file if it exists
	if err := os.Remove(socketPath); err != nil && !os.IsNotExist(err) {
		fmt.Println("Failed to remove existing socket:", err)
		rtm.Runtime.Exit(1)
	}

	// Create Unix domain socket
	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		fmt.Println("Failed to create socket:", err)
		rtm.Runtime.Exit(1)
	}

	d.socketListener = listener

	// Cleanup socket on exit
	rtm.Runtime.OnExit.AddListener(func(code any) {
		if d.socketListener != nil {
			d.socketListener.Close()
			d.socketListener = nil
		}
		os.Remove(socketPath)
	})

	// Start listening for connections in a goroutine
	go func() {
		fmt.Println("Socket listening on:", socketPath)
		for {
			conn, err := listener.Accept()
			if err != nil {
				// Check if listener was closed intentionally
				if d.socketListener == nil {
					return
				}
				return
			}

			// Handle connection in a separate goroutine
			go d.handleSocketConnection(conn)
		}
	}()
}

func (d *dom) handleSocketConnection(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, 4096)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Socket read error:", err)
		return
	}

	// Process the received data
	data := buf[:n]
	fmt.Printf("Received on socket: %s\n", string(data))

	// Send response back
	response := "ACK"
	if _, err := conn.Write([]byte(response)); err != nil {
		fmt.Println("Socket write error:", err)
	}
}
