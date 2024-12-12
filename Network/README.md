##### ***Open a command terminal with in Assignment folder, let's call this terminal as host terminal

```
cd Assignment/

```

##  **************** host terminal ********************

### ------------Register the ca admin for each organization—----------------

##### ***Build the `docker-compose-ca.yaml` in the docker folder

```
docker compose -f docker/docker-compose-ca.yaml up -d

```

```
sudo chmod -R 777 organizations/
```
### ------------Register and enroll the users for each organization—-----------

##### ***Build the `registerEnroll.sh` script file


```
chmod +x registerEnroll.sh

./registerEnroll.sh
```


### —-------------Build the infrastructure—-----------------

##### ***Build the `docker-compose-3org.yaml` in the docker folder

```
docker compose -f docker/docker-compose-3org.yaml up -d
```


### -------------Generate the genesis block—-------------------------------

##### ***Build the `configtx.yaml` file in the config folder

```
export FABRIC_CFG_PATH=./config

export CHANNEL_NAME=credchannel
```

```
configtxgen -profile ThreeOrgsChannel -outputBlock ./channel-artifacts/${CHANNEL_NAME}.block -channelID $CHANNEL_NAME
```


### ------ Create the application channel------

```
export ORDERER_CA=./organizations/ordererOrganizations/cred.com/orderers/orderer.cred.com/msp/tlscacerts/tlsca.cred.com-cert.pem

export ORDERER_ADMIN_TLS_SIGN_CERT=./organizations/ordererOrganizations/cred.com/orderers/orderer.cred.com/tls/server.crt

export ORDERER_ADMIN_TLS_PRIVATE_KEY=./organizations/ordererOrganizations/cred.com/orderers/orderer.cred.com/tls/server.key
```

```
osnadmin channel join --channelID $CHANNEL_NAME --config-block ./channel-artifacts/$CHANNEL_NAME.block -o localhost:7053 --ca-file $ORDERER_CA --client-cert $ORDERER_ADMIN_TLS_SIGN_CERT --client-key $ORDERER_ADMIN_TLS_PRIVATE_KEY
```


```
osnadmin channel list -o localhost:7053 --ca-file $ORDERER_CA --client-cert $ORDERER_ADMIN_TLS_SIGN_CERT --client-key $ORDERER_ADMIN_TLS_PRIVATE_KEY
```

##### ***Open another terminal with in Assignment folder, let's call this terminal as peer0_University terminal.

##  **************** peer0_University terminal ********************

##### ***Build the `core.yaml` in peercfg folder

```
export FABRIC_CFG_PATH=./peercfg
export CHANNEL_NAME=credchannel
export CORE_PEER_LOCALMSPID=UniversityrMSP
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/universityr.cred.com/peers/peer0.universityr.cred.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/universityr.cred.com/users/Admin@universityr.cred.com/msp
export CORE_PEER_ADDRESS=localhost:7051
export ORDERER_CA=${PWD}/organizations/ordererOrganizations/cred.com/orderers/orderer.cred.com/msp/tlscacerts/tlsca.cred.com-cert.pem
export UNIVERSITY=${PWD}/organizations/peerOrganizations/university.cred.com/peers/peer0.university.cred.com/tls/ca.crt
export STUDENT_PEER_TLSROOTCERT=${PWD}/organizations/peerOrganizations/student.cred.com/peers/peer0.student.cred.com/tls/ca.crt
export ver_PEER_TLSROOTCERT=${PWD}/organizations/peerOrganizations/ver.cred.com/peers/peer0.ver.cred.com/tls/ca.crt
```


### —---------------Join peer to the channel—-------------

```
peer channel join -b ./channel-artifacts/$CHANNEL_NAME.block

peer channel list
```


##### ***Open another terminal with in Assignment folder, let's call this terminal as peer0_Student terminal.

##  **************** peer0_Student terminal *****************

```
export FABRIC_CFG_PATH=./peercfg
export CHANNEL_NAME=credchannel 
export CORE_PEER_LOCALMSPID=StudentMSP 
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_ADDRESS=localhost:9051 
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/student.cred.com/peers/peer0.student.cred.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/student.cred.com/users/Admin@student.cred.com/msp
export ORDERER_CA=${PWD}/organizations/ordererOrganizations/cred.com/orderers/orderer.cred.com/msp/tlscacerts/tlsca.cred.com-cert.pem
export UNIVERSITY_PEER_TLSROOTCERT=${PWD}/organizations/peerOrganizations/university.cred.com/peers/peer0.university.cred.com/tls/ca.crt
export STUDENT_PEER_TLSROOTCERT=${PWD}/organizations/peerOrganizations/student.cred.com/peers/peer0.student.cred.com/tls/ca.crt
export ver_PEER_TLSROOTCERT=${PWD}/organizations/peerOrganizations/ver.cred.com/peers/peer0.ver.cred.com/tls/ca.crt
```


