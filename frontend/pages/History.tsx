import { Radio, RadioGroup, Spinner, Stack } from '@chakra-ui/react';
import React, { useEffect, useMemo, useState } from 'react';
import Card from '../components/Card';
import dayjs from 'dayjs';
import Chart from '../components/Chart';
import { useOrderBookMovingAveragesQuery } from '../model/orderBooksMovingAverage';
import { useOrderBooksSnapshotsQuery } from '../model/orderBooksSnapshot';

const History: React.FC = () => {
  return (
    <Card title='履歴'>
      <OrderBooksSnapshotContent />
    </Card>
  );
};

const OrderBooksSnapshotContent: React.FC = () => {
  const [scale, setScale] = useState('0');
  const [now, setNow] = useState(dayjs('2021-08-10 18:00:00').unix());

  const from = useMemo(() => now - 60, [now]);
  const to = now;

  const { data: orderBooksSnapshots, isLoading: historyIsLoading } = useOrderBooksSnapshotsQuery(
    from,
    to,
  );
  const { data: orderBooksMovingAverages, isLoading: averageIsLoading } =
    useOrderBookMovingAveragesQuery(from, to);

  useEffect(() => {
    const timer = setInterval(() => {
      setNow(() => dayjs().unix());
    }, 1000);

    return function cleanup() {
      clearInterval(timer);
    };
  });

  return (
    <div>
      <RadioGroup onChange={setScale} value={scale}>
        <Stack direction='row'>
          <Radio value='0'>5分</Radio>
          <Radio value='1'>1時間</Radio>
          <Radio value='2'>6時間</Radio>
        </Stack>
      </RadioGroup>

      {historyIsLoading && <Spinner />}
      {averageIsLoading && <Spinner />}

      <Chart
        from={from}
        to={to}
        orderBooksSnapshots={orderBooksSnapshots}
        orderBooksMovingAverages={orderBooksMovingAverages}
      />
    </div>
  );
};

export default History;
