#!/bin/bash

# NEC FPD Microservice test script
# Replace these variables with your actual values
MICROSERVICE_URL="your-microservice-url"
DEVICE_FQDN="your-device-fqdn"

echo "Testing NEC FPD Microservice..."
echo "Microservice URL: $MICROSERVICE_URL"
echo "Device FQDN: $DEVICE_FQDN"
echo "----------------------------------------"

# SET Power
echo "Testing SET Power..."
curl -X PUT "http://${MICROSERVICE_URL}/${DEVICE_FQDN}/power" \
     -H "Content-Type: application/json" \
     -d '"on"'
sleep 10

# SET Videoroute
echo "Testing SET Videoroute..."
curl -X PUT "http://${MICROSERVICE_URL}/${DEVICE_FQDN}/videoroute/1" \
     -H "Content-Type: application/json" \
     -d '"1"'
sleep 1

# SET Volume
echo "Testing SET Volume..."
curl -X PUT "http://${MICROSERVICE_URL}/${DEVICE_FQDN}/volume/1" \
     -H "Content-Type: application/json" \
     -d '"30"'
sleep 1

# SET Audiomute
echo "Testing SET Audiomute..."
curl -X PUT "http://${MICROSERVICE_URL}/${DEVICE_FQDN}/audiomute/1" \
     -H "Content-Type: application/json" \
     -d '"false"'
sleep 1

# GET Power
echo "Testing GET Power..."
curl -X GET "http://${MICROSERVICE_URL}/${DEVICE_FQDN}/power"
sleep 1

# GET Videoroute
echo "Testing GET Videoroute..."
curl -X GET "http://${MICROSERVICE_URL}/${DEVICE_FQDN}/videoroute"
sleep 1

# GET Volume
echo "Testing GET Volume..."
curl -X GET "http://${MICROSERVICE_URL}/${DEVICE_FQDN}/volume"
sleep 1

# GET Audiomute
echo "Testing GET Audiomute..."
curl -X GET "http://${MICROSERVICE_URL}/${DEVICE_FQDN}/audiomute"
sleep 1

echo "----------------------------------------"
echo "All tests completed."