### —------------ Join peer to the channel ---------------

```
peer channel join -b ./channel-artifacts/$CHANNEL_NAME.block

peer channel list
```

##### ***Open another terminal with in Assignment folder, let's call this terminal as peer0_Ver terminal.

##  **************** peer0_Ver terminal ******************

```
export FABRIC_CFG_PATH=./peercfg
export CHANNEL_NAME=credchannel 
export CORE_PEER_LOCALMSPID=VerMSP 
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_ADDRESS=localhost:11051 
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/ver.cred.com/peers/peer0.ver.cred.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/ver.cred.com/users/Admin@ver.cred.com/msp
export ORDERER_CA=${PWD}/organizations/ordererOrganizations/cred.com/orderers/orderer.cred.com/msp/tlscacerts/tlsca.cred.com-cert.pem
export UNIVERSITY_PEER_TLSROOTCERT=${PWD}/organizations/peerOrganizations/university.cred.com/peers/peer0.university.cred.com/tls/ca.crt
export STUDENT_PEER_TLSROOTCERT=${PWD}/organizations/peerOrganizations/student.cred.com/peers/peer0.student.cred.com/tls/ca.crt
export ver_PEER_TLSROOTCERT=${PWD}/organizations/peerOrganizations/ver.cred.com/peers/peer0.ver.cred.com/tls/ca.crt
```

### —-------------- Join peer to the channel ----------------------

```
peer channel join -b ./channel-artifacts/$CHANNEL_NAME.block

peer channel list
```

### —-------------------- anchor peer update —---------------------

##  **************** peer0_University terminal ******************

```
peer channel fetch config channel-artifacts/config_block.pb -o localhost:7050 --ordererTLSHostnameOverride orderer.cred.com -c $CHANNEL_NAME --tls --cafile $ORDERER_CA
```

```
cd channel-artifacts
```

```
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
```

```
cd ..
```

```
peer channel update -f channel-artifacts/config_update_in_envelope.pb -c $CHANNEL_NAME -o localhost:7050  --ordererTLSHostnameOverride orderer.cred.com --tls --cafile $ORDERER_CA
```
##  **************** peer0_Student terminal ******************

```
peer channel fetch config channel-artifacts/config_block.pb -o localhost:7050 --ordererTLSHostnameOverride orderer.cred.com -c $CHANNEL_NAME --tls --cafile $ORDERER_CA
```

```
cd channel-artifacts
```

```
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
```

```
cd ..
```

```
peer channel update -f channel-artifacts/config_update_in_envelope.pb -c $CHANNEL_NAME -o localhost:7050  --ordererTLSHostnameOverride orderer.cred.com --tls --cafile $ORDERER_CA
```
##  **************** peer0_Ver terminal ******************

```
peer channel fetch config channel-artifacts/config_block.pb -o localhost:7050 --ordererTLSHostnameOverride orderer.cred.com -c $CHANNEL_NAME --tls --cafile $ORDERER_CA
```

```
cd channel-artifacts
```

```
configtxlator proto_decode --input config_block.pb --type common.Block --output config_block.json
jq '.data.data[0].payload.data.config' config_block.json > config.json
cp config.json config_copy.json

jq '.channel_group.groups.Application.groups.VerMSP.values += {"AnchorPeers":{"mod_policy": "Admins","value":{"anchor_peers": [{"host": "peer0.ver.cred.com","port": 11051}]},"version": "0"}}' config_copy.json > modified_config.json

configtxlator proto_encode --input config.json --type common.Config --output config.pb
configtxlator proto_encode --input modified_config.json --type common.Config --output modified_config.pb
configtxlator compute_update --channel_id $CHANNEL_NAME --original config.pb --updated modified_config.pb --output config_update.pb

configtxlator proto_decode --input config_update.pb --type common.ConfigUpdate --output config_update.json
echo '{"payload":{"header":{"channel_header":{"channel_id":"'$CHANNEL_NAME'", "type":2}},"data":{"config_update":'$(cat config_update.json)'}}}' | jq . > config_update_in_envelope.json
configtxlator proto_encode --input config_update_in_envelope.json --type common.Envelope --output config_update_in_envelope.pb
```

