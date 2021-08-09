import { Table, Thead, Tbody, Tr, Th, Td, AlertIcon, Alert, Skeleton } from '@chakra-ui/react';
import { useOpenOrdersQuery } from '../model/order';
import React from 'react';
import Card from '../components/Card';

const OpenOrders: React.FC = () => {
  return (
    <Card title='未約定の注文一覧'>
      <OpenOrdersContent />
    </Card>
  );
};

const OpenOrdersContent: React.FC = () => {
  const { data: orders, isLoading, isError } = useOpenOrdersQuery();

  if (isError) {
    return (
      <Alert status='error'>
        <AlertIcon />
        エラーが発生しました
      </Alert>
    );
  }
  if (isLoading || !orders) {
    return <Skeleton height={40} />;
  }

  return orders.length === 0 ? (
    <Alert status='info'>
      <AlertIcon />
      注文がありません
    </Alert>
  ) : (
    <Table>
      <Thead>
        <Tr>
          <Th>id</Th>
          <Th>orderType</Th>
          <Th>rate</Th>
          <Th>pair</Th>
          <Th>pendingAmount</Th>
          <Th>pendingMarketBuyAmount</Th>
          <Th>stopLossRate</Th>
          <Th>createdAt</Th>
        </Tr>
      </Thead>
      <Tbody>
        {orders.map((o) => (
          <Tr key={o.id}>
            <Td>{o.id}</Td>
            <Td>{o.orderType}</Td>
            <Td>{o.rate}</Td>
            <Td>{o.pair}</Td>
            <Td>{o.pendingAmount}</Td>
            <Td>{o.pendingMarketBuyAmount}</Td>
            <Td>{o.stopLossRate}</Td>
            <Td>{o.createdAt}</Td>
          </Tr>
        ))}
      </Tbody>
    </Table>
  );
};

export default OpenOrders;
