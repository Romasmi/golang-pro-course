import http from 'k6/http';
import { check } from 'k6';

export const options = {
  scenarios: {
    constant_request_rate: {
      executor: 'constant-arrival-rate',
      rate: 1000,
      timeUnit: '1s',
      duration: '1m',
      preAllocatedVUs: 50,
      maxVUs: 200,
    },
  },
};

const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080';

export default function () {
  const userId = `user-${Math.floor(Math.random() * 1000)}`;
  const eventId = `event-${Math.random().toString(36).substring(2, 11)}`;
  const title = `Load Test Event ${Math.random().toString(36).substring(2, 7)}`;
  const description = `Description ${Math.random().toString(36).substring(2, 15)}`;

  const now = new Date();
  const startAt = new Date(now.getTime() + Math.floor(Math.random() * 1000000)).toISOString();
  const endAt = new Date(new Date(startAt).getTime() + 3600000).toISOString();

  const createPayload = JSON.stringify({
    id: eventId,
    title: title,
    description: description,
    start_at: startAt,
    end_at: endAt,
    user_id: userId,
  });
  const createRes = http.post(`${BASE_URL}/events`, createPayload, {
    headers: { 'Content-Type': 'application/json' },
  });
  check(createRes, { 'create status is 200': (r) => r.status === 200 });

  if (createRes.status !== 200) return;

  const getRes = http.get(`${BASE_URL}/events/${eventId}`);
  check(getRes, { 'get status is 200': (r) => r.status === 200 });

  const updatePayload = JSON.stringify({
    id: eventId,
    title: `Updated ${title}`,
    description: `Updated ${description}`,
    start_at: startAt,
    end_at: endAt,
    user_id: userId,
  });
  const updateRes = http.put(`${BASE_URL}/events/${eventId}`, updatePayload, {
    headers: { 'Content-Type': 'application/json' },
  });
  check(updateRes, { 'update status is 200': (r) => r.status === 200 });

  // List
  const listRes = http.get(`${BASE_URL}/events`);
  check(listRes, { 'list status is 200': (r) => r.status === 200 });
}
