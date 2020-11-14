import axios from 'axios';
import {Project} from "@/model/project";
import {Id} from "@/model/id";

export const MutationNames = {
    ADD_PROJECT: "addProject",
    SELECT_PROJECT: "selectProject",
}

export const ActionNames = {
    LOAD_PROJECTS: "loadProjects",
    ADD_PROJECT: "addProject",
    SELECT_PROJECT: "projectId",
}

interface ProjectState {
    projects: Array<Project>;
    currentProject: Project;
}

type Commit = (s: string, a: any) => void;

export const ProjectModule = {
    state: {
        projects: [],
        currentProject: null,
    } as unknown as ProjectState,

    mutations: {

        [MutationNames.ADD_PROJECT](state: ProjectState, project: Project) {
            for (const projectOld of state.projects) {
                if (project.id === projectOld.id) {
                    return;
                }
            }

            state.projects.push(project);
        },

        [MutationNames.SELECT_PROJECT](state: ProjectState, projectId: Id) {
            const index = state.projects.findIndex(project => project.id === projectId);

            if (index >= 0) {
                state.currentProject = state.projects[index];
            }
        }

    },

    actions: {
        async [ActionNames.ADD_PROJECT]({commit}: {commit: Commit}, payload: Project) {
            const res = await axios.post("/project/create", payload);
            commit(MutationNames.ADD_PROJECT, res.data);
        },

        async [ActionNames.LOAD_PROJECTS]({commit}: {commit: Commit}) {
            const res = await axios.get("/project/list")

            res.data.forEach((project: Project) => commit(MutationNames.ADD_PROJECT, project));
        },

        async [ActionNames.SELECT_PROJECT]({commit, state}: {commit: Commit; state: ProjectState}, payload: Id) {
            if (state.currentProject != null && state.currentProject.id === payload) return;

            const index = state.projects.findIndex(project => project.id === payload);

            if (index === -1) {
                const res = await axios.get(`/project?projectId=${payload}`);
                commit(MutationNames.ADD_PROJECT, res.data);
            }

            commit(MutationNames.SELECT_PROJECT, payload);

        }
    },

    getters: {
        currentProject(state: ProjectState): Project {
            return state.currentProject;
        },

        projects(state: ProjectState): Array<Project> {
            return state.projects;
        }
    }
}