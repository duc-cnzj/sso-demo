import request from '@/utils/request'

export const index = (params) => {
  return request({
    method: 'get',
    url: '/api/admin/permissions',
    params
  })
}

export const store = (data) => {
  return request({
    method: 'post',
    url: `/api/admin/permissions`,
    data
  })
}

export const show = (id) => {
  return request({
    method: 'get',
    url: `/api/admin/permissions/${id}`
  })
}

export const update = (id, data) => {
  return request({
    method: 'put',
    url: `/api/admin/permissions/${id}`,
    data
  })
}

export const destroy = (id) => {
  return request({
    method: 'delete',
    url: `/api/admin/permissions/${id}`
  })
}

export const getByGroups = () => {
  return request({
    method: 'get',
    url: `/api/admin/permissions_by_group`
  })
}

export const getProjects = () => {
  return request({
    method: 'get',
    url: `/api/admin/get_permission_projects`
  })
}
