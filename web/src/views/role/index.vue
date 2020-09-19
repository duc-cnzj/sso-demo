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
      <el-table-column label="角色名" min-width="120px">
        <template slot-scope="{row}">
          <span>{{ row.text }}</span>
        </template>
      </el-table-column>
      <el-table-column label="角色" min-width="120px">
        <template slot-scope="{row}">
          <span>{{ row.name }}</span>
        </template>
      </el-table-column>
      <el-table-column label="权限" min-width="350px">
        <template slot-scope="{row}">
          <span v-for="premission in row.permissions" :key="premission.id" style="margin-right: 5px">
            <el-tag style="margin-bottom: 5px" v-text="`${premission.project}.${premission.text}`" />
          </span>
        </template>
      </el-table-column>
      <el-table-column label="创建时间" min-width="150px">
        <template slot-scope="{row}">
          <span>{{ row.created_at | formatDate }}</span>
        </template>
      </el-table-column>
      <el-table-column label="操作" align="center" width="230" class-name="small-padding fixed-width">
        <template slot-scope="{row,$index}">
          <el-button v-if="row.status!='deleted'" size="mini" type="danger" @click="handleDelete(row,$index)">
            删除
          </el-button>
          <el-button size="mini" type="primary" @click="handleUpdate(row,$index)">
            更新
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-pagination
      v-show="total > 0"
      style="margin-top: 10px"
      background
      layout="prev, pager, next"
      :page-size="listQuery.pageSize"
      :total="total"
      @prev-click="prevPage"
      @current-change="currentChange"
      @next-click="nextPage"
    />
    <el-dialog :title="textMap[dialogStatus]" :visible.sync="dialogFormVisible">
      <el-form ref="dataForm" :rules="rules" :model="temp" label-position="left" label-width="90px" style="width: 80%; margin-left:50px;">
        <el-form-item label="角色名" prop="text">
          <el-input v-model="temp.text" />
        </el-form-item>
        <el-form-item label="角色" prop="name">
          <el-input v-model="temp.name" />
        </el-form-item>
        <el-form-item label="角色权限" prop="permissions">
          <el-select
            v-model="temp.permissions"
            multiple
            filterable
            placeholder="请输入关键词"
            style="width: 100%"
          >
            <el-option-group
              v-for="(group,name) in permissions"
              :key="name"
              :label="name"
            >
              <el-option
                v-for="item in group"
                :key="item.id"
                :label="`${name}.${item.text}`"
                :value="item.id"
              />
            </el-option-group>
          </el-select>
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
  </div>
</template>

<script>
import { index, update, store, destroy } from '@/api/role'
import { getByGroups } from '@/api/permission'
import waves from '@/directive/waves'

export default {
  name: 'Role',
  directives: { waves },
  data() {
    return {
      permissions: [],
      list: null,
      deleteDialogVisible: false,
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
        name: null,
        text: null,
        permissions: []
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
        text: [{ required: true, message: '角色名称必填', trigger: 'change' }],
        name: [{ required: true, message: '角色格式必须为 a-z', trigger: 'change', pattern: /^[a-zA-Z_-]+$/ }]
      }
    }
  },
  created() {
    this.getList()
  },
  methods: {
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

    getPermissions() {
      return getByGroups({}).then(response => {
        const { data } = response
        this.permissions = data
      })
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
        name: this.listQuery.name,
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
      this.getPermissions()

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
            name: this.temp.name,
            text: this.temp.text,
            permission_ids: this.temp.permissions
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
    async handleUpdate(row) {
      await this.getPermissions()
      this.temp.id = row.id
      this.temp.name = row.name
      this.temp.text = row.text
      this.temp.permissions = row.permissions ? row.permissions.map(item => {
        return item.id
      }) : []
      this.dialogStatus = 'update'
      this.dialogFormVisible = true
      this.$nextTick(() => {
        this.$refs['dataForm'].clearValidate()
      })
    },
    updateData() {
      this.$refs['dataForm'].validate((valid) => {
        if (valid) {
          update(this.temp.id, {
            name: this.temp.name,
            text: this.temp.text,
            permission_ids: this.temp.permissions
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
    handleDelete(row, index) {
      this.deleteDialogVisible = true
      this.current = row
    },
    getSortClass: function(key) {
      return this.listQuery.sort === 'asc' ? 'ascending' : 'descending'
    }
  }
}
</script>
