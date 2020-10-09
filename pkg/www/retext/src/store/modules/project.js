import axios from 'axios';

export const ProjectMutations = {
    ADD_PROJECT: "addProject",
    SELECT_PROJECT: "selectProject",
}

export const ProjectActions = {
    LOAD_PROJECTS: "loadProjects",
    ADD_PROJECT: "addProject",
    makeAddProjectPayload(name, description, year, month) {
        return {
            name, description, year, month
        }
    }
}

export const ProjectModule = {
    state: {
        projects: [],
        currentProject: null,
    },

    mutations: {

        [ProjectMutations.ADD_PROJECT](state, project) {
            for (const projectOld of state.projects) {
                if (project.id === projectOld.id) {
                    return;
                }
            }

            state.projects.push(project);
        },

        [ProjectMutations.SELECT_PROJECT](state, index) {
            state.currentProject = state.projects[index];
        },
    },

    actions: {
        async [ProjectActions.ADD_PROJECT]({commit}, payload) {
            const res = await axios.post("/project/create", payload);
            commit(ProjectMutations.ADD_PROJECT, res.data);
        },

        async [ProjectActions.LOAD_PROJECTS]({commit}) {
            const res = await axios.get("/project/list")

            res.data.forEach(project => commit(ProjectMutations.ADD_PROJECT, project));
        }
    },

    getters: {
        currentProject(state) {
            return state.currentProject;
        },

        projects(state) {
            return state.projects;
        }
    }
}