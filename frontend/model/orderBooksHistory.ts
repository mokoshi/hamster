import apiClient from '../lib/apiClient';
import { useSuspenseQuery } from '../hooks/useSuspenceQuery';

export class OrderBooksHistory {
  id: number;
  time: Date;
  lowestAskPrice: number;
  lowestAskQuantity: number;
  highestBidPrice: number;
  highestBidQuantity: number;

  constructor(res: {
    id: number;
    time: Date;
    lowestAskPrice: number;
    lowestAskQuantity: number;
    highestBidPrice: number;
    highestBidQuantity: number;
  }) {
    this.id = res.id;
    this.time = res.time;
    this.lowestAskPrice = res.lowestAskPrice;
    this.lowestAskQuantity = res.lowestAskQuantity;
    this.highestBidPrice = res.highestBidPrice;
    this.highestBidQuantity = res.highestBidQuantity;
  }
}

export function useOrderBookHistoriesQuery(from: string, to: string) {
  const key = `${from}-${to}`;
  return useSuspenseQuery<OrderBooksHistory[]>(key, async () => {
    const { data } = await apiClient.get('/order_books_histories', {
      params: { from, to },
    });

    return data.map((h: any) => new OrderBooksHistory(h));
  });
}
