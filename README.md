# Intro

This utility reads DS18B20 one wire temprature sensors connected to a Rasberry Pi and sends them to Home Assistant via MQTT. 

Temprature sensors connected to the Raspberry Pi are automaticlly discoverred and added. The appropite MQTT messges are sent for each sensor so that Home Assistant MQTT discovery automaticlly adds each sensor.

# Usage

1. clone repository
2. build
3. copy and rename smaple.config.yaml to config.yaml
3. edit config.yaml
4. run



# TODO: 
- [] Expand configuration options
    - [] Gateway name
    - [] additonal options for Home Assistant MQTT auto discovery
    - [] MQTT topic prefix config
- [] implement running as a service
- [] improve documentation
- [] add sample .config file