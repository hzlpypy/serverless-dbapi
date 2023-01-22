## Architecture
![image](../resource/cluster.jpg)

## Introduction

|Component|Use|
|---|---|
|Gateway|Unified access portal of cluster|
|Auth|Safety and certification|
|Etcd Cluster|Data storage and service registry|
|Actuator Sharded Cluster|Sharded routing of requests and scheduling center of Actuators|
|Actuator Cluster|Parse http requests to fill in sql and execute sql|
|Manager Center Cluster|Store data source and api information|
