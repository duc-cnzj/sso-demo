import request from '@/utils/request'

export const index = (params) => {
  return request({
    method: 'get',
    url: '/api/permissions',
    params
  })
}

export const store = (data) => {
  return request({
    method: 'post',
    url: `/api/permissions`,
    data
  })
}

export const show = (id) => {
  return request({
    method: 'get',
    url: `/api/permissions/${id}`
  })
}

export const update = (id, data) => {
  return request({
    method: 'put',
    url: `/api/permissions/${id}`,
    data
  })
}

export const destroy = (id) => {
  return request({
    method: 'delete',
    url: `/api/permissions/${id}`
  })
}

export const getByGroups = () => {
  return request({
    method: 'get',
    url: `/api/permissions_by_group`
  })
}

export const getProjects = () => {
  return request({
    method: 'get',
    url: `/api/get_permission_projects`
  })
}
