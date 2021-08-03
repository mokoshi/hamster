import { useOrderBookHistoriesQuery } from '../model/orderBooksHistory';

export default function History() {
  const { data: orderBooksHistories } = useOrderBookHistoriesQuery('2021-08-01', '2021-08-02');

  if (!orderBooksHistories) {
    return 'loading';
  }

  return (
    <div>
      {orderBooksHistories.map((h) => (
        <div key={h.id}>{h.time}</div>
      ))}
    </div>
  );
}
