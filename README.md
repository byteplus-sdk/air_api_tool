# air_api_tool
## Introduction
The API validation tool is used to reduce the communication costs between clients and API engineers by fixing feedback related to data issues during data preparation. The API validation tool simulates the complete process that the Byteplus team follows when importing customer data, including Data Format validation (Data Parse Check), Schema validation (Data Type Check), and Checker validation (Data Value Check). Customers can use the API validation tool to validate locally exported files and fix issues within the data. This process does not require the involvement of API engineers. This will be efficient for both customers and API engineers. Customers do not need to wait for API engineers to check data files. API engineers also do not need to wait for customers to fix data before rechecking. When customers provide data to the Byteplus team, the data will be of high quality.


## Document
[API Validate Tool Use Guide](https://bytedance.larkoffice.com/wiki/O15AwqMKIiZ4f6kCrp8cE4xPnRg)


## Install
You can install API Validate Tool in two ways.
### 1. Download compiled binary resources
Refer to [API Validate Tool Use Guide](https://bytedance.larkoffice.com/wiki/O15AwqMKIiZ4f6kCrp8cE4xPnRg), and download the binary file and use.


### 2. Download the source code and compile it yourself
First you need to install golang.
* clone the project.
* enter the project root directory.
* install dependencies
* build the binary executable file.

```shell
git clone https://github.com/byteplus-sdk/air_api_tool.git
cd air_api_tool
go mod tidy
./build_linux.sh
```