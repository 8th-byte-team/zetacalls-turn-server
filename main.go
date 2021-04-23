package main

import (
	kp "gopkg.in/alecthomas/kingpin.v2"
)

var (
	RTCRealm = kp.Flag("realm", "STUN/TURN server realm").Required().String()
	RTCPort  = kp.Flag("port", "STUN/TURN server port").Required().Int()
	JWTSign  = kp.Flag("jwt-sign", "Signature for jwt token").Required().String()
)

func init() {
	kp.Parse()
}

func main() {
	if RTCRealm == nil {
		panic("please provide realm for stun/turn server")
	} else if RTCPort == nil {
		panic("please provide port for stun/turn server")
	} else if JWTSign == nil {
		panic("please provide port for stun/turn server")
	} else {
		err := StartServer(
			*RTCPort,
			*RTCRealm,
			ResolveRTCAuth(*JWTSign),
		)
		if err != nil {
			Log.Info(err)
		} else {
			Log.Info("STUN/TURN server is stopped")
		}
	}
}

func ResolveRTCAuth(JWTSignature string) func(JWT string) *string {
	return func(JWT string) *string {
		roomJWT, err := ParseJWT(JWTSignature, JWT)
		if err != nil {
			Log.Info(err)
			return nil
		}

		if roomJWT == nil {
			Log.Info("Err bad refresh token")
			return nil
		}

		password := roomJWT.Id + ":" + roomJWT.Audience

		return &password
	}
}
