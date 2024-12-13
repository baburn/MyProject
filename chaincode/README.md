export CHANNEL_NAME=mychannel
export FABRIC_CFG_PATH=./peercfg
export CORE_PEER_LOCALMSPID=UniversityMSP
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/university.cred.com/peers/peer0.university.cred.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/university.cred.com/users/Admin@university.cred.com/msp
export CORE_PEER_ADDRESS=localhost:7051
export ORDERER_CA=${PWD}/organizations/ordererOrganizations/cred.com/orderers/orderer.cred.com/msp/tlscacerts/tlsca.cred.com-cert.pem
export UNIVERSITY_PEER_TLSROOTCERT=${PWD}/organizations/peerOrganizations/university.cred.com/peers/peer0.university.cred.com/tls/ca.crt
export STUDENT_PEER_TLSROOTCERT=${PWD}/organizations/peerOrganizations/student.cred.com/peers/peer0.student.cred.com/tls/ca.crt
export COMPANY_PEER_TLSROOTCERT=${PWD}/organizations/peerOrganizations/company.cred.com/peers/peer0.company.cred.com/tls/ca.crt



peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.cred.com --tls --cafile $ORDERER_CA -C $CHANNEL_NAME -n Credential-Verification --peerAddresses localhost:7051 --tlsRootCertFiles $UNIVERSITY_PEER_TLSROOTCERT --peerAddresses localhost:9051 --tlsRootCertFiles $STUDENT_PEER_TLSROOTCERT --peerAddresses localhost:11051 --tlsRootCertFiles $COMPANY_PEER_TLSROOTCERT -c '{"function":"CreateResult","Args":["RES1", "Stu1", "100", "90", "90%", "Pass"]}'

Credential-Verification

peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.cred.com --tls --cafile $ORDERER_CA -C $CHANNEL_NAME -n Credential-Verification --peerAddresses localhost:7051 --tlsRootCertFiles $UNIVERSITY_PEER_TLSROOTCERT --peerAddresses localhost:9051 --tlsRootCertFiles $STUDENT_PEER_TLSROOTCERT --peerAddresses localhost:11051 --tlsRootCertFiles $COMPANY_PEER_TLSROOTCERT -c '{"function":"CreateResult","Args":["RES2", "Stu2", "100", "80", "80%", "Pass"]}'

peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.cred.com --tls --cafile $ORDERER_CA -C $CHANNEL_NAME -n Credential-Verification --peerAddresses localhost:7051 --tlsRootCertFiles $UNIVERSITY_PEER_TLSROOTCERT --peerAddresses localhost:9051 --tlsRootCertFiles $STUDENT_PEER_TLSROOTCERT --peerAddresses localhost:11051 --tlsRootCertFiles $COMPANY_PEER_TLSROOTCERT -c '{"function":"CreateResult","Args":["RES3", "Stu2", "100", "95", "95%", "Pass"]}'

peer chaincode query -C $CHANNEL_NAME -n Credential-Verification -c '{"function":"ReadResult","Args":["RES1", "Stu1"]}'

peer chaincode query -C $CHANNEL_NAME -n Credential-Verification -c '{"function":"GetResultHistory","Args":["RES1"]}'

peer chaincode query -C $CHANNEL_NAME -n Credential-Verification -c '{"function":"GetResultsByRange","Args":["RES1","RES3"]}'

peer chaincode query -C $CHANNEL_NAME -n Credential-Verification -c '{"function":"GetResultsWithPagination","Args":[3,"3"]}'


peer chaincode query -C $CHANNEL_NAME -n Credential-Verification -c '{"Args":["GetAllResults"]}'

export CHANNEL_NAME=mychannel
export FABRIC_CFG_PATH=./peercfg
export CORE_PEER_LOCALMSPID=CompanyMSP
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/company.cred.com/peers/peer0.company.cred.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/company.cred.com/users/Admin@company.cred.com/msp
export CORE_PEER_ADDRESS=localhost:7051
export ORDERER_CA=${PWD}/organizations/ordererOrganizations/cred.com/orderers/orderer.cred.com/msp/tlscacerts/tlsca.cred.com-cert.pem
export UNIVERSITY_PEER_TLSROOTCERT=${PWD}/organizations/peerOrganizations/university.cred.com/peers/peer0.university.cred.com/tls/ca.crt
export STUDENT_PEER_TLSROOTCERT=${PWD}/organizations/peerOrganizations/student.cred.com/peers/peer0.student.cred.com/tls/ca.crt
export COMPANY_PEER_TLSROOTCERT=${PWD}/organizations/peerOrganizations/company.cred.com/peers/peer0.company.cred.com/tls/ca.crt


export CTC=$(echo -n "9LPA" | base64 | tr -d \\n)

export DATEOFJOINING=$(echo -n "01/01/2025" | base64 | tr -d \\n)

export DATEOFRELEASE=$(echo -n "19/12/2025" | base64 | tr -d \\n)

export COMPANYNAME=$(echo -n "XXX" | base64 | tr -d \\n)

peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.cred.com --tls --cafile $ORDERER_CA -C $CHANNEL_NAME -n Credential-Verification --peerAddresses localhost:7051 --tlsRootCertFiles $UNIVERSITY_PEER_TLSROOTCERT --peerAddresses localhost:9051 --tlsRootCertFiles $STUDENT_PEER_TLSROOTCERT --peerAddresses localhost:11051 --tlsRootCertFiles $COMPANY_PEER_TLSROOTCERT -c '{"Args":["OfferContract:CreateOffer","Offer1"]}' --transient "{\"ctc\":\"$CTC\",\"dateOfJoining\":\"$DATEOFJOINING\",\"dateOfRelease\":\"$DATEOFRELEASE\",\"companyName\":\"$COMPANYNAME\"}"
