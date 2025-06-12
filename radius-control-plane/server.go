package main

import (
	"fmt"
	"log"
	"time"
	"radius-control-plane/datastore"

	"layeh.com/radius"
	"layeh.com/radius/rfc2865"
	"layeh.com/radius/rfc2866"
)


var (
	/* Server User Name */
	ServerUsername = "testuser"
	/* Server Password */
	ServerPassword = "mypassword"
	/* Radius Server Listening Point*/
	ServerPort = "redis:1813"
)

/* Main Entry Point to the Radius Server Implementation */
func main() {
	
	go initRadiusServer()

	datastore.InitRedisClient()

	select{}
}

/* Initialize the Radius Server. Listens by default on localhost:1813 */
func initRadiusServer() {
	server := radius.PacketServer{
		Addr: ServerPort,
		Handler:      radius.HandlerFunc(processIncomingRequest),
		SecretSource: radius.StaticSecretSource([]byte(`testing123`)),
	}

	log.Printf("Starting Radius Server on : %s", ServerPort)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

/* Process Incoming Radius message and filters the Accounting Request message */
func processIncomingRequest(w radius.ResponseWriter, r *radius.Request) {
	username := rfc2865.UserName_GetString(r.Packet)
        password := rfc2865.UserPassword_GetString(r.Packet)

	log.Printf("Incoming Radius Request Name : %s Password : %s ", username, password)
    var code radius.Code
        if username != ServerUsername {
            code = radius.CodeAccessReject
			log.Printf("Writing %v to %v", code, r.RemoteAddr)
            w.Write(r.Response(code))
			return
        }
	var incoming bool
	if r.Packet.Code == radius.CodeAccountingRequest {
		log.Printf("Incoming Accounting Request of type : %d",  rfc2866.AcctStatusType_Get(r.Packet))
		incoming = true
	}

	// UserName
	record := datastore.AccountingRecord {}
	record.UserName = username

	// NAS-IP-Address
	nasIpAddr := rfc2865.NASIPAddress_Get(r.Packet).String()
	if nasIpAddr != "" {
		record.NASIPAddress = nasIpAddr
	}

	// NAS-Port
	port := rfc2865.NASPort_Get(r.Packet)
	nasPort := fmt.Sprintf("%d", port)
	record.NASPort = nasPort

	// Acct-Status-Type
	acctType := rfc2866.AcctStatusType_Get(r.Packet).String()
	record.AcctStatusType = acctType

	// Acct-Session-Id
	acctSessId := rfc2866.AcctSessionID_GetString(r.Packet)
	record.AcctSessionId = acctSessId

	// Framed-Ip-Address
	framedIpAddr := rfc2865.FramedIPAddress_Get(r.Packet).String()
	record.FramedIpAddress =  framedIpAddr

	// Calling-Station-Id
	cgStationId := rfc2865.CallingStationID_GetString(r.Packet)
	record.CallingStationId = cgStationId

	// Called-Station-Id
	cdStationId := rfc2865.CalledStationID_GetString(r.Packet)
	record.CalledStationId = cdStationId

	// Time Stamp
	timestamp := fmt.Sprintf("%d", time.Now().Unix())
	record.Timestamp = timestamp

	// Client IP
	record.ClientIp = r.RemoteAddr.String()

	// Packet-Type
	if incoming {
		record.PacketType = "Accounting-Request"
	}
	log.Printf("Incoming Request %s", record)

	acctKey :=  fmt.Sprintf("radius:acct:%s:%s:%s",username, acctSessId, timestamp)
	log.Printf("Record Key : %s", acctKey)

	datastore.StoreAccountingRecord(acctKey, record)

	code = radius.CodeAccountingResponse
        log.Printf("Writing %v to %v", code, r.RemoteAddr)
        w.Write(r.Response(code))
}

