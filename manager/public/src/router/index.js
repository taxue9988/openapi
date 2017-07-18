import Vue from 'vue'
import Router from 'vue-router'
import Home from '@/pages/Home'
import Dashboard from '@/pages/Dashboard'
import BasicTable from '@/pages/BasicTable'
import OpenAPI from '@/pages/openapi/OpenAPI'
import Widget from '@/pages/Widget'
import ImageList from '@/pages/ImageList'
import Charts from '@/pages/Charts'
import Login from '@/pages/Login'
import LockScreen from '@/pages/LockScreen'

import ApiStat from '@/pages/openapi/page/ApiStat'
import ApiList from '@/pages/openapi/page/ApiList'

Vue.use(Router)

export default new Router({
    routes: [{
        path: '/',
        component: Home,
        redirect: "/systems/openapi/api-list",
        children: [{
                path: '/root',
                name: 'dashboard',
                component: Dashboard
            },
            {
                path: '/systems/openapi',
                name: 'OpenAPI',
                redirect: '/systems/openapi/api-stat',
                component: OpenAPI,
                children: [{
                        path: 'api-stat',
                        name: 'ApiStat',
                        component: ApiStat
                    },
                    {
                        path: 'api-list',
                        name: 'ApiList',
                        component: ApiList
                    }
                ]
            }
        ]
    }]
})