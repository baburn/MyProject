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

	"University": {
		CertPath:     "../organizations/peerOrganizations/university.cred.com/users/User1@university.cred.com/msp/signcerts/cert.pem",
		KeyDirectory: "../organizations/peerOrganizations/university.cred.com/users/User1@university.cred.com/msp/keystore/",
		TLSCertPath:  "../organizations/peerOrganizations/university.cred.com/peers/peer0.university.cred.com/tls/ca.crt",
		PeerEndpoint: "localhost:7051",
		GatewayPeer:  "peer0.university.cred.com",
		MSPID:        "UniversityMSP",
	},

	"company": {
		CertPath:     "../organizations/peerOrganizations/company.cred.com/users/User1@company.cred.com/msp/signcerts/cert.pem",
		KeyDirectory: "../organizations/peerOrganizations/company.cred.com/users/User1@company.cred.com/msp/keystore/",
		TLSCertPath:  "../organizations/peerOrganizations/company.cred.com/peers/peer0.company.cred.com/tls/ca.crt",
		PeerEndpoint: "localhost:9051",
		GatewayPeer:  "peer0.company.cred.com",
		MSPID:        "CompanyMSP",
	},

	"student": {
		CertPath:     "../organizations/peerOrganizations/student.cred.com/users/User1@student.cred.com/msp/signcerts/cert.pem",
		KeyDirectory: "../organizations/peerOrganizations/student.cred.com/users/User1@student.cred.com/msp/keystore/",
		TLSCertPath:  "../organizations/peerOrganizations/student.cred.com/peers/peer0.student.cred.com/tls/ca.crt",
		PeerEndpoint: "localhost:11051",
		GatewayPeer:  "peer0.student.cred.com",
		MSPID:        "StudentMSP",
	},

	"university2": {
		CertPath:     "../organizations/peerOrganizations/university.cred.com/users/User2@university.cred.com/msp/signcerts/cert.pem",
		KeyDirectory: "../organizations/peerOrganizations/university.cred.com/users/User2@university.cred.com/msp/keystore/",
		TLSCertPath:  "../organizations/peerOrganizations/university.cred.com/peers/peer0.university.cred.com/tls/ca.crt",
		PeerEndpoint: "localhost:7051",
		GatewayPeer:  "peer0.university.cred.com",
		MSPID:        "UniversityMSP",
	},

	"minifab-university": {
		CertPath:     "../Minifab_Network/vars/keyfiles/peerOrganizations/university.cred.com/users/Admin@university.cred.com/msp/signcerts/Admin@university.cred.com-cert.pem",
		KeyDirectory: "../Minifab_Network/vars/keyfiles/peerOrganizations/university.cred.com/users/Admin@university.cred.com/msp/keystore/",
		TLSCertPath:  "../Minifab_Network/vars/keyfiles/peerOrganizations/university.cred.com/peers/peer1.university.cred.com/tls/ca.crt",
		PeerEndpoint: "localhost:7003",
		GatewayPeer:  "peer1.university.cred.com",
		MSPID:        "university-cred-com",
	},

	"minifab-company": {
		CertPath:     "../Minifab_Network/vars/keyfiles/peerOrganizations/company.cred.com/users/Admin@company.cred.com/msp/signcerts/Admin@company.cred.com-cert.pem",
		KeyDirectory: "../Minifab_Network/vars/keyfiles/peerOrganizations/company.cred.com/users/Admin@company.cred.com/msp/keystore/",
		TLSCertPath:  "../Minifab_Network/vars/keyfiles/peerOrganizations/company.cred.com/peers/peer1.company.cred.com/tls/ca.crt",
		PeerEndpoint: "localhost:7004",
		GatewayPeer:  "peer1.company.cred.com",
		MSPID:        "company-cred-com",
	},

	"minifab-student": {
		CertPath:     "../Minifab_Network/vars/keyfiles/peerOrganizations/student.cred.com/users/Admin@student.cred.com/msp/signcerts/Admin@student.cred.com-cert.pem",
		KeyDirectory: "../Minifab_Network/vars/keyfiles/peerOrganizations/student.cred.com/users/Admin@student.cred.com/msp/keystore/",
		TLSCertPath:  "../Minifab_Network/vars/keyfiles/peerOrganizations/student.cred.com/peers/peer1.student.cred.com/tls/ca.crt",
		PeerEndpoint: "localhost:7005",
		GatewayPeer:  "peer1.student.cred.com",
		MSPID:        "student-cred-com",
	},
}