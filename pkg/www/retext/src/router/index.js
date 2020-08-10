import Vue from 'vue'
import VueRouter from 'vue-router'
import Code from '@/views/Code'
import Upload from "@/views/Upload";
import Demo from "@/views/Demo";
import Home from "@/Home";
import Project from "@/views/Project";

Vue.use(VueRouter);

// Set the landing page by changing homePath below.
// homePath gets "injected" into `routes` below programmatically.
const homePath = '/project';

const routes = [
    {
        path: '/',
        name: 'Home',
        component: Home,
        redirect: homePath,
        children: [
            {
                path: '/code/:documentId',
                name: 'Code',
                component: Code
            },
            {
                path: '/demo/:documentId',
                name: 'Demo',
                component: Demo
            },
            {
                path: '/upload',
                name: 'Upload',
                component: Upload,
            },
            {
                path: '/about',
                name: 'About',
                // route level code-splitting
                // this generates a separate chunk (about.[hash].js) for this route
                // which is lazy-loaded when the route is visited.
                component: () => import(/* webpackChunkName: "about" */ '@/views/About.vue')
            },
            {
                path: '/project',
                name: 'Projects',
                component: Project
            }
        ],
    },
];
routes[0].children.forEach(child => {
    if(child.path === homePath) {
        child.alias = '/';
    }
});

const router = new VueRouter({
    mode: 'history',
    base: process.env.BASE_URL,
    routes
});

export default router
