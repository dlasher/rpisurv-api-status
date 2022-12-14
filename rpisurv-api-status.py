# Gruik coded by GuiguiAbloc
import flask
import time
import keyboard
app = flask.Flask(__name__)


@app.route('/')
def index():
  return 'Server Works!'

@app.route('/camera/pause')
def campause():
  try:
      text_file = open("/tmp/rpi.status", "w")
      text_file.write('pause')
      text_file.close()
      text_file = open("/tmp/rpi.camera", "r")
      rcamera = text_file.read()
      text_file.close()
      keyboard.press_and_release('p')
      return {
          "status": "pause",
          "camera": rcamera,
      }
  except:
      print("Error")
      return "Error"

@app.route('/camera/resume')
def camresume():
  try:
      text_file = open("/tmp/rpi.status", "w")
      text_file.write('resume')
      text_file.close()
      text_file = open("/tmp/rpi.camera", "r")
      rcamera = text_file.read()
      text_file.close()
      keyboard.press_and_release('r')
      return {
          "status": "resume",
          "camera": rcamera,
      }
  except:
      print("Error")
      return "Error"

@app.route('/camera/status')
def camstatus():
   try:
      text_file = open("/tmp/rpi.status", "r")
      rstatus = text_file.read()
      text_file.close()
      text_file = open("/tmp/rpi.camera", "r")
      rcamera = text_file.read()
      text_file.close()
      return {
          "status": rstatus,
          "camera": rcamera
      }
   except:
      print("Error")
      return "Error"

@app.route('/camera/1')
def cam1():
  try:
      text_file = open("/tmp/rpi.status", "r")
      rstatus = text_file.read()
      text_file.close()
      text_file = open("/tmp/rpi.camera", "w")
      text_file.write('1')
      text_file.close()
      keyboard.press_and_release('F1')
      return {
          "status": rstatus,
          "camera": "1",
      }
  except:
      print("Error")
      return "Error"

@app.route('/camera/2')
def cam2():
  try:
      text_file = open("/tmp/rpi.status", "r")
      rstatus = text_file.read()
      text_file.close()
      text_file = open("/tmp/rpi.camera", "w")
      text_file.write('2')
      text_file.close()
      keyboard.press_and_release('F2')
      return {
          "status": rstatus,
          "camera": "2",
      }
  except:
      print("Error")
      return "Error"

@app.route('/camera/3')
def cam3():
  try:
      text_file = open("/tmp/rpi.status", "r")
      rstatus = text_file.read()
      text_file.close()
      text_file = open("/tmp/rpi.camera", "w")
      text_file.write('3')
      text_file.close()
      keyboard.press_and_release('F3')
      return {
          "status": rstatus,
          "camera": "3",
      }
  except:
      print("Error")
      return "Error"

@app.route('/camera/4')
def cam4():
  try:
      text_file = open("/tmp/rpi.status", "r")
      rstatus = text_file.read()
      text_file.close()
      text_file = open("/tmp/rpi.camera", "w")
      text_file.write('4')
      text_file.close()
      keyboard.press_and_release('F4')
      return {
          "status": rstatus,
          "camera": "4",
      }
  except:
      print("Error")
      return "Error"

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=5000, debug=False)
