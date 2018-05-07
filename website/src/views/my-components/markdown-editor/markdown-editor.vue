<template>
    <div>
         <Table border stripe :columns="columns" :data="data"></Table>
    </div>
</template>

<script>
export default {
    name: 'showlog',
    data () {
        return {
            columns:
            [
                    {
                        title: '详情',
                        key: 'detail'
                    },
                    {
                        title: 'Uris',
                        key: 'uris'
                    },
                    {
                        title: 'Methods',
                        key: 'methods'
                    },
                    {
                        title: 'UpstreamUrl',
                        key: 'upstreamurl'
                    },
                    {
                        title: 'name',
                        key: 'name'
                    },
                    {
                        title: '消费者',
                        key: 'consumer'
                    }
            ]
        }
    },
    mounted () {
    },
    methods: {
        
    }
};
</script>
