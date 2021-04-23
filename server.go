package main

import (
	"math/rand"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/pion/turn/v2"
)

// StartServer starts stun/trurn server
func StartServer(Port int, Realm string, AuthResolver func(username string) (password *string)) error {
	netConn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return err
	}

	IP := netConn.LocalAddr().(*net.UDPAddr).IP

	udpListener, err := net.ListenPacket("udp4", "0.0.0.0:"+strconv.Itoa(Port))
	if err != nil {
		Log.Info("UDP error: ", err)
		return err
	}

	go runHealthCheck(Port)

	s, err := turn.NewServer(turn.ServerConfig{
		Realm: Realm,
		AuthHandler: func(Username string, realm string, srcAddr net.Addr) ([]byte, bool) {
			//Log.Info("RTC RQ")
			password := AuthResolver(Username)
			if password != nil {
				return turn.GenerateAuthKey(Username, Realm, *password), true
			}

			return nil, false
		},
		PacketConnConfigs: []turn.PacketConnConfig{
			{
				PacketConn: udpListener,
				RelayAddressGenerator: &turn.RelayAddressGeneratorStatic{
					RelayAddress: IP,
					Address:      "0.0.0.0",
				},
			},
		},
	})
	if err != nil {
		return err
	}

	Log.Info("WebRTC STUN/TURN is running on ", IP, ":", Port)

	defer func() {
		Log.Info("UDP defer error: ", err)
		if err = netConn.Close(); err != nil {
			Log.Info(err)
		}

		if err = s.Close(); err != nil {
			Log.Info(err)
		}
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)
	<-sigs

	return nil
}

// HealthCheck is needed for AWS
func runHealthCheck(Port int) error {
	tcpListener, err := net.Listen("tcp4", "0.0.0.0:"+strconv.Itoa(Port))
	if err != nil {
		Log.Info("TCP HC error: ", err)
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			Log.Info("Recovered TSP: ", r)
		}
	}()

	//defer tcpListener.Close()

	rand.Seed(time.Now().Unix())

	Log.Info("WebRTC STUN/TURN HEALTHCHECK is running on port :", Port)

	for {
		c, err := tcpListener.Accept()
		if err != nil {
			Log.Info("TCP HC IN error: ", err)
			return err
		}

		c.Write([]byte("OK"))
		c.Close()
	}
}
