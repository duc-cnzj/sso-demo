<template>
  <div class="app-container">
    <div class="filter-container">
      <el-input v-model="listQuery.name" placeholder="name" style="width: 200px;margin-right: 10px;" class="filter-item" @keyup.enter.native="handleFilter" />
      <el-button v-waves class="filter-item" type="primary" icon="el-icon-search" @click="handleFilter">
        搜索
      </el-button>
      <el-button type="info" style="margin-left: 10px;" @click="handleCreate">
        新增
      </el-button>

    </div>

    <el-table
      key="0"
      v-loading="listLoading"
      :data="list"
      border
      fit
      highlight-current-row
      style="width: 100%;margin-top: 10px;"
      @sort-change="sortChange"
    >
      <el-table-column label="ID" prop="id" sortable="custom" align="center" width="80" :class-name="getSortClass('id')">
        <template slot-scope="{row}">
          <span>{{ row.id }}</span>
        </template>
      </el-table-column>
      <el-table-column label="用户名" min-width="120px">
        <template slot-scope="{row}">
          <span>{{ row.user_name }}</span>
        </template>
      </el-table-column>
      <el-table-column label="邮箱" min-width="150px">
        <template slot-scope="{row}">
          <span>{{ row.email }}</span>
        </template>
      </el-table-column>
      <el-table-column label="角色" min-width="200px" class="roleClass">
        <template slot-scope="{row}">
          <span v-for="r in row.roles" :key="r.id" style="margin-right: 5px;">
            <el-popover
              placement="top-start"
              width="100%"
              trigger="hover"
            >
              <el-tag slot="reference" type="success" v-text="r.name" />
              <template v-if="getPermissionName(r.permissions)!==null">
                <el-tag v-for="name in getPermissionName(r.permissions)" :key="name" style="margin-right: 5px;">
                  {{ name }}
                </el-tag>
              </template>
              <span v-else>该角色下没有权限</span>

            </el-popover>
          </span>
        </template>
      </el-table-column>
      <el-table-column label="最后登录时间" min-width="160px">
        <template slot-scope="{row}">
          <span>{{ row.last_login_at | formatDate }}</span>
        </template>
      </el-table-column>
      <el-table-column label="创建时间" min-width="160px">
        <template slot-scope="{row}">
          <span>{{ row.created_at | formatDate }}</span>
        </template>
      </el-table-column>
      <el-table-column label="操作" align="center" width="380" class-name="small-padding fixed-width">
        <template slot-scope="{row,$index}">
          <el-button size="mini" type="primary" @click="handleUpdate(row,$index)">
            更新
          </el-button>
          <el-button size="mini" type="info" @click="handSyncRoles(row,$index)">
            修改角色
          </el-button>
          <el-button size="mini" type="danger" @click="handleDelete(row,$index)">
            删除
          </el-button>
          <el-button size="mini" type="danger" @click="handleLogout(row,$index)">
            强制登出
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-pagination
      v-show="total > 0"
      background
      layout="prev, pager, next"
      :page-size="listQuery.pageSize"
      :total="total"
      style="margin-top: 10px"
      @prev-click="prevPage"
      @current-change="currentChange"
      @next-click="nextPage"
    />
    <el-dialog :title="textMap[dialogStatus]" :visible.sync="dialogFormVisible">
      <el-form ref="dataForm" :rules="rules" :model="temp" label-position="left" label-width="90px" style="width: 80%; margin-left:50px;">
        <el-form-item label="用户名" prop="user_name">
          <el-input v-model="temp.user_name" />
        </el-form-item>
        <el-form-item label="邮箱" prop="email">
          <el-input v-model="temp.email" />
        </el-form-item>
        <el-form-item v-if="dialogStatus==='create'" label="密码" prop="password">
          <el-input v-model="temp.password" type="password" />
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="dialogFormVisible = false">
          取消
        </el-button>
        <el-button type="primary" @click="dialogStatus==='create'?createData():updateData()">
          确定
        </el-button>
      </div>
    </el-dialog>
    <el-dialog title="修改用户角色" :visible.sync="dialogSyncFormVisible">
      <el-form ref="dataSyncForm" :rules="rules" :model="temp" label-position="left" label-width="90px" style="width: 80%; margin-left:50px;">
        <el-form-item label="角色名称" prop="roleName">
          <el-select
            v-model="temp.role"
            style="width: 80%;"
            filterable
            default-first-option
            multiple
            placeholder="请选择角色"
          >
            <el-option
              v-for="item in temp.roles"
              :key="item.id"
              :label="item.name"
              :value="item.id"
            />
          </el-select>
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="dialogSyncFormVisible = false">
          取消
        </el-button>
        <el-button type="primary" @click="sync()">
          确定
        </el-button>
      </div>
    </el-dialog>

    <el-dialog
      title="确定要删除吗"
      :visible.sync="deleteDialogVisible"
      width="30%"
    >
      <span>删除该数据？</span>
      <span slot="footer" class="dialog-footer">
        <el-button @click="deleteDialogVisible = false">取 消</el-button>
        <el-button type="primary" @click="doDelete">确 定</el-button>
      </span>
    </el-dialog>
    <el-dialog
      title="确定要登出用户吗"
      :visible.sync="logoutDialogVisible"
      width="30%"
    >
      <span>确定要登出用户吗？</span>
      <span slot="footer" class="dialog-footer">
        <el-button @click="logoutDialogVisible = false">取 消</el-button>
        <el-button type="primary" @click="forceLogout">确 定</el-button>
      </span>
    </el-dialog>
  </div>
</template>

