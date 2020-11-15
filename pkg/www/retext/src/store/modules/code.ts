import {Id} from "@/model/id";
import {Code} from "@/model/code";
import {CodeContainer} from "@/model/codeContainer";
import Vue from "vue";
import {getColor, invertColor} from "@/core/Colors";
import {DocumentText} from "@/model/documentText";
import {API} from "@/core/API";
import {Commit} from "vuex";

export const mutations = {
    CLEAR_CONTAINERS: "clearContainers",
    SET_CONTAINERS: "setContainers",
    ADD_CONTAINER: "addContainer",
    ADD_CODE: "addCode",
    ADD_TEXT: "addText",
    DELETE_TEXT: "deleteText",
    SET_PROJECT: "setProject",
    TOGGLE_CLICK: "toggleClick",
}

export const actions = {
    INIT_CONTAINERS: "initContainers",
    CREATE_CONTAINER: "createContainer",
    CREATE_CODE: "createCode",
    ASSOCIATE_TEXT: "associateText",
    DISASSOCIATE_TEXT: "disassociateText",
    SET_COLOR_ACTIVE: "toggleColorActive",
}

type State =  {
    containers: Array<CodeContainerWithMain>;
    idToContainer: {[key: number]: CodeContainerWithMain};
    idToCode: {[key: number]: Code};
    project: number;
}

export const Module = {

    getters: {
        containers(state: State) {
            return state.containers
        },
        idToContainer(state: State) {
            return state.idToContainer;
        },
        idToCode(state: State) {
            return state.idToCode;
        },

        textLength(state: State): (containerId: number) => number {
            return (containerId: number) => {
                if (!(containerId in state.idToContainer)) {
                    return 0;
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
        [mutations.ADD_CONTAINER](state: State, container: CodeContainer) {

            const convertedContainer = prepareContainer(container);

            if(convertedContainer.containerId in state.idToContainer) {
                convertedContainer.colorInfo = {
                    activeClick: false,
                    activeHover: false,
                    bg: getColor(state.containers.length),
                    fg: "",
                };
                convertedContainer.colorInfo.fg = invertColor(convertedContainer.colorInfo.bg, true);
                state.containers.push(convertedContainer);
            } else {
                // if the same container gets added again, replace old with new
                //  note that the old colorInfo is preserved due to how getColor works
                for(const [i, c] of state.containers.entries()) {
                    if(convertedContainer.containerId === c.containerId) {
                        convertedContainer.colorInfo = c.colorInfo;

                        Vue.set(state.containers, i, convertedContainer);

                        break;
                    }
                }
            }

            // override the map values with the values from the new container
            Vue.set(state.idToContainer, convertedContainer.containerId, container);

            if (convertedContainer.main !== null) {
                Vue.set(state.idToCode, convertedContainer.main.id as number, convertedContainer.main)
            }

            for(const code of convertedContainer.subcodes) {
                Vue.set(state.idToCode, code.id as number, code);
            }
        },

        [mutations.CLEAR_CONTAINERS](state: State) {
            while (state.containers.length) {
                state.containers.pop();
            }

            state.idToContainer = {};
            state.idToCode = {};
        },

        [mutations.SET_CONTAINERS] (state: State, containers: Array<CodeContainerWithMain>) {
            Vue.set(state, "containers", containers);
        },

        [mutations.ADD_CODE](state: State, {containerId, code}: {containerId: Id, code: Code}) {
            if (!code.id) {
                console.error("received code without id");
                return;
            }

            state.idToContainer[containerId].subcodes.push(code);
            state.idToCode[code.id] = code;
        },

        [mutations.ADD_TEXT](state: State, {codeId, text}: {codeId: Id, text: DocumentText}) {
            const texts = state.idToCode[codeId].texts as Array<DocumentText>;

            const index = texts.findIndex(t => t.id === text.id);

            if (index < 0) {
                texts.push(text);
            }
        },

        [mutations.DELETE_TEXT](state: State, payload: {textId: Id, codeId: Id}) {
            const texts = state.idToCode[payload.codeId].texts as Array<DocumentText>;

            const index = texts.findIndex(t => t.id === payload.textId);

            if (index >= 0) {
                texts.splice(index, 1);
            }
        },

        [mutations.TOGGLE_CLICK](state: State, {container}: {container: Id}) {
            Vue.set(state.idToContainer[container].colorInfo, "activeClick",
                !state.idToContainer[container].colorInfo.activeClick);
        },

    },
    actions: {
        async [actions.CREATE_CONTAINER]({commit}: {commit: Commit}, payload: {projectId: Id; codeName: string;}) {

            const container = await API.code_container.post(payload.projectId);

            if (!container.id) {
                console.error("could not create container");
                return;
            }

            // TODO this should probably call the create code action below, just not sure how that affects container state, since the initial container would be empty
            const code = await API.code.post({
                name: payload.codeName, container: container.id, texts: []
            });

            container.codes.push(code);

            commit(mutations.ADD_CONTAINER, container);

            return code;
        },

        async [actions.CREATE_CODE]({commit}: {commit: Commit}, payload: {containerId: Id, name: string}) {

            const code = await API.code.post({
                name: payload.name, container: payload.containerId, texts: []
            });

            commit(mutations.ADD_CODE, {containerId: payload.containerId, code});

            return code;
        },

        async [actions.ASSOCIATE_TEXT]({commit}: {commit: Commit}, payload: {codeId: Id, text: DocumentText}) {
            const text = await API.document_text.post(payload.codeId, payload.text);

            commit(mutations.ADD_TEXT, {codeId: payload.codeId, text});
        },

        async [actions.DISASSOCIATE_TEXT]({commit}: {commit: Commit}, payload: {textId: Id, codeId: Id}) {

            await API.document_text.delete(payload.textId);

            commit(mutations.DELETE_TEXT, payload);
        },

        async [actions.INIT_CONTAINERS]({commit}: {commit: Commit}, payload: Id) {
            const containers = await API.code_containers.get(payload);

            // TODO this isn't right, we should just filter the containers in the components instead of clearing them every time
            commit(mutations.CLEAR_CONTAINERS);

            containers.forEach(container => commit(mutations.ADD_CONTAINER, container))
        },
    },
}

export interface CodeContainerWithMain {
    containerId: Id;
    main: Code | null;
    percentage: number;
    subcodes: Array<Code>;
    colorInfo: {
        bg: string;
        fg: string;
        activeHover: boolean;
        activeClick: boolean;
    };
}

function prepareContainer(backendCodeContainer: CodeContainer): CodeContainerWithMain {
    let main = null;

    if (backendCodeContainer.codes.length > 0) {
        main = backendCodeContainer.codes.shift();
    }

    return {
        containerId: backendCodeContainer.id,
        percentage: backendCodeContainer.percentage,
        main,
        subcodes: backendCodeContainer.codes,
    } as CodeContainerWithMain;
}

function makeId(context: any) {
    return `project_id=${context.state.ProjectModule.currentProject.id}`
}