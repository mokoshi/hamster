import apiClient from '../lib/apiClient';
import { useQuery } from 'react-query';
import dayjs from 'dayjs';

export class OrderBooksHistory {
  unix: number;

  constructor(
    public id: number,
    public time: Date,
    public lowestAskPrice: number,
    public lowestAskQuantity: number,
    public highestBidPrice: number,
    public highestBidQuantity: number,
  ) {
    this.unix = dayjs(time).valueOf();
  }
}

export function useOrderBookHistoriesQuery(from: number, to: number) {
  const path = '/order_books/histories';
  const key = `${path}:${from}-${to}`;
  return useQuery<OrderBooksHistory[]>(
    key,
    async () => {
      const { data } = await apiClient.get(path, {
        params: { from, to },
      });

      return data.map(
        (h: any) =>
          new OrderBooksHistory(
            h.id,
            h.time,
            h.lowestAskPrice,
            h.lowestAskQuantity,
            h.highestBidPrice,
            h.highestBidQuantity,
          ),
      );
    },
    { keepPreviousData: true },
  );
}
