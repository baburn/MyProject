#!/bin/bash

echo "------------Register the ca admin for each organization—----------------"

docker compose -f docker/docker-compose-ca.yaml up -d
sleep 3

sudo chmod -R 777 organizations/

echo "------------Register and enroll the users for each organization—-----------"

chmod +x registerEnroll.sh

./registerEnroll.sh
sleep 3

echo "—-------------Build the infrastructure—-----------------"

docker compose -f docker/docker-compose-3org.yaml up -d
sleep 3

echo "-------------Generate the genesis block—-------------------------------"

export FABRIC_CFG_PATH=${PWD}/config

export CHANNEL_NAME=mychannel

configtxgen -profile ThreeOrgsChannel -outputBlock ${PWD}/channel-artifacts/${CHANNEL_NAME}.block -channelID $CHANNEL_NAME
sleep 2

echo "------ Create the application channel------"

export ORDERER_CA=${PWD}/organizations/ordererOrganizations/cred.com/orderers/orderer.cred.com/msp/tlscacerts/tlsca.cred.com-cert.pem

export ORDERER_ADMIN_TLS_SIGN_CERT=${PWD}/organizations/ordererOrganizations/cred.com/orderers/orderer.cred.com/tls/server.crt

export ORDERER_ADMIN_TLS_PRIVATE_KEY=${PWD}/organizations/ordererOrganizations/cred.com/orderers/orderer.cred.com/tls/server.key

osnadmin channel join --channelID $CHANNEL_NAME --config-block ${PWD}/channel-artifacts/$CHANNEL_NAME.block -o localhost:7053 --ca-file $ORDERER_CA --client-cert $ORDERER_ADMIN_TLS_SIGN_CERT --client-key $ORDERER_ADMIN_TLS_PRIVATE_KEY
sleep 2

osnadmin channel list -o localhost:7053 --ca-file $ORDERER_CA --client-cert $ORDERER_ADMIN_TLS_SIGN_CERT --client-key $ORDERER_ADMIN_TLS_PRIVATE_KEY
sleep 2

export FABRIC_CFG_PATH=${PWD}/peercfg
export CORE_PEER_LOCALMSPID=UniversityMSP
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/university.cred.com/peers/peer0.university.cred.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/university.cred.com/users/Admin@university.cred.com/msp
export CORE_PEER_ADDRESS=localhost:7051
export UNIVERSITY_PEER_TLSROOTCERT=${PWD}/organizations/peerOrganizations/university.cred.com/peers/peer0.university.cred.com/tls/ca.crt
export STUDENT_PEER_TLSROOTCERT=${PWD}/organizations/peerOrganizations/student.cred.com/peers/peer0.student.cred.com/tls/ca.crt
export COMPANY_PEER_TLSROOTCERT=${PWD}/organizations/peerOrganizations/company.cred.com/peers/peer0.company.cred.com/tls/ca.crt
sleep 2

echo "—---------------Join Institution peer to the channel—-------------"

echo "-----------------------------------------------------------------------------------------------------------------------------------------------------------"
echo ${FABRIC_CFG_PATH}
echo "-----------------------------------------------------------------------------------------------------------------------------------------------------------"


sleep 2
peer channel join -b ${PWD}/channel-artifacts/${CHANNEL_NAME}.block
sleep 3

echo "-----channel List----"
peer channel list

echo "—-------------Institution anchor peer update—-----------"

peer channel fetch config ${PWD}/channel-artifacts/config_block.pb -o localhost:7050 --ordererTLSHostnameOverride orderer.cred.com -c $CHANNEL_NAME --tls --cafile $ORDERER_CA
sleep 1

cd channel-artifacts

configtxlator proto_decode --input config_block.pb --type common.Block --output config_block.json
jq '.data.data[0].payload.data.config' config_block.json > config.json

cp config.json config_copy.json

jq '.channel_group.groups.Application.groups.UniversityMSP.values += {"AnchorPeers":{"mod_policy": "Admins","value":{"anchor_peers": [{"host": "peer0.university.cred.com","port": 7051}]},"version": "0"}}' config_copy.json > modified_config.json

configtxlator proto_encode --input config.json --type common.Config --output config.pb
configtxlator proto_encode --input modified_config.json --type common.Config --output modified_config.pb
configtxlator compute_update --channel_id ${CHANNEL_NAME} --original config.pb --updated modified_config.pb --output config_update.pb

