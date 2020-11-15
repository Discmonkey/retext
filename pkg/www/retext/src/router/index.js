import Vue from 'vue'
import VueRouter from 'vue-router'
import Code from '@/views/Code'
import Upload from "@/views/Upload";
import Demo from "@/views/Demo";
import Home from "@/Home";
import Project from "@/views/Project";
import ProjectHub from "@/views/ProjectHub";
import NotFound from "@/components/nav/NotFound";
import Buckets from "@/views/Buckets";
import Bucket from "@/views/Bucket";

Vue.use(VueRouter);

// Set the landing page by changing homePath below.
// homePath gets "injected" into `routes` below programmatically.
// const homePath = '/projects';
export const HubName = "ProjectHub";
export const ProjectName = "project";

const routes = [
    {
        path: '/projects',
        name: 'Projects',
        component: Project,
        alias: '/'
    },
    {
        path: '/project/:projectId',
        name: ProjectName,
        component: Home,
        children: [
            {
                path: 'code/:documentId',
                name: 'Code',
                component: Code
            },
            {
                path: 'demo/:documentId',
                name: 'Demo',
                component: Demo
            },

            {
                path: 'view/:codeId',
                component: Bucket,
            },
            {
                path: 'upload',
                name: 'Upload',
                component: Upload,
            },
            {
                path: '',
                name: HubName,
                component: ProjectHub
            },
            {
                path: 'buckets',
                name: 'Coding Buckets',
                component: Buckets
            }
        ],
    },
    {
        path: '/404',
        component: NotFound
    },
    {
        path: '*',
        component: NotFound
    }
];


const router = new VueRouter({
    mode: 'history',
    base: process.env.BASE_URL,
    routes
});

export default router