<script>
import { index, update, store, destroy, syncRoles, forceLogout } from '@/api/user'
import { allRoles } from '@/api/role'
import waves from '@/directive/waves'

export default {
  name: 'Role',
  directives: { waves },
  data() {
    return {
      tooltipId: null,
      list: null,
      deleteDialogVisible: false,
      logoutDialogVisible: false,
      dialogSyncFormVisible: false,
      total: 0,
      current: null,
      listLoading: true,
      listQuery: {
        page: 1,
        pageSize: 15,
        name: '',
        sort: 'desc'
      },
      temp: {
        id: null,
        user_name: null,
        email: null,
        role: [],
        roles: [],
        password: null
      },
      dialogFormVisible: false,
      dialogStatus: '',
      textMap: {
        update: '编辑',
        create: '创建'
      },
      dialogPvVisible: false,
      pvData: [],
      rules: {
        user_name: [{ required: true, message: '用户名必填', trigger: 'change' }],
        email: [{ required: true, message: '邮箱必填', trigger: 'change' }],
        password: [{ required: true, message: '密码必填', trigger: 'change' }, { min: 5, message: '最少5位', trigger: 'blur' }]
      }
    }
  },
  created() {
    this.getList()
  },
  methods: {
    forceLogout() {
      forceLogout(this.current.id).then(res => {
        this.$notify({
          title: '登出成功',
          type: 'success',
          duration: 2000
        })
        this.logoutDialogVisible = false
      })
    },
    getPermissionName(data) {
      if (data && data.length > 0) {
        const res = data.map(item => `${item.project}.${item.name}`)
        return res || null
      }

      return null
    },
    doDelete() {
      destroy(this.current.id).then(res => {
        this.$notify({
          title: '删除成功',
          type: 'success',
          duration: 2000
        })
        this.list.splice(index, 1)
      })
      this.deleteDialogVisible = false
    },
    currentChange(p) {
      this.listQuery.page = p
      this.getList()
    },
    prevPage() {
      this.listQuery.page--
      this.getList()
    },
    nextPage() {
      this.listQuery.page++
      this.getList()
    },
    getList() {
      this.listLoading = true
      index({
        page_size: this.listQuery.pageSize,
        page: this.listQuery.page,
        user_name: this.listQuery.name,
        sort: this.listQuery.sort
      }).then(response => {
        const { data } = response
        this.total = response.total
        this.list = data
        this.listQuery.page = response.page
        this.listQuery.pageSize = response.page_size

        setTimeout(() => {
          this.listLoading = false
        }, 400)
      })
    },
    handleFilter() {
      this.listQuery.page = 1
      this.getList()
    },

    sortChange(data) {
      const { prop, order } = data
      if (prop === 'id') {
        this.sortByID(order)
      }
    },

    sortByID(order) {
      if (this.listQuery.sort === 'desc') {
        this.listQuery.sort = 'asc'
      } else {
        this.listQuery.sort = 'desc'
      }
      this.handleFilter()
    },
    resetTemp() {
      this.temp = {
        id: null,
        name: null,
        permissions: []
      }
    },
    handleCreate() {
      this.resetTemp()
      this.dialogStatus = 'create'
      this.dialogFormVisible = true
      this.$nextTick(() => {
        this.$refs['dataForm'].clearValidate()
      })
    },
    createData() {
      this.$refs['dataForm'].validate((valid) => {
        if (valid) {
          store({
            user_name: this.temp.user_name,
            email: this.temp.email,
            password: this.temp.password
          }).then(res => {
            this.list.unshift(res.data)
            this.dialogFormVisible = false
            this.$notify({
              title: '创建成功',
              type: 'success',
              duration: 2000
            })
          })
        }
      })
    },
    handleUpdate(row) {
      this.temp.id = row.id
      this.temp.user_name = row.user_name
      this.temp.email = row.email
      this.dialogStatus = 'update'
      this.dialogFormVisible = true
      this.$nextTick(() => {
        this.$refs['dataForm'].clearValidate()
      })
    },
    getAllRoles() {
      return allRoles().then(({ data }) => { this.temp.roles = data })
    },
    async handSyncRoles(row) {
      await this.getAllRoles()
      this.temp.id = row.id
      this.temp.role = row.roles ? row.roles.map(item => item.id) : []
      this.dialogSyncFormVisible = true
      this.$nextTick(() => {
        this.$refs['dataSyncForm'].clearValidate()
      })
    },
    updateData() {
      this.$refs['dataForm'].validate((valid) => {
        if (valid) {
          update(this.temp.id, {
            user_name: this.temp.user_name,
            email: this.temp.email
          }).then(({ data }) => {
            this.getList()
            this.dialogFormVisible = false
            this.$notify({
              title: '更新成功',
              message: '用户信息更新成功',
              type: 'success',
              duration: 2000
            })
          })
        }
      })
    },
    sync() {
      this.$refs['dataSyncForm'].validate((valid) => {
        if (valid) {
          syncRoles(this.temp.id, { role_ids: this.temp.role }).then(res => {
            this.getList()
            this.dialogSyncFormVisible = false
            this.$notify({
              title: '同步成功',
              type: 'success',
              duration: 2000
            })
          })
        }
      })
    },
    handleDelete(row, index) {
      this.deleteDialogVisible = true
      this.current = row
    },
    handleLogout(row, index) {
      this.logoutDialogVisible = true
      this.current = row
    },
    getSortClass: function(key) {
      return this.listQuery.sort === 'asc' ? 'ascending' : 'descending'
    }
  }
}
</script>

<style lang="scss" scoped>
  .el-tag:hover {
    cursor: pointer;
  }
</style>
