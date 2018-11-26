import { stringify } from 'qs';
import request, { baseURLV1 } from '../utils/request';

const buRouter = 'roles';

export async function query(params) {
  return request(`${baseURLV1}/${buRouter}?${stringify(params)}`);
}

export async function get(params) {
  return request(`${baseURLV1}/${buRouter}/${params.record_id}`);
}

export async function create(params) {
  return request(`${baseURLV1}/${buRouter}`, {
    method: 'POST',
    body: params,
  });
}

export async function update(params) {
  return request(`${baseURLV1}/${buRouter}/${params.record_id}`, {
    method: 'PUT',
    body: params,
  });
}

export async function del(params) {
  return request(`${baseURLV1}/${buRouter}/${params.record_id}`, {
    method: 'DELETE',
  });
}

export async function delMany(params) {
  return request(`${baseURLV1}/${buRouter}?${stringify(params)}`, {
    method: 'DELETE',
  });
}

export async function enable(params) {
  return request(`${baseURLV1}/${buRouter}/${params.record_id}/enable`, {
    method: 'PATCH',
  });
}

export async function disable(params) {
  return request(`${baseURLV1}/${buRouter}/${params.record_id}/disable`, {
    method: 'PATCH',
  });
}