configtxlator proto_decode --input config_update.pb --type common.ConfigUpdate --output config_update.json
echo '{"payload":{"header":{"channel_header":{"channel_id":"'$CHANNEL_NAME'", "type":2}},"data":{"config_update":'$(cat config_update.json)'}}}' | jq . > config_update_in_envelope.json
configtxlator proto_encode --input config_update_in_envelope.json --type common.Envelope --output config_update_in_envelope.pb

cd ..

peer channel update -f ${PWD}/channel-artifacts/config_update_in_envelope.pb -c $CHANNEL_NAME -o localhost:7050  --ordererTLSHostnameOverride orderer.cred.com --tls --cafile $ORDERER_CA
sleep 1






echo "—---------------package chaincode—-------------"

peer lifecycle chaincode package credentialcred.tar.gz --path ${PWD}/../chaincode/ --lang golang --label credentialcred_1.0
sleep 1

echo "—---------------install chaincode in Institution peer—-------------"

peer lifecycle chaincode install credentialcred.tar.gz
sleep 3

peer lifecycle chaincode queryinstalled
sleep 1

export CC_PACKAGE_ID=$(peer lifecycle chaincode calculatepackageid credentialcred.tar.gz)

echo "—---------------Approve chaincode in Institution peer—-------------"

peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.cred.com --channelID $CHANNEL_NAME --name Credential-Verification --version 1.0 --collections-config ../chaincode/collection-config.json --package-id $CC_PACKAGE_ID --sequence 1 --tls --cafile $ORDERER_CA --waitForEvent
sleep 2




export CORE_PEER_LOCALMSPID=StudentMSP 
export CORE_PEER_ADDRESS=localhost:9051 
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/student.cred.com/peers/peer0.student.cred.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/student.cred.com/users/Admin@student.cred.com/msp

echo "—---------------Join Student peer to the channel—-------------"

peer channel join -b ${PWD}/channel-artifacts/$CHANNEL_NAME.block
sleep 1
peer channel list

echo "—-------------Student anchor peer update—-----------"

peer channel fetch config ${PWD}/channel-artifacts/config_block.pb -o localhost:7050 --ordererTLSHostnameOverride orderer.cred.com -c $CHANNEL_NAME --tls --cafile $ORDERER_CA
sleep 1

cd channel-artifacts

configtxlator proto_decode --input config_block.pb --type common.Block --output config_block.json
jq '.data.data[0].payload.data.config' config_block.json > config.json
cp config.json config_copy.json

jq '.channel_group.groups.Application.groups.StudentMSP.values += {"AnchorPeers":{"mod_policy": "Admins","value":{"anchor_peers": [{"host": "peer0.student.cred.com","port": 9051}]},"version": "0"}}' config_copy.json > modified_config.json

configtxlator proto_encode --input config.json --type common.Config --output config.pb
configtxlator proto_encode --input modified_config.json --type common.Config --output modified_config.pb
configtxlator compute_update --channel_id $CHANNEL_NAME --original config.pb --updated modified_config.pb --output config_update.pb

configtxlator proto_decode --input config_update.pb --type common.ConfigUpdate --output config_update.json
echo '{"payload":{"header":{"channel_header":{"channel_id":"'$CHANNEL_NAME'", "type":2}},"data":{"config_update":'$(cat config_update.json)'}}}' | jq . > config_update_in_envelope.json
configtxlator proto_encode --input config_update_in_envelope.json --type common.Envelope --output config_update_in_envelope.pb

cd ..

peer channel update -f ${PWD}/channel-artifacts/config_update_in_envelope.pb -c $CHANNEL_NAME -o localhost:7050  --ordererTLSHostnameOverride orderer.cred.com --tls --cafile $ORDERER_CA
sleep 1


echo "—---------------install chaincode in Student peer—-------------"

peer lifecycle chaincode install credentialcred.tar.gz
sleep 3

peer lifecycle chaincode queryinstalled

echo "—---------------Approve chaincode in Student peer—-------------"

peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.cred.com --channelID $CHANNEL_NAME --name Credential-Verification --version 1.0 --collections-config ../chaincode/collection-config.json --package-id $CC_PACKAGE_ID --sequence 1 --tls --cafile $ORDERER_CA --waitForEvent
sleep 1

export CORE_PEER_LOCALMSPID=CompanyMSP 
export CORE_PEER_ADDRESS=localhost:11051 
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/company.cred.com/peers/peer0.company.cred.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/company.cred.com/users/Admin@company.cred.com/msp

echo "—---------------Join company peer to the channel—-------------"

peer channel join -b ${PWD}/channel-artifacts/$CHANNEL_NAME.block
sleep 1
peer channel list




echo "—-------------Company anchor peer update—-----------"
peer channel fetch config ${PWD}/channel-artifacts/config_block.pb -o localhost:7050 --ordererTLSHostnameOverride orderer.cred.com -c $CHANNEL_NAME --tls --cafile $ORDERER_CA
sleep 1

cd channel-artifacts

configtxlator proto_decode --input config_block.pb --type common.Block --output config_block.json
jq '.data.data[0].payload.data.config' config_block.json > config.json
cp config.json config_copy.json

jq '.channel_group.groups.Application.groups.CompanyMSP.values += {"AnchorPeers":{"mod_policy": "Admins","value":{"anchor_peers": [{"host": "peer0.company.cred.com","port": 11051}]},"version": "0"}}' config_copy.json > modified_config.json

configtxlator proto_encode --input config.json --type common.Config --output config.pb
configtxlator proto_encode --input modified_config.json --type common.Config --output modified_config.pb
configtxlator compute_update --channel_id $CHANNEL_NAME --original config.pb --updated modified_config.pb --output config_update.pb

configtxlator proto_decode --input config_update.pb --type common.ConfigUpdate --output config_update.json
echo '{"payload":{"header":{"channel_header":{"channel_id":"'$CHANNEL_NAME'", "type":2}},"data":{"config_update":'$(cat config_update.json)'}}}' | jq . > config_update_in_envelope.json
configtxlator proto_encode --input config_update_in_envelope.json --type common.Envelope --output config_update_in_envelope.pb

cd ..

peer channel update -f ${PWD}/channel-artifacts/config_update_in_envelope.pb -c $CHANNEL_NAME -o localhost:7050  --ordererTLSHostnameOverride orderer.cred.com --tls --cafile $ORDERER_CA
sleep 1

peer channel getinfo -c $CHANNEL_NAME

echo "—---------------install chaincode in Company peer—-------------"

peer lifecycle chaincode install credentialcred.tar.gz
sleep 3

peer lifecycle chaincode queryinstalled

echo "—---------------Approve chaincode in Company peer—-------------"

peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.cred.com --channelID $CHANNEL_NAME --name Credential-Verification --version 1.0 --collections-config ../chaincode/collection-config.json --package-id $CC_PACKAGE_ID --sequence 1 --tls --cafile $ORDERER_CA --waitForEvent
sleep 1

echo "—---------------Commit chaincode in Company peer—-------------"

peer lifecycle chaincode checkcommitreadiness --channelID $CHANNEL_NAME --name Credential-Verification --version 1.0 --sequence 1 --collections-config ../chaincode/collection-config.json --tls --cafile $ORDERER_CA --output json

peer lifecycle chaincode commit -o localhost:7050 --ordererTLSHostnameOverride orderer.cred.com --channelID $CHANNEL_NAME --name Credential-Verification --version 1.0 --sequence 1 --collections-config ../chaincode/collection-config.json --tls --cafile $ORDERER_CA --peerAddresses localhost:7051 --tlsRootCertFiles $UNIVERSITY_PEER_TLSROOTCERT --peerAddresses localhost:9051 --tlsRootCertFiles $STUDENT_PEER_TLSROOTCERT --peerAddresses localhost:11051 --tlsRootCertFiles $COMPANY_PEER_TLSROOTCERT
sleep 1

peer lifecycle chaincode querycommitted --channelID $CHANNEL_NAME --name Credential-Verification --cafile $ORDERER_CA


echo "-----------------------------------------------------------------------------------------------------------------------------------------------------------"
echo ${FABRIC_CFG_PATH} ${CHANNEL_NAME}
echo "-----------------------------------------------------------------------------------------------------------------------------------------------------------"

