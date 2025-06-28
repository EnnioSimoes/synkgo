# Synkgo

![Synkgo Logo](./synkgo-logo-2.png)

A command line tool for copy data from one database to other, is very usefull when we need copy data from production to a develop database.

## Steps

* Command to generate config file ✅
* Create docker-compose wich MySQL and PHPMyAdmin Tlest ✅
* Connect to database MySQL ✅
* Create tables em Fake data to test ✅
* Analize tables and count data

synkgo
    init // Start setup configuration to save in synkgo.json ✅

    config // Show config synkgo.json if exists ✅
    config -create // Generate blank template file for configuration ✅

    tables // show config tables
    tables -config // config tables to save in synkgo.json (future) ✅

    start // Start Sync!!