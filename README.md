[![Build Status](https://travis-ci.org/0xC0D3D00D/atossa.svg?branch=master)](https://travis-ci.org/0xC0D3D00D/atossa)
[![Coverage Status](https://coveralls.io/repos/github/0xC0D3D00D/atossa/badge.svg?branch=master)](https://coveralls.io/github/0xC0D3D00D/atossa?branch=master)
[![codebeat badge](https://codebeat.co/badges/92e8c83f-f3a6-4052-a33e-9184e93eda9f)](https://codebeat.co/projects/github-com-0xc0d3d00d-atossa-master)
[![Go Report Card](https://goreportcard.com/badge/github.com/0xc0d3d00d/atossa)](https://goreportcard.com/report/github.com/0xc0d3d00d/atossa)
[![Gitter](https://badges.gitter.im/atossadb/community.svg)](https://gitter.im/atossadb/community?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge)

# Atossa
Fast, polyglot database

# Supported redis commands
Artimis doesn't implement all redis commandset. Instead, it supports the core concepts that will help you to build any type of data models on top of it.

:white_check_mark: Not implemented yet  
:heavy_plus_sign: Partially implemented  
:heavy_check_mark: Implemented  

## Connection
:white_check_mark: `ECHO message`: Echo the given string  
:heavy_plus_sign: `PING [message]`: Ping the server  
:white_check_mark: `QUIT`: Close the connection  

## Administrative
:heavy_check_mark: `COMMAND`: Get array of supported commandset with details  
:white_check_mark: `COMMAND COUNT`: Get total number of supported commands  
:white_check_mark: `COMMAND INFO command-name [command-name ...]`: Get array of specific commands  
:white_check_mark: `CONFIG`: Returns current configuration of the server  
:white_check_mark: `DBSIZE`: Returns the number of keys in the selected database  
:heavy_plus_sign: `INFO`: Get information and statistics about the server  
:white_check_mark: `LOLWUT`: WUT?!  
:white_check_mark: `SHUTDOWN`: Shut down the server  
:white_check_mark: `TIME`: Returns the current server time  

## Keys
:white_check_mark: `DEL key [key ...]`: Delete a key  
:white_check_mark: `EXISTS key [key ...]`: Determine if a key exists  
:white_check_mark: `EXPIRE key seconds`: Set a key's TTL in seconds  
:white_check_mark: `EXPIREAT key timestamp`: Set the expiration for a key as a UNIX timestamp  
:heavy_plus_sign: `KEYS pattern`  
:white_check_mark: `PEXPIRE key milliseconds`: Set a key's TTL in milliseconds  
:white_check_mark: `PEXPIREAT key milliseconds-timestamp`: Set the expiration for a keys as a UNIX timestamp specified in milliseconds  
:white_check_mark: `PTTL key`: Get the TTL for a key in milliseconds  
:white_check_mark: `RANDOMKEY`: Return a random key from the keyspace  
:white_check_mark: `RENAME`: Rename a key  
:white_check_mark: `RENAMENX key newkey`: Rename a key, only if the new key does not exists  
:white_check_mark: `SORT key  [BY pattern] [LIMIT offset count] [GET pattern [GET pattern ...]] [ASC|DESC] [ALPHA] [STORE destination]`: Sort the elements in a list, set or sorted set  
:white_check_mark: `TTL key`: Get the time to live for a key  
:white_check_mark: `TYPE key`: Determine the type stored at the key  


## Strings
:white_check_mark: `APPEND key value`: Append a value to a key  
:white_check_mark: `DECR key`: Decrement the integer value of a key by one  
:white_check_mark: `DECRBY key decrement`: Decrement the integer value of a key by the given number  
:heavy_check_mark: `GET key`: Get the value of a key  
:white_check_mark: `GETRANGE key start end`: Get a substring of the string stored at a key  
:white_check_mark: `GETSET key value`: Set the string value of a key and return its old value  
:white_check_mark: `INCR key`: Increment the integer value of a key by one  
:white_check_mark: `INCRBY key increment`: Increment the integer value of key by the given amount  
:white_check_mark: `INCRBYFLOAT key increment`: Increment the float value of a key by the given amount  
:white_check_mark: `MGET key [key ...]`: Get the values of all the given keys  
:white_check_mark: `MSET key value [key value ...]`: Set multiple keys to multiple values  
:white_check_mark: `MSETNX key value [key value ...]`: Set multiple keys to multiple values, only if none of the keys exist  
:white_check_mark: `PSETEX key milliseconds value`: Set the value and expiration in milliseconds of a key  
:heavy_plus_sign: `SET key value [EX seconds|PX milliseconds] [NX|XX] [KEEPTTL]`: Set the string value of a key  
:white_check_mark: `SETEX key value`: Set the value and expiration of a key  
:white_check_mark: `SETNX key value`: Set the value of a key, only if a key does not exist  
:white_check_mark: `SETRANGE key offset value`: Overwrite part of a string at key starting at the specified offset  
:white_check_mark: `STRLEN key`: Get the length of the value stored in a key  

## Lists
:white_check_mark: `BLPOP key [key ...] timeout`: Remove and get the first element in a list, or block until one is available  
:white_check_mark: `BRPOP key [key ...] timeout`: Remove and get the last element in a list, or block until one is available  
:white_check_mark: `BRPOPLPUSH source destination timeout`: Pop an element from a list, push it to another list and return it; or block until one is available  
:heavy_check_mark: `LINDEX key index`: Get an element from a list by its index  
:white_check_mark: `LINSERT key BEFORE|AFTER pivot element`: Insert an element before or after another element in a list  
:heavy_check_mark: `LLEN key`: Get the length of a list  
:heavy_check_mark: `LPOP key`: Remove and get the first element in a list  
:heavy_check_mark: `LPUSH key element [element ...]`: Prepend one or multiple elements to a list  
:white_check_mark: `LPUSHX key element [element ...]`: Prepend an element to a list, only if the list exists  
:heavy_check_mark: `LRANGE key start stop`: Get a range of elements from a list  
:white_check_mark: `LREM key count element`: Remove elements from a list  
:heavy_check_mark: `LSET key index element`: Set the value of an element in a list by its index  
:white_check_mark: `LTRIM key start stop`: Trim a list to the specified range  
:heavy_check_mark: `RPOP key`: Remove and get the last element in a list  
:white_check_mark: `RPOPLPUSH source destination`: Pop an element from a list, push it to another list and return it  
:heavy_check_mark: `RPUSH key element [element ...]`: Append one or multiple elements to a list  
:white_check_mark: `RPUSHX key element [element ...]`: Append an element to a list, only if the list exists  

# Incompatibility Notes
There is cases that this server behaviour is not compatible with Redis. You can find them listed below:   

- Key must not be an empty string, if an empty key is provided server will return an `ERR NILKEY Key is nil`
