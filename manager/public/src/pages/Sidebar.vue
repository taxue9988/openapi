<template>
<div class="sidebar">
      <Menu ref="menu" theme="dark" width="100%" class="menu" :active-name="activeName" :open-names="openNames" @on-select="handleSelect">
          <!--<Menu-item name="/root">
                <i class="fa fa-dashboard"></i>
                数据大盘  
          </Menu-item>-->
  
          <Submenu name="systems">
              <template slot="title">
                  <i class="fa fa-database"></i>
                  系统列表
              </template>
              <Menu-group title="OpenAPI"  id="sidebar-menu-group">
                 <!--<Menu-item name="/systems/openapi/api-stat">
                     数据统计
                </Menu-item>-->
                <Menu-item name="/systems/openapi/api-list">
                     API管理
                </Menu-item>
             
              </Menu-group>
               
              </Submenu>
              
          </Submenu>
      </Menu>
    </div>
</template>
<script>
    export default {
    name: 'sidebar',
    data () {
      return {
        activeName: '',
        openNames: []
      }
    },
    created () {
      this.update()
    },
    methods: {
      handleSelect (name) {
             this.$router.push(name)
      },
      update (route) {
        const path = route ? route.path : this.$route.path
         const activeName = path
         const openName = path.split('/')[1]

            this.$set(this, 'activeName', activeName)
 
    
         this.$set(this, 'openNames', [openName])
       

        this.$nextTick(() => {
          this.$refs.menu.updateActiveName()
          this.$refs.menu.updateOpened()
        })
      }
    }
  }
</script>

<style rel="stylesheet/scss" lang="less" scoped>
#sidebar-menu-group {
    margin-left:1rem;
   
}


</style>