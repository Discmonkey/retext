class Packet {
    constructor(data, callback) {
        this.data = data;
        this.callback = callback;
    }
}

export class Channel {
    constructor() {
        this._packet = false;
    }

    /**
     * Must be called receive or the data will be ignored/overwritten on the next call to send()
     *
     * @param data
     * @param sendCallback
     */
    send(data, sendCallback) {
        this._packet = new Packet(data, sendCallback);
    }

    /**
     * Must be called after send or it'll just return false. Clears out the package after the first use.
     *
     * @returns {boolean|object}
     */
    receive() {
        let pack = this._packet
        this._packet = false;
        return pack;
    }
}
