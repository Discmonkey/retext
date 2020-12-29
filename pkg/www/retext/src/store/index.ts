import Vue from "vue"
import Vuex from "vuex";
import * as project from "@/store/modules/project";
import * as notification from "@/store/modules/notification";
import * as document from "@/store/modules/file";
import * as code from "@/store/modules/code";

Vue.use(Vuex)

export const mutations = {
    code: code.mutations,
    file: document.mutations,
    notification: notification.mutations,
    project: project.mutations,
}
export const actions = {
    code: code.actions,
    file: document.actions,
    project: project.actions,
}

export const store = new Vuex.Store({
    state: {},
    getters: {},
    mutations: {},
    actions: {},

    modules: {
        ProjectModule: project.Module,
        NotificationModule: notification.Module,
        DocumentModule: document.Module,
        CodeModule: code.Module,
    }
})