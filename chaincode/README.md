### Set environment variables for channel name and configuration path
### Set the local MSP ID (Membership Service Provider) for the University peer
### Enable TLS (Transport Layer Security) and provide the University peer's root certificate
### Set the MSP configuration path for the University peer
### Set the address of the University peer
### Set the Orderer's certificate path
### Set the TLS root certificates for University, Student, and Company peers

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

### Invoke the chaincode function "CreateResult" to create a result for student "Stu1"
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.cred.com --tls --cafile $ORDERER_CA -C $CHANNEL_NAME -n Credential-Verification --peerAddresses localhost:7051 --tlsRootCertFiles $UNIVERSITY_PEER_TLSROOTCERT --peerAddresses localhost:9051 --tlsRootCertFiles $STUDENT_PEER_TLSROOTCERT --peerAddresses localhost:11051 --tlsRootCertFiles $COMPANY_PEER_TLSROOTCERT -c '{"function":"CreateResult","Args":["RES1", "Stu1", "100", "90", "90%", "Pass"]}'

### Invoke the chaincode function "CreateResult" to create a result for student "Stu2" with a "Fail" status
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.cred.com --tls --cafile $ORDERER_CA -C $CHANNEL_NAME -n Credential-Verification --peerAddresses localhost:7051 --tlsRootCertFiles $UNIVERSITY_PEER_TLSROOTCERT --peerAddresses localhost:9051 --tlsRootCertFiles $STUDENT_PEER_TLSROOTCERT --peerAddresses localhost:11051 --tlsRootCertFiles $COMPANY_PEER_TLSROOTCERT -c '{"function":"CreateResult","Args":["RES2", "Stu2", "100", "30", "30%", "Fail"]}'

### Invoke the chaincode function "CreateResult" to create another result for student "Stu2" with a "Pass" status
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.cred.com --tls --cafile $ORDERER_CA -C $CHANNEL_NAME -n Credential-Verification --peerAddresses localhost:7051 --tlsRootCertFiles $UNIVERSITY_PEER_TLSROOTCERT --peerAddresses localhost:9051 --tlsRootCertFiles $STUDENT_PEER_TLSROOTCERT --peerAddresses localhost:11051 --tlsRootCertFiles $COMPANY_PEER_TLSROOTCERT -c '{"function":"CreateResult","Args":["RES3", "Stu3", "100", "95", "95%", "Pass"]}'

### Query the chaincode to read the result for RES1 (student "Stu1")
peer chaincode query -C $CHANNEL_NAME -n Credential-Verification -c '{"function":"ReadResult","Args":["RES1", "Stu1"]}'

### Query the chaincode to get the results in the specified range (RES1 to RES3)
peer chaincode query -C $CHANNEL_NAME -n Credential-Verification -c '{"function":"GetResultsByRange","Args":["RES1","RES3"]}'

### Query the chaincode to get all results
peer chaincode query -C $CHANNEL_NAME -n Credential-Verification -c '{"Args":["GetAllResults"]}'

### Query the chaincode to get the result history for RES1
peer chaincode query -C $CHANNEL_NAME -n Credential-Verification -c '{"function":"GetResultHistory","Args":["RES1"]}'


### Query the chaincode to get paginated results (fetching 3 results per page, starting from page 3)

peer chaincode query -C $CHANNEL_NAME -n Credential-Verification -c '{"function":"GetResultsWithPagination","Args":["3", ""]}'
peer chaincode query -C $CHANNEL_NAME -n Credential-Verification -c '{"function":"GetResultsWithPagination","Args":["3", "someBookmarkValue"]}'


### Switch to the Company peer context by setting relevant environment variables
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

### Encode sensitive data such as CTC, Date of Joining, etc. to base64 format for use in transactions
export CTC=$(echo -n "9LPA" | base64 | tr -d \\n)
export DATEOFJOINING=$(echo -n "01/01/2025" | base64 | tr -d \\n)
export DATEOFRELEASE=$(echo -n "19/12/2025" | base64 | tr -d \\n)
export NAME=$(echo -n "Ram" | base64 | tr -d \\n)
export EMAIL=$(echo -n "ram@gmail.com" | base64 | tr -d \\n)
export COMPANYNAME=$(echo -n "NPCI" | base64 | tr -d \\n)

### Invoke the "CreateOffer" function on the OfferContract chaincode to create a job offer for "Offer1" using transient data
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.cred.com --tls --cafile $ORDERER_CA -C $CHANNEL_NAME -n Credential-Verification --peerAddresses localhost:7051 --tlsRootCertFiles $UNIVERSITY_PEER_TLSROOTCERT --peerAddresses localhost:9051 --tlsRootCertFiles $STUDENT_PEER_TLSROOTCERT --peerAddresses localhost:11051 --tlsRootCertFiles $COMPANY_PEER_TLSROOTCERT -c '{"Args":["OfferContract:CreateOffer","Offer1"]}' --transient "{\"ctc\":\"$CTC\",\"dateOfJoining\":\"$DATEOFJOINING\",\"dateOfRelease\":\"$DATEOFRELEASE\",\"name\":\"$NAME\",\"email\":\"$EMAIL\",\"companyName\":\"$COMPANYNAME\"}"

export CTC=$(echo -n "9LPA" | base64 | tr -d \\n)
export DATEOFJOINING=$(echo -n "01/01/2025" | base64 | tr -d \\n)
export DATEOFRELEASE=$(echo -n "19/12/2025" | base64 | tr -d \\n)
export NAME=$(echo -n "sam" | base64 | tr -d \\n)
export EMAIL=$(echo -n "sam@gmail.com" | base64 | tr -d \\n)
export COMPANYNAME=$(echo -n "NPCI" | base64 | tr -d \\n)

### Invoke the "CreateOffer" function on the OfferContract chaincode to create another job offer for "Offer2" using transient data
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.cred.com --tls --cafile $ORDERER_CA -C $CHANNEL_NAME -n Credential-Verification --peerAddresses localhost:7051 --tlsRootCertFiles $UNIVERSITY_PEER_TLSROOTCERT --peerAddresses localhost:9051 --tlsRootCertFiles $STUDENT_PEER_TLSROOTCERT --peerAddresses localhost:11051 --tlsRootCertFiles $COMPANY_PEER_TLSROOTCERT -c '{"Args":["OfferContract:CreateOffer","Offer2"]}' --transient "{\"ctc\":\"$CTC\",\"dateOfJoining\":\"$DATEOFJOINING\",\"dateOfRelease\":\"$DATEOFRELEASE\",\"name\":\"$NAME\",\"email\":\"$EMAIL\",\"companyName\":\"$COMPANYNAME\"}"

### Query the chaincode to read the offer details for "Offer1"
peer chaincode query -C $CHANNEL_NAME -n Credential-Verification --peerAddresses localhost:11051 --tlsRootCertFiles $COMPANY_PEER_TLSROOTCERT -c '{"Args":["OfferContract:ReadOffer","Offer1"]}'
