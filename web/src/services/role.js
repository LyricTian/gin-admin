import { stringify } from 'qs';
import request, { v1API } from '../utils/request';

const buRouter = 'roles';

export async function query(params) {
  return request(`${v1API}/${buRouter}?${stringify(params)}`);
}

export async function get(params) {
  return request(`${v1API}/${buRouter}/${params.record_id}`);
}

export async function create(params) {
  return request(`${v1API}/${buRouter}`, {
    method: 'POST',
    body: params,
  });
}

export async function update(params) {
  return request(`${v1API}/${buRouter}/${params.record_id}`, {
    method: 'PUT',
    body: params,
  });
}

export async function del(params) {
  return request(`${v1API}/${buRouter}/${params.record_id}`, {
    method: 'DELETE',
  });
}

export async function delMany(params) {
  return request(`${v1API}/${buRouter}?${stringify(params)}`, {
    method: 'DELETE',
  });
}

export async function enable(params) {
  return request(`${v1API}/${buRouter}/${params.record_id}/enable`, {
    method: 'PATCH',
  });
}

export async function disable(params) {
  return request(`${v1API}/${buRouter}/${params.record_id}/disable`, {
    method: 'PATCH',
  });
}
