# Singularity

Singularity is a Iot helper program that is designed to allow local controll and loging while also integrating to other platforms. It's currently is an very early state and right now only has limited feature and device support. 

The plan for this software is to act as a middle agent collecting events from simple IOT devices and then forward these events on to third party platform API's. This software will handle state logging and act as a local proxy to allow the Iot device to be as simple as possible. 

No idea if this is really going to work out at this point. It's also a programming job for me to build my skills with golang, api's and native k8s applications. 

---

I have decided to make this software listen on a MQTT message bus where the IOT devices are publishing metrics and listening for commands. 
