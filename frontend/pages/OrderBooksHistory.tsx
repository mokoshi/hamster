import { useOrderBookHistoriesQuery } from '../model/orderBooksHistory';
import { Alert, AlertIcon, Radio, RadioGroup, Skeleton, Stack } from '@chakra-ui/react';
import React, { useState } from 'react';
import Card from '../components/Card';
import { Line } from '@nivo/line';
import dayjs from 'dayjs';

const OrderBooksHistory: React.FC = () => {
  return (
    <Card title='履歴'>
      <OrderBooksHistoryContent />
    </Card>
  );
};

const now = dayjs('2021-08-09 22:00:00');

const OrderBooksHistoryContent: React.FC = () => {
  const [scale, setScale] = useState('0');
  const from = now.subtract(1, 'minute').unix();
  const to = now.unix();

  const { data: orderBooksHistories, isLoading, isError } = useOrderBookHistoriesQuery(from, to);

  if (isError) {
    return (
      <Alert status='error'>
        <AlertIcon />
        エラーが発生しました
      </Alert>
    );
  }
  if (isLoading || !orderBooksHistories) {
    return <Skeleton height={40} />;
  }

  const lowestAsks = [];
  const highestBids = [];
  let min = 100000000;
  let max = 0;
  for (const h of orderBooksHistories) {
    lowestAsks.push({
      x: h.unix,
      y: h.lowestAskPrice,
    });
    highestBids.push({
      x: h.unix,
      y: h.highestBidPrice,
    });
    if (max < h.lowestAskPrice) {
      max = h.lowestAskPrice;
    }
    if (min > h.highestBidPrice) {
      min = h.highestBidPrice;
    }
  }

  return (
    <div>
      <RadioGroup onChange={setScale} value={scale}>
        <Stack direction='row'>
          <Radio value='0'>5分</Radio>
          <Radio value='1'>1時間</Radio>
          <Radio value='2'>6時間</Radio>
        </Stack>
      </RadioGroup>

      <Line
        width={600}
        height={400}
        margin={{ bottom: 50, left: 60, right: 30 }}
        xScale={{
          type: 'linear',
          min: from * 1000,
          max: to * 1000,
        }}
        yScale={{ type: 'linear', min: min - 1000, max: max + 1000 }}
        xFormat={(v) => dayjs(v).format('MM-DD HH:mm:ss.SSS')}
        pointColor={{ theme: 'background' }}
        curve='linear'
        axisBottom={{
          tickValues: 5,
          format: (v) => dayjs(v).format('HH:mm:ss'),
          tickSize: 5,
          tickPadding: 5,
          tickRotation: 0,
        }}
        animate={false}
        lineWidth={2}
        pointSize={2}
        useMesh={true}
        data={[
          {
            id: 'lowestAsks',
            data: lowestAsks,
          },
          {
            id: 'highestBids',
            data: highestBids,
          },
        ]}
      />
    </div>
  );
};

export default OrderBooksHistory;
