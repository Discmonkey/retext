/**
 * A "channel" that only sends 1 thing at a time. send() and receive() return
 * the same Promise object. The Promise is resolved once both send() and receive
 * have been called (rejected if: see next paragraph). The variable passed to send() is
 * passed to the Promise's resolve().
 *
 * If send() is called twice before receive()ing, the Promise will have it's
 * reject() called and a new Promise will overwrite the old. Same goes for
 * trying to receive() twice before send()ing.
 *
 * todo: the following is valid (but not desired):
 *  1. mousedown and up (click) a category
 *  2. select some text
 *  3. mousedown the selected text
 *  4. as soon as text is mousedown, the category clicked in step 1 is used w/
 *   the selected text to make a /associate post request
 *  Because the receive is only called on mouseup, this can be avoided if the channel.send() is also called on mouseup of dragging text or by fixing this in some other way i haven't thought of cause it's 6am and i woke up at 10am yesterday.
 */
export class Channel {
    static _STATES = {"waiting": "wai", "receiving": "rec", "sending": "sen"};
    static _instanceCounter = 0;

    constructor() {
        this._promise = false;
        this._state = Channel._STATES.waiting;
        this._eventName = "cus-channel-" + ++Channel._instanceCounter;

        // todo: remove counter before final merge. it's a debug counter (used to count # of promises created
        this.counter = 0;

        // internal event-handling helpers
        // https://stackoverflow.com/questions/1530837
        let _eventTarget = document.createTextNode(null);
        this._addEventListener = _eventTarget.addEventListener.bind(_eventTarget);
        this._removeEventListener = _eventTarget.removeEventListener.bind(_eventTarget);
        this._dispatchEvent = _eventTarget.dispatchEvent.bind(_eventTarget);
    }

    send(data) {
        this._sendReceive(Channel._STATES.sending, data);

        return this._promise;
    }

    receive() {
        this._sendReceive(Channel._STATES.receiving);

        return this._promise;
    }

    _sendReceive(sr, data) {
        // only send will have data but, the definitions of send() and receive
        //  were so close i merged them into 1 method
        let rs = false;
        switch (sr) {
            case Channel._STATES.receiving:
                console.log('receiving')
                rs = Channel._STATES.sending;
                break;
            case Channel._STATES.sending:
                console.log('sending')
                rs = Channel._STATES.receiving;
                break;
            default:
                // you can't pass this method _states.waiting
                throw new Error("invalid value")
        }

        console.log(this._state)
        if (this._state === rs) {
            this._state = Channel._STATES.waiting
            this.dispatchChannelEvent(sr, data);
        } else {
            let promiseCount = ++this.counter;
            // send a "sending" event to trigger a reject() on any waiting "sending" event
            let channelE = this.dispatchChannelEvent(sr, data);
            this._state = sr;
            console.log(`created promise ${promiseCount}`)

            this._promise = new Promise((resolve, reject) => {
                let f = e => {
                    // don't want the event we just dispatched to trigger this promise
                    if (e === channelE) {
                        console.log('it happened')
                        return;
                    }

                    let newState = e.detail.state;
                    let passData = data || e.detail.data;
                    console.log(`${promiseCount}|${newState} === ${rs}`)
                    if (passData && newState === rs) {
                        resolve(passData);
                        console.log('resolved ' + promiseCount)
                    } else {
                        reject(passData);
                        console.log('rejected ' + promiseCount)
                    }
                    this.removeChannelListener(f);
                };
                this.addChannelListener(f);
            });
            this._promise.catch(() => {
                // reject() calls are irrelevant (until more than 1 thing in the
                //  channel is supported ... or something)
            });
        }
    }

    addChannelListener(func) {
        this._addEventListener(this._eventName, func);
    }

    removeChannelListener(func) {
        this._removeEventListener(this._eventName, func)
    }

    dispatchChannelEvent(state, data) {
        let e = new CustomEvent(this._eventName, {detail: {state: state, data: data}});
        this._dispatchEvent(e);
        return e;
    }

}

