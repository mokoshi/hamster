import { useQuery, UseQueryResult } from 'react-query';

type QueryResult<TQueryFnData = unknown, TError = unknown, TData = TQueryFnData> = UseQueryResult<
  TData,
  TError
> & {
  data: TData;
};

export function useSuspenseQuery<TQueryFnData = unknown, TError = unknown, TData = TQueryFnData>(
  ...args: Parameters<typeof useQuery>
): QueryResult<TQueryFnData, TError, TData> {
  return useQuery(...args) as QueryResult<TQueryFnData, TError, TData>;
}
