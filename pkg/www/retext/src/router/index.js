import Vue from 'vue'
import VueRouter from 'vue-router'
import Code from '@/views/Code'
import Upload from "@/views/Upload";
import Demo from "@/views/Demo";
import Home from "@/Home";
import Project from "@/views/Project";
import LoadProject from "@/views/LoadProject";

Vue.use(VueRouter);

// `meta: {name: ...}` used in Breadcrumbs.vue
const routes = [
    {
        path: '/project',
        name: 'Projects',
        component: Home,
        redirect: 'project',
        children: [
            {
                path: '',
                name: 'List',
                component: Project,
            },
            {
                path: ':projectId',
                meta: {name: 'projectId'},
                name: 'projectView',
                component: LoadProject,
                redirect: ':projectId',
                children: [
                    {
                        path: 'code',
                        name: 'Code',
                        redirect: 'upload',
                        component: Home,
                        children: [
                            {
                                path: ':documentId',
                                meta: {name: 'documentId'},
                                component: Code
                            },
                        ],
                    },
                    {
                        path: 'demo/:documentId',
                        name: 'Demo',
                        component: Demo
                    },
                    {
                        path: 'upload',
                        alias: '',
                        name: 'Upload',
                        component: Upload,
                    },
                ]
            }
        ],
    },
];

const router = new VueRouter({
    mode: 'history',
    base: process.env.BASE_URL,
    routes
});

export default router
