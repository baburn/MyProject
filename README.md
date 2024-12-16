# Academic Credential Verification


# Installing the Dependencies

Note: If any of the following dependencies are available on your laptop, then no need to install it.

## Update Packages

In case of a fresh Ubuntu 22 installation, use the following commands to update the packages before installing other dependencies.  
```
sudo apt update
```

```
sudo apt upgrade
```

## Visual Studio Code
Download and install the latest version of VS code from here: https://code.visualstudio.com/download


To install, execute the following command from the same folder where VS Code is being downloaded.

Note: Replace file_name with the actual name of the file you've downloaded.
```
sudo dpkg -i file_name
```
eg: sudo dpkg -i code_1.95.2-1730981514_amd64.deb


## cURL
Install curl using the command
```
sudo apt install curl -y
```

```
curl -V
```

## NVM

Install NVM (Node Version Manager), open a terminal and execute the following command.
```
curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.40.0/install.sh | bash
```
Close the current terminal and open a new one.

In the new terminal execute this command to verify nvm has been installed

```
nvm -v
```


## NodeJS (Ver 22.x)

Execute the following command to install NodeJs
```
nvm install 22
```  

Check  the version of nodeJS installed
```
node -v
```

Check  the version of npm installed
```
npm -v
```

## Docker
Step 1: Download the script
```
curl -fsSL https://get.docker.com -o install-docker.sh
```
Step 2: Run the script either as root, or using sudo to perform the installation.
```
sudo sh install-docker.sh
```
Step 2: To manage Docker as a non-root user
```
sudo chmod 777 /var/run/docker.sock
```

```
sudo usermod -aG docker $USER
```

To verify the installtion enter the following commands


```
docker compose version
```

```
docker -v
```

Execute the following command to check whether we can execute docker commands without sudo

```
docker ps -a
```

## JQ
Install JQ using the following command
```
sudo apt install jq -y
```

To verify the installtion enter the following command


```
jq --version
```

## Build Essential
Install Build Essential uisng the commnad
```
sudo apt install build-essential -y
```
To verify the installtion enter the following command


```
dpkg -l | grep build-essential
```

## Go
Step 1: Download Go
```
curl -OL  https://go.dev/dl/go1.22.0.linux-amd64.tar.gz
```
Step 2: Extract
```
sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf go1.22.0.linux-amd64.tar.gz
```

Step 3: Add /usr/local/go/bin to the PATH environment variable. Open the /etc/environment file
```
sudo gedit /etc/environment
```

Step 4: Append the following to the end of `PATH` variable and save
```
:/usr/local/go/bin
```
```
source $HOME/.profile
```

To verify the installtion enter the following command


```
go version
```
Note: If go version is not listed, then restart the system and execute the command again.


# Instructions for setting the network 
Go to the Network directory and using the termianl to execute the commands
```
chmod +x startnetwork.sh 
```

use chmod +x to set scipt file as a exucutable file Note: Use sudo if permisson denied then enter your password

check the containers by using
```
docker ps -a
```

When the network is successfully runing and chaincode is working and installed

Set Up the environmental Variables

```
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
```

Invoke the chaincode using these lines

```
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.cred.com --tls --cafile $ORDERER_CA -C $CHANNEL_NAME -n Credential-Verification --peerAddresses localhost:7051 --tlsRootCertFiles $UNIVERSITY_PEER_TLSROOTCERT --peerAddresses localhost:9051 --tlsRootCertFiles $STUDENT_PEER_TLSROOTCERT --peerAddresses localhost:11051 --tlsRootCertFiles $COMPANY_PEER_TLSROOTCERT -c '{"function":"CreateResult","Args":["RES1", "Stu1", "100", "90", "90%", "Pass"]}'
```

Add more Results by using the same above line by chaging the values

See the Results by using this command
```
peer chaincode query -C $CHANNEL_NAME -n Credential-Verification -c '{"function":"ReadResult","Args":["RES1", "Stu1"]}'
```


## Environment Variable Explanations

```bash
# Channel Name - Defines the blockchain channel for transactions
export CHANNEL_NAME=mychannel

# Fabric Configuration Path
export FABRIC_CFG_PATH=./peercfg

# Local MSP (Membership Service Provider) ID for the peer
export CORE_PEER_LOCALMSPID=CompanyMSP

# Enable TLS (Transport Layer Security)
export CORE_PEER_TLS_ENABLED=true
```

### Peer Configuration Breakdown
- `CHANNEL_NAME`: Identifies the specific channel where transactions occur
- `FABRIC_CFG_PATH`: Specifies the directory containing Fabric configuration files
- `CORE_PEER_LOCALMSPID`: Defines the organization's unique identifier
- `CORE_PEER_TLS_ENABLED`: Enables secure communication between network components

