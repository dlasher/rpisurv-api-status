from flask import Flask, jsonify, request, abort
import keyboard
import ipaddress

app = Flask(__name__)

ALLOWED_IPS = [ipaddress.ip_network('10.4.0.0/16')]

@app.before_request
def limit_remote_addr():
    ip_addr = ipaddress.ip_address(request.remote_addr)
    if ip_addr.version == 6 and ip_addr.ipv4_mapped is not None:
        ip_addr = ip_addr.ipv4_mapped
    if not any(ip_addr in net for net in ALLOWED_IPS):
        abort(403)  # Forbidden

def read_file(path):
    try:
        with open(path, 'r') as file:
            return file.read().strip()
    except IOError:
        return None

def write_file(path, content):
    try:
        with open(path, 'w') as file:
            file.write(content)
        return True
    except IOError:
        return False

@app.route('/')
def index():
    return 'Server Works!'

@app.route('/camera/<action>')
def camera_action(action):
    if action not in ['pause', 'resume', 'status', '1', '2', '3', '4']:
        return jsonify(error='Invalid action'), 400

    if action in ['pause', 'resume']:
        if not write_file("/tmp/rpi.status", action):
            return jsonify(error='File write error'), 500
        key = 'p' if action == 'pause' else 'r'
    elif action in ['1', '2', '3', '4']:
        if not write_file("/tmp/rpi.camera", action):
            return jsonify(error='File write error'), 500
        key = 'F' + action
    else:
        key = None

    if key:
        keyboard.press_and_release(key)

    status = read_file("/tmp/rpi.status")
    camera = read_file("/tmp/rpi.camera")

    if status is None or camera is None:
        return jsonify(error='File read error'), 500

    return jsonify(status=status, camera=camera)

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=5000, debug=False)