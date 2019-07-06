### byfn.sh

该文件在first-network/fabric-samples目录下面



# 操作说明

byfn.sh <mode> [-c <channel name>] [-t <timeout>] [-d <delay>] [-f <docker-compose-file>] [-s <dbtype>] [-l <language>] [-o <consensus-type>] [-i <imagetag>] [-a] [-n] [-v] 

<mode>  例子 ./byfn.sh up 

- up - bring up the network with docker-compose up 通过docker-compose方式启动网络

- down - clear the network with docker-compose down   关掉并移除所有容器并且清理环境（此处需要分析脚本）

- restart - restart the network   重新启动网络

- generate - generate required certificates and genesis block 生成所有的证书和上帝区块

- upgrade  - upgrade the network from version 1.3.x to 1.4.0 将网络从1.3.x版本升级到1.4.0


  -c <channel name> - channel name to use (defaults to \"mychannel\")

  -t <timeout> - CLI timeout duration in seconds (defaults to 10)

  -d <delay> - delay duration in seconds (defaults to 3)

  -f <docker-compose-file> - specify which docker-compose file use (defaults to docker-compose-cli.yaml)

  -s <dbtype> - the database backend to use: goleveldb (default) or couchdb

  -l <language> - the chaincode language: golang (default) or node

  -o <consensus-type> - the consensus-type of the ordering service: solo (default), kafka, or etcdraft

  -i <imagetag> - the tag to be used to launch the network (defaults to \"latest\")
  
  -a - launch certificate authorities (no certificate authorities are launched by default)

  -n - do not deploy chaincode (abstore chaincode is deployed by default)

  -v - verbose mode

byfn.sh -h 打印以上信息


### 通过etcdraft的方式的orderer

```a
./byfn.sh up -o etcdraft
```

核心的compose文件是：docker-compose-etcdraft2.yaml

解读如下：

```file
# Copyright IBM Corp. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#

version: '2'

volumes:
  orderer2.example.com:
  orderer3.example.com:
  orderer4.example.com:
  orderer5.example.com:

networks:
  byfn:

services:


  # yaml 说明文档：https://www.jianshu.com/p/97222440cd08
  # - 代表数组

  orderer2.example.com:
    extends:
      file: base/peer-base.yaml
      service: orderer-base
    container_name: orderer2.example.com
    networks:
    - byfn
    volumes:
        - ./channel-artifacts/genesis.block:/var/hyperledger/orderer/orderer.genesis.block
        - ./crypto-config/ordererOrganizations/example.com/orderers/orderer2.example.com/msp:/var/hyperledger/orderer/msp
        - ./crypto-config/ordererOrganizations/example.com/orderers/orderer2.example.com/tls/:/var/hyperledger/orderer/tls
        - orderer2.example.com:/var/hyperledger/production/orderer
    ports:
    - 8050:7050

  orderer3.example.com:
    extends:
      file: base/peer-base.yaml
      service: orderer-base
    container_name: orderer3.example.com
    networks:
    - byfn
    volumes:
        - ./channel-artifacts/genesis.block:/var/hyperledger/orderer/orderer.genesis.block
        - ./crypto-config/ordererOrganizations/example.com/orderers/orderer3.example.com/msp:/var/hyperledger/orderer/msp
        - ./crypto-config/ordererOrganizations/example.com/orderers/orderer3.example.com/tls/:/var/hyperledger/orderer/tls
        - orderer3.example.com:/var/hyperledger/production/orderer
    ports:
    - 9050:7050

  orderer4.example.com:
    extends:
      file: base/peer-base.yaml
      service: orderer-base
    container_name: orderer4.example.com
    networks:
    - byfn
    volumes:
        - ./channel-artifacts/genesis.block:/var/hyperledger/orderer/orderer.genesis.block
        - ./crypto-config/ordererOrganizations/example.com/orderers/orderer4.example.com/msp:/var/hyperledger/orderer/msp
        - ./crypto-config/ordererOrganizations/example.com/orderers/orderer4.example.com/tls/:/var/hyperledger/orderer/tls
        - orderer4.example.com:/var/hyperledger/production/orderer
    ports:
    - 10050:7050

  orderer5.example.com:
    extends:
      file: base/peer-base.yaml
      service: orderer-base
    container_name: orderer5.example.com
    networks:
    - byfn
    volumes:
        - ./channel-artifacts/genesis.block:/var/hyperledger/orderer/orderer.genesis.block
        # 所挂载的上帝区块需要先生成
        - ./crypto-config/ordererOrganizations/example.com/orderers/orderer5.example.com/msp:/var/hyperledger/orderer/msp
        # msp信息也是其他地方生成的，有待考证
        - ./crypto-config/ordererOrganizations/example.com/orderers/orderer5.example.com/tls/:/var/hyperledger/orderer/tls
        # tls是加密信息？
        - orderer5.example.com:/var/hyperledger/production/orderer
        # orderer?
    ports:
    - 11050:7050

```


