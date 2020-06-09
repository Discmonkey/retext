class Packet {
    constructor(data, callback) {
        this.data = data;
        this.callback = callback;
    }
}

export class Channel {
    /**
     * check if isSending before calling receive to be sure a packet was sent
     */
    constructor() {
        this._packet = null;
        this.isSending = false;
    }

    /**
     * Should be called before receive or the data will be ignored/overwritten on the next call to send()
     *
     * @param data
     * @param sendCallback
     */
    send(data, sendCallback) {
        this._packet = new Packet(data, sendCallback);
        this.isSending = true;
    }

    /**
     * Should be called after send or it'll just return false. Clears out the packet after the first use.
     *
     * Check isSending before calling this to be sure a packet was sent
     *
     * @returns {null|object}
     */
    receive() {
        let pack = this._packet

        this._packet = null;
        this.isSending = false;

        return pack;
    }
}