```
cd ..
```

```
peer channel update -f channel-artifacts/config_update_in_envelope.pb -c $CHANNEL_NAME -o localhost:7050  --ordererTLSHostnameOverride orderer.cred.com --tls --cafile $ORDERER_CA
```
```
peer channel getinfo -c $CHANNEL_NAME
```


### —-----------------Chaincode lifecycle—-------------------

##### ***Build the chaincode (Change MSPId and collection file)

##### ***Make sure that KBA-Automobile chaincode is available in Chaincode folder which is at the same level of Assignment. 

##  **************** peer0_University terminal ******************

```
peer lifecycle chaincode package kbacred.tar.gz --path ../Chaincode/ --lang golang --label kbacred_1.0

```

```
peer lifecycle chaincode install kbacred.tar.gz
```

```
peer lifecycle chaincode queryinstalled

```

```
export CC_PACKAGE_ID=$(peer lifecycle chaincode calculatepackageid kbacred.tar.gz)

```

##  **************** peer0_Student terminal ******************
```
peer lifecycle chaincode install kbacred.tar.gz
```

```
export CC_PACKAGE_ID=$(peer lifecycle chaincode calculatepackageid kbacred.tar.gz)

```

##  **************** peer0_Ver terminal ******************
```
peer lifecycle chaincode install kbacred.tar.gz
```

```
export CC_PACKAGE_ID=$(peer lifecycle chaincode calculatepackageid kbacred.tar.gz)

```

##  **************** peer0_University terminal ******************

```
peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.cred.com --channelID $CHANNEL_NAME --name KBA-Automobile --version 1.0 --collections-config ../Chaincode/collection-credmobile.json --package-id $CC_PACKAGE_ID --sequence 1 --tls --cafile $ORDERER_CA --waitForEvent
```

##  **************** peer0_Student terminal ******************

```
peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.cred.com --channelID $CHANNEL_NAME --name KBA-Automobile --version 1.0 --collections-config ../Chaincode/collection-credmobile.json --package-id $CC_PACKAGE_ID --sequence 1 --tls --cafile $ORDERER_CA --waitForEvent
```

##  **************** peer0_Ver terminal ******************
```
peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.cred.com --channelID $CHANNEL_NAME --name KBA-Automobile --version 1.0 --collections-config ../Chaincode/collection-credmobile.json --package-id $CC_PACKAGE_ID --sequence 1 --tls --cafile $ORDERER_CA --waitForEvent
```

##  **************** peer0_University terminal ******************

```
peer lifecycle chaincode checkcommitreadiness --channelID $CHANNEL_NAME --name KBA-Automobile --version 1.0 --sequence 1 --collections-config ../Chaincode/collection-credmobile.json --tls --cafile $ORDERER_CA --output json
```

```
peer lifecycle chaincode commit -o localhost:7050 --ordererTLSHostnameOverride orderer.cred.com --channelID $CHANNEL_NAME --name KBA-Automobile --version 1.0 --sequence 1 --collections-config ../Chaincode/collection-credmobile.json --tls --cafile $ORDERER_CA --peerAddresses localhost:7051 --tlsRootCertFiles $UNIVERSITY_PEER_TLSROOTCERT --peerAddresses localhost:9051 --tlsRootCertFiles $STUDENT_PEER_TLSROOTCERT --peerAddresses localhost:11051 --tlsRootCertFiles $ver_PEER_TLSROOTCERT
```
```
peer lifecycle chaincode querycommitted --channelID $CHANNEL_NAME --name KBA-Automobile --cafile $ORDERER_CA
```
```
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.cred.com --tls --cafile $ORDERER_CA -C $CHANNEL_NAME -n KBA-Automobile --peerAddresses localhost:7051 --tlsRootCertFiles $UNIVERSITY_PEER_TLSROOTCERT --peerAddresses localhost:9051 --tlsRootCertFiles $STUDENT_PEER_TLSROOTCERT --peerAddresses localhost:11051 --tlsRootCertFiles $ver_PEER_TLSROOTCERT -c '{"function":"CreateCar","Args":["Car-11", "Tata", "Nexon", "White", "Factory-1", "22/07/2024"]}'
```

```
peer chaincode query -C $CHANNEL_NAME -n KBA-Automobile -c '{"Args":["GetAllCars"]}'

```


### --------Invoke Private Transaction----------

##  **************** peer0_Student terminal ******************

