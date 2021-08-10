import React, { useMemo } from 'react';
import { OrderBooksHistory } from '../model/orderBooksHistory';
import { Line, LineProps, Serie } from '@nivo/line';
import dayjs from 'dayjs';

interface Props {
  from: number;
  to: number;
  orderBooksHistories?: OrderBooksHistory[];
}

const chartProps: Partial<LineProps> = {
  margin: {
    top: 20,
    bottom: 50,
    left: 60,
    right: 30,
  },
  xFormat: (v) => dayjs(v).format('MM-DD HH:mm:ss.SSS'),
  axisBottom: {
    tickValues: 5,
    format: (v) => dayjs(v).format('HH:mm:ss'),
    tickSize: 5,
    tickPadding: 5,
    tickRotation: 0,
  },
  pointColor: { theme: 'background' },
  lineWidth: 2,
  pointSize: 2,
};

const Chart: React.FC<Props> = (props) => {
  const { from, to, orderBooksHistories } = props;

  const asks: Serie = useMemo(() => {
    return {
      id: 'asks',
      data:
        orderBooksHistories?.map((h) => ({
          x: h.unix,
          y: h.lowestAskPrice,
        })) ?? [],
    };
  }, [orderBooksHistories]);

  const bids: Serie = useMemo(() => {
    return {
      id: 'bids',
      data:
        orderBooksHistories?.map((h) => ({
          x: h.unix,
          y: h.highestBidPrice,
        })) ?? [],
    };
  }, [orderBooksHistories]);

  const data: Serie[] = useMemo(() => [asks, bids], [asks, bids]);
  const xScale: LineProps['xScale'] = useMemo(
    () => ({
      type: 'linear',
      min: from * 1000,
      max: to * 1000,
    }),
    [from, to],
  );
  const yScale: LineProps['yScale'] = useMemo(
    () => ({ type: 'linear', min: 4000000, max: 6000000 }),
    [],
  );

  return (
    <Line
      {...chartProps}
      width={600}
      height={400}
      xScale={xScale}
      yScale={yScale}
      curve='linear'
      data={data}
    />
  );
};
export default Chart;
