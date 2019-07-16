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
