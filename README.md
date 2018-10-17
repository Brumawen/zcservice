# zcservice
Zeroconf registration and client service written in go.  This service provides access to zeroconf services to third party microservices.


## Installation

Zcservice comes as a single executable file.

Download the latest release file for your operating system from
https://github.com/Brumawen/zcservice/releases

Extract the executable file to a folder and run the following from the command line:

        zcservice -service install
        zcservice -service run

This will install and run the zcservice as a background service on the machine. 


## Configuration

There are a few configuration options available.  To set these, edit the config.json file that zcservice creates when it first runs.  This file contains json text with the following properties
* id: This is a globally unique identifier generated for this installation.
* defaultServiceType: This is the service type that the zcservice registers itself as.  It is also the service type that is used if a web request does not specify a service type.


## API Methods

### Register a service

To register a service with zeroconf, send a POST request to:

        http://127.0.0.1:20404/service/add

with a json document in the request body containing the following properties:

* <b>id</b> : (<i>string</i>) The unique identifier for the service instance.  This is usually a GUID, but can be any value.  If left blank, the zcservice will generate a GUID and return it with the response.
* <b>name</b> : (<i>string</i>) The name of the service.
* <b>serviceType</b> : (<i>string</i>) The service type (e.g. "_microservice._tcp").  If this is left blank then it uses the configured Default Service Type.
* <b>domain</b> : (<i>string</i>) The name of the domain.  Leave this blank for "local."
* <b>portNo</b> : (<i>int</i>) The port number you service is listening on for requests.
* <b>text</b> : (<i>string array</i>) An array of Key=Value text strings that provide additional information about the microservice.

The response will contain a json document with the following properties:

* <b>id</b> : (<i>string</i>) The unique identifier of the registered service.

### Deregister a service

To deregister a service, send a DELETE request to:

        http://127.0.0.1:20404/service/remove/{id}

where {id} is the unique identifier of the service.


### Get a list of services

To get a list of registered services, send a GET or a POST request to:

        http://127.0.0.1:20404/service/get

If a GET request is sent, the result will contain registered services for the configured default service type.  If a POST request is sent, then the request body must contain a json document with the following properties:

* <b>serviceType</b> : (<i>string</i>) The service type to search for.  Leave this blank to use the configured default service type.
* <b>domain</b> : (<i>string</i>) The domain name.  Leave this blank for "local."
* <b>waitTime</b> : (<i>int</i>) The maximum amount of time (in seconds) to wait for a response.  The default is 3 seconds.

The response will contain a json document with the following properties:

* <b>serviceType</b> : (<i>string</i>) The service type.
* <b>domain</b> : (<i>string</i>) The domain name.
* <b>services</b> : (<i>Array</i>) An array containing the details about the registered services found.  Each service will contain the following properties:
    * <b>name</b> : (<i>string</i>) The name of the service instance.
    * <b>port</b> : (<i>int</i>) The port number used by the service.
    * <b>hostname</b> : (<i>string</i>) The hostname of the computer the service is running on.
    * <b>type</b> : (<i>string</i>) The service type.
    * <b>domain</b> : (<i>string</i>) The domain name.
    * <b>text</b> : (<i>string array</i>) An array of Key=Value text strings that provide additional information about the service.
    * <b>ipv4</b> : (<i>string array</i>) An array containing the IPv4 IP address(es) of the service host.
    * <b>ipv6</b> : (<i>string array</i>) An array containing the IPv6 IP address(es) of the service host.


### Check if the service is online

To check if the service is running, send a GET request to:

        http:127.0.0.1:20404/online

This will return the text "true" if the service is online.