## Chaincode Query Commands

### 1. Get Result History
```bash
peer chaincode query -C $CHANNEL_NAME -n Credential-Verification -c '{"function":"GetResultHistory","Args":["RES1"]}'
```
**Derivation**:
- `-C`: Specifies the channel name
- `-n`: Names the chaincode
- `function`: Calls the specific chaincode function
- `Args`: Passes the result ID to retrieve its history
- Purpose: Retrieves the complete transaction history for a specific result

### 2. Get Results by Range
```bash
peer chaincode query -C $CHANNEL_NAME -n Credential-Verification -c '{"function":"GetResultsByRange","Args":["RES1","RES3"]}'
```
**Derivation**:
- Returns results with IDs between RES1 and RES3
- Useful for paginated or segmented result retrieval
- Demonstrates range-based querying in blockchain

### 3. Get Results with Pagination
```bash
peer chaincode query -C $CHANNEL_NAME -n Credential-Verification -c '{"function":"GetResultsWithPagination","Args":[3,"3"]}'
```
**Derivation**:
- First argument `3`: Page size (number of results per page)
- Second argument `3`: Page number
- Enables efficient data retrieval for large datasets
- Prevents overwhelming network with massive result sets

### 4. Get All Results
```bash
peer chaincode query -C $CHANNEL_NAME -n Credential-Verification -c '{"Args":["GetAllResults"]}'
```
**Derivation**:
- Retrieves all results stored in the chaincode
- No specific arguments required
- Useful for initial data loading or comprehensive overview

## Offer Creation Command

```bash
peer chaincode invoke -o localhost:7050 \
    --ordererTLSHostnameOverride orderer.cred.com \
    --tls --cafile $ORDERER_CA \
    -C $CHANNEL_NAME -n Credential-Verification \
    --peerAddresses localhost:7051 --tlsRootCertFiles $UNIVERSITY_PEER_TLSROOTCERT \
    --peerAddresses localhost:9051 --tlsRootCertFiles $STUDENT_PEER_TLSROOTCERT \
    --peerAddresses localhost:11051 --tlsRootCertFiles $COMPANY_PEER_TLSROOTCERT \
    -c '{"Args":["OfferContract:CreateOffer","Offer1"]}' \
    --transient "{\"ctc\":\"$CTC\",\"dateOfJoining\":\"$DATEOFJOINING\",\"dateOfRelease\":\"$DATEOFRELEASE\",\"companyName\":\"$COMPANYNAME\"}"
```

**Command Breakdown**:
- `-o localhost:7050`: Specifies the orderer node
- `--ordererTLSHostnameOverride`: Overrides TLS hostname for orderer
- `--tls`: Enables TLS communication
- `--cafile $ORDERER_CA`: Provides Certificate Authority file
- `-C $CHANNEL_NAME`: Specifies the channel
- `-n Credential-Verification`: Names the chaincode
- `--peerAddresses`: Specifies multiple peer nodes for endorsement
- `--tlsRootCertFiles`: Provides TLS certificates for each peer
- `-c`: JSON-formatted chaincode invocation
- `--transient`: Passes sensitive data not stored on the blockchain

### Transient Data Encoding
```bash
# Base64 encoding of sensitive information
export CTC=$(echo -n "9LPA" | base64 | tr -d \\n)
export DATEOFJOINING=$(echo -n "01/01/2025" | base64 | tr -d \\n)
export DATEOFRELEASE=$(echo -n "19/12/2025" | base64 | tr -d \\n)
export COMPANYNAME=$(echo -n "XXX" | base64 | tr -d \\n)
```
**Encoding Purpose**:
- Protects sensitive information
- Prevents direct storage of raw data on the blockchain
- Allows secure transmission of confidential details

## Offer Reading Command

```bash
peer chaincode query -C $CHANNEL_NAME -n Credential-Verification \
    --peerAddresses localhost:11051 --tlsRootCertFiles $COMPANY_PEER_TLSROOTCERT \
    -c '{"Args":["OfferContract:ReadOffer","Offer1"]}'
```

**Command Explanation**:
- Queries a specific offer by its ID
- Uses Company peer for verification
- Retrieves offer details without modifying blockchain state

## Key Considerations

1. Multi-endorsement ensures transaction integrity
2. TLS provides secure communication
3. Transient data protection for sensitive information
4. Flexible querying mechanisms

## Recommended Practices

- Always use TLS
- Implement multi-peer endorsement
- Use base64 encoding for sensitive data
- Implement proper access controls
```
