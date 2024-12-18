package main

// Config represents the configuration for a role.
type Config struct {
	CertPath     string `json:"certPath"`
	KeyDirectory string `json:"keyPath"`
	TLSCertPath  string `json:"tlsCertPath"`
	PeerEndpoint string `json:"peerEndpoint"`
	GatewayPeer  string `json:"gatewayPeer"`
	MSPID        string `json:"mspID"`
}

// Create a Profile map
var profile = map[string]Config{

	"university": {
		CertPath:     "../Network/organizations/peerOrganizations/university.cred.com/users/User1@university.cred.com/msp/signcerts/cert.pem",
		KeyDirectory: "../Network/organizations/peerOrganizations/university.cred.com/users/User1@university.cred.com/msp/keystore/",
		TLSCertPath:  "../Network/organizations/peerOrganizations/university.cred.com/peers/peer0.university.cred.com/tls/ca.crt",
		PeerEndpoint: "localhost:7051",
		GatewayPeer:  "peer0.university.cred.com",
		MSPID:        "UniversityMSP",
	},

	"student": {
		CertPath:     "../Network/organizations/peerOrganizations/student.cred.com/users/User1@student.cred.com/msp/signcerts/cert.pem",
		KeyDirectory: "../Network/organizations/peerOrganizations/student.cred.com/users/User1@student.cred.com/msp/keystore/",
		TLSCertPath:  "../Network/organizations/peerOrganizations/student.cred.com/peers/peer0.student.cred.com/tls/ca.crt",
		PeerEndpoint: "localhost:9051",
		GatewayPeer:  "peer0.student.cred.com",
		MSPID:        "StudentMSP",
	},

	"company": {
		CertPath:     "../Network/organizations/peerOrganizations/company.cred.com/users/User1@company.cred.com/msp/signcerts/cert.pem",
		KeyDirectory: "../Network/organizations/peerOrganizations/company.cred.com/users/User1@company.cred.com/msp/keystore/",
		TLSCertPath:  "../Network/organizations/peerOrganizations/company.cred.com/peers/peer0.company.cred.com/tls/ca.crt",
		PeerEndpoint: "localhost:11051",
		GatewayPeer:  "peer0.company.cred.com",
		MSPID:        "StudentMSP",
	},

	"university2": {
		CertPath:     "../Network/organizations/peerOrganizations/university.cred.com/users/User2@university.cred.com/msp/signcerts/cert.pem",
		KeyDirectory: "../Network/organizations/peerOrganizations/university.cred.com/users/User2@university.cred.com/msp/keystore/",
		TLSCertPath:  "../Network/organizations/peerOrganizations/university.cred.com/peers/peer0.university.cred.com/tls/ca.crt",
		PeerEndpoint: "localhost:7051",
		GatewayPeer:  "peer0.university.cred.com",
		MSPID:        "UniversityMSP",
	},

}