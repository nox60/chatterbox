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
