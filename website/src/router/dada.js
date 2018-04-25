export const test = [
    {
        path: '/component',
        icon: 'social-buffer',
        name: 'component',
        title: '机构管理',
        component: Main,
        children: [
            {
                path: 'text-editor',
                icon: 'compose',
                name: 'text-editor',
                title: '机构管理',
                component: () => import('@/views/tables/editable-table.vue')
            }
            ,
            {
                path: 'md-editor',
                icon: 'pound',
                name: 'md-editor',
                title: 'Markdown编辑器',
                component: () => import('@/views/my-components/markdown-editor/markdown-editor.vue')
            }
        ]
    }]