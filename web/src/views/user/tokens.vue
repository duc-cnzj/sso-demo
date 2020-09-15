<template>
  <div class="app-container">
    <div class="filter-container">
      <el-input v-model="listQuery.name" placeholder="name" style="width: 200px;margin-right: 10px;" class="filter-item" @keyup.enter.native="handleFilter" />
      <el-button v-waves class="filter-item" type="primary" icon="el-icon-search" @click="handleFilter">
        搜索
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
    >
      <el-table-column label="ID" prop="id" sortable="custom" align="center" width="80" :class-name="getSortClass('id')">
        <template slot-scope="{row}">
          <span>{{ row.id }}</span>
        </template>
      </el-table-column>
      <el-table-column label="用户名" min-width="120px">
        <template slot-scope="{row}">
          <span>{{ row.user.user_name }}</span>
        </template>
      </el-table-column>
      <el-table-column label="token" min-width="120px">
        <template slot-scope="{row}">
          <span style="cursor: pointer;" v-text="row.api_token" />
        </template>
      </el-table-column>
      <el-table-column label="最后使用时间" min-width="160px">
        <template slot-scope="{row}">
          <span>{{ row.last_use_at | formatDate }}</span>
        </template>
      </el-table-column>
      <el-table-column label="创建时间" min-width="160px">
        <template slot-scope="{row}">
          <span>{{ row.created_at | formatDate }}</span>
        </template>
      </el-table-column>
      <!-- <el-table-column label="操作" align="center" width="380" class-name="small-padding fixed-width">
        <template slot-scope="{row,$index}">
          <el-button size="mini" type="danger" @click="handleLogout(row,$index)">
            销毁
          </el-button>
        </template>
      </el-table-column> -->
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
  </div>
</template>

<script>
import { index as tokenindex } from '@/api/token'
import waves from '@/directive/waves'

export default {
  name: 'Role',
  directives: { waves },
  data() {
    return {
      list: null,
      //   deleteDialogVisible: false,
      total: 0,
      current: null,
      listLoading: true,
      listQuery: {
        page: 1,
        pageSize: 15
      }
    }
  },
  created() {
    this.getList()
  },
  methods: {
    doDelete() {
    //   deleteToken(this.current.id).then(res => {
    //     this.$notify({
    //       title: '删除成功',
    //       type: 'success',
    //       duration: 2000
    //     })
    //     this.list.splice(index, 1)
    //   })
    //   this.deleteDialogVisible = false
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

      tokenindex({
        page_size: this.listQuery.pageSize,
        page: this.listQuery.page
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

<style lang="scss" scoped>
  .el-tag:hover {
    cursor: pointer;
  }
</style>
