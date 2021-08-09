import { useOrderBookHistoriesQuery } from '../model/orderBooksHistory';
import { Spinner } from '@chakra-ui/react';

export default function History() {
  const { data: orderBooksHistories, isLoading } = useOrderBookHistoriesQuery(
    '2021-08-01',
    '2021-08-02',
  );

  if (isLoading || !orderBooksHistories) {
    return <Spinner />;
  }

  return (
    <div>
      {orderBooksHistories.map((h) => (
        <div key={h.id}>{h.time}</div>
      ))}
    </div>
  );
}
