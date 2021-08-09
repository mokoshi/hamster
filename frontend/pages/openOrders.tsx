import {
  Spinner,
  Table,
  Thead,
  Tbody,
  Tr,
  Th,
  Td,
  TableCaption,
  AlertIcon,
  Alert,
} from '@chakra-ui/react';
import { useOpenOrdersQuery } from '../model/order';

export default function OpenOrders() {
  const { data: orders, isLoading } = useOpenOrdersQuery();

  if (isLoading || !orders) {
    return <Spinner />;
  }

  if (orders.length === 0) {
    return (
      <Alert status='info'>
        <AlertIcon />
        注文がありません
      </Alert>
    );
  }

  return (
    <Table>
      <TableCaption>未約定の注文一覧</TableCaption>
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
}
