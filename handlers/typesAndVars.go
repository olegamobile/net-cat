package handlers

import (
	"net"
	"sync"
)

type Client struct {
	conn net.Conn
	name string
}

// Request represents a client request with a timestamp.
type Request struct {
	client Client
	data   string
}

type Users map[string]*Client

var LogFile = CreateLogFile()
var UserList = CreateUserList()
var MsgHistory string
var MessagePipe = make(chan Request, 100)
var BrodcastPipe = make(chan Request, 100)
var Lock sync.Mutex

var (
	Green = "\033[32m" // Green text color
	Red   = "\033[31m" // Red text color
	White = "\033[37m" // White text color (default)
	Reset = "\033[0m"  // Reset to default terminal color
)

const RoomSize = 10
const LogFileDir = "logs"

const welcomeIcon = `
         _nnnn_
        dGGGGMMb
       @p~qp~~qMb
       M|@||@) M|
       @,----.JM|
      JS^\__/  qKL
     dZP        qKRb
    dZP          qKKb
   fZP            SMMb
   HZM            MMMM
   FqM            MMMM
 __| ".        |\dS"qML
|    ` + "`.       | `' \\Zq\n" +
	"_)      \\.___.,|     .'\n" +
	"\\____   )MMMMMP|   .'\n" +
	"     `-'       `--'\n"
