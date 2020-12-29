import Vue from "vue";

export const mutations = {
    add: "addNotification",
    ignore: "ignoreNotification",
    dismiss: "dismissNotification",
    resolve: "resolveNotification",
}

// export const mutations = {}
interface Notification {
    id: number,
    title?: string,
    text: string,
    dismissed?: boolean,
    okText?: string,
    cancelText?: string,
    // notification has been hidden but not yet ignored nor resolved
    dismissCallback?: (n: Notification) => void, // this may be called multiple times
    ignoreCallback?: (n: Notification) => void,
    resolveCallback?: (n: Notification) => void,
}

interface State {
    notifications: Array<Notification>,
    errNum: number,
}

function findNotificationIndex(nots: Array<{id: number}>, id: number) {
    return nots.findIndex(n => n.id === id)
}

export const Module = {
    state: {
        notifications: [],
        errNum: 0,
    } as State,
    mutations: {
        [mutations.add](state: State, notification: Notification) {
            state.errNum++;

            notification.id = state.errNum;
            Vue.set(notification, "dismissed", false);

            state.notifications.push(notification);
        },
        [mutations.dismiss](state: State, id: number) {
            const i = findNotificationIndex(state.notifications, id);

            if(i === -1) {
                return;
            }

            const n = state.notifications[i];

            n.dismissed = true;
            if(n.dismissCallback) {
                n.dismissCallback(n);
            }
        },
        [mutations.ignore](state: State, id: number) {
            const i = findNotificationIndex(state.notifications, id);

            if(i === -1) {
                return;
            }

            const n = state.notifications[i];

            if(n.ignoreCallback) {
                n.ignoreCallback(n);
            }

            state.notifications.splice(i, 1);
        },
        [mutations.resolve](state: State, id: number) {
            const i = findNotificationIndex(state.notifications, id);

            if(i === -1) {
                return;
            }

            const n = state.notifications[i];

            if(n.resolveCallback) {
                n.resolveCallback(n);
            }

            state.notifications.splice(i, 1);
        },
    },
    getters: {
        notifications(state: State): Array<Notification> {
            return state.notifications;
        },
    },
}
