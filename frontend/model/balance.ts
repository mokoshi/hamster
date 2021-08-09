import apiClient from '../lib/apiClient';
import { useQuery } from 'react-query';

export class Balance {
  constructor(
    public jpy: number,
    public btc: number,
    public jpyReserved: number,
    public btcReserved: number,
    public jpyLendInUse: number,
    public btcLendInUse: number,
    public jpyLent: number,
    public btcLent: number,
    public jpyDebt: number,
    public btcDebt: number,
  ) {}
}

export function useBalanceQuery() {
  const path = '/account/balance';
  return useQuery<Balance>(path, async () => {
    const { data } = await apiClient.get(path);

    return new Balance(
      data.jpy,
      data.btc,
      data.jpyReserved,
      data.btcReserved,
      data.jpyLendInUse,
      data.btcLendInUse,
      data.jpyLent,
      data.btcLent,
      data.jpyDebt,
      data.btcDebt,
    );
  });
}
