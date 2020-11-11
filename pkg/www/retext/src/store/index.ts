import Vue from "vue"
import Vuex from "vuex";
import {ProjectModule} from "@/store/modules/project";
import {getColor, invertColor} from "@/core/Colors";
import {Id} from "@/model/id";
import {CodeContainer} from "@/model/codeContainer";
import {Code} from "@/model/code";

Vue.use(Vuex)

export const getters = {
    CONTAINERS: "containers",
    ID_TO_CONTAINER: "idToContainer",
    ID_TO_CODE: "idToCode",
    GET_CODE: "getCode",
    GET_TEXTS_LENGTH: "getTextsLength",
}

export const mutations = {
    CLEAR_CONTAINERS: "clearContainers",
    SET_CONTAINERS: "setContainers",
    ADD_CONTAINER: "addContainer",
    ADD_CODE: "addCode",
    ADD_TEXT: "addText",
    SET_PROJECT: "setProject"
}
export const actions = {
    INIT_CONTAINERS: "initContainers",
    CREATE_CONTAINER: "createContainer",
    CREATE_CODE: "createCode",
    ASSOCIATE_TEXT: "associateText",
    DISASSOCIATE_TEXT: "disassociateText",
    SET_COLOR_ACTIVE: "toggleColorActive",
}
type Commit = (s: string, a: any) => void;

/*
 * Takes a container object returned from the backend and turns it into an object usable by the front-end
 */

export interface CodeContainerWithMain {
    containerId: Id;
    main: Code | null;
    subcodes: Array<Code>;
    colorInfo: {
        bg: string;
        fg: string;
        active: boolean;
    };
}

function prepareContainer(backendCodeContainer: CodeContainer): CodeContainerWithMain {
    let main = null;

    if (backendCodeContainer.codes.length > 0) {
        main = backendCodeContainer.codes.shift();
    }

    return {
        containerId: backendCodeContainer.id,
        main,
        subcodes: backendCodeContainer.codes,
    } as CodeContainerWithMain;
}
async function createCode(containerId: Id, name: string): Promise<Code> {
    return (await Vue.axios.post("/code/create", {
                code: name,
                containerId
            })).data
}

function makeId(context: any) {
    return `projectId=${context.state.ProjectModule.currentProject.id}`
}

interface StoreState {
    containers: Array<CodeContainerWithMain>;
    idToContainer: {[key: number]: CodeContainerWithMain};
    idToCode: {[key: number]: Code};
    project: number;
}

export const store = new Vuex.Store({
    state: {
        containers: [],
        idToContainer: {},
        idToCode: {},
        project: -1,
    } as StoreState,

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
        [getters.GET_TEXTS_LENGTH] (state: StoreState) {
            return (containerId: Id) => {
                if (!(containerId in state.idToContainer)) {
                    return;
                }

                const container = state.idToContainer[containerId];

                if (container.main === null) {
                    return 0;
                }

                let length = container.main.texts.length;
                if (container.subcodes.length > 0) {
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
        [mutations.ADD_CONTAINER](state: StoreState, container: CodeContainerWithMain) {
            if(!(container.containerId in state.idToContainer)) {
                container.colorInfo = {
                    active: false,
                    bg: getColor(state.containers.length),
                    fg: "",
                };
                container.colorInfo.fg = invertColor(container.colorInfo.bg, true);
                state.containers.push(container);
            } else {
                // if the same container gets added again, replace old with new
                //  note that the old colorInfo is preserved due to how getColor works
                for(const [i, c] of state.containers.entries()) {
                    if(container.containerId === c.containerId) {
                        container.colorInfo = c.colorInfo;
                        state.containers[i] = container;
                        break;
                    }
                }
            }

            // override the map values with the values from the new container
            Vue.set(state.idToContainer, container.containerId, container);

            if (container.main !== null) {
                Vue.set(state.idToCode, container.main.id as number, container.main)
            }

            for(const code of container.subcodes) {
                Vue.set(state.idToCode, code.id as number, code);
            }
        },

        [mutations.CLEAR_CONTAINERS](state) {
            while (state.containers.length) {
                state.containers.pop();
            }

            state.idToContainer = {};
            state.idToCode = {};
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
        }
    },
    actions: {
        async [actions.CREATE_CONTAINER](context, {name}) {
            const res = await Vue.axios.post(`/code/container/create?${makeId(context)}`);
            const containerId = res.data.ContainerId as Id;

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
            if(
                words.anchor.paragraph > words.last.paragraph ||
                words.anchor.paragraph === words.last.paragraph && words.anchor.sentence > words.last.sentence ||
                words.anchor.paragraph === words.last.paragraph && words.anchor.sentence === words.last.sentence && words.anchor.word > words.last.word
            ) {
                // swap the order to make sure the first is first
                let t = words.last;
                words.last = words.anchor;
                words.anchor = t;
            }

            return Vue.axios.post("/code/associate", {
                key: parseInt(words.documentId),
                codeId: codeId,
                text: words.text,
                anchor: words.anchor,
                last: words.last,
            }).then(() => {
                // TODO: need the textId here?
                context.commit(mutations.ADD_TEXT, {codeId, text: words})
            });
        },
        async [actions.DISASSOCIATE_TEXT](context, {codedTexts}) {
            let allTextIds = [] as Array<Id>;
            for(let ti of Object.values(codedTexts)) {
                for(let textId of ti as Array<Id>) {
                    allTextIds.push(textId);
                }
            }

            return Vue.axios.delete("/code/disassociate", {
                data: {textIds: allTextIds}
            }).then(() => {
                for (const [codeId, textIds] of Object.entries(codedTexts)) {
                    let code = context.getters[getters.ID_TO_CODE][codeId];

                    let newTexts = [];
                    for(let text of code.texts) {
                        if (!(textIds as Array<Id>).includes(text.id)) {
                            newTexts.push(text);
                        }
                    }
                    code.texts = newTexts;
                }
            });
        },
        /**
         * if `toggleTo` is not provided, whatever the current "status" is will be flipped
         *
         * only 1 container can be "color active" at a time
         */
        async [actions.SET_COLOR_ACTIVE] (context, {containerId, toggleTo}) {
            let c = context.getters[getters.ID_TO_CONTAINER][containerId];
            toggleTo = toggleTo === undefined ? !c.colorInfo.active : toggleTo;

            for(let c of context.getters[getters.CONTAINERS]) {
                c.colorInfo.active = false;
            }

            c.colorInfo.active = toggleTo;
        },
        async [actions.INIT_CONTAINERS](context) {
            Vue.axios.get(`/code/list?${makeId(context)}`).then((res) => {
                const containers = res.data;

                context.commit(mutations.CLEAR_CONTAINERS);

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