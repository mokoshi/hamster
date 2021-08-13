import React, { useMemo } from 'react';
import { Line, LineProps, Serie } from '@nivo/line';
import dayjs from 'dayjs';
import { OrderBooksHistory } from '../model/orderBooksHistory';
import { OrderBooksMovingAverage } from '../model/orderBooksMovingAverage';

interface Props {
  from: number;
  to: number;
  orderBooksHistories?: OrderBooksHistory[];
  orderBooksMovingAverages?: OrderBooksMovingAverage[];
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
  const { from, to, orderBooksHistories, orderBooksMovingAverages } = props;

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

  const mvAsks: Serie = useMemo(() => {
    return {
      id: 'mvAsks',
      data:
        orderBooksMovingAverages?.map((a) => ({
          x: a.unix,
          y: a.askPrice,
        })) ?? [],
    };
  }, [orderBooksMovingAverages]);

  const mvBids: Serie = useMemo(() => {
    return {
      id: 'mvBids',
      data:
        orderBooksMovingAverages?.map((a) => ({
          x: a.unix,
          y: a.bidPrice,
        })) ?? [],
    };
  }, [orderBooksMovingAverages]);

  const data: Serie[] = useMemo(() => [asks, bids, mvAsks, mvBids], [asks, bids, mvAsks, mvBids]);
  const xScale: LineProps['xScale'] = useMemo(
    () => ({
      type: 'linear',
      min: from * 1000,
      max: to * 1000,
    }),
    [from, to],
  );
  const yScale: LineProps['yScale'] = useMemo(() => {
    let min = Number.MAX_SAFE_INTEGER;
    let max = 0;
    for (let h of orderBooksHistories ?? []) {
      if (max < h.lowestAskPrice) {
        max = h.lowestAskPrice;
      }
      if (min > h.highestBidPrice) {
        min = h.highestBidPrice;
      }
    }
    return { type: 'linear', min: min - 1000, max: max + 1000 };
  }, [orderBooksHistories]);

  return (
    <Line
      {...chartProps}
      width={600}
      height={400}
      xScale={xScale}
      yScale={yScale}
      curve='linear'
      data={data}
      animate={false}
      useMesh
      isInteractive
    />
  );
};
export default Chart;
