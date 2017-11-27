#!/bin/bash

cd support/config-server/

keytool -genkeypair -alias goblogkey -keyalg RSA -dname "CN=Go Blog,OU=Unit,O=Organization,L=City,S=State,C=SE" -keypass changeme -keystore server.jks -storepass letmein -validity 730
