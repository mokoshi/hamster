import apiClient from '../lib/apiClient';
import { useQuery } from 'react-query';
import dayjs from 'dayjs';

export class RateHistory {
  unix: number;

  constructor(
    public id: number,
    public orderType: string,
    public pair: string,
    public rate: number,
    public createdAt: number,
  ) {
    this.unix = dayjs(createdAt).valueOf();
  }
}

export function useRateHistoriesQuery(from: number, to: number) {
  const path = '/exchange/rate_histories';
  const key = `${path}:${from}-${to}`;
  return useQuery<RateHistory[]>(
    key,
    async () => {
      const { data } = await apiClient.get(path, {
        params: { from, to },
      });

      return data.map((h: any) => new RateHistory(h.id, h.orderType, h.pair, h.rate, h.createdAt));
    },
    { keepPreviousData: true },
  );
}
