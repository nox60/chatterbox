pki

PKI是Public Key Infrastructure的首字母缩写，翻译过来就是公钥基础设施；PKI是一种遵循标准的利用公钥加密技术为电子商务的开展提供一套安全基础平台的技术和规范。

证明书的文件构造是一种叫做 X.509 的协议规定的。另一方面，认证机关也其实就是一个网络应用程序。

构成PKI的三大要素：

1. 证明书

2. 认证机关

3. 证书库


PKI核心流程： 

https://blog.csdn.net/liuhuiyi/article/details/7776825


a PKI provides a list of identities, 

and an MSP says which of these are members of a given organization that participates in the network.

核心单词：

tamper

integrity

impersonator

disseminated

Intermediate

exclusive

affiliation

organizational units (OUs) 超级账本专用名词

https://hyperledger-fabric-ca.readthedocs.io/en/release-1.4/#

reading
https://hyperledger-fabric.readthedocs.io/en/latest/membership/membership.html

local MSP

channel MSP

The exclusive relationship between an organization and its MSP makes it sensible to name the MSP after the organization, a convention you’ll find adopted in most policy configurations. 


In these cases it makes sense to have multiple MSPs and name them accordingly, e.g., ORG2-MSP-NATIONAL and ORG2-MSP-GOVERNMENT, reflecting the different membership roots of trust within ORG2 in the NATIONAL sales channel compared to the GOVERNMENT regulatory channel.
