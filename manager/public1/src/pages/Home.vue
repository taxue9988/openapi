<template>
  <div class="home">
    <Row class="header" type="flex" align="middle">
      <div class="logo">
        <img src="../assets/img/logo.png" height="30" alt="">
        <span>Open</span>Api <Tag>管理后台</Tag>
      </div>
      <VmMsgPush  :data="msgPushData"></VmMsgPush>
      <Dropdown class="login-info" placement="bottom-end">
        <Button type="ghost">
            <img src="../assets/img/photo.jpg" height="30" alt="">Admin
            <Icon type="arrow-down-b"></Icon>
        </Button> 
        <Dropdown-menu slot="list">
            <Dropdown-item divided><i class="fa fa-key"></i>Log out</Dropdown-item>
            <Dropdown-item divided><i class="fa fa-cog"></i>Setting</Dropdown-item>
        </Dropdown-menu>
      </Dropdown>
    </Row>   
    <Sidebar ref="sidebar" />

    
    <div class="main-content">
      <router-view></router-view> 
      <!--<p class="vm-author" style="bottom:0;">vue-manager designed by <a href="https://github.com/luosijie" target="_blank">Jesse Luo</a></p> -->
    </div>
  
  </div>
</template>
<script>
import VmMsgPush from '@/components/vm-msg-push.vue'
import Sidebar from '@/pages/Sidebar'

export default {
  name: 'home',
  components: {
    VmMsgPush,
    Sidebar
  },
  beforeRouteUpdate (route, from, next) {
     this.$nextTick(() => {
        this.$refs.sidebar.update(route)
      })
      next()
    },
  data () {
    return {
      msgPushData: [
        {
          image: require('@/assets/img/photo.jpg'),
          from: 'Sunface',
          time: '2017-7-15',
          message: 'OpenApi V0.3.0发布'
        }
      ]
    }
  }
}

</script>