```
export MAKE=$(echo -n "Tata" | base64 | tr -d \\n)

export MODEL=$(echo -n "Tiago" | base64 | tr -d \\n)

export COLOR=$(echo -n "White" | base64 | tr -d \\n)

export STUDENT_NAME=$(echo -n "XXX" | base64 | tr -d \\n)
```

```
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.cred.com --tls --cafile $ORDERER_CA -C $CHANNEL_NAME -n KBA-Automobile --peerAddresses localhost:7051 --tlsRootCertFiles $UNIVERSITY_PEER_TLSROOTCERT --peerAddresses localhost:9051 --tlsRootCertFiles $STUDENT_PEER_TLSROOTCERT --peerAddresses localhost:11051 --tlsRootCertFiles $ver_PEER_TLSROOTCERT -c '{"Args":["OrderContract:CreateOrder","ORD201"]}' --transient "{\"make\":\"$MAKE\",\"model\":\"$MODEL\",\"color\":\"$COLOR\",\"studentName\":\"$STUDENT_NAME\"}"
```
```
peer chaincode query -C $CHANNEL_NAME -n KBA-Automobile -c '{"Args":["OrderContract:ReadOrder","ORD201"]}'
```


### --------- Stop the Assignment --------------

##  **************** Host terminal ******************

```
docker compose -f docker/docker-compose-3org.yaml down
```
```
docker compose -f docker/docker-compose-ca.yaml down
```
```
docker volume rm $(docker volume ls -q)
```
```
sudo rm -rf organizations/
```
```
sudo rm -rf channel-artifacts/
```
```
sudo rm kbacred.tar.gz
```
```
docker ps -a
```

##### ***if there still exists the containers then execute the following commands.

```
docker rm $(docker container ls -q) --force
```
```
docker container prune
```
```
docker system prune
```
```
docker volume prune
```
```
docker network prune
```


### Run using startAutomobileNetwork.sh script

##### Build `startAutomobileNetwork.sh` script file

```
chmod +x startAutomobileNetwork.sh
```

```
./startAutomobileNetwork.sh
```


##### To submit transaction as UniversityMSP
```
export CHANNEL_NAME=credchannel
export FABRIC_CFG_PATH=./peercfg
export CORE_PEER_LOCALMSPID=UniversityMSP
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/university.cred.com/peers/peer0.university.cred.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/university.cred.com/users/Admin@university.cred.com/msp
export CORE_PEER_ADDRESS=localhost:7051
export ORDERER_CA=${PWD}/organizations/ordererOrganizations/cred.com/orderers/orderer.cred.com/msp/tlscacerts/tlsca.cred.com-cert.pem
export UNIVERSITY_PEER_TLSROOTCERT=${PWD}/organizations/peerOrganizations/university.cred.com/peers/peer0.university.cred.com/tls/ca.crt
export STUDENT_PEER_TLSROOTCERT=${PWD}/organizations/peerOrganizations/student.cred.com/peers/peer0.student.cred.com/tls/ca.crt
export ver_PEER_TLSROOTCERT=${PWD}/organizations/peerOrganizations/ver.cred.com/peers/peer0.ver.cred.com/tls/ca.crt
```

```
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.cred.com --tls --cafile $ORDERER_CA -C $CHANNEL_NAME -n KBA-Automobile --peerAddresses localhost:7051 --tlsRootCertFiles $UNIVERSITY_PEER_TLSROOTCERT --peerAddresses localhost:9051 --tlsRootCertFiles $STUDENT_PEER_TLSROOTCERT --peerAddresses localhost:11051 --tlsRootCertFiles $ver_PEER_TLSROOTCERT -c '{"function":"CreateCar","Args":["Car-101", "Tata", "Nexon", "White", "Factory-1", "22/07/2024"]}'
```

```
peer chaincode query -C $CHANNEL_NAME -n KBA-Automobile -c '{"Args":["GetAllCars"]}'
```

### To stop the network using script file

##### Build `stopAutomobileNetwork.sh` script file
```
chmod +x stopAutomobileNetwork.sh
```
```
./stopAutomobileNetwork.sh
```

### To stop the network without clearing the data using docker

```
docker compose -f ./docker/docker-compose-ca.yaml -f ./docker/docker-compose-3org.yaml stop
```
### To start the network with existing docker containers
```
docker compose -f ./docker/docker-compose-ca.yaml -f ./docker/docker-compose-3org.yaml start
```
#OUTPUT SCREENSHOT

![alt text](<Screenshot from 2024-12-06 08-37-42.png>)