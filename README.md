# watch
watch message  received from lowan server with gateway and application 

### [usage]
watch.exe --mqtt-server tcp://127.0.0.1:1883 --gw-topic  gateway/+/stats  --app-topic application/+/node/+/rx

### [gw-topic]
1. gateway/[MAC]/stats  --Topic for gateway statistics
2. gateway/[MAC]/rx     --Topic for received packets (from nodes)
3. gateway/[MAC]/tx     --Topic for publishing packets to be transmitted by the given gateway.

### [app-topic]
1. application/[AppEUI]/node/[DevEUI]/tx     --Applications are able to send data to the nodes
2. application/[AppEUI]/node/[DevEUI]/rx     --To receive data from your node
3. application/[AppEUI]/node/[DevEUI]/join   --Topic for join notifications
4. application/[AppEUI]/node/[DevEUI]/ack    --Topic for ACK notifications.
5. application/[AppEUI]/node/[DevEUI]/error  --Topic for error notifications. 
6. application/[AppEUI]/node/[DevEUI]/rxinfo --Topic for for rx information of received packets (e.g. frequency, bandwidth, ADR, ...)

### [Network-Control(TODO)]
1. application/[AppEUI]/node/[DevEUI]/rxinfo  --Topic on which RX related information is published for each received packet.
2. application/[AppEUI]/node/[DevEUI]/mac/rx  --Topic for received MAC commands (from the nodes).
3. application/[appEUI]/node/[DevEUI]/mac/error  --Topic for error notifications
4. application/[AppEUI]/node/[DevEUI]/mac/tx   --Topic for sending MAC commands to the node

### [Thanks]
1. [loraserver](https://github.com/brocaar/loraserver)
2. [lora-gateway-bridge](https://github.com/brocaar/lora-gateway-bridge)
3. [loraserverdoc](https://docs.loraserver.io/loraserver/sending-data/)