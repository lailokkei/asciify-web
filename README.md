# Asciify Web
A realtime web GUI for [Asciify](https://github.com/toodemhard/asciify). 
<br><br>
The text generation is computed on the server because I want to reuse my go library and js is probably too slow. Image uploads and requests are sent using a websocket api in order to persist data and reduce network usage.
<br><br>
Deployed on [https://asciify.dev](https://asciify.dev). When it no longer works it's because my aws free trial expired :(

![image](https://github.com/toodemhard/asciify-web/assets/100080774/8954c5cf-6aac-430a-82c2-da695bc377bb)

## Install
```sh
git clone https://github.com/toodemhard/asciify-web
cd asciify-web
go run .
```
