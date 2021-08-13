import apiClient from '../lib/apiClient';
import { useQuery } from 'react-query';
import dayjs from 'dayjs';

export class OrderBooksMovingAverage {
  unix: number;

  constructor(
    public id: number,
    public time: Date,
    public duration: number,
    public askPrice: number,
    public bidPrice: number,
  ) {
    this.unix = dayjs(time).valueOf();
  }
}

export function useOrderBookMovingAveragesQuery(from: number, to: number) {
  const path = '/order_books/moving_averages';
  const key = `${path}:${from}-${to}`;
  return useQuery<OrderBooksMovingAverage[]>(
    key,
    async () => {
      const { data } = await apiClient.get(path, {
        params: { from, to },
      });

      return data.map(
        (h: any) => new OrderBooksMovingAverage(h.id, h.time, h.duration, h.askPrice, h.bidPrice),
      );
    },
    { keepPreviousData: true },
  );
}
