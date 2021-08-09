import apiClient from '../lib/apiClient';
import { useQuery } from 'react-query';

export class Order {
  id: number;
  orderType: string;
  rate: number;
  pair: string;
  pendingAmount: number;
  pendingMarketBuyAmount: number;
  stopLossRate: number;
  createdAt: Date;

  constructor(res: {
    id: number;
    orderType: string;
    rate: number;
    pair: string;
    pendingAmount: number;
    pendingMarketBuyAmount: number;
    stopLossRate: number;
    createdAt: Date;
  }) {
    this.id = res.id;
    this.orderType = res.orderType;
    this.rate = res.rate;
    this.pair = res.pair;
    this.pendingAmount = res.pendingAmount;
    this.pendingMarketBuyAmount = res.pendingMarketBuyAmount;
    this.stopLossRate = res.stopLossRate;
    this.createdAt = res.createdAt;
  }
}

export function useOpenOrdersQuery() {
  return useQuery<Order[]>('', async () => {
    const { data } = await apiClient.get('/orders/open');

    return data.map((h: any) => new Order(h));
  });
}
