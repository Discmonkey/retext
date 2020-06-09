export class Channel {
    constructor() {
        this._state = this._STATES.waiting;
        this._sendCallback = () => {};
        this._sendData = false;
    }

    send(data, sendCallback) {
        // send a "sending" event to trigger a reject() on any waiting "sending" event
        this._state = this._STATES.sending;
        this._sendData = data;
        this._sendCallback = sendCallback;
    }

    receive(receiveCallback) {
        if (this.isSending()) {
            this._state = this._STATES.waiting;
            receiveCallback(this._sendData, this._sendCallback);
        }
    }

    isSending() {
        return this._state === this._STATES.sending;
    }
}

Channel.prototype._STATES = {"waiting": "wai", "receiving": "rec", "sending": "sen"};
