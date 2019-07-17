https://hyperledger-fabric-ca.readthedocs.io/en/release-1.4

fabric ca 文档

TLS

当客户端连接到支持TLS协议的服务器要求创建安全连接并列出了受支持的密码组合（加密密码算法和加密哈希函数），握手开始。

服务器从该列表中决定加密和散列函数，并通知客户端。

服务器发回其数字证书，此证书通常包含服务器的名称、受信任的证书颁发机构（CA）和服务器的公钥。

客户端确认其颁发的证书的有效性。

为了生成会话密钥用于安全连接，客户端使用服务器的公钥加密随机生成的密钥，并将其发送到服务器，只有服务器才能使用自己的私钥解密。

利用随机数，双方生成用于加密和解密的对称密钥。这就是TLS协议的握手，握手完毕后的连接是安全的，直到连接（被）关闭。如果上述任何一个步骤失败，TLS握手过程就会失败，并且断开所有的连接。 [1] 

Certificate Signing Request (CSR) 

CSR 文件：
CSR是Certificate Signing Request的英文缩写，即证书请求文件，也就是证书申请者在申请数字证书时由CSP(加密服务提供者)在生成私钥的同时也生成证书请求文件，证书申请者只要把CSR文件提交给证书颁发机构后，证书颁发机构使用其根证书私钥签名就生成了证书公钥文件，也就是颁发给用户的证书。

用途：
通常CSR文件是在拿到参考码、授权码进行证书签发和下载时，通过网页提交给CA的（也可以由浏览器自动生成）。

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

consortium

instantiate

premise

entrant

strain

payroll

maturity 票据到期

premium 优质的，高昂的，附加费

creditworthiness

creditworthy

obligations

issuance

mandatory

sign off 

elaborate

allegiance

efficient and effective solutions

very descriptive 生动

EuroBond, DollarBond, YenBond
 
organizational units (OUs) 超级账本专用名词

https://hyperledger-fabric-ca.readthedocs.io/en/release-1.4/#

reading
https://hyperledger-fabric.readthedocs.io/en/latest/membership/membership.html

local MSP

channel MSP

vectors 向量

scalars 标量

generic container

intuitive 直观的：it’s a really good idea to use this kind of DNS name because well-chosen names will make your blockchain designs intuitive to other people. This idea applies equally well to smart contract names.

analogous

collisions

granularity

representation

key features

overall

remainder 

descriptive ： See how the list has a descriptive name: org.papernet.papers;

concatenation ： The key for a PaperNet commercial paper is formed by a concatenation of the Issuer and paper properties

This puts a short term strain on its finances – it will require an extra 5M USD each month to pay these new employees.

The exclusive relationship between an organization and its MSP makes it sensible to name the MSP after the organization, a convention you’ll find adopted in most policy configurations. 

In these cases it makes sense to have multiple MSPs and name them accordingly, e.g., ORG2-MSP-NATIONAL and ORG2-MSP-GOVERNMENT, reflecting the different membership roots of trust within ORG2 in the NATIONAL sales channel compared to the GOVERNMENT regulatory channel.

A solid understanding of the features will help you design and implement efficient and effective solutions.

MSP定义：

The power of an MSP goes beyond simply listing who is a network participant or member of a channel. An MSP can identify specific roles an actor might play either within the scope of the organization the MSP represents (e.g., admins, or as members of a sub-organization group), and sets the basis for defining access privileges in the context of a network and channel (e.g., channel admins, readers, writers).

MSPs appear in two places in a blockchain network: channel configuration (channel MSPs), and locally on an actor’s premise (local MSP). Local MSPs are defined for clients (users) and for nodes (peers and orderers). Node local MSPs define the permissions for that node (who the peer admins are, for example). The local MSPs of the users allow the user side to authenticate itself in its transactions as a member of a channel (e.g. in chaincode transactions), or as the owner of a specific role into the system (an org admin, for example, in configuration transactions).

Notice how each state is self-describing; each property has a name and a value. Although all our commercial papers currently have the same properties, this need not be the case for all time, as Hyperledger Fabric supports different states having different properties. This allows the same ledger world state to contain different forms of the same asset as well as different types of asset. It also makes it possible to update a state’s structure; imagine a new regulation that requires an additional data field. Flexible state properties support the fundamental requirement of data evolution over time.

To make these kinds of search tasks possible, it’s helpful to group all related papers together in a logical list（这里说的是逻辑列表）. The PaperNet design incorporates the idea of a commercial paper list – a logical container which is updated whenever commercial papers are issued or otherwise changed.

At the heart of a blockchain network is a smart contract. 

Usually there will only be one smart contract per file – contracts tend to have different lifecycles, which makes it sensible to separate them. However, in some cases, multiple smart contracts might provide syntactic help for applications, e.g. EuroBond, DollarBond, YenBond, but essentially provide the same function. In such cases, smart contracts and transactions can be disambiguated.

disambiguated

issue, buy and redeem 

commercial paper.

To solidify your understanding of the structure of a smart contract transaction, locate the buy and redeem transaction definitions, and see if you can see how they map to their corresponding commercial paper transactions.

An application can interact with a blockchain network by submitting transactions to a ledger or querying ledger content. 

You’re going to see how a typical application performs these six steps using the Fabric SDK. You’ll find the application code in the issue.js file. View it in your browser, or open it in your favourite editor if you’ve downloaded it. Spend a few moments looking at the overall structure of the application; even with comments and spacing, it’s only 100 lines of code!

Think of a wallet holding the digital equivalents of your government ID, driving license or ATM card. The X.509 digital certificates within it will associate the holder with a organization, thereby entitling them to rights in a network channel. 

For example, Isabella might be an administrator in MagnetoCorp, and this could give her more privileges than a different user – Balaji from DigiBank. Moreover, a smart contract can retrieve this identity during smart contract processing using the transaction context.

That’s it! In this topic you’ve understood how to call a smart contract from a sample application by examining how MagnetoCorp’s application issues a new commercial paper in PaperNet. Now examine the key ledger and smart contract data structures are designed by in the architecture topic behind them.

In the Developing Applications topic, we can see how the Fabric SDKs provide high level programming abstractions which help application and smart contract developers to focus on their business problem, rather than the low level details of how to interact with a Fabric network.

When a chaincode is installed and instantiated, all the smart contracts within it are made available to the corresponding channel.

All Hyperledger Fabric CA servers in a cluster share the same database for keeping track of identities and certificates. If LDAP is configured, the identity information is kept in LDAP rather than the database.

A server may contain multiple CAs. Each CA is either a root CA or an intermediate CA. Each intermediate CA has a parent CA which is either a root CA or another intermediate CA.


go install

https://segmentfault.com/q/1010000004044176

定义	
协议	年份
SSL 1.0	未知
SSL 2.0	1995
SSL 3.0	1996
TLS 1.0	1999
TLS 1.1	2006
TLS 1.2	2008
TLS 1.3	2018
