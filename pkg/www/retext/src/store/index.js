import Vue from "vue"
import Vuex from "vuex";
import {ProjectModule} from "@/store/modules/project";

Vue.use(Vuex)

export const getters = {
    CONTAINERS: "containers",
    ID_TO_CONTAINER: "idToContainer",
    ID_TO_CODE: "idToCode",
    GET_CODE: "getCode",
    GET_TEXTS_LENGTH: "getTextsLength",
}

export const mutations = {
    SET_CONTAINERS: "setContainers",
    ADD_CONTAINER: "addContainer",
    ADD_CODE: "addCode",
    ADD_TEXT: "addText",
}
export const actions = {
    INIT_CONTAINERS: "initContainers",
    CREATE_CONTAINER: "createContainer",
    CREATE_CODE: "createCode",
    ASSOCIATE_TEXT: "associateText",
}

/*
 * Takes a container object returned from the backend and turns it into an object usable by the front-end
 */
function prepareContainer(backendCodeContainer) {
    const main = backendCodeContainer.subcodes.shift();
    return {
        containerId: backendCodeContainer.containerId,
        main,
        subcodes: backendCodeContainer.subcodes
    }
}
async function createCode(containerId, name) {
    return (await Vue.axios.post("/code/create", {
                code: name,
                containerId
            })).data
}

export const store = new Vuex.Store({
    state: {
        containers: [],
        idToContainer: {},
        idToCode: {},
    },
    getters: {
        [getters.CONTAINERS](state) {
            return state.containers
        },
        [getters.ID_TO_CONTAINER](state) {
            return state.idToContainer;
        },
        [getters.ID_TO_CODE](state) {
            return state.idToCode;
        },
        [getters.GET_TEXTS_LENGTH] (state) {
            return (containerId) => {
                if (!(containerId in state.idToContainer)) {
                    return;
                }
                let container = state.idToContainer[containerId];

                let length = container.main.texts == null ? 0 : container.main.texts.length;

                if (container.subcodes) {
                    for (let subCode of container.subcodes) {
                        if (subCode.texts != null) {
                            length += subCode.texts.length;
                        }
                    }
                }

                return length;
            }
        }
    },
    mutations: {
        [mutations.ADD_CONTAINER](state, container) {
            if(!(container.containerId in state.idToContainer)) {
                state.containers.push(container);
            } else {
                // if the same container gets added again, replace old with new
                for(const [i, container] of state.containers.entries()) {
                    if(container.containerId === container.containerId) {
                        state.containers[i] = container;
                        break;
                    }
                }
            }

            // override the map values with the values from the new container
            state.idToContainer[container.containerId] = container;
            state.idToCode[container.main.id] = container.main;
            for(const code of container.subcodes) {
                state.idToCode[code.id] = code;
            }
        },
        [mutations.SET_CONTAINERS] (state, containers) {
            Vue.set(state, "containers", containers);
        },
        [mutations.ADD_CODE](state, {containerId, code}) {
            state.idToContainer[containerId].subcodes.push(code);
            state.idToCode[code.id] = code;
        },
        [mutations.ADD_TEXT](state, {codeId, text}) {
            state.idToCode[codeId].texts.push(text);
        },
    },
    actions: {
        async [actions.CREATE_CONTAINER](context, {name}) {
            const res = await Vue.axios.post("/code/container/create");
            const containerId = res.data.ContainerId;

            const code = await createCode(containerId, name);

            const newContainer = {
                containerId: containerId,
                main: code,
                subcodes: [],
            }

            context.commit(mutations.ADD_CONTAINER, newContainer);

            return code;
        },
        async [actions.CREATE_CODE](context, {containerId, name}) {
            const code = await createCode(containerId, name);
            code.texts = [];

            context.commit(mutations.ADD_CODE, {containerId, code});

            return code;
        },
        async [actions.ASSOCIATE_TEXT](context, {codeId, words}) {
            return Vue.axios.post("/code/associate", {
                key: parseInt(words.documentId),
                codeId: codeId,
                text: words.text
            }).then(() => {
                context.commit(mutations.ADD_TEXT, {codeId, text: words.text})
            });
        },
        async [actions.INIT_CONTAINERS](context) {
            Vue.axios.get("/code/list").then((res) => {
                const containers = res.data;

                for(const c of containers) {
                    const container = prepareContainer(c)
                    context.commit(mutations.ADD_CONTAINER, container);
                }
            })
        },
    },

    modules: {
        ProjectModule
    }
})