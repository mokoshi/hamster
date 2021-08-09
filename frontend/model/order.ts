import apiClient from '../lib/apiClient';
import { useQuery } from 'react-query';

export class Order {
  constructor(
    public id: number,
    public orderType: string,
    public rate: number,
    public pair: string,
    public pendingAmount: number,
    public pendingMarketBuyAmount: number,
    public stopLossRate: number,
    public createdAt: Date,
  ) {}
}

export function useOpenOrdersQuery() {
  const path = '/exchange/open_orders';
  return useQuery<Order[]>(path, async () => {
    const { data } = await apiClient.get(path);

    return data.map(
      (r: any) =>
        new Order(
          r.id,
          r.orderType,
          r.rate,
          r.pair,
          r.pendingAmount,
          r.pendingMarketBuyAmount,
          r.stopLossRate,
          r.createdAt,
        ),
    );
  });
}
