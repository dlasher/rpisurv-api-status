# rpisurv-api-status
API for rpisurv (https://github.com/SvenVD/rpisurv)

Simple API to simulate keyboard press for changing screen, with changes to show status and output as JSON

Request :
* pip3 install flask
* pip3 install keyboard

sudo nohup python3 rpisurv-api-status.py > log.txt 2>&1 &

http call example:

* pause rotation

curl http://api-ip:5000/camera/pause 

* resume

curl http://api-ip:5000/camera/resume

* camera 1

curl http://api-ip:5000/camera/1

* camera 2

curl http://api-ip:5000/camera/2

* camera 3

curl http://api-ip:5000/camera/3

* camera 4

curl http://api-ip:5000/camera/4

* status

curl http://api-ip:5000/camera/status


output is now JSON

#curl -q http://api-ip:5000/camera/status
{"camera":"1","status":"resume"}

#curl -q http://api-ip:5000/camera/status | jq
{
  "camera": "1",
  "status": "resume"
}


